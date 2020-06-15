package cli

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewListCmd(configPath string) *cobra.Command {
	var searchCmd = &cobra.Command{
		Use:   "list",
		Short: "list details from zot",
		Long:  `list details from zot`,
	}

	searchCmd.AddCommand(NewCveCommand(NewCveSearchService(), configPath))

	return searchCmd
}

func setupConfig(configPath string) error {
	file, err := os.OpenFile(configPath, os.O_RDONLY|os.O_CREATE, 0644)

	if err != nil {
		return err
	}

	file.Close()
	viper.SetConfigFile(configPath)
	viper.SetConfigType("properties")
	viper.SetDefault("showSpinner", true)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
