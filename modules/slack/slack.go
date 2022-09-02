package slack

import (
	"bytes"
	"encoding/json"
	"net/http"
	// "fmt"
)

type SlackPayload struct {
	Text     string `json:"text"`
	Username string `json:"username"`
}

func PostMessage(url string, payload SlackPayload) error {
	// body, err := json.Marshal(map[string]string{"text": message, "username": style})
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	err = response.Body.Close()
	if err != nil {
		return err
	}

	return nil
}
