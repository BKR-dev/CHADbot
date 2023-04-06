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

	for {
		responses, botUserId, err := api.GetLastPost()
		if err != nil {
			scribe.Errorf("Error getting responses from Topic: ", err)
		}

		highestPostNumber = responses.HighestPostNumber

		for _, p := range responses.PostStream.Posts {

			if p.PostNumber == highestPostNumber {

				if p.UserID == botUserId {
					scribe.Infof("Last Post is from Bot - 3 sec sleepy time")
					time.Sleep(3 * time.Second)
					break
				}

				scribe.Infof("Post: %v", p.Cooked)
				response, err = answer.GetResponse(p.Cooked, p.Username, p.UserTitle)
				scribe.Infof("Response: ", response)
				if err != nil {
					scribe.Errorf("Error getting response: ", err)
				}

			}
		}

		err = api.PostResponseToTopic(response)
		if err != nil {
			scribe.Error(err)
		}

		time.Sleep(3 * time.Second)
	}
}
