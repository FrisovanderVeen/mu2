package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	logging "github.com/op/go-logging"
	"github.com/spf13/cobra"

	"github.com/fvdveen/bf"
	"github.com/fvdveen/mu2/commands"
	_ "github.com/fvdveen/mu2/commands/help"
	_ "github.com/fvdveen/mu2/commands/info"
	_ "github.com/fvdveen/mu2/commands/pingpong"
	_ "github.com/fvdveen/mu2/commands/sound"
)

var runCmd = &cobra.Command{
	Use:   "run [flags]",
	Short: "Runs the bot",
	Long:  "Run runs the bot",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runBot(token, prefix); err != nil {
			log.Critical("Could not run bot: %v", err)
		}
	},
}

var (
	token  string
	prefix string
)

func init() {
	runCmd.Flags().StringVar(&token, "token", "DGTOKEN", "The environment variable containing the discord token")
	runCmd.Flags().StringVar(&prefix, "prefix", "DGPREFIX", "The environment variable containing the discord prefix")
	rootCmd.AddCommand(runCmd)
}

func runBot(tokenEnv, prefixEnv string) error {
	logging.SetBackend(logging.NewBackendFormatter(logging.NewLogBackend(os.Stderr, "", 0), logging.MustStringFormatter(`%{color}%{time:15:04:05.000} %{module} ▶ %{level:.4s} %{id:03x} %{message} %{color:reset}`)))

	bot, err := bf.NewBot(getEnvVars(tokenEnv, prefixEnv))
	if err != nil {
		log.Critical("could not make bot: %v", err)
		return err
	}

	bot.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		s.UpdateStatus(0, fmt.Sprintf("%shelp", bot.Prefix))
	})

	if err := bot.AddCommand(commands.Commands...); err != nil {
		log.Errorf("could not add command: %v", err)
		return err
	}

	if err := bot.Open(); err != nil {
		log.Critical("Could not open session: %v", err)
		return err
	}

	log.Info("Logged in as")
	log.Info(bot.Session.State.User.Username)
	log.Info(bot.Session.State.User.ID)
	log.Info("Bot is now running. Press CTRL-C to exit.")
	log.Info("-------------------")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	if err := bot.Close(); err != nil {
		log.Errorf("could not close session: %v", err)
		return err
	}

	return nil
}

func getEnvVars(tokenEnv string, prefixEnv string) bf.OptionFunc {
	token := os.Getenv(tokenEnv)
	prefix := os.Getenv(prefixEnv)
	return func(b *bf.Bot) error {
		b.Token = token
		b.Prefix = prefix

		return nil
	}
}
