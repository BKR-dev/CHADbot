package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"terminator-shitpost/logging"
	"time"
)

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

// target thread, apiKey and User
var thread string
var apiKey string
var apiUser string

func init() {
	thread = "1118"
	apiKey = "5634da9f596ecc2740440a75499176a3b8181752aa418696b61ed08b982c3a43"
	apiUser = "terminator"
}

func GetLastPost(logging.Logger) (string, error) {
	logger := logging.NewLogger(logging.DebugLevel, os.Stdout)
	var lastPost string
	for {
		hp := getHighestPost()

		if hp == (postCount{}) || hp.HighestPost == 0 {
			logger.Info("Empty, sleeping it off...")

			time.Sleep(2 * time.Minute)
			continue
		}

		lp := getLatestPost(hp.HighestPost)

		if lp == "none" {
			log.Println("latestPost: empty result")
			time.Sleep(5 * time.Second)
			continue
		}
		lastPost = lp

		return lastPost, nil
	}
}

func PostResponseToTopic(log logging.Logger, message string) error {

	jsonBody := []byte(`{"topic_id": "%s", "raw": "%s"}`)
	s := fmt.Sprintf(string(jsonBody), thread, message)
	body := bytes.NewReader([]byte(s))

	url := fmt.Sprintf("https://forum.pixelspace.xyz/posts.json")
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Api-Username", apiUser)

	client := http.Client{Timeout: 30 * time.Second}
	_, err = client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func getHighestPost() postCount {
	url := "https://forum.pixelspace.xyz/t/" + thread + ".json"
	var result postCount

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Failed new request")
		fmt.Println(err)
		return result
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Api-Username", apiUser)

	res, err := client.Do(req)
	if err != nil {
		log.Println("Failed client request")
		return result
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Failed res body")
		return result
	}

	err = json.Unmarshal(body, &result)

	if err != nil {
		log.Println("Failed unmarshal")
		return result
	}

	res.Body.Close()

	return result
}

func getLatestPost(highestPost int) string {
	url := "https://forum.pixelspace.xyz/t/" + thread + "/" + strconv.Itoa(highestPost) + ".json"
	client := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Failed new request")
		return "none"
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Api-Username", apiUser)

	res, err := client.Do(req)
	if err != nil {
		log.Println("Failed client request")
		return "none"
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Failed response")
		return "none"
	}

	var latest latestPost
	err = json.Unmarshal(body, &latest)

	if err != nil {
		log.Println("Failed unmarshal")
		return "none"
	}

	if len(latest.PostStream.Posts) == 0 {
		log.Println("PostStream.Posts is empty")
		return "none"
	}

	currentPost := latest.PostStream.Posts[len(latest.PostStream.Posts)-1].Cooked
	res.Body.Close()
	return currentPost
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
