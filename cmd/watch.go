package cmd

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/raojinlin/apollo-client/apollo"
	"log"
	"net/http"
	"os"
	"os/exec"
	path2 "path"
	"strings"
	"sync"
)

type Notify struct {
	Script string
	Url    string
}

var scriptPath = "/tmp"

func generateScriptName(notify string) string {
	s := md5.New().Sum(bytes.NewBufferString(notify).Bytes())
	return path2.Join(scriptPath, "apollo-notify-"+hex.EncodeToString(s)+".sh")
}

func createNotifyScript(content, script string) error {
	file, err := os.OpenFile(script, os.O_CREATE|os.O_RDWR, os.ModePerm)
	defer file.Close()

	if errors.Is(err, os.ErrExist) {
		return nil
	}

	if err != nil {
		return err
	}

	_, err = file.WriteString("#!/usr/bin/env bash\n\n")
	if err != nil {
		return err
	}
	_, err = file.WriteString(content)
	return err
}

func execNotifyScript(script string, response []apollo.NotificationResponse, config []*apollo.Response) error {
	if script == "" {
		return nil
	}

	scriptName := generateScriptName(script)
	err := createNotifyScript(script, scriptName)
	if err != nil {
		return err
	}

	log.Println("exec notify script: " + scriptName)
	cmd := exec.Command("bash", scriptName)
	names := make([]string, len(response))
	for i, r := range response {
		names[i] = r.NamespaceName
	}

	cmd.Env = []string{
		"PATH=" + os.Getenv("PATH"),
		"NAMESPACES=" + strings.Join(names, ","),
		"CACHEDIR=" + cacheDir,
		"APPID=" + appId,
		"CLUSTER=" + cluster,
		"SERVER=" + server,
	}
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(output))
	}

	return err
}

func pushToServer(url string, response []apollo.NotificationResponse, config []*apollo.Response) error {
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

func notified(notify *Notify, response []apollo.NotificationResponse, config []*apollo.Response) error {
	var err error
	wg := sync.WaitGroup{}

	if notify.Script != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = execNotifyScript(notify.Script, response, config)
		}()
	}

	if notify.Url != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = pushToServer(notify.Url, response, config)
			if err != nil {
				log.Println("Push to ", notify.Url, "failed", err.Error())
			}
		}()
	}

	wg.Wait()
	return err
}

func watch(server, appId, cluster, cacheDir string, namespaces []string, notify *Notify) error {
	var ns = make([]*apollo.NotificationRequestPayload, len(namespaces))
	for i, namespace := range namespaces {
		ns[i] = &apollo.NotificationRequestPayload{
			Cluster:        cluster,
			NamespaceName:  namespace,
			NotificationId: -1,
		}
	}

	return apollo.Subscribe(server, appId, cluster, cacheDir, ns, func(err error, response []apollo.NotificationResponse, config []*apollo.Response) {
		err = notified(notify, response, config)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			log.Println("Notified.")
		}
	})
}
