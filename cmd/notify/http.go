package notify

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/raojinlin/apollo-client/apollo"
	"log"
	"net/http"
)

func pushToServer(url string, config []*apollo.Response) error {
	j, err := json.Marshal(config)
	if err != nil {
		return err
	}

	log.Println("Push change to server: ", url)
	body := bytes.NewReader(j)
	res, err := http.Post(url, "application/json", body)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New("Invalid response status code" + res.Status)
	}

	return nil
}
