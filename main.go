package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

var client *http.Client

func main() {
	log.Println("Starting Application")
	client = &http.Client{}

	// Authentication
	username := os.Getenv("FOM_USER")
	password := os.Getenv("FOM_PWD")
	fmt.Println("REMOVE: User Data", username, password)
	session := getLoginCookie(username, password)

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
	fmt.Println(response.Body)
	text, _ := httputil.DumpResponse(response, true)
	fmt.Println(string(text))
}

// getLoginCookie creates a Session Cookie with FOM-OC Login.do Endpoint
func getLoginCookie(user, pwd string) []*http.Cookie {
	// Make HTTP GET request
	endpoint := "https://campus.bildungscentrum.de"
	params := "?crt=19453&assl=&iehack=%C3%A2%CB%9C%C2%A0&quelle=LoginForm-FOM&i=fom"
	url := endpoint + "/nfcampus/Login.do" + params + "&name=" + user + "&password=" + pwd
	fmt.Println("Connecting to ", url)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	// Return Set-Cookie from Response
	loginCookie := response.Cookies()
	if len(loginCookie) != 0 {
		return loginCookie
	}
	return nil
}
