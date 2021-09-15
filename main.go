package main

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/robfig/cron"
)

var client *http.Client

const endpoint string = "https://campus.bildungscentrum.de"

func main() {
	log.Println("Starting FOM-OC Discord Bot")
	// Check if Env-Vars are present
	checkEnvVars()

	// HTTP Client Setup with global cookie storage
	jar, _ := cookiejar.New(nil)
	client = &http.Client{
		Jar: jar,
	}

	// Setup execution every 15m for periodicly downloading the lastest OC-News
	processOCData()
	c := cron.New()
	cErr := c.AddFunc("@every 15m", processOCData)
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
		c.Stop()
	}()
}

//checkEnvVars tests if all enviroment variables are correctly set
func checkEnvVars() {
	if os.Getenv("FOM_USER") == "" || os.Getenv("FOM_PWD") == "" {
		log.Fatal("User or PWD Env-Var is empty. Please provided login credentials")
	}
	if os.Getenv("FOM_WEBHOOK") == "" {
		log.Fatal("Discord WebHook Env-Var is not set. Cancelling..")
	}
}

// processOCData parses Blackboard News and Course notification data
func processOCData() {
	// Create Authorization Context
	context := createLoginContext()
	user, pwd := os.Getenv("FOM_USER"), os.Getenv("FOM_PWD")
	getLoginCookie(user, pwd, context)

	// Parsing new Blackboard Messages
	log.Println("Requesting Blackboard Data")
	news := getDashboardBlackboard()
	parseBlackBoardData(news)

	log.Println("Finished working on Blackboard Data")
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
