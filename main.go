package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var client *http.Client

const endpoint string = "https://campus.bildungscentrum.de"

var msgQueue []blackBoardMsg // Global Queue for storing parsed Items
var bindChannel string
var d *discordgo.Session

func main() {
	log.Println("Starting FOM-OC Discord Bot")

	// Init Discord Bot
	d = initBot()
	//welcomeMessage(bindChannel) // Currently only to register that the bot is online

	// HTTP Client Setup with global cookie storage
	jar, _ := cookiejar.New(nil)
	client = &http.Client{
		Jar: jar,
	}

	// Authenticate Session
	username := os.Getenv("FOM_USER")
	password := os.Getenv("FOM_PWD")
	context := getLoginContext()
	getLoginCookie(username, password, context)

	// Parsing new OC-Messages
	//news := getDashboardBlackboard()
	news := loadSampleBlackboard("samples/api.html")
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
