package main

import (
	"os"

	"errors"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

const QUOTE_API = "http://ron-swanson-quotes.herokuapp.com/v2/quotes"

var (
	log = logrus.New()
)

type Response struct {
	Quote string `json:"quote"`
}

type Request struct {
	ChannelID string `json:"channelID"`
}

func HandleLambdaEvent(request Request) (Response, error) {
	if len(request.ChannelID) == 0 {
		return Response{}, errors.New("required parameter not set: ChannelID")
	}

	token, ok := os.LookupEnv("SLACK_TOKEN")
	if !ok {
		return Response{}, errors.New("required configuration missing: SLACK_TOKEN")
	}

	api := slack.New(token, slack.OptionDebug(true))

	var quotes []string
	client := resty.New()
	_, err := client.R().
		EnableTrace().
		SetResult(&quotes).
		Get(QUOTE_API)

	if err != nil {
		return Response{}, err
	}

	if len(quotes) == 0 {
		log.Fatal("found no quotes")
		return Response{}, err
	}

	quote := quotes[0]

	if _, _, _, err := api.JoinConversation(request.ChannelID); err != nil {
		return Response{}, err
	}

	channelID, timestamp, err := api.PostMessage(request.ChannelID, slack.MsgOptionText(quote, true))
	if err != nil {
		return Response{}, err
	}

	log.Infof("message successfully sent to channel %s at %s", channelID, timestamp)

	return Response{Quote: quote}, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
