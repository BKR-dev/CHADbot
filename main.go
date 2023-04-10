package main

import (
	"math/rand"
	api "terminator-shitpost/apihandler"
	scribe "terminator-shitpost/logging"
	answer "terminator-shitpost/responses"
	"time"
)

func main() {

	scribe, err := scribe.NewLogger()
	if err != nil {
		panic(err)
	}

	var response string
	var highestPostNumber int
	botTimeout := 3

	for {
		responses, botUserId, err := api.GetLastPost()
		if err != nil {
			scribe.Errorf("Error getting responses from Topic: ", err)
		}

		highestPostNumber = responses.HighestPostNumber

	postresponse:
		for _, p := range responses.PostStream.Posts {
			// more randomness
			rand.Seed(time.Now().UnixNano())
			// only highest Post aka latest is getting looked at
			if p.PostNumber == highestPostNumber {

				// if the last user was the bot, sleepy time
				if p.UserID == botUserId {
					scribe.Infof("Last Post is from Bot ")
					time.Sleep(time.Duration(botTimeout) * time.Second)
					break postresponse
				}

				// FIXME: when no new post is created or last post is from bot no resposne should be sent
				// get a response for the last post
				scribe.Infof("Post: %v", p.Cooked)
				response, err = answer.GetResponse(p.Cooked, p.Username, p.UserTitle)
				scribe.Infof("Response: %v", response)
				if err != nil {
					scribe.Errorf("Error getting response: ", err)
				}
				// if the response is empty, no keyword, leepy time
				if response == "" {
					scribe.Infof("Sleeping for %d", botTimeout)
					time.Sleep(time.Duration(botTimeout) * time.Second)
				}
			}
		}

		// actually post response to topic
		scribe.Infof("Sending response: %v", response)
		err = api.PostResponseToTopic(response)
		if err != nil {
			scribe.Error(err)
		}
		scribe.Infof("Sleeping for 5 seconds")
		time.Sleep(5 * time.Second)
	}
}
