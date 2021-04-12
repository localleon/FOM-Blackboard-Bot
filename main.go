package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
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

	// Setup execution every 30m for periodicly downloading the lastest OC-News
	getLatestOCNews()
	c := cron.New()
	cErr := c.AddFunc("@every 15m", getLatestOCNews)
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

func getLatestOCNews() {
	log.Print("Requesting new FOM-OC Blackboard Data")

	// Decode Credentials to cleartext
	context := createLoginContext()
	user, pwd := createLoginCredentials("FOM_USER", "FOM_PWD")

	getLoginCookie(user, pwd, context)
	// Parsing new OC-Messages
	news := getDashboardBlackboard()
	parseBlackBoardData(news)
	// // Check notification for courses
	getCourseNotification()
}

//createLoginCredentials reads the ENV-Vars out and decodes the credentials from base64
func createLoginCredentials(userEnv, pwdEnv string) (string, string) {
	envUser, uErr := base64.StdEncoding.DecodeString(os.Getenv(userEnv))
	envPwd, pErr := base64.StdEncoding.DecodeString(os.Getenv(pwdEnv))
	if uErr != nil || pErr != nil {
		log.Fatalf("Error decoding base64 values of user/password values")
	}

	return string(envUser), string(envPwd)
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
