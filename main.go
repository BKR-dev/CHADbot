package main

import (
	api "terminator-shitpost/apihandler"
	conf "terminator-shitpost/conf"
	scribe "terminator-shitpost/logging"
	answer "terminator-shitpost/responses"
	"time"
)

// TODO: add config file for credentials
// TODO: att a frontend for user interaction (start / stop / last logg / add insults)
// TODO: add go routines just for the heck of it
func main() {

	scribe, err := scribe.NewLogger()
	if err != nil {
		panic(err)
	}

	conf.GetSettings()

	var response string
	var highestPostNumber int
	botTimeout := 3

	for {
		responses, botUserId, err := api.GetLastPost()
		if err != nil {
			scribe.Errorf("Error getting responses from Topic: ", err)
		}

		highestPostNumber = responses.HighestPostNumber

		//for _, p := range responses.PostStream.Posts {
		// only highest Post aka latest is getting looked at
		currentPost := responses.PostStream.Posts[len(responses.PostStream.Posts)-1].PostNumber
		currentResp := responses.PostStream.Posts[len(responses.PostStream.Posts)-1]
		if currentPost == highestPostNumber {
			//if p.PostNumber == highestPostNumber {

			// if the last user was the bot, sleepy time
			if currentResp.UserID == botUserId {
				scribe.Infof("Last Post is from Bot ")
				time.Sleep(time.Duration(botTimeout) * time.Second)
				botTimeout++
				scribe.Infof("Sleeping for %d", botTimeout)
				continue
			}

			// get a response for the last post
			scribe.Infof("Post: %v", currentResp.Cooked)
			response, err = answer.GetResponse(currentResp.Cooked, currentResp.Username, currentResp.UserTitle)
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
		botTimeout++

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
