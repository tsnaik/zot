package accesscontrol

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/anuvu/zot/pkg/log"
	jsoniter "github.com/json-iterator/go"
	strcas "github.com/qiangmzsx/string-adapter/v2"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
)

func IsAuthorized(username, requestMethod, requestURI, configPath string, logger log.Logger) bool {
	config, err := readConfig(configPath, logger)

	// should not interfere if the path is not specified in the config. For backwards compatibility
	if err != nil && errors.Is(err, errNoPath) {
		return true
	}

	policy := configToPolicy(config)
	policyAdapter := strcas.NewAdapter(policy)
	m := model.Model{}
	err = m.LoadModelFromText(casbinModel)

	if err != nil {
		logger.Panic().Err(err).Msg("Error parsing access control model: " + err.Error())
	}

	e, err := casbin.NewEnforcer(m, policyAdapter)

	if err != nil {
		logger.Panic().Err(err).Msg("Error parsing access control policy: " + err.Error())
	}

	ok, err := e.Enforce(username, requestURI, requestMethod)

	if err != nil {
		logger.Panic().Err(err).Msg("Error enforcing access control: " + err.Error())
	}

	return ok
}

type accessConfig struct {
	Repositories []struct {
		Repository string `json:"name"`
		Users      []struct {
			Username  string `json:"username"`
			AllowAPIs []struct {
				Path   string `json:"path"`
				Method string `json:"method"`
			} `json:"allowAPIs"`
		} `json:"users"`
	} `json:"repositories"`
}

func readConfig(configPath string, logger log.Logger) (accessConfig, error) {
	config := accessConfig{}

	if configPath == "" {
		return config, errNoPath
	}

	f, err := ioutil.ReadFile(configPath)
	if err != nil {
		logger.Panic().Err(err).Msg("Error parsing access control config file: " + err.Error())
	}

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(f, &config)

	if err != nil {
		logger.Panic().Err(err).Msg("Error parsing access control config file: " + err.Error())
	}

	return config, nil
}

func configToPolicy(config accessConfig) string {
	var builder strings.Builder

	for _, repo := range config.Repositories {
		for _, user := range repo.Users {
			for _, allow := range user.AllowAPIs {
				processedPath := strings.ReplaceAll(allow.Path, "{name}", repo.Repository)
				builder.WriteString(fmt.Sprintf("p, %s, %s, %s\n", user.Username, processedPath, allow.Method))
			}
		}
	}

	return builder.String()
}

const (
	casbinModel = `
	[request_definition]
	r = sub, obj, act
	
	[policy_definition]
	p = sub, obj, act
	
	[policy_effect]
	e = some(where (p.eft == allow))
	
	[matchers]
	m = r.sub == p.sub && keyMatch3(r.obj, p.obj) && regexMatch(r.act, p.act)
`
)

var (
	errNoPath = errors.New("access-control: config file not specified")
)
