package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/ovh/go-ovh/ovh"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type config struct {
	Endpoint          string   // OVH Endpoint api
	ApplicationKey    string   // Application Key from OVH Api
	ApplicationSecret string   // Application Secret from OVH Api
	ConsumerKey       string   // Consumer Key from OVH Api
	Domain            []string // Your personal Domains
}

var (
	cfgFile       string      // configuration filename
	Domain        string      // Current Domain
	currentConfig config      // current configuration
	OvhClient     *ovh.Client // OVH Client
	fromFlag      string      // Email of redirection
	toFlag        string      // Email of destination
)

var RootCmd = &cobra.Command{
	Use:                   "govh-mrm [list|add|remove] ([redirection mail]) ([destination mail]) (--config=[FILENAME]) (--domain=[DOMAIN])",
	Short:                 "OVH Mail redirection manager",
	Long:                  `This application manage mail redirection for OVH ¨Provider.`,
	DisableFlagsInUseLine: true,
	CompletionOptions: cobra.CompletionOptions{
		DisableDescriptions: true,
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Config file (default is $HOME/.config/ovh/ovh.yaml)")
	RootCmd.PersistentFlags().StringVarP(&Domain, "domain", "d", "", "Domain to use (default is the first in config file)")
	viper.SetDefault("author", "Julien Vinet <julien@vinet.dev>")
	viper.SetDefault("license", "GNU GENERAL PUBLIC LICENSE")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home + "/.config/ovh/")
		viper.AddConfigPath("/etc/ovh/")
		viper.AddConfigPath(".")
		viper.SetConfigName("ovh")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config file not found.")
		} else {
			fmt.Println("Fatal error config file: %w", err)
		}
		os.Exit(1)
	}

	if err := viper.Unmarshal(&currentConfig); err != nil {
		fmt.Println("Error reading configuration : %w", err)
		os.Exit(1)
	}

	OvhClient, _ = ovh.NewClient(
		currentConfig.Endpoint,
		currentConfig.ApplicationKey,
		currentConfig.ApplicationSecret,
		currentConfig.ConsumerKey,
	)

	if len(currentConfig.Domain) == 0 && Domain == "" {
		fmt.Println("Error reading configuration : no domain found !")
		os.Exit(1)
	}
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
