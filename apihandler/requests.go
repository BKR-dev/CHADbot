package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"terminator-shitpost/conf"
	"terminator-shitpost/logging"
	"time"
	"unicode"
)

// TopicResponse
type TopicResponse struct {
	ID                int    `json:"id"`
	Title             string `json:"title"`
	PostCount         int    `json:"posts_count"`
	HighestPostNumber int    `json:"highest_post_number"`
	PostStream        struct {
		Posts []struct {
			ID         int    `json:"id"`
			UserID     int    `json:"user_id"`
			Username   string `json:"username"`
			UserTitle  string `json:"user_title"`
			PostNumber int    `json:"post_number"`
			Cooked     string `json:"cooked"`
		} `json:"posts"`
	} `json:"post_stream"`
}

// latestPost struct
type LatestPost struct {
	ID                int    `json:"id"`
	Title             string `json:"title"`
	HighestPostNumber int    `json:"highest_post_number"`
	PostStream        struct {
		Posts []struct {
			ID         int    `json:"id"`
			UserID     int    `json:"user_id"`
			Username   string `json:"username"`
			UserTitle  string `json:"user_title"`
			PostNumber int    `json:"post_number"`
			Cooked     string `json:"cooked"`
		} `json:"posts"`
	} `json:"post_stream"`
}

var scribe *logging.Logger

// var settings *conf.Config

// input needed from config file
var topicId string
var apiKey string
var apiUser string
var apiUserId int
var url string
var highestPost int

// initializing vars
func init() {
	scribe, err := logging.NewLogger()
	if err != nil {
		fmt.Println("Error creating logger in responses.go")
	}
	scribe.Infof("Lets go")
}

// gets the topic object and returns the latest post number from it
func getPostsFromTopic() (int, error) {
	settings, err := conf.GetSettings()
	if err != nil {
		scribe.Errorf("Could not obtain config", err)
	}
	// fmt.Println(settings)
	topicId = settings.TopicId
	apiKey = settings.ApiKey
	apiUser = settings.ApiUser
	apiUserId, _ = strconv.Atoi(settings.ApiUserId)
	url = settings.Url
	highestPost = 1
	reqUrl := createUrlString(url, topicId)
	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return 0, fmt.Errorf("Error creating request: %v", err)
	}
	// good to know: if the request cannot be created req.Header.Set returns a SIGSEGV
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Api-Username", apiUser)

	res, err := client.Do(req)
	if !(res.StatusCode >= 200 && res.StatusCode <= 204) {
		fmt.Println(reqUrl)
		fmt.Println(req.Header)
		return 0, errors.New("httpStatusCode is worrysome: " + fmt.Sprint(res.StatusCode))
	}

	if err != nil {
		scribe.Infof("Request Status: %v", res.StatusCode)
		return 0, err
	}

	defer res.Body.Close()

	var topicResponse TopicResponse

	err = json.NewDecoder(res.Body).Decode(&topicResponse)
	if err != nil {
		return 0, err
	}
	highestPost = topicResponse.HighestPostNumber
	// time.Sleep(3 * time.Second)
	return highestPost, nil
}

// Returns the complete TopicResponse as 1:1 mapped struct, Posts are inside
// the struct as Slice of Posts in PostStream. Also returns error.
// will return empty Struct and error if Response StatusCode is NOT 200 - 204
func GetLastPost() (LatestPost, int, error) {

	postNumber, err := getPostsFromTopic()
	if err != nil {
		return LatestPost{}, apiUserId, err
	}

	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", url+topicId+"/"+fmt.Sprint(postNumber)+".json", nil)
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Api-User", apiUser)

	if err != nil {
		return LatestPost{}, apiUserId, err
	}

	res, err := client.Do(req)
	if !(res.StatusCode >= 200 && res.StatusCode <= 204) {
		scribe.Warning("the request in question: %v", req)
		return LatestPost{}, apiUserId,
			errors.New("httpStatusCode is worrysome: " + fmt.Sprint(res.StatusCode))
	}

	if err != nil {
		scribe.Infof("Request Status: %v", res.StatusCode)
		return LatestPost{}, apiUserId, err
	}

	defer res.Body.Close()

	var lastPost LatestPost

	err = json.NewDecoder(res.Body).Decode(&lastPost)
	if err != nil {
		return LatestPost{}, apiUserId, err
	}
	scribe.Infof("Got last post, quick 1 second sleep")
	time.Sleep(1 * time.Second)
	return lastPost, apiUserId, nil
}

// posts response to discourse api
func PostResponseToTopic(message string) error {

	if message == "" {
		return errors.New("Empty Message String")
	}

	jsonBody := []byte(`{"topic_id": "%s", "raw": "%s"}`)
	s := fmt.Sprintf(string(jsonBody), topicId, message)
	body := bytes.NewReader([]byte(s))

	url := "https://forum.pixelspace.xyz/posts.json"
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Api-Username", apiUser)

	client := http.Client{Timeout: 30 * time.Second}
	httpCode, err := client.Do(req)
	if !(httpCode.StatusCode >= 200 && httpCode.StatusCode <= 204) {
		return errors.New("httpStatusCode is worrysome: " + fmt.Sprint(httpCode.StatusCode))
	}
	if err != nil {
		return err
	}
	return nil
}

// fetches a random bible verse from labs.bible.org and returns a verse or error
func GetRandomBibleVerse() (string, error) {
	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", "https://labs.bible.org/api/?passage=random", nil)
	if err != nil {
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	return string(body), nil
}

// creates satinized url free from qoutation marks
func createUrlString(url string, topicId string) string {
	jst := ".json"
	t := "/t/"
	var sb strings.Builder
	var newUrl string
	sb.WriteString(url)
	sb.WriteString(t)
	sb.WriteString(topicId)
	sb.Write([]byte(jst))
	oldUrl := sb.String()
	sb.Reset()
	for _, r := range oldUrl {
		if !unicode.Is(unicode.Quotation_Mark, r) {
			sb.WriteRune(r)
		}
	}
	newUrl = sb.String()
	return newUrl
}
