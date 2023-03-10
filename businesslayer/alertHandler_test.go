package businesslayer

import (
	"os"
	"testing"

	"github.com/arizon-dread/status-checker-api/models"
)

func TestSendStatusHappyPath(t *testing.T) {
	//{"name":"status-checker-api", "callUrl": "http://localhost:8080/healthz", "httpMethod": "GET", "ResponseMatch": "Healthy",
	//"alertBody": "shit went south", "alertUrl": "http://dev.null.com", "alertEmail": "erik.j.svensson@gmail.com"}'
	body := "message body"
	system := models.Systemstatus{
		ID:               1,
		Name:             "google",
		CallUrl:          "https://google.com",
		CallStatus:       "OK",
		HttpMethod:       "GET",
		Message:          "Message",
		ResponseMatch:    "<html>",
		AlertBody:        "payload={\"text\": \"Google is down\"}",
		AlertHasBeenSent: false,
		AlertUrl:         os.Getenv("SLACK_URL"),
	}
	err := sendStatus(&system, body)
	if err != nil {
		t.Fatalf("Happy path test, err should be nil, %v", err)
	}
}

func TestSendStatusEmptyBody(t *testing.T) {
	body := ""
	system := models.Systemstatus{
		ID:               1,
		Name:             "google",
		CallUrl:          "https://google.com",
		CallStatus:       "OK",
		HttpMethod:       "GET",
		Message:          "Message",
		ResponseMatch:    "<html>",
		AlertBody:        "payload={\"text\": \"Google is down\"}",
		AlertHasBeenSent: false,
		AlertUrl:         os.Getenv("SLACK_URL"),
	}
	err := sendStatus(&system, body)
	if err == nil {
		t.Fatal("test failed, expected error when sending an empty body but got nil")
	}
}

func TestSendStatusEmptyModel(t *testing.T) {
	body := "message body"
	system := models.Systemstatus{}
	err := sendStatus(&system, body)
	if err == nil {
		t.Fatal("expected err, got nil")
	}
}
