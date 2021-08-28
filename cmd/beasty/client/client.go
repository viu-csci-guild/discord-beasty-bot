package client

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-yaml/yaml"
	"github.com/viu-csci-guild/beasty/cmd/beasty/responses"
)

type Beasty struct {
	token         string
	Connection    *discordgo.Session
	config        map[interface{}]interface{}
	responses     *responses.Responses
	stuRoleID     string
	watchRoomID   string
	startupRoomID string
	serverID      string
}

var BeastyHandle *Beasty = nil

// Provides instance of beasty bot with valid configuration
// determined by yml files
// singleton
func NewBeasty(t string, srid string, wrid string, surid string, sid string) *Beasty {
	rand.Seed(time.Now().UnixNano())
	if BeastyHandle != nil {
		return BeastyHandle
	}
	configMap := make(map[interface{}]interface{})
	data, err := ioutil.ReadFile("./client/config.yaml")
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
	r := responses.NewResponses(responseData)

	dgo, err := discordgo.New("Bot " + t)
	if err != nil {
		log.Fatalf("error creating Discord session: %s", err)
	}

	b := &Beasty{
		token:         t,
		config:        c,
		responses:     r,
		Connection:    dgo,
		stuRoleID:     srid,
		watchRoomID:   wrid,
		startupRoomID: surid,
		serverID:      sid,
	}
	BeastyHandle = b
	return b
}

// randomly fetch a response based on a lookup string, for personality
func (b *Beasty) Response(lookup string) string {
	return b.responses.GenerateResponse(lookup)
}

// configures intent and applies callbacks to discord session
func (b *Beasty) bootstrapConnection() {
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
func (b *Beasty) MessageCallback(s *discordgo.Session, m *discordgo.MessageCreate) {
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
	case "student":
		gid := b.serverID
		rid := b.stuRoleID
		err := b.Connection.GuildMemberRoleAdd(gid, m.Author.ID, rid)
		if err != nil {
			log.Printf("Role add error: %s", err)
		}
		b.Connection.ChannelMessageSend(m.ChannelID, b.Response("student"))
	case "joke":
		b.Connection.ChannelMessageSend(m.ChannelID, b.Response("joke"))
	default:
		b.Connection.ChannelMessageSend(m.ChannelID, b.Response("unkown"))
	}
}

// callback for handling guild joins, such as a greeting
func (b *Beasty) GuildJoinCallback(s *discordgo.Session, m *discordgo.MessageCreate) {}

// if valid message provided, branch to desired action and return true
// if no action performed, return false
func (b *Beasty) MessageParse(m string) ([]string, bool) {
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
func (b *Beasty) Start() {
	// if start ends for any reason, close redis
	b.bootstrapConnection()
	b.Connection.ChannelMessageSend(b.startupRoomID, b.Response("startup"))
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	// make unbuffered blocking channel to wait for user kill signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	// // this is my life now
	// for {
	// }
	b.Connection.ChannelMessageSend(b.startupRoomID, b.Response("shutdown"))
	// Cleanly close down the Discord session.
	b.Connection.Close()

}
