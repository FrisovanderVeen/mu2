package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/fvdveen/mu2/bot"
	"github.com/fvdveen/mu2/commands"
	"github.com/fvdveen/mu2/config"
	"github.com/fvdveen/mu2/db"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	// register all commands
	_ "github.com/fvdveen/mu2/commands/all"

	// register all dbs
	_ "github.com/fvdveen/mu2/db/all"
)

var (
	cfgFile string
	conf    config.Config
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "mu2",
	Short: "A discord music bot",
	Long: `Mu2 is a discord music bot.

To configure the bot either use environment variables or use a config file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := db.Get(conf.Database)
		if err != nil {
			return fmt.Errorf("create db: %v", err)
		}

		b, err := bot.New(bot.WithConfig(conf.Discord), bot.WithDB(store))
		if err != nil {
			return fmt.Errorf("create bot: %v", err)
		}

		if err := b.AddCommand(commands.All()...); err != nil {
			return fmt.Errorf("add commands: %v", err)
		}

		if err := b.Open(); err != nil {
			return fmt.Errorf("open session: %v", err)
		}

		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		logrus.Info("Bot is now running press CRTL-C to exit")
		<-sc

		if err := b.Close(); err != nil {
			return fmt.Errorf("close session: %v", err)
		}
		return nil
	},
	SilenceUsage: true,
}

// Execute runs the cli
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("MU2")

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file location")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err != nil {
			logrus.Errorf("Reading in config: %v", err)
		}
	}

	defaults := map[string]interface{}{
		"log": map[string]interface{}{
			"level": "",
			"discord": map[string]interface{}{
				"level":   "",
				"webhook": "",
			},
		},
		"discord": map[string]interface{}{
			"token":  "",
			"prefix": "$",
		},
		"database": map[string]interface{}{
			"host":     "",
			"user":     "",
			"password": "",
			"ssl":      "",
			"type":     "postgres",
		},
	}

	for k, v := range defaults {
		viper.SetDefault(k, v)
	}

	viper.AutomaticEnv() // read in environment variables that match

	for _, key := range viper.AllKeys() {
		val := viper.Get(key)
		viper.Set(key, val)
	}

	if err := viper.Unmarshal(&conf); err != nil {
		logrus.Fatalf("Unmarshalling config: %v", err)
		return
	}

	loadLogger(conf.Log, logrus.StandardLogger())
	logrus.Debugf("Using config file: %s", viper.ConfigFileUsed())
}
