package apollo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	path2 "path"
	"regexp"
	"strings"
	"sync"
)

func pullConfig(option Option) (error, *Response) {
	if option.Cluster == "" {
		option.Cluster = "default"
	}

	res, err := http.Get(fmt.Sprintf("%s/configs/%s/%s/%s", option.Server, option.AppId, option.Cluster, option.Namespaces[0]))
	if err != nil {
		return err, nil
	}

	defer func() {
		res.Body.Close()
	}()

	resp, err := io.ReadAll(res.Body)
	if err != nil {
		return err, nil
	}

	var response Response
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return err, nil
	}

	return nil, &response
}

func PullConfigBatch(option Option) ([]*Response, error) {
	var result = make([]*Response, len(option.Namespaces))
	var wg sync.WaitGroup
	for i, namespace := range option.Namespaces {
		wg.Add(1)
		go func(j int, ns string) {
			defer func() {
				wg.Done()
			}()
			pn := option
			pn.Namespaces = []string{ns}
			err, r := pullConfig(pn)
			if err != nil {
				result[j] = nil
			} else {
				result[j] = r
			}
		}(i, namespace)
	}

	wg.Wait()
	return result, nil
}

func getOutputFile(response *Response) string {
	return strings.Join([]string{response.AppId, response.Cluster, response.NamespaceName}, "-") + ".env"
}

func save(path string, response *Response) error {
	path = path2.Join(path, getOutputFile(response))
	info, _ := os.Stat(path)

	if info != nil && info.IsDir() {
		return fmt.Errorf("invalid file type, %s not file", path)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}

	err = file.Truncate(0)
	if err != nil {
		return err
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	defer func() {
		file.Close()
	}()

	dotReg, err := regexp.Compile("\\.")
	if err != nil {
		return err
	}

	reg, err := regexp.Compile("\"")

	for key, val := range response.Configurations {
		key = strings.ToUpper(dotReg.ReplaceAllString(key, "_"))
		val = reg.ReplaceAllStringFunc(val, func(s string) string {
			return "\\\""
		})
		_, err = io.WriteString(file, key+"=\""+val+"\"\n")
	}

	if err != nil {
		return err
	}

	//_, err = file.Write(r)
	return err
}

func PullConfigAndSave(option Option) (r []*Response, err error) {
	r, err = PullConfigBatch(option)
	if err != nil {
		return
	}

	for _, item := range r {
		err = save(option.CacheDir, item)
		if err != nil {
			return
		}
	}

	return
}
