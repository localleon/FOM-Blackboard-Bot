package main

import (
	"html"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type PrivateMessage struct {
	Subject    string
	Date       string
	NotifyDate time.Time
	Link       string
	Text       string
}

func parseBlackBoardData(d blackboardRes) {
	if d.Status == 200 {
		output := html.UnescapeString(d.HTML)
		html := replaceUmlauts(output)

		// Parse HTML
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			log.Println("Couldn't parse html of BlackBoardData")
			return
		}

		// Find the news items
		doc.Find("#cell_blackboardtype1").Each(func(i int, s *goquery.Selection) {
			// For each item found, parse the message
			s.Find("li").Each(parseMessageHTML)
		})
		// Find the news items
		doc.Find("#cell_mPrio").Each(func(i int, s *goquery.Selection) {
			// For each item found, parse the message
			s.Find("li").Each(parseMessageHTML)
		})
		return // skip
	}
	log.Println("Detected error message in blackBoardRes while trying to parse. Seems like an API Error. aborting.. ")
}

func parseMessageHTML(i int, s *goquery.Selection) {
	// Only parse msgs with content in it
	if !s.Is(":empty") {
		log.Println("Got new Data in Blackboard. Starting parsing process..... ")
		// Find all Values in HTML Doc
		Subject := s.Find(".titel").Text()
		date := s.Find(".date").Text()
		body := s.Find(".abstract").Text()
		link, state := s.Find(".abstract").Find("a").Attr("href")
		if !state {
			log.Println("Message", Subject, "does not contain an Hyperlink for more information")
		}

		// Cleanup and create message object
		body = replaceUmlauts(body)
		richBody := parseMessageBodyFromRef(link)
		if richBody != "" {
			body = richBody
		}

		// Send Notification to discord
		sendWebHook(os.Getenv("FOM_WEBHOOK"), "FOM-OC", Subject, link, "Am "+date+":", body)
		return
	}
	log.Print("Couldn't find any new articles. \n")
}

func parseMessageBodyFromRef(ref string) string {
	var msgString string

	url := endpoint + ref
	// Prepare new HTTP request
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Add("Content-Type", "charset=UTF-8")

	// Send HTTP request and move the response to the variable
	res, errR := client.Do(request)
	if errR != nil {
		log.Println("Cant get document from link, seems invalid", errR.Error())
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		doc, _ := goquery.NewDocumentFromReader(res.Body)
		// Parse each <p> Tag in the content div where the message is displayed.
		doc.Find("#content").Find("p").Each(func(i int, s *goquery.Selection) {
			if !s.Is(":empty") {
				txt := s.Text()
				if strings.Contains(txt, "Übersicht") { // Skip if its the "Übersicht" Dialog that is not relevant for the message
					return
				}
				msgString += txt + "\n"
			}
		})
	}

	// Cleanup first few D
	msgString = strings.Replace(msgString, "\n", "", 4)
	return msgString
}
