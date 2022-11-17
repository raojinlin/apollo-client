package notify

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
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

func execNotifyScript(option *apollo.Option, script string, response []apollo.NotificationResponse) error {
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
		"CACHEDIR=" + option.CacheDir,
		"APPID=" + option.AppId,
		"CLUSTER=" + option.Cluster,
		"SERVER=" + option.Server,
	}
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(output))
	}

	return err
}
