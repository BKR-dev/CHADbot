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

/*
go routine getNewestPost(
	get topicRespone -> struct
	push struct -> channel
	get topicRespone -> newStruct
	if newStruct has higher PostCount then struct
		push struct -> channel
)


*/

// TopicResponse
type TopicResponse struct {
	ID                int    `json:"id"`
	Title             string `json:"title"`
	PostCount         int    `json:"posts_count"`
	HighestPostNumber int    `json:"highest_post_number"`
	PostStream        struct {
		Posts []struct {
			ID        int    `json:"id"`
			Username  string `json:"username"`
			UserID    int    `json:"user_id"`
			UserTitle string `json:"user_title"`
			Cooked    string `json:"cooked"`
		} `json:"posts"`
	} `json:"post_stream"`
}

// postCount struct
type postCount struct {
	HighestPost int `json:"highest_post_number"`
}

// latestPost struct
type latestPost struct {
	PostStream struct {
		Posts []struct {
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

func init() {
	topicId = "1118"
	apiKey = "5634da9f596ecc2740440a75499176a3b8181752aa418696b61ed08b982c3a43"
	apiUser = "terminator"
	apiUserId = 9
	url = "https://forum.pixelspace.xyz/t/"

	var err error
	Scribe, err = logging.NewLogger()
	if err != nil {
		fmt.Println("Error creating logger in responses.go")
	}
}

func GetLastPost(dataCh chan<- TopicResponse, errCh chan<- error) {

	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", url+topicId+".json", nil)
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Api-User", apiUser)

	if err != nil {
		errCh <- err
		return
	}
	fmt.Println(req)
	res, err := client.Do(req)

	if err != nil {
		errCh <- err
		return
	}

	defer res.Body.Close()

	var topicRespone TopicResponse

	err = json.NewDecoder(res.Body).Decode(&topicRespone)
	if err != nil {
		errCh <- err
		return
	}

	dataCh <- topicRespone

}

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
