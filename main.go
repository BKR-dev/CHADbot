package main

import (
	"os"
	"terminator-shitpost/handler"
	"terminator-shitpost/logging"
	"terminator-shitpost/responses"
)

func main() {

	var response string

	logger := logging.NewLogger(logging.DebugLevel, os.Stdout)

	post, err := handler.GetLastPost(logger)
	if err != nil {
		logger.Error("Could not get last Post")
	}

	response, err = responses.GetResponse(logger, post)
	if err != nil {
		logger.Error(err.Error())
	}

	err = handler.PostResponseToTopic(logger, response)
	if err != nil {
		logger.Error(err.Error())
	}
}
