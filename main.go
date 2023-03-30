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

const thread = "1118"
const apiKey = "5634da9f596ecc2740440a75499176a3b8181752aa418696b61ed08b982c3a43"
const apiUser = "terminator"

func main() {
	file, err := os.OpenFile("shitpost.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.SetOutput(file)
	now := time.Now()
	now.Format("2006-01-02 15:04:05")
	log.Print("\n\n\n\nStarting up at %v", now)

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
		// TODO: FIX html tags
		userPost := strings.Split(lp, " ")
		now = time.Now()
		log.Printf("UserPost: %v\n", userPost)
		botResponse := getRandomResponse(userPost)
		log.Printf("response from Bot: %v\n", botResponse)
		callback(botResponse)
		time.Sleep(10 * time.Second)

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

func getRandomResponse(keyword []string) string {
	log.Println("starting to find a response")
	// all trigger for a response
	trigger := []string{"weeb", "anime", "glock", "1911", "bible", "shill", "crypto", "bitcoin",
		"etherium", "vegan", "keto", "linux", "macos", "windows", "fed", "fbi", "cia", "atf",
		"aft", "vision", "(((", ")))", "propaganda", "1984", "innawoods", "inna woods", "gaming", "gayming"}

	// response struct
	type TopicResponses struct {
		topic     string
		keywords  []string
		responses []string
	}

	var match string
	var snarckyResponse string

	// all responses per topic
	weebAnime := TopicResponses{
		topic:    "weebAnime",
		keywords: []string{"anime", "weeb"},
		responses: []string{"weebs are retarded and anime is trash", "2D women are for brainlets and retards",
			"anime is cringe and fake touch some grass"},
	}

	innawoods := TopicResponses{
		topic:    "innawoods",
		keywords: []string{"innawoods", "inna woods"},
		responses: []string{"forever alone in the woods", "pissing in jars to keep some company",
			"getting buttfucked by the local wendingo", "starving in the cold is better than buying starbucks"},
	}

	nineteen11 := TopicResponses{
		topic:     "1911",
		keywords:  []string{"1911"},
		responses: []string{"two world wars", "chinese made glocks are more versatile"},
	}

	shillCrypto := TopicResponses{
		topic:    "shill",
		keywords: []string{"shill", "crypto", "bitcoin", "etherium", "NFT"},
		responses: []string{"Thanks to Coinbase resiliency and UFC NFTs, crypto is now linked directly to my Wells Fargo account :chris_party:",
			"My NTF ETF i just shorted got me the platinum AMEX so i get paid for spending money i dont have :think_about_it:",
			"bro invest in my hyper value adding NFT web4.2 finTech renewable ecommerce gaming startup, bro"},
	}

	diet := TopicResponses{
		topic:    "diet",
		keywords: []string{"vegan", "keto"},
		responses: []string{"just eat some real food and stop being a cunt",
			"stop pretending you are not being brainwashed by some cucked incels to buy there supplements"},
	}

	opSys := TopicResponses{
		topic:     "opSys",
		keywords:  []string{"linux", "macos", "windows"},
		responses: []string{"stop being such a poor and use a real OS", "you got your programming socks already?"},
	}

	feds := TopicResponses{
		topic:     "feds",
		keywords:  []string{"fed", "FBI", "CIA", "ATF", "AFT"},
		responses: []string{"We've had enough, time to blow this fucker up!", "Wow you're so cool! Go, commit a crime :hugs:"},
	}
	// responses slice of structs
	allResponses := []TopicResponses{weebAnime, innawoods, nineteen11, shillCrypto, diet, opSys, feds}

	// match trigger and keyword
findMatch:
	for _, trig := range trigger {
		log.Println("trigger: " + trig)
		for _, key := range keyword {
			log.Println("key: " + key)
			if trig == strings.ToLower(key) {
				match = trig
				log.Println("FOUND A MATCH!!!!")
				break findMatch
			}
		}
	}

	// match keywords
findResponse:
	for _, strc := range allResponses {
		for _, key := range strc.keywords {
			if key == match {
				log.Println(snarckyResponse)
				snarckyResponse = returnResponseFromSlice(strc.responses)
				break findResponse
			}
		}
	}

	return snarckyResponse
}

func returnResponseFromSlice(responses []string) string {
	return responses[rand.Intn((len(responses)-1)+1)]
}
