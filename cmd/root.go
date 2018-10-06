package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/fvdveen/mu2-config/consul"
	"github.com/fvdveen/mu2/commands"
	"github.com/fvdveen/mu2-config/events"
	"github.com/fvdveen/mu2/log"
	"github.com/fvdveen/mu2/watch"
	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/fvdveen/mu2/services/search"
	"github.com/fvdveen/mu2/commands/play"
	"github.com/spf13/viper"

	// register all commands
	_ "github.com/fvdveen/mu2/commands/all"

	// register all dbs
	_ "github.com/fvdveen/mu2/db/all"
)

var (
	logLvl string
	conf   struct {
		Consul struct {
			Address string `mapstructure:"address"`
		} `mapstructure:"consul"`
		Log struct {
			Level string `mapstructure:"level"`
		} `mapstructure:"log"`
		Config struct {
			Path string `mapstructure:"path"`
			Type string `mapstructure:"type"`
		} `mapstructure:"config"`
	}
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "mu2",
	Short: "A discord music bot",
	Long:  `Mu2 is a discord music bot.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logrus.WithField("type", "main").Debug("Starting mu2...")

		cc := api.DefaultConfig()
		if conf.Consul.Address != "" {
			cc.Address = conf.Consul.Address
		}

		c, err := api.NewClient(cc)
		if err != nil {
			return fmt.Errorf("create consul client: %v", err)
		}

		p, err := consul.NewProvider(c, conf.Config.Path, conf.Config.Type, nil)
		if err != nil {
			return fmt.Errorf("create provider: %v", err)
		}

		var ch, b, l, db, ssch <-chan *events.Event

		ch = events.Watch(p.Watch())
		logrus.WithField("type", "main").Debug("Created config provider")
		b, ch = events.Bot(ch)
		l, ch = events.Log(ch)
		db, ch = events.Database(ch)
		ssch, ch = events.SearchService(ch)
		events.Null(ch)

		var wg sync.WaitGroup
		
		wg.Add(1)
		ld := watch.Log(logrus.StandardLogger(), l, &wg)
		logrus.WithField("type", "main").Debug("Created log watcher")
		
		wg.Add(1)
		s, dbd := watch.DB(db, &wg)
		logrus.WithField("type", "main").Debug("Created db watcher")

		ss, ssd := watch.SearchService(ssch, cc)
		go addCommands(ss)

		wg.Add(1)
		check := make(chan interface{})
		bd, _ := watch.Bot(b, check, s, &wg) 
		logrus.WithField("type", "main").Debug("Created bot watcher")

		wg.Wait()
		logrus.WithField("type", "main").Debug("Created watchers")

		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		logrus.WithField("type", "main").Info("Bot is now running press CRTL-C to exit")

		<-sc
		logrus.WithField("type", "main").Info("Shutting down...")
		p.Close()
		<-ld
		<-dbd
		<-bd
		<-ssd
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

	rootCmd.PersistentFlags().StringVar(&conf.Log.Level, "log-level", "", "log level")
	rootCmd.PersistentFlags().StringVar(&conf.Consul.Address, "consul-addr", "", "consul address")
	rootCmd.PersistentFlags().StringVar(&conf.Config.Path, "config-path", "bot/config", "config path on the kv store")
	rootCmd.PersistentFlags().StringVar(&conf.Config.Type, "config-type", "json", "config type on the kv store")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	defaults := map[string]interface{}{
		"log": map[string]interface{}{
			"level": logLvl,
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
		logrus.WithField("type", "main").Fatalf("Unmarshalling config: %v", err)
		return
	}

	var lvl logrus.Level

	if conf.Log.Level != "" {
		lvl = log.GetLevel(conf.Log.Level)
	} else if viper.IsSet("log.level") {
		lvl = log.GetLevel(viper.GetString("log.level"))
	} else {
		lvl = logrus.InfoLevel
	}

	logrus.SetLevel(lvl)
}

func addCommands(ss search.Service) {
	pc := play.New(ss)
	commands.Register(pc)
}