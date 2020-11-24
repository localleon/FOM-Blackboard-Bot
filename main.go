package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"golang.org/x/net/html/charset"
)

var client *http.Client

const endpoint string = "https://campus.bildungscentrum.de"

type blackboardRes struct {
	Status      int    `json:"status"`
	HTML        string `json:"html"`
	NewElements int    `json:"newelements"`
	TotalRows   int    `json:"total_rows"`
}

type blackBoardMsg struct {
	Title   string
	Date    string
	Message string
	Link    string
}

var msgQueue []blackBoardMsg // Global Queue for storing parsed Items

func main() {
	log.Println("Starting Application")
	jar, _ := cookiejar.New(nil)
	client = &http.Client{
		Jar: jar,
	}

	// Authentication
	username := os.Getenv("FOM_USER")
	password := os.Getenv("FOM_PWD")
	fmt.Println("LoginData: ", username, password)
	// Authenticate Session
	context := getLoginContext()
	getLoginCookie(username, password, context)
	// // Parsing
	//news := getDashboardBlackboard()
	news := loadSampleBlackboard("samples/api.html")
	parseBlackBoardData(news)
	sendQueueMessages()
}

// Works on the global Message queue
func sendQueueMessages() {
	for _, msg := range msgQueue {
		fmt.Println("----------------------------")
		fmt.Println("Working on:")
		fmt.Println("- Title:", msg.Title)
		fmt.Println("- Posted on:", msg.Date)
		fmt.Println("- Link:", msg.Link)
		fmt.Println(msg.Message)
	}
	// Newly initalize queue
	msgQueue = []blackBoardMsg{}
}

func parseBlackBoardData(d blackboardRes) {
	if d.Status == 200 {
		log.Println("------Starting parsing of Blackboard API Data------")
		output := html.UnescapeString(d.HTML)
		html := replaceUmlauts(output)

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			log.Println("Couldn't parse html")
			return
		}
		// Find the news  items
		doc.Find("#cell_blackboardtype1").Each(func(i int, s *goquery.Selection) {
			// For each item found, parse the message
			s.Find("ul").Each(parseMessageHTML)
		})
		// Find the news  items
		doc.Find("#cell_mPrio").Each(func(i int, s *goquery.Selection) {
			// For each item found, parse the message
			s.Find("ul").Each(parseMessageHTML)
		})
	}
}

func parseMessageHTML(i int, s *goquery.Selection) {
	// Only parse msgs with content in it
	if !s.Is(":empty") {
		// Find all Values in HTML Doc
		title := s.Find(".titel").Text()
		date := s.Find(".date").Text()
		body := s.Find(".abstract").Text()
		link, state := s.Find(".abstract").Find("a").Attr("href")
		if state != true {
			log.Println("Message", title, "doesn not contain an Hyperlink for more information")
		}

		// Cleanup and create message object
		body = replaceUmlauts(body)
		richBody := parseMessageBodyFromRef(link)
		if richBody != "" {
			body = richBody
		}

		msg := blackBoardMsg{
			Title:   title,
			Date:    date,
			Message: body,
			Link:    link,
		}
		msgQueue = append(msgQueue, msg) // Add Item to queue to be parsed
	}
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
	res, err := client.Do(request)
	if err != nil {
		log.Println("Cant get document from link, seems invalid", err.Error())
	}
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
	return msgString
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

func getDashboardBlackboard() blackboardRes {
	params := "/nfcampus/startapi/blackboard"
	//params := "/nfcampus/Node.do?n=5003"
	url := endpoint + params

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
	if response.StatusCode == 200 {
		data := blackboardRes{}
		json.NewDecoder(response.Body).Decode(&data)
		return data
	}
	fmt.Println("Error while getting Dashboard", response.Status)
	return blackboardRes{Status: response.StatusCode}
}

func detectContentCharset(body io.Reader) string {
	r := bufio.NewReader(body)
	if data, err := r.Peek(1024); err == nil {
		if _, name, ok := charset.DetermineEncoding(data, ""); ok {
			return name
		}
	}
	return "utf-8"
}

// getLoginCookie creates a Session Cookie with FOM-OC Login.do Endpoint
func getLoginCookie(user, pwd string, ctx []*http.Cookie) []*http.Cookie {
	log.Println("Authenticating Session.....")

	resource := "/nfcampus/Login.do"
	// Emulate Form Data of Login Page
	data := url.Values{}
	data.Set("crt", "19453")
	data.Set("assl", "")
	data.Set("iehack", "%C3%A2%CB%9C%C2%A0")
	data.Set("quelle", "LoginForm-BCW")
	data.Set("i", "bcw")
	data.Set("name", user)
	data.Set("password", pwd)

	// Build request
	u, _ := url.ParseRequestURI(endpoint)
	u.Path = resource
	urlStr := u.String()                                                                       // "https://api.com/user/"
	request, err := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	request.Header.Add("Referer", "https://campus.bildungscentrum.de/nfcampus/pages/login.jsp")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36")

	// REad out request body
	// bodybytes, _ := ioutil.ReadAll(request.Body)
	// fmt.Println(string(bodybytes))

	// Add SessionID from Login Context
	for _, c := range ctx {
		request.AddCookie(c)
	}

	// ERROR: Die Benutzerrerkennung wird benötigt.
	// Send HTTP request and move the response to the variable
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	// Return Set-Cookie from Response
	loginCookie := response.Cookies()
	if len(loginCookie) != 0 {
		return loginCookie
	}
	return nil
}

func getLoginContext() []*http.Cookie {
	log.Println("Getting Login-Form Auth Context")
	params := "/nfcampus/Login.do"
	url := endpoint + params
	// Prepare new HTTP request
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36")

	// Send HTTP request and move the response to the variable
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	// Return login context sessionid
	if response.StatusCode == 200 {
		return response.Cookies()
	}
	return nil
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
