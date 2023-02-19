package main

import (
	"github.com/bwmarrin/discordgo"
	"god-of-right-go/cmd_handlers"
	"log"
	"os"
	"os/signal"
)

var botToken = os.Getenv("BOT_TOKEN")
var s *discordgo.Session

func init() {
	var err error
	s, err = discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalf("Invalid bot token: %s", err)
	}
}

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "ksm",
		Description: "Kim Sung-mo",
	},
	{
		Name:        "stock",
		Description: "Stock",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "symbol",
				Description: "Stock symbol",
				Required:    true,
			},
		},
	},
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"ksm":   cmd_handlers.KsmHandler,
	"stock": cmd_handlers.StockHandler,
}

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Bot is running as %s", r.User.Username)
	})

	err := s.Open()
	if err != nil {
		log.Fatalf("Failed to open session: %s", err)
	}

	log.Printf("Registering commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, command := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, s.State.Guilds[0].ID, command)
		if err != nil {
			log.Fatalf("Failed to register command: %s", err)
		}
		registeredCommands[i] = cmd
	}

	defer func(s *discordgo.Session) {
		err := s.Close()
		if err != nil {

		}
	}(s)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Bot is running. Press CTRL-C to exit.")
	<-stop
}
