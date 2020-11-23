package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"gopkg.in/resty.v1"
)

func parseCopyTemplate(directory string, outDirectory string, relizaHubUri string, environment string, tagKey string,
	tagVal string, apiKeyId string, apiKey string, instance string, namespace string) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Println("Error opening parse directory = " + directory)
		fmt.Println(err)
		os.Exit(1)
	}
	for _, f := range files {
		// fmt.Println(f.Name())
		// open read file
		fullFile, fileOpenErr := os.Open(directory + "/" + f.Name())
		if fileOpenErr != nil {
			fmt.Println("Error opening source file for parse = " + directory + "/" + f.Name())
			fmt.Println(fileOpenErr)
			os.Exit(1)
		}
		// open write file
		writeFile, writeFileCreateErr := os.Create(outDirectory + "/" + f.Name())
		if writeFileCreateErr != nil {
			fmt.Println("Error creating parse output file = " + outDirectory + "/" + f.Name())
			fmt.Println(writeFileCreateErr)
			os.Exit(1)
		}

		s := bufio.NewScanner(fullFile)
		for s.Scan() {
			line := s.Text()
			if strings.Contains(line, "<%PROJECT__") {
				// reset all params to avoid incorrect results from old templates
				branch = ""
				project = ""
				product = ""
				// extract project id to request latest applicable release from reliza hub
				//match, _ := regexp.MatchString("<%PROJECT__(.+)%>", line)
				r, _ := regexp.Compile("<%PROJECT__([-_a-zA-Z0-9\\s]+)%>")
				subLineToReplace := r.FindStringSubmatch(line)[0]
				projectProductTemplate := r.FindStringSubmatch(line)[1]
				projectProductArray := strings.Split(projectProductTemplate, "PRODUCT__")
				productId := ""
				if len(projectProductArray) == 2 {
					// product is present
					productIdWBranchArr := strings.Split(projectProductArray[1], "__")
					productId = productIdWBranchArr[0]
					if len(productIdWBranchArr) == 2 {
						branch = productIdWBranchArr[1]
					}
				}

				projectIdWBranchArr := strings.Split(projectProductArray[0], "__")
				projectId := projectIdWBranchArr[0]
				if len(branch) < 1 && len(productId) < 1 && len(projectIdWBranchArr) == 2 {
					branch = projectIdWBranchArr[1]
				}
				//fmt.Println(subLineToReplace)
				//fmt.Println(projectId)
				//fmt.Println(branch)
				// call Reliza Hub with specified project id
				body := getLatestReleaseFunc("false", relizaHubUri, projectId, productId, branch, environment,
					tagKey, tagVal, apiKeyId, apiKey, instance, namespace)

				// parse body json
				var bodyJson map[string]interface{}
				json.Unmarshal(body(), &bodyJson)
				// assume only one artifact - should be configured by tags - later add type selector - TODO
				// for now only use first digest - TODO
				var pickedArtifact string
				artifactsArr := bodyJson["artifactDetails"].([]interface{})
				zeroArtifact := artifactsArr[0].(map[string]interface{})
				artifactsIdentifier := zeroArtifact["identifier"].(string)
				artifactDigests := zeroArtifact["digests"].([]interface{})
				pickedArtifact = artifactsIdentifier + "@" + artifactDigests[0].(string)
				//fmt.Println(pickedArtifact)

				// perform string replacement
				line = strings.Replace(line, subLineToReplace, pickedArtifact, -1)

				// fmt.Println("2nd print")
				// fmt.Println(string(body()))
			}
			// output line to out directory
			writeFile.WriteString(line + "\n")

			// fmt.Println(s.Text())
		}
		writeFile.Close()
	}
}

func substituteCopyBasedOnMap(inFile string, outFile string, substitutionMap map[string]string) {
	fmt.Println("Opening in file...")
	// open read file
	inFileOpened, fileOpenErr := os.Open(inFile)
	if fileOpenErr != nil {
		fmt.Println("Error opening inFile = " + inFile)
		fmt.Println(fileOpenErr)
		os.Exit(1)
	}
	fmt.Println("Opening output file...")
	// open write file
	outFileOpened, fileOpenErr := os.Create(outFile)
	if fileOpenErr != nil {
		fmt.Println("Error opening outFile = " + outFile)
		fmt.Println(fileOpenErr)
		os.Exit(1)
	}
	inScanner := bufio.NewScanner(inFileOpened)
	for inScanner.Scan() {
		line := inScanner.Text()
		// check if line contains any key of substitution map
		for k, v := range substitutionMap {
			if strings.Contains(line, k) {
				line = strings.ReplaceAll(line, k, v)
				break
			}
		}
		outFileOpened.WriteString(line + "\n")
	}
}

func scanTagFile(tagSourceFile string, typeVal string) map[string]string {
	tagFile, fileOpenErr := os.Open(tagSourceFile)
	if fileOpenErr != nil {
		fmt.Println("Error opening tagSourceFile = " + tagSourceFile)
		fmt.Println(fileOpenErr)
		os.Exit(1)
	}

	tagSourceMap := map[string]string{}

	if typeVal == "cyclonedx" {
		cycloneBytes, ioReadErr := ioutil.ReadAll(tagFile)
		if ioReadErr != nil {
			fmt.Println("Error opening tagFile = " + tagSourceFile)
			fmt.Println(ioReadErr)
			os.Exit(1)
		}
		var bomJson map[string]interface{}
		json.Unmarshal(cycloneBytes, &bomJson)
		// go to components in cyclone dx
		bomComponents := bomJson["components"].([]interface{})
		for _, bomc := range bomComponents {
			// check that type is container
			if bomc.(map[string]interface{})["type"] == "container" {

				resolvedImage := false

				// name is a mandatory field in cyclonedx spec
				contName := bomc.(map[string]interface{})["name"].(string)
				// 1st try to parse purl if present

				if bomc.(map[string]interface{})["purl"] != nil {
					purl := bomc.(map[string]interface{})["purl"].(string)
					// sample purl pkg:docker/test-cont@sha256:testsha256hash?repository_url=123.dkr.ecr.us-east-1.amazonaws.com
					// remove pkg:docker/ thing first - must be there
					purl = strings.ReplaceAll(purl, "pkg:docker/", "")
					// check if the image is not on docker hub - split by ?repository_url= if present, otherewise repository is docker hub (ignore)
					purlImageName := purl
					if strings.Contains(purl, "?repository_url=") {
						purlImageName = strings.Split(purl, "?repository_url=")[1] + "/" + strings.Split(purl, "?repository_url=")[0]
					}
					parseImageNameIntoMap(purlImageName, tagSourceMap)
					resolvedImage = true
				}
				if !resolvedImage && bomc.(map[string]interface{})["hashes"] != nil {
					// if purl is not set - use name and hash if present, but only if hashes contain SHA-256 algorithm
					for _, hashEntry := range bomc.(map[string]interface{})["hashes"].([]interface{}) {
						alg := hashEntry.(map[string]interface{})["alg"].(string)
						if 0 == strings.Compare(alg, "SHA-256") {
							// take name and attach hash
							fullImageName := stripImageHashTag(contName) + "@sha256:" + hashEntry.(map[string]interface{})["content"].(string)
							parseImageNameIntoMap(fullImageName, tagSourceMap)
							resolvedImage = true
							break
						}
					}
				}
				if !resolvedImage {
					// if both purl and hashes are not set - use only name and treat it same as text file case
					parseImageNameIntoMap(contName, tagSourceMap)
				}
			}
		}
	} else if typeVal == "text" {
		tagScanner := bufio.NewScanner(tagFile)
		for tagScanner.Scan() {
			line := tagScanner.Text()
			parseImageNameIntoMap(line, tagSourceMap)
		}
	}
	return tagSourceMap
}

/**
* This adds value into tag source map based on image name
 */
func parseImageNameIntoMap(imageName string, tagSourceMap map[string]string) {
	strippedImageName := stripImageHashTag(imageName)
	tagSourceMap[strippedImageName] = imageName
}

func stripImageHashTag(imageName string) string {
	strippedImageName := imageName
	if strings.Contains(imageName, "@") {
		sourceTagSplit := strings.Split(imageName, "@")
		strippedImageName = sourceTagSplit[0]
	} else if strings.Contains(imageName, ":") {
		sourceTagSplit := strings.SplitN(imageName, ":", 2)
		strippedImageName = sourceTagSplit[0]
	}
	return strippedImageName
}

func getLatestReleaseFunc(debug string, relizaHubUri string, project string, product string, branch string, environment string,
	tagKey string, tagVal string, apiKeyId string, apiKey string, instance string, namespace string) func() []byte {
	if debug == "true" {
		fmt.Println("Using Reliza Hub at", relizaHubUri)
	}

	path := relizaHubUri + "/api/programmatic/v1/release/getLatestProjectRelease"
	body := map[string]string{"project": project}
	if len(environment) > 0 {
		body["environment"] = environment
	}

	if len(product) > 0 {
		body["product"] = product
	}

	if len(tagKey) > 0 && len(tagVal) > 0 {
		body["tags"] = tagKey + "____" + tagVal
	}

	if len(branch) > 0 {
		body["branch"] = branch
	}

	if len(instance) > 0 {
		body["instance"] = instance
	}

	if len(namespace) > 0 {
		body["namespace"] = namespace
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "Reliza Go Client").
		SetHeader("Accept-Encoding", "gzip, deflate").
		SetBody(body).
		SetBasicAuth(apiKeyId, apiKey).
		Post(path)

	printResponse(err, resp)
	return resp.Body

	// path := relizaHubUri + "/api/programmatic/v1/release/getLatestProjectRelease/" + project + "/" + branch
	// if len(environment) > 0 {
	// 	path = path + "/" + environment
	// }

	// if len(tagKey) > 0 && len(tagVal) > 0 {
	// 	path = path + "?tag=" + tagKey + "____" + tagVal
	// }

	// if debug == "true" {
	// 	fmt.Println(path)
	// }

	// client := resty.New()
	// resp, err := client.R().
	// 	SetHeader("Content-Type", "application/json").
	// 	SetHeader("User-Agent", "Reliza Go Client").
	// 	SetHeader("Accept-Encoding", "gzip, deflate").
	// 	SetBasicAuth(apiKeyId, apiKey).
	// 	Get(path)

}
