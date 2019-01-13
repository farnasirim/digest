package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	gdrive "google.golang.org/api/drive/v3"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	homedir "github.com/mitchellh/go-homedir"

	yaml "gopkg.in/yaml.v2"

	"github.com/farnasirim/digest/diff"
	"github.com/farnasirim/digest/drive"
	"github.com/farnasirim/digest/smtp"
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

func getDigestDir() string {
	return path.Join(getHomeDir(), ".digest")
}

func getConfigDir() string {
	return getDigestDir()
}

func persistConfigs() {
	err := os.MkdirAll(getConfigDir(), 0755)
	if err != nil {
		log.Fatalln(err.Error())
	}

	viperSettings := viper.AllSettings()
	fileName := path.Join(getConfigDir(), defaultConfigName+".yaml")

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

func initConfig() {
	if configFile == "" {
		viper.AddConfigPath(getConfigDir())
		viper.SetConfigName(defaultConfigName)
	} else {
		viper.SetConfigFile(configFile)
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("While reading config file: %s", err.Error())
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	viper.SetDefault("folder", "notes")
	viper.SetDefault("auth-dir", path.Join(getDigestDir(), "auth"))
	viper.SetDefault("data-dir", path.Join(getDigestDir(), "data"))
	viper.SetDefault("smtp-server-host", "smtp.gmail.com")
	viper.SetDefault("smtp-server-port", "587")

	rootCmd.Flags().StringVar(&configFile, "config", "",
		fmt.Sprintf(
			"Path to config file. Will look it up in %s if not specified",
			getConfigDir()))

	rootCmd.Flags().String("folder", "",
		"Name of folder in Google Drive that contains your google docs")
	rootCmd.Flags().String("auth-dir", "",
		"Where to look for/store google oauth2 credentials. Default: $HOME/.digest/auth")
	rootCmd.Flags().String("data-dir", "",
		"Where to keep downloaded google docs. Default: $HOME/.digest/data")
	rootCmd.Flags().String("smtp-server-host", "",
		"SMTP server host address to use for sending emails. Default: smtp.gmail.com")
	rootCmd.Flags().String("smtp-server-port", "",
		"SMTP server port to use for sending emails. Default: 587")
	rootCmd.Flags().String("smtp-user", "",
		`SMTP username to login with. Will also be used as "from" address`)
	rootCmd.Flags().String("smtp-pass", "",
		"SMTP password to login with")
	rootCmd.Flags().String("smtp-to", "",
		"Address to send the diff email to. Defaults to smtp-user")

	rootCmd.Flags().BoolVar(&persistConfs, "persist-confs", false,
		`Overwrite the default config file with config from the current run
WARNING: may write sensitive information e.g. smtp password to file.
Use at own risk.`)

	viper.BindPFlags(rootCmd.Flags())
	viper.SetEnvPrefix("DIG")
	viper.AutomaticEnv()
	rootCmd.Run = digestFunc
}

func digestFunc(cmd *cobra.Command, args []string) {
	fromAddr := viper.GetString("smtp-user")
	password := viper.GetString("smtp-pass")
	smtpHost := viper.GetString("smtp-server-host")
	smtpPort := viper.GetString("smtp-server-port")
	smtpTo := viper.GetString("smtp-to")

	if smtpTo == "" {
		smtpTo = fromAddr
	}

	secretDir := viper.GetString("auth-dir")
	dataDir := viper.GetString("data-dir")

	googleDocsFolder := viper.GetString("folder")

	googleAuthor := drive.NewGoogleAuthenticator(secretDir)
	authClient, err := googleAuthor.GetOrCreateClient()

	if err != nil {
		log.Fatalln(err.Error())
	}

	svc, err := gdrive.New(authClient)
	if err != nil {
		log.Fatalln(err.Error())
	}

	driveService := drive.NewDriveService(dataDir, svc)

	smtpServer := smtp.NewSimpleSMTP(smtpHost, smtpPort, fromAddr, password)

	driveService.TakeAndPersistTimedSnapshot(googleDocsFolder)

	older, newer, err := driveService.LastTwoDirs()
	if err == drive.ErrLessThanTwoSnapshots {
		log.Println("Not enough directories to send the diff")
		log.Println("Exiting successfully")
		return
	}

	differ := diff.NewPlainTextDiff()
	diff := differ.DiffDirsHtml(older, newer)
	if strings.TrimSpace(diff) == "" {
		diff = `
<div style="word-wrap: break-word; width:700px; font-family: monospace;">
<font size=4px color="darkred">
		No new notes! *LOUD GASP* ⊙▃⊙ 
</font>
</div>
		`
	}

	diff = `
<html>

    <body bgcolor="#e6e6fa">

` + diff + `

    </body>
</html>
	`

	now := time.Now()
	subject := fmt.Sprintf("%d-%d-%d", now.Year(), now.Month(), now.Day())
	if err := smtpServer.SendMailHtml(smtpTo, subject, []byte(diff)); err != nil {
		log.Fatalln(err.Error())
		return
	}

	log.Println("Email sent successfully")

	if persistConfs {
		log.Println("Persisting configurations")
		persistConfigs()
	}
}
