package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

const thread = "1113"
const apiKey = "5634da9f596ecc2740440a75499176a3b8181752aa418696b61ed08b982c3a43"
const apiUser = "terminator"

func main() {
	file, err := os.OpenFile("shitpost.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.SetOutput(file)

	for {
		hp := getHighestPost()

		if hp == (postCount{}) || hp.HighestPost == 0 {
			log.Println("Empty, sleeping it off...")
			time.Sleep(2 * time.Minute)
			continue
		}

		lp := getLatestPost(hp.HighestPost)
		if err != nil {
			fmt.Println(err)
			return
		}

		if lp == "none" {
			log.Println("latestPost: empty result")
			time.Sleep(5 * time.Second)
			continue
		}

		keyword := strings.ToLower(lp)

		if strings.Contains(keyword, "weeb") {
			log.Println("Responding to weeb")
			msg := "Anime is trash"
			// adding an insult to the message
			callback(msg)
		} else if strings.Contains(keyword, "inna woods") {
			log.Println("Responding to inna woods")
			msg := convertText("forever alone in the woods")
			callback(msg)
		} else if strings.Contains(keyword, "1911") {
			log.Println("Responding to 1911")
			msg := convertText("two world wars")
			callback(msg)
		} else if strings.Contains(keyword, "bible") && strings.Contains(keyword, "verse") {
			log.Println("Responding to bible/verse")
			msg := getRandomBibleVerse()
			callback(msg)
		} else if strings.Contains(keyword, "shill") || strings.Contains(keyword, "profit") {
			log.Println("Responding to shill/profit")
			msg := "Thanks to Coinbase resiliency and UFC NFTs, crypto is now linked directly to my Wells Fargo account :chris_party:"
			callback(msg)
			// new keyword "vegan" and "keto"
		} else if strings.Contains(keyword, "vegan") || strings.Contains(keyword, "keto") {
			log.Println("Responding to vegan/keto")
			msg := "just eat some real food and stop being a cunt"
			callback(msg)
			// new keyword "linux"
		} else if strings.Contains(keyword, "linux") {
			log.Println("Responding to linux")
			msg := "stop being such a poor and use a real OS"
			callback(msg)
		} else if strings.Contains(keyword, "the fed") {
			log.Println("Responding to the fed")
			msg := "We've had enough, time to blow this fucker up"
			callback(msg)
		} else if strings.Contains(keyword, "vision") {
			log.Println("Responding to vision")
			msg := "Fatwood is a cunt"
			callback(msg)
		} else {
			time.Sleep(3 * time.Second)
		}
	}
}

func getRandomBibleVerse() string {
	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", "https://labs.bible.org/api/?passage=random", nil)
	if err != nil {
		log.Println(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	return string(body)
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
		log.Println("Error with new request")
		fmt.Println(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Api-Username", apiUser)

	client := http.Client{Timeout: 30 * time.Second}
	_, err = client.Do(req)
	if err != nil {
		log.Println("Error with client")
		fmt.Println(err)
	}
}
