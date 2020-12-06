package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron"
)

var client *http.Client

const endpoint string = "https://campus.bildungscentrum.de"

var msgQueue []blackBoardMsg // Global Queue for storing parsed Items
var bindChannel string       // #blackboard channel can be set static
var d *discordgo.Session

func main() {
	log.Println("Starting FOM-OC Discord Bot")
	checkEnvVars()

	// Init Discord Bot
	d = initBot()
	//welcomeMessage(bindChannel) // Currently only to register that the bot is online

	// HTTP Client Setup with global cookie storage
	jar, _ := cookiejar.New(nil)
	client = &http.Client{
		Jar: jar,
	}

	// Setup execution every 30m for periodicly downloading the lastest OC-News
	getLatestOCNews()
	c := cron.New()
	cErr := c.AddFunc("@every 30m", getLatestOCNews)
	if cErr != nil {
		log.Println("Can't setup cron handler")
		os.Exit(5)
	}
	c.Start()

	// Wait for shutdown via control-c
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch

	// Close various things
	defer func() {
		log.Println("Received shutdown signal, exiting gracefully........")
		d.Close()
		c.Stop()
	}()
}

//checkEnvVars tests if all enviroment variables are correctly set
func checkEnvVars() {
	if os.Getenv("FOM_USER") == "" || os.Getenv("FOM_PWD") == "" {
		log.Fatal("User or PWD Env-Var is empty. Please provided login credentials")
	}
	if os.Getenv("FOM_DTOKEN") == "" {
		log.Fatal("Discord Token Env-Var is not set. Cancelling..")
	}
	if os.Getenv("FOM_CHANNEL") == "" {
		log.Fatal("Discord Channel ID is empty. We need a channel to write to for the bot")
	}
}

func getLatestOCNews() {
	log.Print("Requesting new FOM-OC Blackboard Data")
	// Authenticate Session
	username := os.Getenv("FOM_USER")
	password := os.Getenv("FOM_PWD")

	context := getLoginContext()
	getLoginCookie(username, password, context)
	// Parsing new OC-Messages
	news := getDashboardBlackboard()
	parseBlackBoardData(news)
	// Working on the Messages and sending to Discord
	sendQueueMessages()
}

// 	prints the msg to stdout for debug purposes
func printBlackboardMSG(msg blackBoardMsg) {
	fmt.Println("----------------------------")
	fmt.Println("Working on:")
	fmt.Println("- Title:", msg.Title)
	fmt.Println("- Posted on:", msg.Date)
	fmt.Println("- Link:", msg.Link)
	fmt.Println(msg.Message)
}

//loadSampleBlackboards loads an API response from a local file so we dont generate to much network traffic while developing
func loadSampleBlackboard(path string) blackboardRes {
	// load parse.html
	data, _ := ioutil.ReadFile(path)

	b := blackboardRes{
		Status:      200,
		NewElements: 1,
		TotalRows:   1,
		HTML:        string(data),
	}
	return b
}

func replaceUmlauts(s string) string {
	// Common German Umlauts replacment
	s = strings.Replace(s, "ä", "ae", -1)
	s = strings.Replace(s, "ö", "oe", -1)
	s = strings.Replace(s, "ü", "ue", -1)
	s = strings.Replace(s, "ß", "ss", -1)
	s = strings.Replace(s, "\n", "'", -1)

	return s
}
