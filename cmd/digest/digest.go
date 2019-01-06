package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	homedir "github.com/mitchellh/go-homedir"

	yaml "gopkg.in/yaml.v2"
)

var (
	googleDriveFolder string
	configFile        string
	persistConfs      bool
)

var (
	defaultConfigName = "config"
)

func getHomeDir() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		log.Fatalln(err.Error())
	}
	return homeDir
}

func getConfigDir() string {
	return path.Join(getHomeDir(), ".digest")
}

func persistConfigs() {
	err := os.MkdirAll(getConfigDir(), 0755)
	if err != nil {
		log.Fatalln(err.Error())
	}

	viperSettings := viper.AllSettings()
	fileName := path.Join(getConfigDir(), defaultConfigName)

	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = yaml.NewEncoder(file).Encode(viperSettings)
	if err != nil {
		log.Fatalf("unable to marshal config to YAML: %v", err)
	}
	log.Println("Current config persisted at: " + fileName)
}

func digestFunc(cmd *cobra.Command, args []string) {
	// TODO:

	if persistConfs {
		persistConfigs()
	}
}

func initConfig() {
	if configFile == "" {
		viper.AddConfigPath(getConfigDir())
		viper.SetConfigName(defaultConfigName)
	} else {
		viper.SetConfigFile(configFile)
	}

	viper.ReadInConfig()
}

func init() {
	cobra.OnInitialize(initConfig)

	viper.SetDefault("folder", "def")

	rootCmd.Flags().StringVar(&configFile, "config", "",
		fmt.Sprintf(
			"Path to config file. Will look it up in %s if not specified",
			getConfigDir()))

	rootCmd.Flags().String("folder", "",
		"Name of Google Drive folder containing your google docs")

	rootCmd.Flags().BoolVar(&persistConfs, "persist-confs", false,
		`Overwrite the default config file with config from the current run
WARNING: may write sensitive information e.g. smtp password to file.
Use at own risk.`)

	viper.BindPFlags(rootCmd.Flags())
	viper.SetEnvPrefix("DIG")
	viper.AutomaticEnv()
	rootCmd.Run = digestFunc
}
