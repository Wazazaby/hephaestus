package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/vrischmann/envconfig"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	var conf Config
	if err := envconfig.Init(&conf); err != nil {
		logger.Fatal("unable to parse config", zap.Error(err))
	}

	logger.Info("parsed config from env", zap.String("token", conf.Token))

	session, _ := discordgo.New("Bot " + conf.Token)

	// status handler
	session.AddHandler(func(_ *discordgo.Session, _ *discordgo.Ready) {
		logger.Info("hephaestus is up and running")
	})

	// will respond to "ayo"
	session.AddHandler(ayoMessageCreate)

	if err := session.Open(); err != nil {
		logger.Fatal("unable to open connection", zap.Error(err))
	}
	defer session.Close()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	logger.Info("received close signal, stopping")
}

// for status purpose
func ayoMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content != "ayo" {
		return
	}

	s.ChannelMessageSend(m.ChannelID, "ayo")
}
