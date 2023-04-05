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

	for {
		responses, err := api.GetLastPost()
		if err != nil {
			scribe.Errorf("Error getting responses from Topic: ", err)
		}

		for _, p := range responses.PostStream.Posts {
			if p.PostNumber == responses.HighestPostNumber {

				scribe.Infof("Post: %v", p.Cooked)
				response, err = answer.GetResponse(p.Cooked)
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
