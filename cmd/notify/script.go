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

type ScriptNotification struct {
	Script string
	Path   string
}

func (s *ScriptNotification) generateScriptName() string {
	content := md5.New().Sum(bytes.NewBufferString(s.Script).Bytes())
	return path2.Join(s.Path, "apollo-notify-"+hex.EncodeToString(content)+".sh")
}

func (s *ScriptNotification) createNotifyScript(scriptName string) error {
	file, err := os.OpenFile(scriptName, os.O_CREATE|os.O_RDWR, os.ModePerm)

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
	_, err = file.WriteString(s.Script)
	return err
}

func (s *ScriptNotification) Notify(opt *apollo.Option, response []apollo.NotificationResponse, config []*apollo.Response) error {
	script := s.Script
	if script == "" {
		return nil
	}

	scriptName := s.generateScriptName()
	err := s.createNotifyScript(scriptName)
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
		"CACHEDIR=" + opt.CacheDir,
		"APPID=" + opt.AppId,
		"CLUSTER=" + opt.Cluster,
		"SERVER=" + opt.Server,
	}
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(output))
	}

	return err
}

func NewScriptNotification(script string) *ScriptNotification {
	return &ScriptNotification{
		Script: script,
		Path:   "/tmp/",
	}
}
