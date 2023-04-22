package notify

import (
	"github.com/raojinlin/apollo-client/apollo"
	"log"
	"sync"
)

type Notify struct {
	Script string
	Url    string
}

func Push(opt *apollo.Option, notify *Notify, response []apollo.NotificationResponse, config []*apollo.Response) error {
	var err error
	wg := sync.WaitGroup{}

	if notify.Script != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = NewScriptNotification(notify.Script).Notify(opt, response, config)
		}()
	}

	if notify.Url != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = NewHttpNotification(notify.Url).Notify(opt, response, config)
			if err != nil {
				log.Println("Push to ", notify.Url, "failed", err.Error())
			}
		}()
	}

	wg.Wait()
	return err
}
