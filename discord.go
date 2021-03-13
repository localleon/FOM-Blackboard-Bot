package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
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

func sendWebHook(hook, hookName, title, eURL, fieldName, msg string) {
	c := &http.Client{}

	// Construct Webhook

	e := Embeds{
		Title: title,
		URL:   eURL,
		Color: 3066993, // Green
		Fields: []Fields{
			{
				Name:   fieldName,
				Value:  msg,
				Inline: false,
			},
		},
	}

	w := Webhook{
		Username: hookName,
		Embeds:   []Embeds{e},
	}

	webhookReq, err := json.Marshal(w)
	if err != nil {
		log.Println("Couldn't parse struct into WebHook Request")
	}

	req, rErr := http.NewRequest("POST", hook, bytes.NewBuffer(webhookReq))
	if rErr != nil {
		log.Println("Couldn't create discord webhook request out of data")
	}
	req.Header.Add("Content-Type", "application/json")

	res, qErr := c.Do(req)
	if qErr != nil {
		log.Println("Error while sending http discord webhook")
	}
	defer res.Body.Close()
}
