package responses

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	api "terminator-shitpost/apihandler"
	"terminator-shitpost/logging"

	"golang.org/x/net/html"
)

// response struct
// TODO: add bool "mock" for extra mocking
// TODO: add bool "insult" for extra spice
type InsultingResponse struct {
	topic     string
	keywords  []string
	responses []string
}

var Scribe logging.Logger

func init() {
	var err error
	Scribe, err = logging.NewLogger()
	if err != nil {
		fmt.Println("Error creating logger in responses.go")
	}
}

// From Post, looks up keywords and matches unto existing responses
func GetResponse(post string) (string, error) {

	Scribe.Infof("Cleaning up Post: %v", post)
	cleanPost := extractTextFromHTML(post)
	keywords := strings.Split(cleanPost, " ")
	Scribe.Infof("Looking for keywords: %v", keywords)

	answer, err := getRandomResponse(keywords)
	if err != nil {
		return "", nil
	}
	if answer == "no match found" {
		Scribe.Infof("No keyword from post matched: ", answer)
		return "", errors.New("No match between post content and keywords")
	}
	Scribe.Infof("Using response: %v", answer)
	return answer, nil
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

func getRandomResponse(keyword []string) (string, error) {
	var err error
	// all trigger for a response
	trigger := []string{"weeb", "anime", "glock", "1911", "bible", "shill", "crypto", "bitcoin",
		"etherium", "vegan", "keto", "linux", "macos", "windows", "fed", "fbi", "cia", "atf",
		"aft", "vision", "(((", ")))", "propaganda", "1984", "innawoods", "inna woods", "gaming",
		"gayming"}

	var match string
	var snarckyResponse string

	// all responses per topic
	weebAnime := InsultingResponse{
		topic:    "weebAnime",
		keywords: []string{"anime", "weeb"},
		responses: []string{"weebs are retarded and anime is trash", "2D women are for brainlets and retards",
			"anime is cringe and fake touch some grass"},
	}

	innawoods := InsultingResponse{
		topic:    "innawoods",
		keywords: []string{"innawoods", "inna woods"},
		responses: []string{"forever alone in the woods", "pissing in jars to keep some company",
			"getting buttfucked by the local wendingo", "starving in the cold is better than buying starbucks"},
	}

	nineteen11 := InsultingResponse{
		topic:     "1911",
		keywords:  []string{"1911"},
		responses: []string{"two world wars", "chinese made glocks are more versatile"},
	}

	shillCrypto := InsultingResponse{
		topic:    "shill",
		keywords: []string{"shill", "crypto", "bitcoin", "etherium", "NFT"},
		responses: []string{
			"Thanks to Coinbase resiliency and UFC NFTs, crypto is now linked directly to my Wells Fargo account :chris_party:",
			"My NTF ETF i just shorted got me the platinum AMEX so i get paid for spending money i dont have :think_about_it:",
			"bro invest in my hyper value adding NFT web4.2 finTech renewable ecommerce gaming startup, bro"},
	}

	diet := InsultingResponse{
		topic:    "diet",
		keywords: []string{"vegan", "keto"},
		responses: []string{"just eat some real food and stop being a cunt",
			"stop pretending you are not being brainwashed by some cucked incels to buy there supplements"},
	}

	opSys := InsultingResponse{
		topic:     "opSys",
		keywords:  []string{"linux", "macos", "windows"},
		responses: []string{"stop being such a poor and use a real OS", "you got your programming socks already?"},
	}

	feds := InsultingResponse{
		topic:     "feds",
		keywords:  []string{"fed", "FBI", "CIA", "ATF", "AFT"},
		responses: []string{"We've had enough, time to blow this fucker up!", "Wow you're so cool! Go, commit a crime :hugs:"},
	}
	// responses slice of structs
	allResponses := []InsultingResponse{weebAnime, innawoods, nineteen11, shillCrypto, diet, opSys, feds}

	// match trigger and keyword
findMatch:
	for _, trig := range trigger {
		for _, key := range keyword {
			if trig == strings.ToLower(key) {
				match = trig
				break findMatch
			} else if strings.ToLower(key) == "bible" {
				match = "bible"
			} else {
				match = "null"
			}
		}
	}

	// match keywords
findResponse:
	for _, strc := range allResponses {
		for _, key := range strc.keywords {
			if key == match {
				snarckyResponse = returnResponseFromSlice(strc.responses)
				break findResponse
			} else if match == "bible" {
				snarckyResponse, err = api.GetRandomBibleVerse()
				if err != nil {
					return "", err
				}
				break findResponse
			} else if match == "null" {
				snarckyResponse = "no match found"
			}
		}
	}

	return snarckyResponse, nil
}

func returnResponseFromSlice(responses []string) string {
	return responses[rand.Intn((len(responses)-1)+1)]
}

func extractTextFromHTML(htmlString string) string {
	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		Scribe.Errorf("Error removing HTML: %v", err)
	}
	var f func(*html.Node)
	var text string
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			text += n.Data
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return text
}
