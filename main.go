package main

import (
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

		for _, p := range responses.PostStream.Posts {

			if p.PostNumber == highestPostNumber {

				if p.UserID == botUserId {
					scribe.Infof("Last Post is from Bot ")
					time.Sleep(time.Duration(botTimeout) * time.Second)
					botTimeout++
					scribe.Infof("Sleeping for %d", botTimeout)
					break
				}

				scribe.Infof("Post: %v", p.Cooked)
				response, err = answer.GetResponse(p.Cooked, p.Username, p.UserTitle)
				scribe.Infof("Response: %v", response)
				if err != nil {
					scribe.Errorf("Error getting response: ", err)
				}
				if response == "" {
					time.Sleep(time.Duration(botTimeout) * time.Second)
					scribe.Infof("Sleeping for %d", botTimeout)
				}
			}
			botTimeout++
		}

		err = api.PostResponseToTopic(response)
		if err != nil {
			scribe.Error(err)
		}

		botTimeout = 3
		time.Sleep(5 * time.Second)
	}
}
