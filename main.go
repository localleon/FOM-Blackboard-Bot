package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

var client *http.Client

func main() {
	log.Println("Starting Application")
	client = &http.Client{}

	// Authentication
	username := os.Getenv("FOM_USER")
	password := os.Getenv("FOM_PWD")
	// Authenticate Session
	context := getLoginContext()
	session := getLoginCookie(username, password, context)

	return
	// Parsing
	fmt.Println("Working with the following cookies:", session)
	getDashboardBlackboard(session)
}

func getDashboardBlackboard(s []*http.Cookie) {
	endpoint := "https://campus.bildungscentrum.de"
	//params := "/nfcampus/startapi/blackboard"
	params := "/nfcampus/Node.do?n=5003"
	url := endpoint + params
	fmt.Println("Connecting to ", url)

	// Prepare new HTTP request
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Add all cookies from sessino to request
	for _, c := range s {
		request.AddCookie(c)
	}
	// Send HTTP request and move the response to the variable
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response.Status)
	// text, _ := httputil.DumpResponse(response, true)
	// fmt.Println(string(text))
}

// getLoginCookie creates a Session Cookie with FOM-OC Login.do Endpoint
func getLoginCookie(user, pwd string, ctx []*http.Cookie) []*http.Cookie {
	log.Println("Authenticating Session.....")
	// Make HTTP GET request
	// params := "crt=19453&assl=&iehack=%C3%A2%CB%9C%C2%A0&quelle=LoginForm-BCW&i=bcw"
	// body := params + "&name=" + user + "&password=" + pwd
	// fmt.Println("Composing body with value:", body)
	//
	// fmt.Println("Connecting to ", url)

	apiUrl := "https://campus.bildungscentrum.de"
	resource := "/nfcampus/Login.do"
	data := url.Values{}
	data.Set("crt", "19453")
	data.Set("quelle", "LoginForm-BCW")
	data.Set("i", "bcw")
	data.Set("name", user)
	data.Set("password", pwd)

	u, _ := url.ParseRequestURI(apiUrl)
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

	// Debug Output
	// text, _ := httputil.DumpRequest(request, true)
	// fmt.Println("Request:")
	// fmt.Println(string(text))

	fmt.Println("Response:")
	text2, _ := httputil.DumpResponse(response, true)
	fmt.Println(string(text2))

	// Return Set-Cookie from Response
	loginCookie := response.Cookies()
	if len(loginCookie) != 0 {
		return loginCookie
	}
	return nil
}

func getLoginContext() []*http.Cookie {
	log.Println("Getting Login-Form Auth Context")
	endpoint := "https://campus.bildungscentrum.de"
	params := "/nfcampus/pages/login.jsp"
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
