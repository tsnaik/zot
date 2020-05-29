package cveinfo

import (
	"time"

	"github.com/anuvu/zot/pkg/log"
	integration "github.com/aquasecurity/trivy/integration"
	config "github.com/aquasecurity/trivy/integration/config"
)

// UpdateCVEDb ...
func UpdateCVEDb(dbDir string, log log.Logger, interval time.Duration, isTest bool) error {
	config, err := config.NewDbConfig(dbDir)
	if err != nil {
		log.Error().Err(err).Msg("Unable to get config")
	}

	if isTest {
		err = integration.RunDb(config)
		if err != nil {
			log.Error().Err(err).Msg("Unable to update DB ")
		}

		return err
	}

	for {
		err = integration.RunDb(config)
		if err != nil {
			log.Error().Err(err).Msg("Unable to update DB ")
			return err
		}

		time.Sleep(interval * time.Hour)
	}
}
