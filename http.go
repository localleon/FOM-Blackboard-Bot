package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
)

//blackboardRes represents the API Response of the /startapi/blackboard endpoint
type blackboardRes struct {
	Status      int    `json:"status"`
	HTML        string `json:"html"`
	NewElements int    `json:"newelements"`
	TotalRows   int    `json:"total_rows"`
}

func getDashboardBlackboard() blackboardRes {
	params := "/nfcampus/startapi/blackboard"
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
	log.Println("Error while getting Dashboard", response.Status)
	return blackboardRes{Status: response.StatusCode}
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

	if user == "" || pwd == "" {
		log.Fatal("User or PWD Env-Var is empty. Please provided login credentials")
	}
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

	// Add SessionID from Login Context
	for _, c := range ctx {
		request.AddCookie(c)
	}
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
