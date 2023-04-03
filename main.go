package main

import (
	"fmt"
	"sync"
	api "terminator-shitpost/apihandler"
)

var wg sync.WaitGroup
var mu sync.Mutex

func main() {

	dataCh := make(chan api.TopicResponse)
	errCh := make(chan error)
	stopCh := make(chan bool)

	go api.GetLastPost(dataCh, errCh)

	for {

		select {

		// case data := <-dataCh:
		// 	fmt.Println("TopicResponse: ", data)

		case err := <-errCh:
			fmt.Println("Error: ", err)

		case <-stopCh:
			return

		default:
			data := <-dataCh
			fmt.Println(data)

		}

	}

	// logger, err := logging.NewLogger()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// logger.Infof("Started Terminator Bot")

	// logger.Infof("Getting last post")
	// post, err := api.GetLastPost()
	// if err != nil {
	// 	logger.Errorf("Could not get last Post: %v", err)
	// }

	// logger.Infof("Getting a response for post: %v", post)
	// response, err = responses.GetResponse(post)
	// if err != nil {
	// 	logger.Errorf("Could not get a response: %v", err)
	// }
	// logger.Infof("Sending out response: %v", response)
	// err = api.PostResponseToTopic(response)
	// if err != nil {
	// 	logger.Errorf("Could not post response to topic: %v", err)
	// }

	// err = logger.Close()
	// if err != nil {
	// 	fmt.Println(err)
	// }
}
