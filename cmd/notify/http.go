package notify

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/raojinlin/apollo-client/apollo"
	"log"
	"net/http"
)

type HTTPNotification struct {
	NotifyUrl string
}

func (h *HTTPNotification) Notify(opt *apollo.Option, response []apollo.NotificationResponse, config []*apollo.Response) error {
	j, err := json.Marshal(config)
	if err != nil {
		return err
	}

	log.Println("Push change to server: ", h.NotifyUrl)
	body := bytes.NewReader(j)
	res, err := http.Post(h.NotifyUrl, "application/json", body)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New("Invalid response status code" + res.Status)
	}

	return nil
}

func NewHttpNotification(url string) *HTTPNotification {
	return &HTTPNotification{NotifyUrl: url}
}
