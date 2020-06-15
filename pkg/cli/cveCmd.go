package cli

import (
	"fmt"
	"time"

	zotErrors "github.com/anuvu/zot/errors"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCveCommand(searchService CveSearchService, configPath string) *cobra.Command {
	searchCveParams := make(map[string]*string)

	var servURL string

	var user string

	var cveCmd = &cobra.Command{
		Use:   "cve",
		Short: "Find CVEs",
		Long:  `Find CVEs (Common Vulnerabilities and Exposures)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			spin := spinner.New(spinner.CharSets[39], spinnerDuration, spinner.WithWriter(cmd.ErrOrStderr()))
			spin.Prefix = "Searching... "

			err := setupConfig(configPath)
			if err != nil {
				return err
			}

			if servURL == "" {
				if viper.InConfig("url") {
					servURL = viper.GetString("url")
				} else {
					return zotErrors.ErrInvalidArgs
				}
			}

			if viper.GetBool("showSpinner") {
				spin.Start()
			}

			result, err := searchCve(searchCveParams, searchService, &servURL, &user)

			if viper.GetBool("showSpinner") {
				spin.Stop()
			}

			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), result)

			return nil
		},
	}

	setupCmdFlags(cveCmd, searchCveParams, &servURL, &user)

	return cveCmd
}

func setupCmdFlags(cveCmd *cobra.Command, searchCveParams map[string]*string, servURL *string, user *string) {
	searchCveParams["imageName"] = cveCmd.Flags().StringP("image-name", "I", "", "Find by image name for affected CVEs")
	searchCveParams["cveID"] = cveCmd.Flags().StringP("cve-id", "i", "", "Find images affected by a CVE")

	cveCmd.Flags().StringVar(servURL, "url", "", "Specify zot server URL if not configured")
	cveCmd.Flags().StringVarP(user, "user", "u", "", `User Credentials of zot server in "username:password" format`)
}

func searchCve(params map[string]*string, service CveSearchService, servURL, user *string) (string, error) {
	for _, searcher := range getSearchers() {
		results, err := searcher.search(params, service, servURL, user)
		if err != nil {
			if err == ErrCannotSearch {
				continue
			} else {
				return "", err
			}
		} else {
			return results, nil
		}
	}

	return "", zotErrors.ErrInvalidFlagsCombination
}

const (
	spinnerDuration = 150 * time.Millisecond
)
