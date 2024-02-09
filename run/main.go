package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/tae2089/bob-logging/logger"
	web "github.com/tae2089/webpage-check"
	"go.uber.org/zap"
)

func main() {
	r, err := web.GetWebCheck("/scripts/check.yaml")
	if err != nil {
		logger.Error(err)
		return
	}
	slackWebHookURL := os.Getenv("SLACK_WEB_HOOK_URL")
	if slackWebHookURL == "" {
		logger.Error(errors.New("SLACK_WEB_HOOK_URL is not set"))
		return
	}
	for _, homePage := range r.HomePage {
		for _, page := range homePage.Path {
			url := fmt.Sprintf("%s%s", homePage.Host, page)
			resp, err := http.Get(url)
			message := ""
			if err != nil {
				logger.Error(err, zap.String("url", url))
				message = fmt.Sprintf("this page is not healthy - url:%s  | error:%s ", url, err.Error())
				sendSlackMessage(slackWebHookURL, message)
				continue
			}
			if resp.StatusCode >= 400 {
				message = fmt.Sprintf("this page is not healthy - %s and status code is %d", url, resp.StatusCode)
				logger.Error(errors.New(message))
				sendSlackMessage(slackWebHookURL, message)
			}
		}
	}
	logger.Info("home urls is healthy")
}

func sendSlackMessage(slackWebHookURL string, message string) error {
	slackMessage := fmt.Sprintf(`{"text": "%s"}`, message)
	reqBody := bytes.NewBufferString(slackMessage)
	http.Post(slackWebHookURL, "application/json", reqBody)
	return nil
}
