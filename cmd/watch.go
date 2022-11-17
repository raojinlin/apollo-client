package cmd

import (
	"fmt"
	"github.com/raojinlin/apollo-client/apollo"
	n "github.com/raojinlin/apollo-client/cmd/notify"
	"log"
)

func watch(opt apollo.Option, notify *n.Notify) error {
	var ns = make([]*apollo.NotificationRequestPayload, len(opt.Namespaces))
	for i, namespace := range opt.Namespaces {
		ns[i] = &apollo.NotificationRequestPayload{
			Cluster:        cluster,
			NamespaceName:  namespace,
			NotificationId: -1,
		}
	}

	return apollo.Subscribe(opt, ns, func(err error, response []apollo.NotificationResponse, config []*apollo.Response) {
		err = n.Push(&opt, notify, response, config)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			log.Println("Notified.")
		}
	})
}
