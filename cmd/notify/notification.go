package notify

import (
	"github.com/raojinlin/apollo-client/apollo"
)

type Notification interface {
	Notify(opt *apollo.Option, response []apollo.NotificationResponse, config []*apollo.Response) error
}

func NewNotification(notify *Notify) Notification {
	if notify.Script != "" {
		return NewScriptNotification(notify.Script)
	}

	return NewHttpNotification(notify.Url)
}