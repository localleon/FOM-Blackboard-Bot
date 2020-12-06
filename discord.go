package main

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

func initBot() *discordgo.Session {
	token := os.Getenv("FOM_DTOKEN")
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println("Error connecting with Token to Discord API.")
		return nil
	}
	// Open the websocket and begin listening.
	err = discord.Open()
	if err != nil {
		log.Println("Error opening Discord session: ", err)
	}

	// Set channel to bind to
	bindChannel = os.Getenv("FOM_CHANNEL")

	return discord
}

func sendMessageToDiscord(msg blackBoardMsg) {
	c, err := d.Channel(bindChannel)
	if err != nil {
		log.Println("Error while trying to bind to #blackboard channel", err.Error())
	}
	// Format msg Body nicley and send to discord
	var text string
	text += "**" + msg.Title + "**\n"
	text += "*Date: " + msg.Date + "*\n\n"
	text += msg.Message + "\n\n"
	text += "*" + endpoint + msg.Link + "*"

	_, mErr := d.ChannelMessageSend(
		c.ID,
		text,
	)
	if mErr != nil {
		log.Println("Error while sending BlackBoardMessage to Discord Channel")
	}
}

func welcomeMessage(channelID string) {
	// Send WelcomeMessage to bindChannel
	c, err := d.Channel(channelID)
	if err != nil {
		log.Println("Error while trying to write welcome message", err.Error())
	}
	_, mErr := d.ChannelMessageSend(
		c.ID,
		"FOM-OC Bot is ready to rock!",
	)
	if mErr != nil {
		log.Println("Error while sending BlackBoardMessage to Discord Channel")
	}
}

// Works on the global Message queue
func sendQueueMessages() {
	for _, msg := range msgQueue {
		//printBlackboardMSG(msg)
		sendMessageToDiscord(msg)
	}
	// Newly initalize queue
	msgQueue = []blackBoardMsg{}
}
