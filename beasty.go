package main

import (
	"io/ioutil"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/go-yaml/yaml"
)

type beasty struct {
	token      string
	Connection *discordgo.Session
	config     map[interface{}]interface{}
	responses  *responses
	storage    struct{}
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
func (b beasty) Response(lookup string) string {
	return b.responses.generateResponse(lookup)
}

// configures intent and applies callbacks to discord session
func (b beasty) bootstrapConnection() {}

// initial point for handling discord events
// will call parse to determine action and if valid
func (b beasty) MessageEventCallback(s *discordgo.Session, m *discordgo.MessageCreate) {}

// if valid message provided, branch to desired action and return true
// if no action performed, return false
func (b beasty) MessageParse(msg string) bool {
	return true
}

// call bootstrap and initiate blocking unbuffered channel to handle
// discord message events
func (b beasty) Start() {
	return
}
