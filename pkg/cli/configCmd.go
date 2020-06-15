package cli

import (
	"fmt"
	"io/ioutil"
	"strings"

	zotErrors "github.com/anuvu/zot/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewConfigCommand(configPath string) *cobra.Command {
	var isListing bool

	var configCmd = &cobra.Command{
		Use:     "config <parameter> [value]",
		Example: examples,
		Short:   "Configure zot CLI",
		Long:    `Configure default parameters for CLI`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if isListing {
				res, err := getAllConfig(configPath)
				if err != nil {
					return err
				}
				fmt.Fprint(cmd.OutOrStdout(), res)

				return nil
			}

			switch len(args) {
			case 0:
				return zotErrors.ErrInvalidArgs
			case 1:
				res, err := getConfigValue(configPath, args[0])
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), res)

			case 2:
				if err := setConfigValue(configPath, args[0], args[1]); err != nil {
					return err
				}
			default:
				return zotErrors.ErrInvalidArgs
			}

			return nil
		},
	}

	configCmd.Flags().BoolVarP(&isListing, "list", "l", false, "List current configuration")
	configCmd.SetUsageTemplate(configCmd.UsageTemplate() + supportedOptions)

	return configCmd
}

func getConfigValue(configPath, key string) (string, error) {
	err := setupConfig(configPath)
	if err != nil {
		return "", err
	}

	if err := viper.ReadInConfig(); err != nil {
		return "", err
	}

	return viper.GetString(key), nil
}

func setConfigValue(configPath, key, value string) error {
	err := setupConfig(configPath)
	if err != nil {
		return err
	}

	viper.Set(key, value)

	text, err := getAllConfig(configPath)

	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(configPath, []byte(text), 0644); err != nil {
		return err
	}

	return nil
}

func getAllConfig(configPath string) (string, error) {
	err := setupConfig(configPath)
	if err != nil {
		return "", err
	}

	var builder strings.Builder

	values := viper.AllSettings()
	for key, val := range values {
		fmt.Fprintf(&builder, "%s = %v\n", key, val)
	}

	return builder.String(), nil
}

const (
	examples = `  zot config url https://zot-foo.com:8080
  zot config url
  zot config --list`

	supportedOptions = `
Useful parameters:
  url		zot server URL
  showSpinner	show spinner while loading data [true/false]`
)
