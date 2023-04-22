package apollo

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	url2 "net/url"
	"os"
	"strconv"
)

func readNotificationId(opt Option, namespace string) int {
	path := fmt.Sprintf("/tmp/%s-%s-%s.nid", opt.AppId, opt.Cluster, namespace)
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return -1
	}

	s, err := io.ReadAll(file)
	id, err := strconv.Atoi(string(s))
	if err != nil {
		return -1
	}

	return id
}

func saveNotificationId(opt Option, namespace string, id int) error {
	path := fmt.Sprintf("/tmp/%s-%s-%s.nid", opt.AppId, opt.Cluster, namespace)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}

	_, err = io.WriteString(file, strconv.Itoa(id))
	return err
}

func Subscribe(opt Option, subject []*NotificationRequestPayload, handler func(error, []NotificationResponse, []*Response)) error {
	var stop bool
	for !stop {
		url := fmt.Sprintf("%s/notifications/v2?appId=%s&cluster=%s&", opt.Server, opt.AppId, opt.Cluster)
		for _, payload := range subject {
			payload.NotificationId = readNotificationId(opt, payload.NamespaceName)
		}

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
						err = saveNotificationId(opt, item.NamespaceName, item2.NotificationId)
						if err != nil {
							log.Printf("Error, save %s notification id %d, %s", item.NamespaceName, item.NotificationId, err.Error())
						}
						break
					}
				}
			}
			log.Println("Configuration changed. update local config file.")
			newConfig, err := PullConfigAndSave(opt)
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
