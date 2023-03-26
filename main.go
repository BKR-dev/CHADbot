package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type postCount struct {
	HighestPost int `json:"highest_post_number"`
}

type latestPost struct {
	PostStream struct {
		Posts []struct {
			PostNumber int    `json:"post_number"`
			Cooked     string `json:"cooked"`
		} `json:"posts"`
	} `json:"post_stream"`
}

const thread = "1016"
const apiKey = "5634da9f596ecc2740440a75499176a3b8181752aa418696b61ed08b982c3a43"
const apiUser = "terminator"

func main() {
	temp := 0

	for {
		hp := getHighestPost()
		if hp.HighestPost == temp {
			time.Sleep(5 * time.Minute)
			continue
		}
		if hp == (postCount{}) {
			fmt.Println("See error: empty result")
			return
		}

		lp := getLatestPost(hp)
		if lp == "none" {
			fmt.Println("See error: empty result")
			return
		}

		keyword := strings.ToLower(lp)
		x := len(keyword)
		if strings.Contains(lp, "\n") {
			x = strings.Index(lp, "\n")
		}

		abridgedPost := lp[:x]

		if strings.Contains(keyword, "weeb") {
			msg := "Weebs are trash"
			callback(msg)
		} else if strings.Contains(keyword, "terminator") {
			msg := convertText(abridgedPost)
			callback(msg)
		} else if strings.Contains(keyword, "inna woods") {
			msg := convertText(abridgedPost)
			callback(msg)
		} else if strings.Contains(keyword, "1911") {
			msg := convertText("two world wars")
			callback(msg)
		} else if strings.Contains(keyword, "shill") || strings.Contains(keyword, "profit") {
			msg := "Thanks to Coinbase resiliency and UFC NFTs, crypto is now linked directly to my Wells Fargo account"
			callback(msg)
		} else {
			msg := convertText(abridgedPost)
			msg += "\n\nUgandan baboon"
			callback(msg)
		}
		temp = hp.HighestPost + 1
		time.Sleep(300 * time.Minute)
	}
}

func getHighestPost() postCount {
	url := "https://forum.pixelspace.xyz/t/" + thread + ".json"
	var result postCount

	client := http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return result
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Api-Username", apiUser)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return result
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return result
	}

	err = json.Unmarshal(body, &result)

	if err != nil {
		fmt.Println(err)
		return result
	}

	res.Body.Close()

	return result
}

func getLatestPost(pc postCount) string {
	url := "https://forum.pixelspace.xyz/t/" + thread + "/" + strconv.Itoa(pc.HighestPost) + ".json"
	client := http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return "none"
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Api-Username", apiUser)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "none"
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "none"
	}

	var latest latestPost
	err = json.Unmarshal(body, &latest)

	if err != nil {
		fmt.Println(err)
		return "none"
	}

	currentPost := latest.PostStream.Posts[len(latest.PostStream.Posts)-1].Cooked
	return currentPost
}

func convertText(statement string) string {
	var newStatement []string
	var msg string

	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	symbols := "'\"@!.<>"

	for ind, letter := range statement {
		if strings.Contains(letters, string(letter)) {
			if ind%2 == 0 {
				newStatement = append(newStatement, strings.ToLower(string(letter)))
			} else {
				newStatement = append(newStatement, strings.ToUpper(string(letter)))
			}
		} else if strings.Contains(symbols, string(letter)) {
			newStatement = append(newStatement, string(letter))
		} else if strings.Contains(" ", string(letter)) {
			newStatement = append(newStatement, " ")
		}
	}
	for i := 0; i < len(newStatement); i++ {
		msg += newStatement[i]
	}

	return msg
}

func callback(message string) {
	jsonBody := []byte(`{"topic_id": "%s", "raw": "%s"}`)
	s := fmt.Sprintf(string(jsonBody), thread, message)
	body := bytes.NewReader([]byte(s))

	url := fmt.Sprintf("https://forum.pixelspace.xyz/posts.json")
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Api-Username", apiUser)

	client := http.Client{Timeout: 30 * time.Second}
	_, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
}
