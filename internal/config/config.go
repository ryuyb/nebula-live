package config

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var (
	cfgFile       string
	globalConfig  Config
	isInitialized bool
)

func GetConfig() *Config {
	return &globalConfig
}

func InitConfig(rootCmd *cobra.Command) {
	if isInitialized {
		return
	}
	isInitialized = true
	cobra.OnInitialize(loadConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
	bindFlags(&globalConfig, rootCmd, "")
}

func loadConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("NEBULA")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("[LoadConfig] Failed to read configuration file: %v\n", err)
	}
	if err := viper.Unmarshal(&globalConfig); err != nil {
		fmt.Printf("[LoadConfig] Failed to unmarshal configuration: %v\n", err)
	}
}
