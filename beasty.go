package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/go-yaml/yaml"
)

type beasty struct {
	token      string
	Connection *discordgo.Session
	config     map[interface{}]interface{}
	responses  *responses
	storage    struct{}
	localUse   bool
}

// Provides instance of beasty bot with valid configuration
// determined by yml files
func NewBeasty(t string) *beasty {
	configMap := make(map[interface{}]interface{})
	data, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("Error: %s", err)
	} else {
		err := yaml.Unmarshal(data, &configMap)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
	}
	c := configMap

	responseMap := make(map[interface{}]interface{})
	responseLoc, valid := c["responses_file"].(string)
	if !valid {
		log.Fatalf("Error: incorrect interface type in response file")
	} else {
		data, err = ioutil.ReadFile(responseLoc)
		if err != nil {
			log.Fatalf("Error: %s", err)
		} else {
			err := yaml.Unmarshal(data, &responseMap)
			if err != nil {
				log.Fatalf("Error: %s", err)
			}
		}
	}
	responseData := responseMap
	r := newResponses(responseData)

	dgo, err := discordgo.New("Bot " + t)
	if err != nil {
		log.Fatalf("error creating Discord session: %s", err)
	}

	b := &beasty{
		token:      t,
		config:     c,
		responses:  r,
		Connection: dgo,
	}
	return b
}

// randomly fetch a response based on a lookup string, for personality
func (b *beasty) Response(lookup string) string {
	return b.responses.generateResponse(lookup)
}

// enables user prompts and blocking if called with CLI arguments
func (b *beasty) SetLocalUse(f bool) {
	log.Println("Note: set bot to provide local CLI prompt")
	b.localUse = f
}

// configures intent and applies callbacks to discord session
func (b *beasty) bootstrapConnection() {
	b.Connection.AddHandler(b.MessageCallback)
	b.Connection.AddHandler(b.GuildJoinCallback)
	b.Connection.Identify.Intents = discordgo.MakeIntent(
		discordgo.IntentsGuildMessages |
			discordgo.IntentsDirectMessages)
	err := b.Connection.Open()
	if err != nil {
		log.Fatalf("Socket connection error: %s", err)
	}
}

// initial point for handling discord message events
// will call parse to determine action and if valid
func (b *beasty) MessageCallback(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore beasty messages
	if m.Author.ID == s.State.User.ID {
		return
	}
	arguments, forBot := b.MessageParse(m.Content)
	log.Printf("Note: parsed: %s, forBot: %t, from channel: %s and user: %s", arguments, forBot, m.ChannelID, m.Author.ID)
	if !forBot {
		return
	}
	command := arguments[0]
	switch command {
	case "help":
		b.Connection.ChannelMessageSend(m.ChannelID, b.Response("help"))
	case "code":
		// fetch hash code from redis here and check ts
		gid, gidv := b.config["guild_id"].(string)
		rid, ridv := b.config["student_role_id"].(string)
		if !gidv || !ridv {
			log.Printf("Error: failed to get ids for role change")
		}
		b.Connection.GuildMemberRoleAdd(gid, m.Author.ID, rid)
		b.Connection.ChannelMessageSend(m.ChannelID, b.Response("student_role"))
	case "student":
		if len(arguments) == 2 && b.CheckValidEmail(arguments[1]) {
			b.Connection.ChannelMessageSend(m.ChannelID, b.Response("student"))
			// insert hash with ts here
			// email code here
		} else {
			b.Connection.ChannelMessageSend(m.ChannelID, b.Response("email_invalid"))
		}
	case "joke":
		b.Connection.ChannelMessageSend(m.ChannelID, b.Response("joke"))
	default:
		b.Connection.ChannelMessageSend(m.ChannelID, b.Response("unkown"))
	}
}

// callback for handling guild joins, such as a greeting
func (b *beasty) GuildJoinCallback(s *discordgo.Session, m *discordgo.MessageCreate) {}

func (b *beasty) CheckValidEmail(email string) bool {
	return true
}

// if valid message provided, branch to desired action and return true
// if no action performed, return false
func (b *beasty) MessageParse(m string) ([]string, bool) {
	forBot := false
	signalWord, valid := b.config["signal_word"].(string)
	if !valid {
		log.Fatalf("Error: could not retreive signal word from config")
	}
	parts := strings.Split(m, " ")
	if parts[0] == signalWord {
		forBot = true
	}
	arguments := parts[1:]
	return arguments, forBot
}

// call bootstrap and initiate blocking unbuffered channel to handle
// discord message events
func (b *beasty) Start() {
	b.bootstrapConnection()
	startChannel, valid := b.config["startup_channel_id"].(string)
	if !valid {
		log.Fatalf("Could not retrieve startup channel ID")
	}
	b.Connection.ChannelMessageSend(startChannel, b.Response("startup"))
	if b.localUse == true {
		fmt.Println("Bot is now running.  Press CTRL-C to exit.")
		// make unbuffered blocking channel to wait for user kill signal
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc
	} else {
		// this is my life now
		for {
		}
	}
	b.Connection.ChannelMessageSend(startChannel, b.Response("shutdown"))
	// Cleanly close down the Discord session.
	b.Connection.Close()

}
