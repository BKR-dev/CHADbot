package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"io/ioutil"
	"net/http"
	"terminator-shitpost/logging"
	"time"
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

var Scribe logging.Logger

// target thread, apiKey and User
var topicId string
var apiKey string
var apiUser string
var apiUserId int
var url string
var highestPost int

func init() {
	topicId = "1127"
	apiKey = "0e0f0211b74c488045ffde571cb152ffdc65a96f4078535b632f7f05fc5c27da"
	apiUser = "CHADbot"
	apiUserId = 101
	url = "https://forum.pixelspace.xyz/t/"
	highestPost = 1

	// 101
	// 0e0f0211b74c488045ffde571cb152ffdc65a96f4078535b632f7f05fc5c27da
	// CHADbot
	// burying feds innawoods after making them sleep with akimbo 1911 because their vegan or  keto weirdo gaming fags who love anime and slurp ramen all day

	var err error
	Scribe, err = logging.NewLogger()
	if err != nil {
		fmt.Println("Error creating logger in responses.go")
	}
}

// gets the topic object and returns the latest post number from it
func getPostsFromTopic() (int, error) {

	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", url+topicId+".json", nil)
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Api-User", apiUser)

	if err != nil {
		return 0, err
	}

	res, err := client.Do(req)
	if !(res.StatusCode >= 200 && res.StatusCode <= 204) {
		return 0, errors.New("httpStatusCode is worrysome: " + fmt.Sprint(res.StatusCode))
	}

	if err != nil {
		Scribe.Infof("Request Status: %v", res.StatusCode)
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
		return LatestPost{}, apiUserId,
			errors.New("httpStatusCode is worrysome: " + fmt.Sprint(res.StatusCode))
	}

	if err != nil {
		Scribe.Infof("Request Status: %v", res.StatusCode)
		return LatestPost{}, apiUserId, err
	}

	defer res.Body.Close()

	var lastPost LatestPost

	err = json.NewDecoder(res.Body).Decode(&lastPost)
	if err != nil {
		return LatestPost{}, apiUserId, err
	}
	Scribe.Infof("Got last post, quick 1 second sleep")
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
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	return string(body), nil
}
