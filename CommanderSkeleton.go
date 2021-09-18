package main

import (
	"flag"
	//"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	//"time"
	"syscall"
	"github.com/bwmarrin/discordgo"
)


var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
//	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var s *discordgo.Session

func init() { flag.Parse() }

func init() {
	token, isset := os.LookupEnv("TOKEN")
	if (!isset) {
		log.Fatalf("TOKEN is not set")
	}
	var err error
	s, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

var(
	commands = []*discordgo.ApplicationCommand{
		{
			Name: "basic-command",
			Description: "Basic command",
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"basic-command": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hey there! Congratulations, you just executed your first slash command",
				},
			})
		},
  }
)

func message_handler(s *discordgo.Session, m *discordgo.MessageCreate){
	log.Println(m.Author.ID)
	if m.Author.ID == s.State.User.ID {
        return
  }

	 var tokenized []string = strings.Split(m.Content, " ")
	 // Define Commands here
	 switch tokenized[0] {
	 case "!Hi":
		 s.ChannelMessageSend(m.ChannelID, "Hello! This is a test, and will not be working in the future.")
	 }
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
		log.Println("Bot is up!")
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	s.AddHandler(message_handler)
	//log.Println(len(commands))
	for _, v := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
	}

	defer s.Close()

	sc := make(chan os.Signal, 1)
  signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
  <-sc
	log.Println("Gracefully shutdowning")
}
