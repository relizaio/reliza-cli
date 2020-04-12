package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/go-resty/resty"
)

func parseCopyTemplate(directory string, relizaHubUri string, environment string, tagKey string, tagVal string, apiKeyId string, apiKey string) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, f := range files {
		fmt.Println(f.Name())
		fullFile, fileOpenErr := os.Open(directory + "/" + f.Name())
		if fileOpenErr != nil {
			fmt.Println(fileOpenErr)
			os.Exit(1)
		}
		s := bufio.NewScanner(fullFile)
		for s.Scan() {
			line := s.Text()
			if strings.Contains(line, "<%PROJECT__") {
				// extract project id to request latest applicable release from reliza hub
				//match, _ := regexp.MatchString("<%PROJECT__(.+)%>", line)
				r, _ := regexp.Compile("<%PROJECT__([-_a-zA-Z0-9]+)%>")
				projectIdWBranch := r.FindStringSubmatch(line)[1]
				projectIdWBranchArr := strings.Split(projectIdWBranch, "__")
				projectId := projectIdWBranchArr[0]
				branch := projectIdWBranchArr[1]
				fmt.Println(projectId)
				fmt.Println(branch)

				// call Reliza Hub with specified project id
				body := getLatestReleaseFunc("false", relizaHubUri, projectId, branch, environment, tagKey, tagVal, apiKeyId, apiKey)
				fmt.Println("2nd print")
				fmt.Println(string(body()))
			}
			// fmt.Println(s.Text())
		}
	}
}

func getLatestReleaseFunc(debug string, relizaHubUri string, project string, branch string, environment string,
	tagKey string, tagVal string, apiKeyId string, apiKey string) func() []byte {
	if debug == "true" {
		fmt.Println("Using Reliza Hub at", relizaHubUri)
	}

	path := relizaHubUri + "/api/programmatic/v1/release/getLatestProjectRelease/" + project + "/" + branch
	if len(environment) > 0 {
		path = path + "/" + environment
	}

	if len(tagKey) > 0 && len(tagVal) > 0 {
		path = path + "?tag=" + tagKey + "____" + tagVal
	}

	if debug == "true" {
		fmt.Println(path)
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "Reliza Go Client").
		SetHeader("Accept-Encoding", "gzip, deflate").
		SetBasicAuth(apiKeyId, apiKey).
		Get(path)

	printResponse(err, resp)
	return resp.Body
}
