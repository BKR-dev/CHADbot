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

type RandomInsult struct {
	mock   bool // trigger for mocking statement, randomly set
	insult bool // trigger for personal attack, randomly set
}

var scribe logging.Logger

func init() {
	var err error
	scribe, err = logging.NewLogger()
	if err != nil {
		fmt.Println("Error creating logger in responses.go")
	}
}

// From Post, looks up keywords and matches unto existing responses
// adds user name and user title to randomly add them for a more personal insult
func GetResponse(post string, username string, usertitle string) (string, error) {

	scribe.Infof("Cleaning up Post: %v", post)
	// remove html tags from post
	cleanPost := extractTextFromHTML(post)
	// split post into words for being used at akeywords
	keywords := strings.Split(cleanPost, " ")
	scribe.Infof("Looking for keywords: %v", keywords)
	// get
	answer, err := getRandomResponse(keywords, username, usertitle)
	if err != nil {
		return "", nil
	}

	if answer == nil {
		scribe.Infof("No keyword from post matched: ", answer)
		return "", errors.New("No match between post content and keywords")
	}

	botAnswer := strings.Join(answer, ", ")

	scribe.Infof("Using response: %v", answer)
	return botAnswer, nil
}

// return a random response when trigger is detected
func getRandomResponse(keyword []string, username string, usertitle string) ([]string, error) {
	var err error
	var responseFull []string
	var allMatches []string
	var snarckyResponse string
	var rndmInsult RandomInsult

	if len(keyword) < 1 {
		return nil, errors.New("keywords empty")
	}

	allResponses, insults, trigger := ProvideResponsesAndInsults()
	scribe.Infof("keywords from responses: %v", trigger)
	// match trigger and keyword and add to slice for all matches

	for _, trig := range trigger {
		for _, key := range keyword {
			if trig == strings.ToLower(key) {
				allMatches = append(allMatches, trig)
				scribe.Infof("Found match: %v", trig)
			} else if strings.ToLower(key) == "bible" {
				scribe.Infof("Found match: bible")
				allMatches = append(allMatches, "bible")
			}
		}
	}

	scribe.Infof("allMatches found: %v", allMatches)

	// if len allMatches is empty, we is done
	if len(allMatches) < 1 {
		return nil, errors.New("found no match")
		// bible case
	} else if allMatches[0] == "bible" {
		snarckyResponse, err = api.GetRandomBibleVerse()
		if err != nil {
			return nil, err
		}
	} else if len(allMatches) > 1 {
		// random match from allMatches
		match := randomStringFromSlice(allMatches)
	findResponse:
		for _, r := range allResponses {
			for _, k := range r.keywords {
				if k == match {
					snarckyResponse = randomStringFromSlice(r.responses)
					scribe.Infof("Found response %v for keyword %v", snarckyResponse, k)
					break findResponse
				}
			}
		}
	}

	responseFull = append(responseFull, snarckyResponse)

	rndmInsult = randomBoolSet(rndmInsult)
	// response + mock + insult
	if rndmInsult.mock && rndmInsult.insult {
		responseFull[0] = convertText(snarckyResponse)
		responseFull = append(responseFull, randomStringFromSlice(insults.personalAttack))
	}
	//  response + mock
	if rndmInsult.mock && !rndmInsult.insult {
		responseFull[0] = convertText(snarckyResponse)
	}
	// response + insult
	if !rndmInsult.mock && rndmInsult.insult {
		responseFull = append(responseFull, randomStringFromSlice(insults.personalAttack))
		responseFull = append(responseFull, username)
	}

	return responseFull, nil
}

// returns a random response from response slice
func randomStringFromSlice(responses []string) string {
	return responses[rand.Intn((len(responses)-1)+1)]
}

// extracts all html tags from input and returns HTML tag free string
func extractTextFromHTML(htmlString string) string {
	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		scribe.Errorf("Error removing HTML: %v", err)
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

// transform a message to A mEsSaGe
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

// randomly sets bool fields to randomInsult
func randomBoolSet(rndmI RandomInsult) RandomInsult {
	rndmI.mock = rand.Intn(2) == 1
	rndmI.insult = rand.Intn(2) == 1
	return rndmI
}
