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

// InsultResponse
type InsultingResponse struct {
	topic     string   // name of the InsultResponse
	keywords  []string // slice of keywords that yield a respone
	responses []string // slice of spicy responses
}

// Insults
type Insults struct {
	personalAttack      []string // you sodding tiktak
	personalTitleAttack []string // insults you for/with your title
}

type RandomInsult struct {
	mock   bool // trigger for mocking statement
	insult bool // trigger for personal attack
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
// adds user name and user title to randomly add them for a more personal insult
func GetResponse(post string, username string, usertitle string) (string, error) {

	Scribe.Infof("Cleaning up Post: %v", post)
	// remove html tags from post
	cleanPost := extractTextFromHTML(post)
	// split post into words for being used at akeywords
	keywords := strings.Split(cleanPost, " ")
	Scribe.Infof("Looking for keywords: %v", keywords)
	// get
	answer, err := getRandomResponse(keywords, username, usertitle)
	if err != nil {
		return "", nil
	}

	if answer == nil {
		Scribe.Infof("No keyword from post matched: ", answer)
		return "", errors.New("No match between post content and keywords")
	}

	botAnswer := strings.Join(answer, " ")

	Scribe.Infof("Using response: %v", answer)
	return botAnswer, nil
}

// return a random response when trigger is detected
func getRandomResponse(keyword []string, username string, usertitle string) ([]string, error) {
	var err error
	var responseFull []string
	// insults to add to the response
	insults := Insults{
		personalAttack: []string{
			"you sodding tiktak", "absolute muppet",
			"complete retard", "cum guzzling lunatic", "beyond meat enjoyer", "incel cuck",
			"vegan soy brainlet", "just shut the fuck up fag", "first grade degenrate"},
		personalTitleAttack: []string{"%s is your title!? Fitting...", "%s, yeah right, my ass!"},
	}
	// all trigger for a response
	trigger := []string{
		"weeb", "anime", "glock", "1911", "bible", "shill", "crypto", "bitcoin",
		"etherium", "vegan", "keto", "linux", "macos", "windows", "fed", "fbi", "cia", "atf",
		"aft", "vision", "(((", ")))", "propaganda", "1984", "innawoods", "inna woods",
		"gaming", "gayming"}

	var match string
	var snarckyResponse string
	var rndmInsult RandomInsult

	allResponses := provideResponses()

	// match trigger and keyword
findMatch:
	for _, trig := range trigger {
		for _, key := range keyword {
			if trig == strings.ToLower(key) {
				match = trig
				Scribe.Infof("Found match: %v", match)
				break findMatch
			} else if strings.ToLower(key) == "bible" {
				Scribe.Infof("Found match: bible")
				match = "bible"
			} else {
				match = ""
			}
		}
	}

	if len(match) < 1 {
		return nil, errors.New("found no match")
	}

	// match keywords
findResponse:
	for _, respo := range allResponses {
		for _, key := range respo.keywords {
			if key == match {
				snarckyResponse = randomStringFromSlice(respo.responses)
				break findResponse
				// special case for bible as response comes from an API
			} else if match == "bible" {
				snarckyResponse, err = api.GetRandomBibleVerse()
				if err != nil {
					return nil, err
				}
				break findResponse
			} else if match == "" {
				return nil, errors.New("found no match")
			}
		}
	}

	responseFull = append(responseFull, snarckyResponse)

	/*
		Anatomy if a Response:
		when mock is true - convertText
		snarckyResponse = fmt.Sprintf("%v, %v", convertText(snarckyResponse))
		when insult is true - response + personal insult
		snarckyResponse = fmt.Sprintf("%v, %v", snarckyResponse, personalAttack)
	*/

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

// stores all responses
func provideResponses() []InsultingResponse {

	/* InsultResponse Struct
	//type InsultingResponse struct {
	// 	topic     string   // name of the InsultResponse
	// 	keywords  []string // slice of keywords that yield a respone
	// 	responses []string // slice of spicy responses
	// 	mock      bool     // trigger for mocking statement
	// 	insult    bool     // trigger for personal attack }
	*/

	// all responses per topic
	weebAnime := InsultingResponse{
		topic:    "weebAnime",
		keywords: []string{"anime", "weeb"},
		responses: []string{"weebs are retarded and anime is trash", "2D women are for brainlets and retards",
			"anime is cringe and fake, go and touch some grass"},
	}

	innawoods := InsultingResponse{
		topic:    "innawoods",
		keywords: []string{"innawoods", "inna woods"},
		responses: []string{"forever alone in the woods", "pissing in jars to keep some company",
			"getting buttfucked by the local wendingo", "starving in the cold is better than buying starbucks soy caramel faggoccino"},
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
			"My NFT ETF i just shorted got me the platinum AMEX so i get paid for spending money i dont have :think_about_it:",
			"bro invest in my hyper value adding NFT web4.2 finTech renewable ecommerce gaming startup, bro"},
	}

	diet := InsultingResponse{
		topic:    "diet",
		keywords: []string{"vegan", "keto"},
		responses: []string{"just eat some real food and stop being a cunt",
			"stop pretending you are not being brainwashed by some cucked incels to buy their supplements"},
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

	mockingParantheses := InsultingResponse{
		topic:     "mockingParantheses",
		keywords:  []string{"(((", ")))"},
		responses: []string{},
	}

	vidja := InsultingResponse{
		topic:    "videogames",
		keywords: []string{"gaming", "gayming", "vidja", "videogames"},
		responses: []string{"grow up already, you are not a child anymore, just let it go",
			"you played that!? wow amazing, wayyyyy better than watching netflix to numb your brain right!?"},
	}
	// responses slice of structs
	allResponses := []InsultingResponse{weebAnime, innawoods, nineteen11, shillCrypto, diet, opSys, feds, mockingParantheses, vidja}

	return allResponses
}

// randomly sets bool fields
func randomBoolSet(rndmI RandomInsult) RandomInsult {
	rndmI.mock = rand.Intn(2) == 1
	rndmI.insult = rand.Intn(2) == 1
	return rndmI
}
