package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-yaml/yaml"
	"github.com/gomodule/redigo/redis"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type beasty struct {
	token      string
	Connection *discordgo.Session
	config     map[interface{}]interface{}
	responses  *responses
	Storage    redis.Conn
	localUse   bool
}

var codeBase int = 1000000000000000
var codeTop int = 1999999999999999
var beastyHandle *beasty = nil

// Provides instance of beasty bot with valid configuration
// determined by yml files
// singleton
func NewBeasty(t string) *beasty {
	rand.Seed(time.Now().UnixNano())
	if beastyHandle != nil {
		return beastyHandle
	}
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

	redisServer := os.Getenv("REDIS_SERVER")
	redisPort := os.Getenv("REDIS_PORT")
	conn, err := redis.Dial("tcp", redisServer+":"+redisPort)
	if err != nil {
		log.Fatal(err)
	}

	b := &beasty{
		token:      t,
		config:     c,
		responses:  r,
		Connection: dgo,
		Storage:    conn,
	}
	beastyHandle = b
	return b
}

// randomly fetch a response based on a lookup string, for personality
func (b *beasty) Response(lookup string) string {
	return b.responses.generateResponse(lookup)
}

// enables user prompts and blocking if called with CLI arguments
func (b *beasty) SetLocalUse(f bool) {
	if f {
		log.Println("Note: set bot to provide local CLI prompt")
	}
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
		// type cast redis value, as it's safer than user string input
		storedCode, err := redis.String(b.Storage.Do("GET", m.Author.ID))
		if err != nil {
			log.Printf("Note: user ID: %s name: %s submitted code %s and redis GET failed", m.Author.ID, m.Author.Username, arguments[1])
		} else {
			if storedCode == arguments[1] {
				gid, gidv := b.config["guild_id"].(string)
				rid, ridv := b.config["student_role_id"].(string)
				if !gidv || !ridv {
					log.Fatalf("Error: failed to get ids for role change")
				}
				b.Connection.GuildMemberRoleAdd(gid, m.Author.ID, rid)
				b.Connection.ChannelMessageSend(m.ChannelID, b.Response("student_role"))
			} else {
				b.Connection.ChannelMessageSend(m.ChannelID, b.Response("code_invalid"))
				log.Printf("Note: user ID: %s name: %s submitted code %s compared to stored %s", m.Author.ID, m.Author.Username, arguments[1], storedCode)
			}
		}

	case "student":
		if len(arguments) == 2 && b.CheckValidEmail(arguments[1]) {
			// TODO: make sure !exist in redis already
			newCode := codeBase + rand.Intn(codeTop-codeBase)
			userName := m.Author.Username

			emailContent, vc := b.config["sender_message"].(string)
			emailSubject, vs := b.config["email_subject"].(string)
			emailSrc, vsrc := b.config["sender_email"].(string)
			apiKey := os.Getenv("SENDGRID_API_KEY")
			if !vc || !vs || !vsrc {
				log.Fatalf("Error: could not retrieve email configurations")
			}
			from := mail.NewEmail("viu-csci-bot", emailSrc)
			to := mail.NewEmail("Csci Student", arguments[1])
			plainTextContent := fmt.Sprintf(emailContent, userName, newCode)
			// We don't need html markup so we add plain twice to support both user modes
			message := mail.NewSingleEmail(from, emailSubject, to, plainTextContent, plainTextContent)
			client := sendgrid.NewSendClient(apiKey)
			response, err := client.Send(message)
			if err != nil {
				log.Fatalf("Error: email %s", err)
			} else {
				log.Printf("Email response: %d %s %v", response.StatusCode, response.Body, response.Headers)
			}
			b.Storage.Do("SET", m.Author.ID, newCode)
			b.Connection.ChannelMessageSend(m.ChannelID, b.Response("student"))
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
	// if start ends for any reason, close redis
	defer b.Storage.Close()
	b.bootstrapConnection()
	startChannel, valid := b.config["startup_channel_id"].(string)
	if !valid {
		log.Fatalf("Could not retrieve startup channel ID")
	}
	noWakeup, _ := strconv.ParseBool(os.Getenv("DISABLE_WAKEUP"))
	if !noWakeup {
		b.Connection.ChannelMessageSend(startChannel, b.Response("startup"))
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	// make unbuffered blocking channel to wait for user kill signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	if !noWakeup {
		b.Connection.ChannelMessageSend(startChannel, b.Response("shutdown"))
	}
	// Cleanly close down the Discord session.
	b.Connection.Close()

}
