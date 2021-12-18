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
	Endpoint          string
	ApplicationKey    string
	ApplicationSecret string
	ConsumerKey       string
	Domain            string
}

var (
	cfgFile       string
	OptFrom       string
	OptTo         string
	Domain        string
	currentConfig config
	OvhClient     *ovh.Client
)

var RootCmd = &cobra.Command{
	Use:   "govh-mrm [list|add|remove] --from=[MAIL] (--to=[MAIL])",
	Short: "OVH Mail redirection manager",
	Long:  `This application manage mail redirection for OVH ¨Provider.`,
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default is $HOME/.config/ovh/ovh.yaml)")
	//RootCmd.PersistentFlags().StringVar(&OptFrom, "from", "", "Name of redirection")
	//RootCmd.PersistentFlags().StringVar(&OptTo, "to", "", "Email of redirection target")

	viper.SetDefault("author", "Julien Vinet <contact@julienvinet.dev>")
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
	Domain = currentConfig.Domain
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
