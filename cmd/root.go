package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/fvdveen/mu2/bot"
	"github.com/fvdveen/mu2/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var conf config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "backend",
	Short: "Mu2 is a Discord music bot",
	Run: func(cmd *cobra.Command, args []string) {
		b, err := bot.New(&conf.Bot)
		if err != nil {
			logrus.Fatal(err)
			return
		}

		err = b.Open()
		if err != nil {
			logrus.Fatal(err)
			return
		}
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		logrus.Info("Bot is now running press CRTL-C to exit")
		<-sc
		err = b.Close()
		if err != nil {
			logrus.Fatal(err)
			return
		}
	},
}

// Execute runs the CLI
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")
}

// initConfi reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in current directory with name "config" (without extension).
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Error(err)
		return
	}

	// Umarshal the config
	err = viper.Unmarshal(&conf)
	if err != nil {
		logrus.Error(err)
		return
	}

	logrus.Infof("Using config file: %s", viper.ConfigFileUsed())
}
