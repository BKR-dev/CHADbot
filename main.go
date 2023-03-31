package main

import (
	"fmt"
	"terminator-shitpost/handler"
	"terminator-shitpost/logging"
	"terminator-shitpost/responses"
)

func main() {

	var response string

	logger, err := logging.NewLogger()
	if err != nil {
		fmt.Println(err)
	}
	logger.Infof("Started Terminator Bot")

	logger.Infof("Getting last post")
	post, err := handler.GetLastPost()
	if err != nil {
		logger.Errorf("Could not get last Post: %v", err)
	}

	logger.Infof("Getting a response for post: %v", post)
	response, err = responses.GetResponse(post)
	if err != nil {
		logger.Errorf("Could not get a response: %v", err)
	}
	logger.Infof("Sending out response: %v", response)
	err = handler.PostResponseToTopic(response)
	if err != nil {
		logger.Errorf("Could not post response to topic: %v", err)
	}

	err = logger.Close()
	if err != nil {
		fmt.Println(err)
	}
}
