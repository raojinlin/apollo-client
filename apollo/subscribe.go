package apollo

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	url2 "net/url"
)

func Subscribe(server, appId, cluster, cacheDir string, subject []*NotificationRequestPayload, handler func(error, []NotificationResponse, []*Response)) error {
	var stop bool
	for !stop {
		url := fmt.Sprintf("%s/notifications/v2?appId=%s&cluster=%s&", server, appId, cluster)
		notifications, err := json.Marshal(subject)
		if err != nil {
			return err
		}

		url += "notifications=" + url2.QueryEscape(string(notifications))
		log.Println(url)
		resp, err := http.Get(url)
		if err != nil {
			return err
		}

		if resp.StatusCode == 200 {
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			var notificationResponse []NotificationResponse
			err = json.Unmarshal(data, &notificationResponse)
			if err != nil {
				return err
			}

			var namespaces = make([]string, len(subject))
			for i, item := range subject {
				namespaces[i] = item.NamespaceName
				for _, item2 := range notificationResponse {
					if item2.NamespaceName == item.NamespaceName {
						item.NotificationId = item2.NotificationId
						break
					}
				}
			}
			log.Println("Configuration changed. update local config file.")
			newConfig, err := PullConfigAndSave(cacheDir, server, appId, cluster, namespaces)
			if err != nil {
				log.Println("pull and save config failed: ", err)
			}

			handler(err, notificationResponse, newConfig)
			resp.Body.Close()
		} else if resp.StatusCode == 304 {
			log.Println("No configuration changed, continue check.")
			continue
		} else {
			return fmt.Errorf("got http error, status_code='%d', status='%s'", resp.StatusCode, resp.Status)
		}
	}

	return nil
}
