package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
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

const thread = "1016"
const apiKey = "5634da9f596ecc2740440a75499176a3b8181752aa418696b61ed08b982c3a43"
const apiUser = "terminator"

func main() {
	temp := 0

	file, err := os.OpenFile("shitpost.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.SetOutput(file)

	for {
		hp := getHighestPost()
		if hp.HighestPost == temp {
			log.Println("No new posts...")
			time.Sleep(5 * time.Second)
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
			log.Println("Responding to weeb")
			msg := "Weebs are trash"
			// adding an insult to the message
			msg = add_insult(msg)
			callback(msg)
			temp = hp.HighestPost + 1
		} else if strings.Contains(keyword, "terminator") {
			log.Println("Responding to terminator")
			msg := convertText(abridgedPost)
			callback(msg)
			temp = hp.HighestPost + 1
		} else if strings.Contains(keyword, "inna woods") || strings.Contains(keyword, "innawoods") {
			log.Println("Responding to inna woods")
			msg := convertText(abridgedPost)
			callback(msg)
			temp = hp.HighestPost + 1
		} else if strings.Contains(keyword, "1911") {
			log.Println("Responding to 1911")
			msg := convertText("two world wars")
			callback(msg)
			temp = hp.HighestPost + 1
		} else if strings.Contains(keyword, "shill") || strings.Contains(keyword, "profit") {
			log.Println("Responding to shill/profit")
			msg := "Thanks to Coinbase resiliency and UFC NFTs, crypto is now linked directly to my Wells Fargo account :chris_party:"
			callback(msg)
			temp = hp.HighestPost + 1
			// new keyword "vegan" and "keto"
		} else if strings.Contains(keyword, "vegan") || strings.Contains(keyword, "keto") {
			log.Println("Responding to vegan/keto")
			msg := "just eat some real food and stop being a douche"
			callback(msg)
			temp = hp.HighestPost + 1
			// new keyword "linux"
		} else if strings.Contains(keyword, "linux") {
			log.Println("Responding to linux")
			msg := "stop being such a poor and use a real OS"
			callback(msg)
			temp = hp.HighestPost + 1
		} else if strings.Contains(keyword, "img") {
			log.Println("Image detected...")
			continue
		} else {
			time.Sleep(5 * time.Second)
		}
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

func add_insult(message string) string {
	// setup slice of insults
	insult := []string{"listen up retard ", " piece of fuck", " you absolute muppet", " you hopeless degenerate"}
	// getting random int 0-3
	rndm := rand.Intn((4 - 1) + 1)
	// adding the insults to the message
	if rndm < 1 {
		message = insult[rndm] + message
	} else {
		message = message + insult[rndm]
	}

	return message
}
