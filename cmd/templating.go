package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

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

	// Copy data from input-file to output-file, with tags replaced according to substitution map.
	inScanner := bufio.NewScanner(inFileOpened)
	for inScanner.Scan() {
		line := inScanner.Text()
		// check if line contains any key of substitution map
		for k, v := range substitutionMap {
			// have a version stripping docker.io and docker.io/library if it's present
			// for this establish base image text
			baseImageText := ""

			if strings.Contains(line, k+":") || strings.Contains(line, k+"@") || strings.HasSuffix(line, k) {
				baseImageText = k
			}

			if len(baseImageText) < 1 {
				// try stripping docker.io
				contText := strings.Replace(k, "docker.io/", "", 1)
				if strings.Contains(line, contText+":") || strings.Contains(line, contText+"@") || strings.HasSuffix(line, contText) {
					baseImageText = contText
				}
			}

			if len(baseImageText) < 1 {
				// try stripping docker.io/library
				contText := strings.Replace(k, "docker.io/library/", "", 1)
				// note that exact match is too loose here, so instead we only look for image: pattern
				if !strings.Contains(line, "//"+contText) &&
					(strings.Contains(line, contText+":") || strings.Contains(line, contText+"@") || strings.Contains(strings.ToLower(line), "image: "+contText) ||
						strings.Contains(strings.ToLower(line), "image:"+contText)) {
					baseImageText = contText
				}
			}

			if len(baseImageText) > 0 && !strings.HasSuffix(line, ":") {
				// found a match - do substitution
				//split line before image name and concat with substitution map value
				parts := strings.Split(line, baseImageText)

				// remove beginning quotes if present
				startLine := parts[0]
				re := regexp.MustCompile("\"$")
				startLine = re.ReplaceAllLiteralString(startLine, "")
				re = regexp.MustCompile("'$")
				startLine = re.ReplaceAllLiteralString(startLine, "")

				line = startLine + v
				break
			}
		}
		outFileOpened.WriteString(line + "\n")
	}
}

/*
	This function adds some extra info in the form of comments to the top of the outfile
	specified (adds two lines to exisiting content, does not overwrite). The first line
	notes the version of reliza-cli that ran the command to generate the outfile, as
	well as the date the file was generated. The second line contains info about where
	the replaced tags were sourced from.
*/
func addProvenanceToReplaceTagsOutfile(outfile string, apiKeyId string, apiKey string, tagSourceFile string, environment string, instance string, instanceURI string, revision string, definitionReferenceFile string, typeVal string, version string, bundle string) {
	// Read outfile, store in string, then clear outfile
	outFileOpened, fileOpenErr := os.Open(outfile)
	if fileOpenErr != nil {
		fmt.Println("Error opening outFile = " + outfile)
		fmt.Println(fileOpenErr)
		os.Exit(1)
	}
	// Store generated outfile in temp variable
	// Either need to store lines with '\n' character or else add it when putting lines back into file
	var outFileLines []string
	inScanner := bufio.NewScanner(outFileOpened)
	for inScanner.Scan() {
		line := inScanner.Text()
		outFileLines = append(outFileLines, line)
	}
	// Close file and check for errors
	fileCloseErr := outFileOpened.Close()
	if fileCloseErr != nil {
		//fmt.Fprintf(os.Stderr, "error closing outfile: %v\n", fileCloseErr)
		fmt.Println("Error closing outFile = " + outfile)
		fmt.Println(fileCloseErr)
		os.Exit(1)
	}
	// Reopen file, but this type creating a new file to replace the old outfile
	outFileCreated, fileCreateErr := os.Create(outfile)
	if fileCreateErr != nil {
		fmt.Println("Error opening outFile = " + outfile)
		fmt.Println(fileCreateErr)
		os.Exit(1)
	}

	// Add some provenance to top of file (as comments)
	var provenanceLine1 string
	var provenanceLine2 string

	// First line: current reliza-cli version and current datetime
	currentDateTimeFormatted := time.Now().UTC().Format(time.RFC3339)
	relizaCliCurrentVersion := Version // This is not really getting updated atm??
	provenanceLine1 = "# Tags replaced with Reliza CLI version " + relizaCliCurrentVersion + " on " + currentDateTimeFormatted

	// Second line: where tags come from, either:
	// (tagsource file) or (environment) or (bundle+version)
	// or (instance+revision) or (instanceuri+revision) , revision is optional, otherwise uses latest revision
	// or (apiKeyId suffix, if using apiKeyId+apiKey pair from instance)
	if tagSourceFile != "" {
		provenanceLine2 = "# According to tag source file " + tagSourceFile
	} else if len(environment) > 0 {
		//provenanceLine2 = "# According to the latest approved images in environment " + environment + " for instance " + "XXX"
		provenanceLine2 = "# According to the latest approved images in  " + environment + " environment."
	} else if len(bundle) > 0 && len(version) > 0 {
		provenanceLine2 = "# According to bundle " + bundle + " version " + version
	} else if len(instance) > 0 {
		if len(revision) > 0 {
			provenanceLine2 = "# According to revision " + revision + " of the instance " + instance
		} else {
			// no revision specified, using latest
			provenanceLine2 = "# According to latest approved images for the instance " + instance
		}
	} else if len(instanceURI) > 0 {
		if len(revision) > 0 {
			provenanceLine2 = "# According to revision " + revision + " of the instance at " + instanceURI
		} else {
			// no revision specified, using latest
			provenanceLine2 = "# According to latest approved images for the instance at " + instanceURI
		}
	} else if strings.HasPrefix(apiKeyId, "INSTANCE__") { // Unessecary
		instUUIDFromAPIKeyId := apiKeyId[10:37] // remove first 10 chars
		if len(revision) > 0 {
			provenanceLine2 = "# According to revision " + revision + " of the instance " + instUUIDFromAPIKeyId
		} else {
			// no revision specified, using latest
			provenanceLine2 = "# According to latest approved images for the instance " + instUUIDFromAPIKeyId
		}
	} else {
		// should have at least one of those things
		provenanceLine2 = "missing replacetags input"
	}

	outFileCreated.WriteString(provenanceLine1 + "\n")
	outFileCreated.WriteString(provenanceLine2 + "\n")

	// Add back the rest of the outfile
	for _, line := range outFileLines {
		outFileCreated.WriteString(line + "\n")
	}

	// Close file and check for errors
	fileCreatedCloseErr := outFileCreated.Close()
	if fileCreatedCloseErr != nil {
		//fmt.Fprintf(os.Stderr, "error closing outfile: %v\n", fileCloseErr)
		fmt.Println("Error closing outfile = " + outfile)
		fmt.Println(fileCloseErr)
		os.Exit(1)
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
		var bomJSON map[string]interface{}
		json.Unmarshal(cycloneBytes, &bomJSON)
		extractComponentsFromCycloneJSON(bomJSON, tagSourceMap)

	} else if typeVal == "text" {
		tagScanner := bufio.NewScanner(tagFile)
		for tagScanner.Scan() {
			line := tagScanner.Text()
			parseImageNameIntoMap(line, tagSourceMap)
		}
	}
	return tagSourceMap
}

func extractComponentsFromCycloneJSON(bomJSON map[string]interface{}, tagSourceMap map[string]string) map[string]string {
	// go to components in cyclone dx
	var bomComponents []interface{}
	if components, ok := bomJSON["components"]; ok {
		bomComponents = components.([]interface{})
	} else {
		fmt.Println("Error: CycloneDX BOM componenets are empty!")
		os.Exit(1)
	}

	// bomComponents := bomJSON["components"].([]interface{})
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

func scanTags(tagSourceFile string, typeVal string, apiKeyId string, apiKey string, instance string, revision string, instanceURI string, bundle string, version string, environment string) map[string]string {
	tagSourceMap := map[string]string{} // 1st - scan tag source file and construct a map of generic tag to actual tag
	if tagSourceFile != "" {
		tagSourceMap = scanTagFile(tagSourceFile, typeVal)
	} else if len(environment) > 0 {
		cycloneBytes := getEnvironmentCycloneDxExportV1(apiKeyId, apiKey, environment)
		// fmt.Println("res", tagSourceRes)
		var bomJSON map[string]interface{}
		json.Unmarshal(cycloneBytes.Body(), &bomJSON)
		fmt.Println(string(cycloneBytes.Body()))
		extractComponentsFromCycloneJSON(bomJSON, tagSourceMap)
	} else if len(bundle) > 0 && len(version) > 0 {
		cycloneBytes := getBundleVersionCycloneDxExportV1(apiKeyId, apiKey, bundle, version)
		// fmt.Println("res", tagSourceRes)
		var bomJSON map[string]interface{}
		json.Unmarshal(cycloneBytes.Body(), &bomJSON)
		extractComponentsFromCycloneJSON(bomJSON, tagSourceMap)
	} else if len(instance) > 0 || len(instanceURI) > 0 || strings.HasPrefix(apiKeyId, "INSTANCE__") {
		// tagSourceMap = getInstanceRevisionCycloneDxExportV1(apiKeyId, apiKey, instance, revision)
		cycloneBytes := getInstanceRevisionCycloneDxExportV1(apiKeyId, apiKey, instance, revision, instanceURI)
		// fmt.Println("res", tagSourceRes)
		var bomJSON map[string]interface{}
		json.Unmarshal(cycloneBytes.Body(), &bomJSON)
		extractComponentsFromCycloneJSON(bomJSON, tagSourceMap)
	} else {
		fmt.Println("Scan Tags Failed! specify either tagsource or instance or bundle and version")
		os.Exit(1)
	}
	return tagSourceMap
}

func getInstanceRevisionCycloneDxExportV1(apiKeyId string, apiKey string, instance string, revision string, instanceURI string) *resty.Response {

	if len(instance) <= 0 && len(instanceURI) <= 0 && !strings.HasPrefix(apiKeyId, "INSTANCE__") {
		//throw error and exit
		fmt.Println("instance or instanceURI not specified!")
		os.Exit(1)
	}

	if "" == revision {
		revision = "-1"
	}

	path := relizaHubUri + "/api/programmatic/v1/instanceRevision/cyclonedxExport?revision=" + revision

	if len(instance) > 0 {
		path += "&uuid=" + instance
	}

	if len(instanceURI) > 0 {
		path += "&uri=" + instanceURI
	}
	client := resty.New()
	resp, err := client.R().
		SetHeader("User-Agent", "Reliza Go Client").
		SetHeader("Accept-Encoding", "gzip, deflate").
		SetBasicAuth(apiKeyId, apiKey).
		Get(path)

	if resp.StatusCode() != 200 {
		fmt.Println("Error Response Info:")
		fmt.Println("Error      :", err)
		fmt.Println("Status Code:", resp.StatusCode())
		fmt.Println("Status     :", resp.Status())
		fmt.Println("Time       :", resp.Time())
		fmt.Println("Received At:", resp.ReceivedAt())
		os.Exit(1)
	}
	return resp
}

func getBundleVersionCycloneDxExportV1(apiKeyId string, apiKey string, bundle string, version string) *resty.Response {

	if len(bundle) <= 0 && len(version) <= 0 {
		//throw error and exit
		fmt.Println("instance or instanceURI not specified!")
		os.Exit(1)
	}

	// if "" == revision {
	// 	revision = "-1"
	// }

	path := relizaHubUri + "/api/programmatic/v1/productRelease/exportAsBom?bundle_name=" + url.QueryEscape(bundle) + "&bundle_version=" + url.QueryEscape(version)

	client := resty.New()
	resp, err := client.R().
		SetHeader("User-Agent", "Reliza Go Client").
		SetHeader("Accept-Encoding", "gzip, deflate").
		SetBasicAuth(apiKeyId, apiKey).
		Get(path)

	if resp.StatusCode() != 200 {
		fmt.Println("Error Response Info:")
		fmt.Println("Error      :", err)
		fmt.Println("Status Code:", resp.StatusCode())
		fmt.Println("Status     :", resp.Status())
		fmt.Println("Time       :", resp.Time())
		fmt.Println("Received At:", resp.ReceivedAt())
		os.Exit(1)
	}
	return resp
}

func getEnvironmentCycloneDxExportV1(apiKeyId string, apiKey string, environment string) *resty.Response {

	if len(environment) <= 0 {
		//throw error and exit
		fmt.Println("environment not specified!")
		os.Exit(1)
	}

	path := relizaHubUri + "/api/programmatic/v1/environmentRelease/exportAsBom?environment=" + url.QueryEscape(environment)
	client := resty.New()
	resp, err := client.R().
		SetHeader("User-Agent", "Reliza Go Client").
		SetHeader("Accept-Encoding", "gzip, deflate").
		SetBasicAuth(apiKeyId, apiKey).
		Get(path)

	if resp.StatusCode() != 200 {
		fmt.Println("Error Response Info:")
		fmt.Println("Error      :", err)
		fmt.Println("Status Code:", resp.StatusCode())
		fmt.Println("Status     :", resp.Status())
		fmt.Println("Time       :", resp.Time())
		fmt.Println("Received At:", resp.ReceivedAt())
		os.Exit(1)
	}
	return resp
}
