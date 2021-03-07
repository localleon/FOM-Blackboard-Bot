package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Thumbnail is the embeded contentview of a discord message
type Thumbnail struct {
	URL string `json:"url"`
}

//Fields is discords field object in embeded content
type Fields struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

//Footer is discords footer information
type Footer struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url"`
}

//Embeds is discord any embedded content
type Embeds struct {
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	Color     int       `json:"color"`
	Timestamp string    `json:"timestamp"`
	Thumbnail Thumbnail `json:"thumbnail"`
	Fields    []Fields  `json:"fields"`
	Footer    Footer    `json:"footer"`
}

//Webhook represents the json struct the discord API expects
type Webhook struct {
	Username  string   `json:"username"`
	AvatarURL string   `json:"avatar_url"`
	Embeds    []Embeds `json:"embeds"`
}

func sendCourseNotification(confLink, title string) {
	c := &http.Client{}

	// Construct Webhook

	e := Embeds{
		Title: title,
		Color: 3066993, // Green
		Fields: []Fields{
			Fields{
				Name:   "Conference Notification",
				Value:  confLink,
				Inline: false,
			},
		},
	}

	w := Webhook{
		Username: "FOM-OC",
		Embeds:   []Embeds{e},
	}

	webhookReq, err := json.Marshal(w)
	if err != nil {
		log.Println("Couldn't parse Blackboard Message into Discord Embeded struct")
	}

	// Send the webhook to the discord api
	url := os.Getenv("FOM_WEBHOOK_COURSES")
	fmt.Println(url)
	req, rErr := http.NewRequest("POST", url, bytes.NewBuffer(webhookReq))
	if rErr != nil {
		log.Println("Couldn't create discord request out of course notification")
	}
	req.Header.Add("Content-Type", "application/json")

	res, qErr := c.Do(req)
	if qErr != nil {
		log.Println("Error while sending http discord webhook")
	}
	defer res.Body.Close()
}

func sendMessageToDiscord(msg blackBoardMsg) {
	c := &http.Client{}

	// Emebed Content only supports 1024 characters in total
	if len(msg.Message) >= 1023 {
		msg.Message = msg.Message[:1000] + ".........."
	}

	// Construct Webhook
	nURL := "https://campus.bildungscentrum.de/" + msg.Link

	e := Embeds{
		Title: msg.Title,
		URL:   nURL,
		Color: 3066993, // Green
		Fields: []Fields{
			Fields{
				Name:   "Am " + msg.Date + ":",
				Value:  msg.Message,
				Inline: true,
			},
		},
	}

	w := Webhook{
		Username: "FOM-OC",
		Embeds:   []Embeds{e},
	}

	webhookReq, err := json.Marshal(w)
	if err != nil {
		log.Println("Couldn't parse Blackboard Message into Discord Embeded struct")
	}

	// Send the webhook to the discord api
	url := os.Getenv("FOM_WEBHOOK")
	req, rErr := http.NewRequest("POST", url, bytes.NewBuffer(webhookReq))
	if rErr != nil {
		log.Println("Couldn't create discord request out of blackboard messag")
	}
	req.Header.Add("Content-Type", "application/json")

	res, qErr := c.Do(req)
	if qErr != nil {
		log.Println("Error while sending http discord webhook")
	}
	defer res.Body.Close()
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
