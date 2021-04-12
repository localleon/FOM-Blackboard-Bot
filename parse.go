package main

import (
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var notifyBuffer []string // saves sent notifications on runtime. The notification feature makes the application somewhat stateful because we need to remeber sent messages

type PrivateMessage struct {
	Subject    string
	Date       string
	NotifyDate time.Time
	Link       string
	Text       string
}

func processOCMailbox(doc *goquery.Document) {
	msgs := extractMessagesFromNode(doc)

	for _, m := range msgs {
		// Check if we already sent the notification

		if !contains(notifyBuffer, m.Subject) {
			notifyBuffer = append(notifyBuffer, m.Subject)
		} else {
			processPrivateMessage(m) // Send Discord Notification out
		}
	}
}

//extractMessagesFromNode extracts all Zoom-Notify messages from the HTML-Node
func extractMessagesFromNode(doc *goquery.Document) []PrivateMessage {
	found := []PrivateMessage{}

	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		// Keep date/time for element
		dateMatch := false
		day := ""

		s.Find("td").Each(func(h int, row *goquery.Selection) {
			switch h {
			case 2: // Check if Notification is from today
				day = time.Now().In(loc).Format("02.01.2006")

				if strings.Contains(row.Text(), day) {
					dateMatch = true
				}
			case 5: // Extract Notify Date
				rowT := row.Text()
				headerDate := extractDateFromHeading(rowT)

				// Does the heading match all checks?
				if dateMatch && isInTimeWindow(day) && isZoomNotification(rowT) {

					// Extract Data
					subject := strings.ReplaceAll(rowT, "'", "")
					subject = strings.TrimSpace(subject)
					msgLink, _ := row.Find("a").Attr("href")

					// Add message to all found messages
					msg := PrivateMessage{
						Subject:    subject,
						Date:       day,
						NotifyDate: headerDate,
						Link:       msgLink,
					}
					found = append(found, msg)
				}
			}
		})
	})

	return found
}

func extractDateFromHeading(heading string) time.Time {
	// TODO: Add function
	// fmt.Println(heading) //
	return time.Now()
}

func isZoomNotification(heading string) bool {
	return strings.Contains(heading, "Ihre Videokonferenz startet in Kuerze um")
}

func isInTimeWindow(day string) bool {
	hour := time.Now().In(loc)

	// Time windows where a notification could be sent , UTC time in strings is converted to CET
	t1, _ := time.Parse("02.01.2006 15:04:05", day+" 06:30:00")
	t2, _ := time.Parse("02.01.2006 15:04:05", day+" 07:30:00")
	t3, _ := time.Parse("02.01.2006 15:04:05", day+" 10:30:00")
	t4, _ := time.Parse("02.01.2006 15:04:05", day+" 11:30:00")
	// Check if our message is in window
	return hour.After(t1.In(loc)) && hour.Before(t2.In(loc)) || hour.After(t3.In(loc)) && hour.Before(t4.In(loc))
}

// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func processPrivateMessage(msg PrivateMessage) {
	url := endpoint + msg.Link

	// Prepare new HTTP request
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Add("Content-Type", "charset=UTF-8")

	// Send HTTP request and move the response to the variable
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		data, errB := ioutil.ReadAll(response.Body)
		if errB != nil {
			log.Println("Error decoding Body of Notifications")
		}
		// Create HTML Document
		output := html.UnescapeString(string(data))
		html := replaceUmlauts(output)
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
		s := doc.Find("table").Find("fieldset")

		s.Contents().Each(func(i int, s *goquery.Selection) {
			if !s.Is("br") {
				r := regexp.MustCompile("^http|https+://*$") // This matches a line that contains only a link
				if r.Match([]byte(s.Text())) {
					// Send Webhook to discord endpoint
					sendWebHook(os.Getenv("FOM_WEBHOOK_COURSES"), "FOM-Notify", "", msg.Subject, "Zoom Notification", s.Text())
				}
			}
		})

	}
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

// stringToHMTL converts a string to a parsable goquery Document
func stringToHTML(data string) *goquery.Document {
	// Create HTML Document
	output := html.UnescapeString(data)
	html := replaceUmlauts(output)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		log.Println("Couldn't parse html of Notifications")
		return nil
	}
	return doc
}
