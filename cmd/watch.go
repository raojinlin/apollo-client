package cmd

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/raojinlin/apollo-client/apollo"
	"log"
	"os"
	"os/exec"
	path2 "path"
	"strings"
)

var scriptPath = "/tmp"

func generateScriptName(notify string) string {
	s := md5.New().Sum(bytes.NewBufferString(notify).Bytes())
	return path2.Join(scriptPath, "apollo-notify-"+hex.EncodeToString(s)+".sh")
}

func createNotifyScript(content, script string) error {
	file, err := os.OpenFile(script, os.O_CREATE|os.O_RDWR, os.ModePerm)
	defer file.Close()
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

func watch(server, appId, cluster, notify, cacheDir string, namespaces []string) error {
	var ns = make([]*apollo.NotificationRequestPayload, len(namespaces))
	for i, namespace := range namespaces {
		ns[i] = &apollo.NotificationRequestPayload{
			Cluster:        cluster,
			NamespaceName:  namespace,
			NotificationId: -1,
		}
	}

	var script string
	if notify != "" {
		script := generateScriptName(notify)
		err := createNotifyScript(notify, script)
		if err != nil {
			return err
		}
	}

	return apollo.Subscribe(server, appId, cluster, cacheDir, ns, func(err error, response []apollo.NotificationResponse) {
		log.Println("updated", response, err)
		if script != "" {
			cmd := exec.Command("bash", script)
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
		}
	})
}
