/*
The MIT License (MIT)

Copyright (c) 2020 Reliza Incorporated (Reliza (tm), https://reliza.io)

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"),
to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

*/
package cmd

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/machinebox/graphql"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/resty.v1"

	"github.com/spf13/viper"
)

var action string
var apiKeyId string
var apiKey string
var artBuildId []string
var artBuildUri []string
var artCiMeta []string
var artDigests []string
var artId []string
var artType []string
var artVersion []string
var artPublisher []string
var artGroup []string
var artPackage []string
var approvalType string
var branch string
var bundle string
var cfgFile string
var commit string
var commitMessage string
var commits string // base64-encoded list of commits obtained with: git log $LATEST_COMMIT..$CURRENT_COMMIT --date=iso-strict --pretty='%H|||%ad|||%s' | base64 -w 0
var dateActual string
var dateStart []string
var dateEnd []string
var debug string
var disapprove bool // approve (default) or disapprove
var endpoint string
var environment string
var featureBranchVersioning string
var hash string
var imageFilePath string
var imageString string
var imageStyle string
var includeApi bool
var instance string
var instanceURI string
var manual bool
var metadata string
var modifier string
var namespace string
var onlyVersion bool
var outDirectory string
var parseDirectory string
var inDirectory string
var infile string
var outfile string
var tagSourceFile string
var definitionReferenceFile string
var provenance bool // add provenance (default), or do not add provenance
var releaseId string
var releaseVersion string
var relizaHubUri string
var revision string
var product string
var project string
var projectName string
var projectType string
var senderId string
var status string
var tagKey string
var tagKeyArr []string
var tagVal string
var tagValArr []string
var typeVal string
var version string
var versionSchema string
var vcsName string
var vcsTag string
var vcsType string
var vcsUri string
var vcsUuid string

const (
	defaultConfigFilename = ".reliza"
	envPrefix             = ""
	configType            = "env"
)

type ErrorBody struct {
	Timestamp string
	Status    int
	Error     string
	Message   string
	Path      string
}

const RELEASE_GQL_DATA = `
	uuid
	createdType
	lastUpdatedBy
	createdDate
	version
	status
	org
	project
	branch
	parentReleases {
		timeSent
		release
		artifact
		type
		namespace
		properties
		state
		replicas {
			id
			state
		}
	}
	optionalReleases {
		timeSent
		release
		artifact
		type
		namespace
		properties
		state
		replicas {
			id
			state
		}
	}
	sourceCodeEntry
	artifacts
	type
	notes
	approvals
	timing {
		lifecycle
		dateFrom
		dateTo
		environment
		instanceUuid
		event
		duration
	}
	endpoint
	commits
`

const FULL_RELEASE_GQL_DATA = RELEASE_GQL_DATA + `
	sourceCodeEntryDetails {
		uuid
		branchUuid
		vcsUuid
		vcsBranch
		commit
		commits
		commitMessage
		vcsTag
		notes
		org
		dateActual
	}
	vcsRepository {
		uuid
		name
		org
		uri
		type
	}
	artifactDetails {
		uuid
		identifier
		org
		branch
		buildId
		buildUri
		cicdMeta
		digests
		isInternal
		artifactType {
			name
			aliases
		}
		notes
		tags
		dateFrom
		dateTo
		buildDuration
		packageType
		version
		publisher
		group
		dependencies
	}
	projectName
	namespace
`

const PROJECT_GQL_DATA = `
	uuid
	name
	org
	type
	versionSchema
	vcsRepository
	featureBranchVersioning
	integrations {
		projectIntegrationUuid
		type
		active
		instance
		vcsUuid:
		eventTypes
		parameters
	}
	envBranchMap
	repositoryEnabled
	status
	apiKeyId
	apiKey
`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "reliza-cli",
	Short: "CLI client for programmatic actions on Reliza Hub",
	Long:  `CLI client for programmatic actions on Reliza Hub.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initConfig(cmd)
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// 	},
}

var printversionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints current version of the CLI",
	Long:  `Prints current version of the CLI`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("reliza-cli version " + Version)
	},
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Persisits API Key Id and API Key Secret",
	Long:  "This CLI command takes API Key Id and API Key Secret and writes them to a configuration file in home directory",
	Run: func(cmd *cobra.Command, args []string) {

		home, err := homedir.Dir()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		configPath := filepath.Join(home, defaultConfigFilename+"."+configType)
		if _, err := os.Stat(configPath); err == nil {
			// config file already exists, it will be overwritten
		} else if os.IsNotExist(err) {
			//create new config file
			if _, err := os.Create(configPath); err != nil { // perm 0666
				fmt.Println(err)
				os.Exit(1)
			}
		}

		viper.Set("apikey", apiKey)
		viper.Set("apikeyid", apiKeyId)
		viper.Set("uri", relizaHubUri)

		if err := viper.WriteConfigAs(configPath); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var addreleaseCmd = &cobra.Command{
	Use:   "addrelease",
	Short: "Creates release on Reliza Hub",
	Long: `This CLI command would create new releases on Reliza Hub
			for authenticated project.`,
	Run: func(cmd *cobra.Command, args []string) {
		if debug == "true" {
			fmt.Println("Using Reliza Hub at", relizaHubUri)
		}

		body := map[string]interface{}{"branch": branch, "version": version}
		if len(status) > 0 {
			body["status"] = strings.ToUpper(status)
		}
		if len(endpoint) > 0 {
			body["endpoint"] = endpoint
		}
		if len(project) > 0 {
			body["project"] = project
		}
		if len(artId) > 0 {
			// use artifacts, construct artifact array
			artifacts := make([]map[string]interface{}, len(artId), len(artId))
			for i, aid := range artId {
				artifacts[i] = map[string]interface{}{"identifier": aid}
			}

			// now do some length validations and add elements
			if len(artBuildId) > 0 && len(artBuildId) != len(artId) {
				fmt.Println("number of --artbuildid flags must be either zero or match number of --artid flags")
				os.Exit(2)
			} else if len(artBuildId) > 0 {
				for i, abid := range artBuildId {
					artifacts[i]["buildId"] = abid
				}
			}

			if len(artBuildUri) > 0 && len(artBuildUri) != len(artId) {
				fmt.Println("number of --artbuildUri flags must be either zero or match number of --artid flags")
				os.Exit(2)
			} else if len(artBuildUri) > 0 {
				for i, aburi := range artBuildUri {
					artifacts[i]["buildUri"] = aburi
				}
			}

			if len(artCiMeta) > 0 && len(artCiMeta) != len(artId) {
				fmt.Println("number of --artcimeta flags must be either zero or match number of --artid flags")
				os.Exit(2)
			} else if len(artCiMeta) > 0 {
				for i, acm := range artCiMeta {
					artifacts[i]["cicdMeta"] = acm
				}
			}

			if len(artType) > 0 && len(artType) != len(artId) {
				fmt.Println("number of --arttype flags must be either zero or match number of --artid flags")
				os.Exit(2)
			} else if len(artType) > 0 {
				for i, at := range artType {
					artifacts[i]["type"] = at
				}
			}

			if len(artDigests) > 0 && len(artDigests) != len(artId) {
				fmt.Println("number of --artdigests flags must be either zero or match number of --artid flags")
				os.Exit(2)
			} else if len(artDigests) > 0 {
				for i, ad := range artDigests {
					adSpl := strings.Split(ad, ",")
					artifacts[i]["digests"] = adSpl
				}
			}

			if len(dateStart) > 0 && len(dateStart) != len(artId) {
				fmt.Println("number of --datestart flags must be either zero or match number of --artid flags")
				os.Exit(2)
			} else if len(dateStart) > 0 {
				for i, ds := range dateStart {
					artifacts[i]["dateFrom"] = ds
				}
			}

			if len(dateEnd) > 0 && len(dateEnd) != len(artId) {
				fmt.Println("number of --dateEnd flags must be either zero or match number of --artid flags")
				os.Exit(2)
			} else if len(dateEnd) > 0 {
				for i, de := range dateEnd {
					artifacts[i]["dateTo"] = de
				}
			}

			if len(artVersion) > 0 && len(artVersion) != len(artId) {
				fmt.Println("number of --artversion flags must be either zero or match number of --artid flags")
				os.Exit(2)
			} else if len(artVersion) > 0 {
				for i, av := range artVersion {
					artifacts[i]["artifactVersion"] = av
				}
			}

			if len(artPublisher) > 0 && len(artPublisher) != len(artId) {
				fmt.Println("number of --artpublisher flags must be either zero or match number of --artid flags")
				os.Exit(2)
			} else if len(artPublisher) > 0 {
				for i, ap := range artPublisher {
					artifacts[i]["publisher"] = ap
				}
			}

			if len(artPackage) > 0 && len(artPackage) != len(artId) {
				fmt.Println("number of --artpackage flags must be either zero or match number of --artid flags")
				os.Exit(2)
			} else if len(artPackage) > 0 {
				for i, ap := range artPackage {
					artifacts[i]["packageType"] = strings.ToUpper(ap)
				}
			}

			if len(artGroup) > 0 && len(artGroup) != len(artId) {
				fmt.Println("number of --artgroup flags must be either zero or match number of --artid flags")
				os.Exit(2)
			} else if len(artGroup) > 0 {
				for i, ag := range artGroup {
					artifacts[i]["group"] = ag
				}
			}

			if len(tagKeyArr) > 0 && len(tagKeyArr) != len(artId) {
				fmt.Println("number of --tagkey flags must be either zero or match number of --artid flags")
				os.Exit(2)
			} else if len(tagValArr) > 0 && len(tagValArr) != len(artId) {
				fmt.Println("number of --tagval flags must be either zero or match number of --artid flags")
				os.Exit(2)
			} else if len(tagKeyArr) > 0 && len(tagValArr) < 1 {
				fmt.Println("number of --tagval and --tagkey flags must be the same and must match number of --artid flags")
				os.Exit(2)
			} else if len(tagKeyArr) > 0 {
				for i, key := range tagKeyArr {
					tagKeys := strings.Split(key, ",")
					tagVals := strings.Split(tagValArr[i], ",")
					if len(tagKeys) != len(tagVals) {
						fmt.Println("number of keys and values per each --tagval and --tagkey flag must be the same")
						os.Exit(2)
					}
					k := map[string]string{}
					for j := range tagKeys {
						k[tagKeys[j]] = tagVals[j]
					}
					artifacts[i]["tags"] = k
				}
			}

			body["artifacts"] = artifacts
		}

		if commit != "" {
			commitMap := map[string]string{"uri": vcsUri, "type": vcsType, "commit": commit, "commitMessage": commitMessage}
			if vcsTag != "" {
				commitMap["vcsTag"] = vcsTag
			}
			if dateActual != "" {
				commitMap["dateActual"] = dateActual
			}
			body["sourceCodeEntry"] = commitMap
		}

		if len(commits) > 0 {
			// fmt.Println(commits)
			plainCommits, err := base64.StdEncoding.DecodeString(commits)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			indCommits := strings.Split(string(plainCommits), "\n")
			commitsInBody := make([]map[string]interface{}, len(indCommits)-1, len(indCommits)-1)
			for i := range indCommits {
				if len(indCommits[i]) > 0 {
					singleCommitEl := map[string]interface{}{}
					commitParts := strings.Split(indCommits[i], "|||")
					singleCommitEl["commit"] = commitParts[0]
					singleCommitEl["dateActual"] = commitParts[1]
					singleCommitEl["commitMessage"] = commitParts[2]
					commitsInBody[i] = singleCommitEl

					// if commit is not present but we are here, use first line as commit
					if len(commit) < 1 && i == 0 {
						commitMap := map[string]string{"commit": commitParts[0], "dateActual": commitParts[1], "commitMessage": commitParts[2]}
						if vcsTag != "" {
							commitMap["vcsTag"] = vcsTag
						}
						if vcsUri != "" {
							commitMap["uri"] = vcsUri
						}
						if vcsType != "" {
							commitMap["type"] = vcsType
						}
						body["sourceCodeEntry"] = commitMap
					}
				}
			}
			body["commits"] = commitsInBody
		}

		// 		fmt.Println(body)
		jsonBody, _ := json.Marshal(body)
		fmt.Println(string(jsonBody))

		client := graphql.NewClient(relizaHubUri + "/graphql")
		req := graphql.NewRequest(`
			mutation ($releaseInputProg: ReleaseInputProg) {
				addReleaseProg(release:$releaseInputProg) {` + RELEASE_GQL_DATA + `}
			}`,
		)
		req.Var("releaseInputProg", body)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Reliza Go Client")
		req.Header.Set("Accept-Encoding", "gzip, deflate")

		if len(apiKeyId) > 0 && len(apiKey) > 0 {
			auth := base64.StdEncoding.EncodeToString([]byte(apiKeyId + ":" + apiKey))
			req.Header.Add("Authorization", "Basic "+auth)
		}

		var respData map[string]interface{}
		if err := client.Run(context.Background(), req, &respData); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		jsonResponse, _ := json.Marshal(respData["addReleaseProg"])
		fmt.Println(string(jsonResponse))
	},
}

var approveReleaseCmd = &cobra.Command{
	Use:   "approverelease",
	Short: "Programmatic approval of releases using valid API key",
	Long: `This CLI command would connect to Reliza Hub and submit approval for a release using valid API key.
			The API key used must be valid and also must be authorized
			to perform requested approval.`,
	Run: func(cmd *cobra.Command, args []string) {
		if debug == "true" {
			fmt.Println("Using Reliza Hub at", relizaHubUri)
		}
		body := map[string]interface{}{}
		approvalMap := map[string]bool{approvalType: !disapprove}
		body["approvals"] = approvalMap
		if len(releaseId) > 0 {
			body["uuid"] = releaseId
		}
		if len(releaseVersion) > 0 {
			body["version"] = releaseVersion
		}
		if len(project) > 0 {
			body["project"] = project
		}
		if len(instance) > 0 {
			body["instance"] = instance
		}
		if len(namespace) > 0 {
			body["namespace"] = namespace
		}

		client := graphql.NewClient(relizaHubUri + "/graphql")
		req := graphql.NewRequest(`
			mutation ($ApproveReleaseInput: ApproveReleaseInput) {
				approveReleaseProg(release:$ApproveReleaseInput) {` + RELEASE_GQL_DATA + `}
			}
		`)
		req.Var("ApproveReleaseInput", body)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Reliza Go Client")
		req.Header.Set("Accept-Encoding", "gzip, deflate")

		if len(apiKeyId) > 0 && len(apiKey) > 0 {
			auth := base64.StdEncoding.EncodeToString([]byte(apiKeyId + ":" + apiKey))
			req.Header.Add("Authorization", "Basic "+auth)
		}

		var respData map[string]interface{}
		if err := client.Run(context.Background(), req, &respData); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		jsonResponse, _ := json.Marshal(respData["approveReleaseProg"])
		fmt.Println(string(jsonResponse))
	},
}

var isApprovalNeededCmd = &cobra.Command{
	Use:   "isapprovalneeded",
	Short: "Check if a release needs to be approvid using valid API key",
	Long: `This CLI command would connect to Reliza Hub and check if a specific release needs to be approved.
			It no longer needs to be approved, if it has been previously approved or rejected.`,
	Run: func(cmd *cobra.Command, args []string) {
		if debug == "true" {
			fmt.Println("Using Reliza Hub at", relizaHubUri)
		}
		body := map[string]interface{}{}
		body["type"] = approvalType
		if len(releaseId) > 0 {
			body["uuid"] = releaseId
		}
		if len(releaseVersion) > 0 {
			body["version"] = releaseVersion
		}
		if len(project) > 0 {
			body["project"] = project
		}
		if len(instance) > 0 {
			body["instance"] = instance
		}
		if len(namespace) > 0 {
			body["namespace"] = namespace
		}

		client := graphql.NewClient(relizaHubUri + "/graphql")
		req := graphql.NewRequest(`
			query ($IsApprovalNeededInput: IsApprovalNeededInput) {
				isApprovalNeeded(release:$IsApprovalNeededInput)
			}
		`)
		req.Var("IsApprovalNeededInput", body)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Reliza Go Client")
		req.Header.Set("Accept-Encoding", "gzip, deflate")

		if len(apiKeyId) > 0 && len(apiKey) > 0 {
			auth := base64.StdEncoding.EncodeToString([]byte(apiKeyId + ":" + apiKey))
			req.Header.Add("Authorization", "Basic "+auth)
		}

		var respData map[string]interface{}
		if err := client.Run(context.Background(), req, &respData); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		jsonResponse, _ := json.Marshal(respData["isApprovalNeeded"])
		fmt.Println(string(jsonResponse))
	},
}

var instDataCmd = &cobra.Command{
	Use:   "instdata",
	Short: "Sends instance data to Reliza Hub",
	Long:  `This CLI command would stream agent data from instance to Reliza Hub`,
	Run: func(cmd *cobra.Command, args []string) {
		if debug == "true" {
			fmt.Println("Using Reliza Hub at", relizaHubUri)
		}

		body := map[string]interface{}{}
		// if imageString (--images flag) is supplied, image File path is ignored
		if imageString != "" {
			// only non-k8s images supported
			body["images"] = strings.Fields(imageString)
		} else {
			imageBytes, err := ioutil.ReadFile(imageFilePath)
			if err != nil {
				fmt.Println("Error when reading images file")
				fmt.Print(err)
				os.Exit(1)
			}
			if imageStyle == "k8s" {
				var k8sjson []map[string]interface{}
				errJson := json.Unmarshal(imageBytes, &k8sjson)
				if errJson != nil {
					fmt.Println("Error unmarshalling k8s images")
					fmt.Println(errJson)
					os.Exit(1)
				}
				body["type"] = "k8s"
				body["images"] = k8sjson
			} else {
				body["images"] = strings.Fields(string(imageBytes))
			}
		}
		body["timeSent"] = time.Now().String()
		if len(namespace) > 0 {
			body["namespace"] = namespace
		}
		if len(senderId) > 0 {
			body["senderId"] = senderId
		}

		if debug == "true" {
			fmt.Println(body)
		}

		client := graphql.NewClient(relizaHubUri + "/graphql")
		req := graphql.NewRequest(`
			mutation ($InstanceDataInput: InstanceDataInput) {
				instData(instance:$InstanceDataInput)
			}
		`)
		req.Var("InstanceDataInput", body)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Reliza Go Client")
		req.Header.Set("Accept-Encoding", "gzip, deflate")

		if len(apiKeyId) > 0 && len(apiKey) > 0 {
			auth := base64.StdEncoding.EncodeToString([]byte(apiKeyId + ":" + apiKey))
			req.Header.Add("Authorization", "Basic "+auth)
		}

		var respData map[string]interface{}
		if err := client.Run(context.Background(), req, &respData); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		jsonResponse, _ := json.Marshal(respData["instData"])
		fmt.Println(string(jsonResponse))
	},
}

var matchBundleCmd = &cobra.Command{
	Use:   "matchbundle",
	Short: "Match images to bundle version",
	Long:  `This CLI command would stream list of images with sha256 digests to Reliza Hub and attempt to match it to product release`,
	Run: func(cmd *cobra.Command, args []string) {
		if debug == "true" {
			fmt.Println("Using Reliza Hub at", relizaHubUri)
		}

		body := map[string]interface{}{}
		// if imageString (--images flag) is supplied, image File path is ignored
		if imageString != "" {
			// only non-k8s images supported
			body["images"] = strings.Fields(imageString)
		} else {
			imageBytes, err := ioutil.ReadFile(imageFilePath)
			if err != nil {
				fmt.Println("Error when reading images file")
				fmt.Print(err)
				os.Exit(1)
			}
			if imageStyle == "k8s" {
				var k8sjson []map[string]interface{}
				errJson := json.Unmarshal(imageBytes, &k8sjson)
				if errJson != nil {
					fmt.Println("Error unmarshalling k8s images")
					fmt.Println(errJson)
					os.Exit(1)
				}
				body["type"] = "k8s"
				body["images"] = k8sjson
			} else {
				body["images"] = strings.Fields(string(imageBytes))
			}
		}
		body["timeSent"] = time.Now().String()
		if len(namespace) > 0 {
			body["namespace"] = namespace
		}

		if debug == "true" {
			fmt.Println(body)
		}

		client := graphql.NewClient(relizaHubUri + "/graphql")
		req := graphql.NewRequest(`
			mutation ($InstanceDataInput: InstanceDataInput) {
				matchToProductRelease(release:$InstanceDataInput)
			}
		`)
		req.Var("InstanceDataInput", body)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Reliza Go Client")
		req.Header.Set("Accept-Encoding", "gzip, deflate")

		if len(apiKeyId) > 0 && len(apiKey) > 0 {
			auth := base64.StdEncoding.EncodeToString([]byte(apiKeyId + ":" + apiKey))
			req.Header.Add("Authorization", "Basic "+auth)
		}

		var respData map[string]interface{}
		if err := client.Run(context.Background(), req, &respData); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		jsonResponse, _ := json.Marshal(respData["matchToProductRelease"])
		fmt.Println(string(jsonResponse))
	},
}

var createProjectCmd = &cobra.Command{
	Use:   "createproject",
	Short: "Create new project",
	Long:  `This CLI command would connect to Reliza Hub which would create a new project `,
	Run: func(cmd *cobra.Command, args []string) {
		if debug == "true" {
			fmt.Println("Using Reliza Hub at", relizaHubUri)
		}

		body := map[string]interface{}{"name": projectName}
		if len(projectType) > 0 {
			body["type"] = strings.ToUpper(projectType)
			if strings.ToUpper(projectType) == "BUNDLE" {
				body["type"] = "PRODUCT"
			}
		}
		if len(versionSchema) > 0 {
			body["versionSchema"] = versionSchema
		}
		if len(featureBranchVersioning) > 0 {
			body["featureBranchVersioning"] = featureBranchVersioning
		}
		if len(vcsUuid) > 0 {
			body["vcsRepositoryUuid"] = vcsUuid
		}

		if len(vcsUri) > 0 {
			vcsRepository := map[string]string{"uri": vcsUri}
			if len(vcsName) > 0 {
				vcsRepository["name"] = vcsName
			}
			if len(vcsType) > 0 {
				vcsRepository["type"] = vcsType
			}
			body["vcsRepository"] = vcsRepository
		}

		body["includeApi"] = includeApi

		client := graphql.NewClient(relizaHubUri + "/graphql")
		req := graphql.NewRequest(`
			mutation ($CreateProjectInput: CreateProjectInput) {
				createProjectProg(project:$CreateProjectInput) {` + PROJECT_GQL_DATA + `}
			}
		`)
		req.Var("CreateProjectInput", body)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Reliza Go Client")
		req.Header.Set("Accept-Encoding", "gzip, deflate")

		if len(apiKeyId) > 0 && len(apiKey) > 0 {
			auth := base64.StdEncoding.EncodeToString([]byte(apiKeyId + ":" + apiKey))
			req.Header.Add("Authorization", "Basic "+auth)
		}

		var respData map[string]interface{}
		if err := client.Run(context.Background(), req, &respData); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		jsonResponse, _ := json.Marshal(respData["createProjectProg"])
		fmt.Println(string(jsonResponse))
	},
}

var getVersionCmd = &cobra.Command{
	Use:   "getversion",
	Short: "Get next version for branch for a particular project",
	Long: `This CLI command would connect to Reliza Hub which would generate next Atomic version for particular project.
			Project would be identified by the API key that is used`,
	Run: func(cmd *cobra.Command, args []string) {
		if debug == "true" {
			fmt.Println("Using Reliza Hub at", relizaHubUri)
		}

		body := map[string]interface{}{"branch": branch}
		if len(project) > 0 {
			body["project"] = project
		}
		if len(modifier) > 0 {
			body["modifier"] = modifier
		}
		if len(metadata) > 0 {
			body["metadata"] = metadata
		}
		if len(action) > 0 {
			body["action"] = action
		}

		if len(versionSchema) > 0 {
			body["versionSchema"] = versionSchema
		}

		if commit != "" {
			commitMap := map[string]string{"uri": vcsUri, "type": vcsType, "commit": commit}
			if vcsTag != "" {
				commitMap["vcsTag"] = vcsTag
			}
			if dateActual != "" {
				commitMap["dateActual"] = dateActual
			}
			body["sourceCodeEntry"] = commitMap
		}
		if manual {
			body["status"] = "draft"
		}

		body["onlyVersion"] = onlyVersion

		client := graphql.NewClient(relizaHubUri + "/graphql")
		req := graphql.NewRequest(`
			mutation ($GetNewVersionInput: GetNewVersionInput) {
				getNewVersion(project:$GetNewVersionInput)
			}
		`)
		req.Var("GetNewVersionInput", body)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Reliza Go Client")
		req.Header.Set("Accept-Encoding", "gzip, deflate")

		if len(apiKeyId) > 0 && len(apiKey) > 0 {
			auth := base64.StdEncoding.EncodeToString([]byte(apiKeyId + ":" + apiKey))
			req.Header.Add("Authorization", "Basic "+auth)
		}

		var respData map[string]interface{}
		if err := client.Run(context.Background(), req, &respData); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		jsonResponse, _ := json.Marshal(respData["getNewVersion"])
		fmt.Println(string(jsonResponse))
	},
}

var checkReleaseByHashCmd = &cobra.Command{
	Use:   "checkhash",
	Short: "Checks whether artifact with this hash is present for particular project",
	Long: `This CLI command would connect to Reliza Hub which would check if the artifact was already submitted as a part of some
			existing release of the current project.
			Project would be identified by the API key that is used`,
	Run: func(cmd *cobra.Command, args []string) {
		if debug == "true" {
			fmt.Println("Using Reliza Hub at", relizaHubUri)
		}

		client := graphql.NewClient(relizaHubUri + "/graphql")
		req := graphql.NewRequest(`
			query ($hash: String!) {
				getReleaseByHash(hash: $hash) {` + RELEASE_GQL_DATA + `}
			}
		`)
		req.Var("hash", hash)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Reliza Go Client")
		req.Header.Set("Accept-Encoding", "gzip, deflate")

		if len(apiKeyId) > 0 && len(apiKey) > 0 {
			auth := base64.StdEncoding.EncodeToString([]byte(apiKeyId + ":" + apiKey))
			req.Header.Add("Authorization", "Basic "+auth)
		}

		var respData map[string]interface{}
		if err := client.Run(context.Background(), req, &respData); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		jsonResponse, _ := json.Marshal(respData["getReleaseByHash"])
		fmt.Println(string(jsonResponse))
	},
}

var getLatestReleaseCmd = &cobra.Command{
	Use:   "getlatestrelease",
	Short: "Obtains latest release for Project or Product",
	Long: `This CLI command would connect to Reliza Hub and would obtain latest release for specified Project and Branch
			or specified Product and Feature Set.`,
	Run: func(cmd *cobra.Command, args []string) {
		getLatestReleaseFunc(debug, relizaHubUri, project, product, branch, environment, tagKey, tagVal, apiKeyId, apiKey, instance, namespace)
	},
}

var getMyReleaseCmd = &cobra.Command{
	Use:   "getmyrelease",
	Short: "Get releases to be deployed on this instance",
	Long: `This CLI command is to be used by programmatic access from instance. 
			It would connect to Reliza Hub which would return release and artifacts versions that should be used on this instance.
			Instance would be identified by the API key that is used`,
	Run: func(cmd *cobra.Command, args []string) {
		client := graphql.NewClient(relizaHubUri + "/graphql")
		req := graphql.NewRequest(`
			query ($namespace: String) {
				getMyRelease(namespace: $namespace) {` + FULL_RELEASE_GQL_DATA + `}
			}
		`)
		req.Var("namespace", namespace)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Reliza Go Client")
		req.Header.Set("Accept-Encoding", "gzip, deflate")

		if len(apiKeyId) > 0 && len(apiKey) > 0 {
			auth := base64.StdEncoding.EncodeToString([]byte(apiKeyId + ":" + apiKey))
			req.Header.Add("Authorization", "Basic "+auth)
		}

		var respData map[string]interface{}
		if err := client.Run(context.Background(), req, &respData); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		jsonResponse, _ := json.Marshal(respData["getMyRelease"])
		fmt.Println(string(jsonResponse))
	},
}

var parseCopyTemplatesCmd = &cobra.Command{
	Use:   "parsetemplate",
	Short: "Parses Reliza templates and copies to target",
	Long: `This CLI command parses template files, replacing project codes styled as 
			<%PROJECT__55fe88fc-5621-4c90-b006-db94ea1d8a08%> - and replaces them with the latest or target
			versions of those projects from Reliza Hub as defined by target environment and tags`,
	Run: func(cmd *cobra.Command, args []string) {

		parseCopyTemplate(parseDirectory, outDirectory, relizaHubUri, environment, tagKey, tagVal, apiKeyId, apiKey,
			instance, namespace)
	},
}

// Modern way to parse templates (re-write over parse copy template)
var replaceTagsCmd = &cobra.Command{
	Use:   "replacetags",
	Short: "Replaces tags in k8s, helm or compose files",
	Long:  `Modern version of parse copy template`,
	Run: func(cmd *cobra.Command, args []string) {
		// v1 - takes inFile = inFile var, outFile = outfile, source txt file, definition reference file - i.e. result of helm template
		// inFile, outFile, tagSourceFile, definitionReferenceFile
		// type - typeVal: options - text, cyclonedx

		// 1st - scan tag source file and construct a map of generic tag to actual tag
		tagSourceMap := scanTags(tagSourceFile, typeVal, apiKeyId, apiKey, instance, revision, instanceURI, bundle, version, environment)

		// 2nd - scan definition reference file and identify all used tags (scan by "image:" pattern)
		substitutionMap := map[string]string{}
		if definitionReferenceFile != "" {
			fmt.Println("Scanning definition references...")
			defFile, fileOpenErr := os.Open(definitionReferenceFile)
			if fileOpenErr != nil {
				fmt.Println(fileOpenErr)
				os.Exit(1)
			}

			// map to store definition images to their replacements -> will be applied on source files
			defScanMap := map[string]string{}

			defScanner := bufio.NewScanner(defFile)
			// input files must be utf-8 !!!
			for defScanner.Scan() {
				line := defScanner.Text()
				if strings.Contains(strings.ToLower(line), "image: ") {
					// extract actual image
					imageLineArray := strings.Split(strings.ToLower(line), "image: ")
					image := imageLineArray[1]
					// remove beginning and ending quotes if present
					re := regexp.MustCompile("^\"")
					image = re.ReplaceAllLiteralString(image, "")
					re = regexp.MustCompile("\"$")
					image = re.ReplaceAllLiteralString(image, "")
					re = regexp.MustCompile("^'")
					image = re.ReplaceAllLiteralString(image, "")
					re = regexp.MustCompile("'$")
					image = re.ReplaceAllLiteralString(image, "")
					// parse and add to map
					if strings.Contains(image, "@") {
						tagSplit := strings.Split(image, "@")
						defScanMap[tagSplit[0]] = image
					} else if strings.Contains(line, ":") {
						tagSplit := strings.SplitN(image, ":", 2)
						defScanMap[tagSplit[0]] = image
					} else {
						defScanMap[image] = image
					}
				}
			}

			// combine 2 maps and come up with substitution map to apply to source (i.e. to source helm chart)
			// traverse defScanMap, map to tagSourceMap and put to substitution map
			for k := range defScanMap {
				// https://stackoverflow.com/questions/2050391/how-to-check-if-a-map-contains-a-key-in-go
				if tagVal, ok := tagSourceMap[k]; ok {
					substitutionMap[k] = tagVal
				}
			}
		} else {
			substitutionMap = tagSourceMap
		}

		// Check if input is infile or inDirectory (operating on directory or file?)
		if len(infile) > 0 && len(inDirectory) == 0 {
			// Make sure infile is a file and not a directory
			fileInfo, err := os.Stat(infile)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			} else if fileInfo.IsDir() {
				fmt.Println("Error: infile must be a path to a file, not a directory!")
				os.Exit(1)
			}
			// Open infile if not directory:
			var inFileOpened *os.File
			var inFileOpenedError error
			inFileOpened, inFileOpenedError = os.Open(infile)
			if inFileOpenedError != nil {
				fmt.Println("Error opening infile: " + infile)
				fmt.Println(inFileOpenedError)
				os.Exit(1)
			}

			// Create outFile to write to, if outfile not specified, write to stdout
			var outFileOpened *os.File
			var outFileOpenedError error
			if len(outfile) > 0 {
				//fmt.Println("Opening output file...")
				outFileOpened, outFileOpenedError = os.Create(outfile)
				if outFileOpenedError != nil {
					fmt.Println("Error opening outfile: " + outfile)
					fmt.Println(outFileOpenedError)
					os.Exit(1)
				}
			}

			// need to add provenance first, beacuse can only write to stdout sequentially
			// check for argument --no-provenance, if true, add provenance to file
			if provenance == true {
				addProvenanceToReplaceTagsOutput(outFileOpened, apiKeyId, apiKey, tagSourceFile, environment, instance, instanceURI, revision, definitionReferenceFile, typeVal, version, bundle)
			}

			// Parse infile and write to outfile with replace tags (or stdout if no outfile)
			substituteCopyBasedOnMap(inFileOpened, outFileOpened, substitutionMap)

			// Remeber to close outfile+infile when done, and check for errors
			//fmt.Println("Closing output file...")
			if outFileOpened != nil { // outfile might not exist if writing to stdout
				outFileCloseError := outFileOpened.Close()
				if outFileCloseError != nil {
					fmt.Println("Error closing outfile: " + outfile)
					fmt.Println(outFileCloseError)
					os.Exit(1)
				}
			}
			// Close infile
			inFileCloseError := inFileOpened.Close()
			if inFileCloseError != nil {
				fmt.Println("Error closing infile: " + infile)
				fmt.Println(inFileCloseError)
				os.Exit(1)
			}
			// No infile input present, operate on inDirectory instead.
		} else if len(infile) == 0 && len(inDirectory) > 0 {
			// If parsing files from input directory, an output directory path should be provided, not an output file path.
			if len(outfile) > 0 {
				fmt.Println("Warning: please only provide '--outdirectory' flag (no '--outfile') when using '--indirectory' as input instead of '--infile'.")
			}
			// Check that outDirectory has value. Cannot write to stdout when parsing multiple files from a directory.
			if len(outDirectory) == 0 {
				fmt.Println("Error: '--outdirectory' is empty. Must supply a path to an output directory when using --indirectory flag.")
				os.Exit(1)
			}
			// Check that inDirectory and out dir end in '/' or '\'
			if string(outDirectory[len(outDirectory)-1:]) != "\\" && string(outDirectory[len(outDirectory)-1:]) != "/" {
				outDirectory = outDirectory + "\\"
			}
			if string(inDirectory[len(inDirectory)-1:]) != "\\" && string(inDirectory[len(inDirectory)-1:]) != "/" {
				inDirectory = inDirectory + "\\"
			}
			// check that outDirectory is a real directory (no stdout output for inDirectory)
			dirInfo, err := os.Stat(outDirectory)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			} else if !dirInfo.IsDir() {
				fmt.Println("Error: outdirectory must be a path to a valid directory!")
				os.Exit(1)
			}
			// Open
			var fileNames []string
			files, err := ioutil.ReadDir(inDirectory)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			// Get slice of names of each file in inDirectory
			for _, f := range files {
				//fmt.Println(f.Name())
				fileNames = append(fileNames, f.Name())
			}

			// replacetags based on substitutionMap for each file in directory
			for _, fileName := range fileNames {
				// Generate path of next output file (same as input file name, but in outDirectory)
				outFilePath := outDirectory + fileName
				// Create outFile to write to inside outDirectory
				var outFileOpened *os.File
				var err error
				outFileOpened, err = os.Create(outFilePath)
				if err != nil {
					fmt.Println("Error opening outfile: " + outFilePath)
					fmt.Println(err)
					os.Exit(1)
				}
				// open infile
				var inFileOpened *os.File
				inFileOpened, err = os.Open(inDirectory + fileName)
				if err != nil {
					fmt.Println("Error opening infile: " + inDirectory + fileName)
					fmt.Println(err)
					os.Exit(1)
				}

				// check for argument --no-provenance, if true, add provenance to file
				if provenance == true {
					addProvenanceToReplaceTagsOutput(outFileOpened, apiKeyId, apiKey, tagSourceFile, environment, instance, instanceURI, revision, definitionReferenceFile, typeVal, version, bundle)
				}

				// Parse infile and write to outfile with replace tags (or stdout if no outfile)
				substituteCopyBasedOnMap(inFileOpened, outFileOpened, substitutionMap)

				// Remeber to close outfile when done, and check for errors
				if outFileOpened != nil {
					outFileCloseError := outFileOpened.Close()
					if outFileCloseError != nil {
						fmt.Println("Error closing outfile: " + outfile)
						fmt.Println(outFileCloseError)
						os.Exit(1)
					}
				}
			}

		} else {
			// either infile and inDirectory provided (too many inputs), or neither provided
			fmt.Println("Error: Must supply either infile or indirectory (but not both)!")
		}
	},
}

var exportInstCmd = &cobra.Command{
	Use:   "exportinst",
	Short: "Outputs the Cyclone DX spec of your instance",
	Long:  `Outputs the Cyclone DX spec of your instance`,
	Run: func(cmd *cobra.Command, args []string) {
		cycloneBytes := getInstanceRevisionCycloneDxExportV1(apiKeyId, apiKey, instance, revision, instanceURI)
		fmt.Println(string(cycloneBytes))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.reliza-cli.yaml)")
	rootCmd.PersistentFlags().StringVarP(&relizaHubUri, "uri", "u", "https://app.relizahub.com", "FQDN of Reliza Hub server")
	rootCmd.PersistentFlags().StringVarP(&apiKey, "apikey", "k", "", "API Key Secret")
	rootCmd.PersistentFlags().StringVarP(&apiKeyId, "apikeyid", "i", "", "API Key ID")
	rootCmd.PersistentFlags().StringVarP(&debug, "debug", "d", "false", "If set to true, print debug details")

	// flags for addrelease command
	addreleaseCmd.PersistentFlags().StringVarP(&branch, "branch", "b", "", "Name of VCS Branch used")
	addreleaseCmd.PersistentFlags().StringVarP(&version, "version", "v", "", "Release version")
	addreleaseCmd.MarkPersistentFlagRequired("version")
	addreleaseCmd.MarkPersistentFlagRequired("branch")
	addreleaseCmd.PersistentFlags().StringVar(&endpoint, "endpoint", "", "Test endpoint for this release")
	addreleaseCmd.PersistentFlags().StringVar(&project, "project", "", "Project UUID for this release if org-wide key is used")
	addreleaseCmd.PersistentFlags().StringVar(&vcsUri, "vcsuri", "", "URI of VCS repository")
	addreleaseCmd.PersistentFlags().StringVar(&vcsType, "vcstype", "", "Type of VCS repository: git, svn, mercurial")
	addreleaseCmd.PersistentFlags().StringVar(&commit, "commit", "", "Commit id")
	addreleaseCmd.PersistentFlags().StringVar(&commitMessage, "commitmessage", "", "Commit message or subject (optional)")
	addreleaseCmd.PersistentFlags().StringVar(&commits, "commits", "", "Base64-encoded list of commits associated with this release, can be obtained with 'git log --date=iso-strict --pretty='%H|||%ad|||%s' | base64 -w 0' command (optional)")
	addreleaseCmd.PersistentFlags().StringVar(&vcsTag, "vcstag", "", "VCS Tag")
	addreleaseCmd.PersistentFlags().StringVar(&dateActual, "date", "", "Commit date and time in iso strict format, use git log --date=iso-strict (optional).")
	addreleaseCmd.PersistentFlags().StringArrayVar(&artId, "artid", []string{}, "Artifact ID (multiple allowed)")
	addreleaseCmd.PersistentFlags().StringArrayVar(&artBuildId, "artbuildid", []string{}, "Artifact Build ID (multiple allowed)")
	addreleaseCmd.PersistentFlags().StringArrayVar(&artBuildUri, "artbuilduri", []string{}, "Artifact Build URI (multiple allowed)")
	addreleaseCmd.PersistentFlags().StringArrayVar(&artCiMeta, "artcimeta", []string{}, "Artifact CI Meta (multiple allowed)")
	addreleaseCmd.PersistentFlags().StringArrayVar(&artType, "arttype", []string{}, "Artifact Type (multiple allowed)")
	addreleaseCmd.PersistentFlags().StringArrayVar(&artDigests, "artdigests", []string{}, "Artifact Digests (multiple allowed, separate several digests for one artifact with commas)")
	addreleaseCmd.PersistentFlags().StringArrayVar(&tagKeyArr, "tagkey", []string{}, "Artifact Tag Keys (multiple allowed, separate several tag keys for one artifact with commas)")
	addreleaseCmd.PersistentFlags().StringArrayVar(&tagValArr, "tagval", []string{}, "Artifact Tag Values (multiple allowed, separate several tag values for one artifact with commas)")
	addreleaseCmd.PersistentFlags().StringArrayVar(&dateStart, "datestart", []string{}, "Artifact Build Start date and time (optional, multiple allowed)")
	addreleaseCmd.PersistentFlags().StringArrayVar(&dateEnd, "dateend", []string{}, "Artifact Build End date and time (optional, multiple allowed)")
	addreleaseCmd.PersistentFlags().StringArrayVar(&artVersion, "artversion", []string{}, "Artifact version, if different from release (multiple allowed)")
	addreleaseCmd.PersistentFlags().StringArrayVar(&artPackage, "artpackage", []string{}, "Artifact package type (multiple allowed)")
	addreleaseCmd.PersistentFlags().StringArrayVar(&artPublisher, "artpublisher", []string{}, "Artifact publisher (multiple allowed)")
	addreleaseCmd.PersistentFlags().StringArrayVar(&artGroup, "artgroup", []string{}, "Artifact group (multiple allowed)")

	addreleaseCmd.PersistentFlags().StringVar(&status, "status", "", "Status of release - set to 'rejected' for failed releases, otherwise 'completed' is used (optional).")

	// flags for approve release command
	approveReleaseCmd.PersistentFlags().StringVar(&releaseId, "releaseid", "", "UUID of release to be approved (either releaseid or releaseversion and project must be set)")
	approveReleaseCmd.PersistentFlags().StringVar(&releaseVersion, "releaseversion", "", "Version of release to be approved (either releaseid or releaseversion and project must be set)")
	approveReleaseCmd.PersistentFlags().StringVar(&project, "project", "", "UUID of project or product which release should be approved (either instance and project or releaseid or releaseversion and project must be set)")
	approveReleaseCmd.PersistentFlags().StringVar(&instance, "instance", "", "UUID or URI of instance for which release should be approved (either instance and project or releaseid or releaseversion and project must be set)")
	approveReleaseCmd.PersistentFlags().StringVar(&namespace, "namespace", "", "Namespace of the instance for which release should be approved (optional, only considered if instance is specified")
	approveReleaseCmd.PersistentFlags().StringVar(&approvalType, "approval", "", "Name of approval to set")
	approveReleaseCmd.PersistentFlags().BoolVar(&disapprove, "disapprove", false, "(Optional) Set --disapprove flag to indicate disapproval instead of approval")
	approveReleaseCmd.MarkPersistentFlagRequired("approval")

	// flags for is approval needed check command
	isApprovalNeededCmd.PersistentFlags().StringVar(&releaseId, "releaseid", "", "UUID of release to be checked (either releaseid or releaseversion and project must be set)")
	isApprovalNeededCmd.PersistentFlags().StringVar(&releaseVersion, "releaseversion", "", "Version of release to be checked (either releaseid or releaseversion and project must be set)")
	isApprovalNeededCmd.PersistentFlags().StringVar(&project, "project", "", "UUID of project or product which release should be checked (either instance and project or releaseid or releaseversion and project must be set)")
	isApprovalNeededCmd.PersistentFlags().StringVar(&instance, "instance", "", "UUID or URI of instance for which release should be checked (either instance and project or releaseid or releaseversion and project must be set)")
	isApprovalNeededCmd.PersistentFlags().StringVar(&namespace, "namespace", "", "Namespace of the instance for which release should be checked (optional, only considered if instance is specified")
	isApprovalNeededCmd.PersistentFlags().StringVar(&approvalType, "approval", "", "Name of approval type to check")
	isApprovalNeededCmd.MarkPersistentFlagRequired("approval")

	// flags for instance data command
	instDataCmd.PersistentFlags().StringVarP(&imageFilePath, "imagefile", "f", "/resources/images", "Path to image file, ignored if --images parameter is supplied")
	instDataCmd.PersistentFlags().StringVar(&imageString, "images", "", "Whitespace separated images with digests or simply digests, if supplied takes precedence over imagefile")
	instDataCmd.PersistentFlags().StringVar(&imageStyle, "imagestyle", "", "Image format style (optional); set to 'k8s' for k8s style formatting, otherwise default string array of digests is assumed")
	instDataCmd.PersistentFlags().StringVar(&namespace, "namespace", "default", "Namespace to submit instance data to")
	instDataCmd.PersistentFlags().StringVar(&senderId, "sender", "default", "Namespace to submit instance data to")

	// flags for match bundle command
	matchBundleCmd.PersistentFlags().StringVarP(&imageFilePath, "imagefile", "f", "/resources/images", "Path to image file, ignored if --images parameter is supplied")
	matchBundleCmd.PersistentFlags().StringVar(&imageString, "images", "", "Whitespace separated images with digests or simply digests, if supplied takes precedence over imagefile")
	matchBundleCmd.PersistentFlags().StringVar(&namespace, "namespace", "default", "Namespace (Optional, exists for compatibility with instance data command).")

	// flags for getmyrelease command
	getMyReleaseCmd.PersistentFlags().StringVar(&namespace, "namespace", "", "Namespace to submit instance data to")

	// flags for createproject command
	createProjectCmd.PersistentFlags().StringVar(&projectName, "name", "", "Name of project to create")
	createProjectCmd.MarkPersistentFlagRequired("name")
	createProjectCmd.PersistentFlags().StringVar(&projectType, "type", "", "Specify to create either a project or bundle")
	createProjectCmd.MarkPersistentFlagRequired("type")
	createProjectCmd.PersistentFlags().StringVar(&versionSchema, "versionschema", "semver", "Version schema of project, default set to semver. Available version schemas: https://github.com/relizaio/versioning")
	createProjectCmd.PersistentFlags().StringVar(&featureBranchVersioning, "featurebranchversioning", "Branch.Micro", "Feature branch version schema of project (Optional, default set to Branch.Micro")
	createProjectCmd.PersistentFlags().StringVar(&vcsUuid, "vcsuuid", "", "Vcs repository UUID (if retreiving existing vcs repository, either vcsuuid or vcsuri must be set)")
	createProjectCmd.PersistentFlags().StringVar(&vcsUri, "vcsuri", "", "Vcs repository URI, if existing repository with uri does not exist and vcsname and vcstype are not set, will attempt to autoparse github, gitlab, and bitbucket uri's")
	createProjectCmd.PersistentFlags().StringVar(&vcsName, "vcsname", "", "Name of vcs repository (Optional - required if creating new vcs repository and uri cannot be parsed)")
	createProjectCmd.PersistentFlags().StringVar(&vcsType, "vcstype", "", "Type of vcs type (Optional - required if creating new vcs repository and uri cannot be parsed)")
	createProjectCmd.PersistentFlags().BoolVar(&includeApi, "includeapi", false, "(Optional) Set --includeapi flag to create and return api key and id for created project during command")

	// flags for get version command
	getVersionCmd.PersistentFlags().StringVarP(&branch, "branch", "b", "", "Name of VCS Branch used")
	getVersionCmd.MarkPersistentFlagRequired("branch")
	getVersionCmd.PersistentFlags().StringVar(&project, "project", "", "Project UUID for this release if org-wide key is used")
	getVersionCmd.PersistentFlags().StringVar(&action, "action", "", "Bump action name: bump | bumppatch | bumpminor | bumpmajor | bumpdate")
	getVersionCmd.PersistentFlags().StringVar(&metadata, "metadata", "", "Version metadata")
	getVersionCmd.PersistentFlags().StringVar(&modifier, "modifier", "", "Version modifier")
	getVersionCmd.PersistentFlags().StringVar(&versionSchema, "pin", "", "Version pin if creating new branch")
	getVersionCmd.PersistentFlags().StringVar(&vcsUri, "vcsuri", "", "URI of VCS repository")
	getVersionCmd.PersistentFlags().StringVar(&vcsType, "vcstype", "", "Type of VCS repository: git, svn, mercurial")
	getVersionCmd.PersistentFlags().StringVar(&commit, "commit", "", "Commit id")
	getVersionCmd.PersistentFlags().StringVar(&vcsTag, "vcstag", "", "VCS Tag")
	getVersionCmd.PersistentFlags().StringVar(&dateActual, "date", "", "Commit date and time in iso strict format, use git log --date=iso-strict (optional).")
	getVersionCmd.PersistentFlags().BoolVar(&manual, "manual", false, "(Optional) Set --manual flag to indicate a manual release.")
	getVersionCmd.PersistentFlags().BoolVar(&onlyVersion, "onlyversion", false, "(Optional) Set --onlyVersion flag to retrieve next version only and not create a release.")

	// flags for check release by hash command
	checkReleaseByHashCmd.PersistentFlags().StringVar(&hash, "hash", "", "Hash of artifact to check")

	// flags for latest project or product release
	getLatestReleaseCmd.PersistentFlags().StringVar(&project, "project", "", "Project or Product UUID from Reliza Hub of project or product from which to obtain latest release")
	getLatestReleaseCmd.PersistentFlags().StringVar(&product, "product", "", "Product UUID from Reliza Hub to condition project release to this product bundle (optional)")
	getLatestReleaseCmd.PersistentFlags().StringVarP(&branch, "branch", "b", "", "Name of branch or Feature Set from Reliza Hub for which latest release is requested, if not supplied UI mapping is used (optional)")
	getLatestReleaseCmd.PersistentFlags().StringVar(&environment, "env", "", "Environment to obtain approvals details from (optional)")
	getLatestReleaseCmd.PersistentFlags().StringVar(&instance, "instance", "", "Instance ID for which to check release (optional)")
	getLatestReleaseCmd.PersistentFlags().StringVar(&namespace, "namespace", "", "Namespace within instance for which to check release, only matters if instance is supplied (optional)")
	getLatestReleaseCmd.PersistentFlags().StringVar(&tagKey, "tagkey", "", "Tag key to use for picking artifact (optional)")
	getLatestReleaseCmd.PersistentFlags().StringVar(&tagVal, "tagval", "", "Tag value to use for picking artifact (optional)")

	// flags for parse template command
	parseCopyTemplatesCmd.PersistentFlags().StringVar(&environment, "env", "", "Environment to obtain approvals details from (optional)")
	parseCopyTemplatesCmd.PersistentFlags().StringVar(&tagKey, "tagkey", "", "Tag key to use for picking artifact (optional)")
	parseCopyTemplatesCmd.PersistentFlags().StringVar(&tagVal, "tagval", "", "Tag value to use for picking artifact (optional)")
	parseCopyTemplatesCmd.PersistentFlags().StringVar(&parseDirectory, "indirectory", "/indir", "Input directory to parse template files from")
	parseCopyTemplatesCmd.PersistentFlags().StringVar(&outDirectory, "outdirectory", "/outdir", "Output directory to output resulting files with substitutions")
	parseCopyTemplatesCmd.PersistentFlags().StringVar(&instance, "instance", "", "Instance ID for which to check release (optional)")
	parseCopyTemplatesCmd.PersistentFlags().StringVar(&namespace, "namespace", "", "Namespace within instance for which to check release (optional)")

	// flags for get tags
	// Now that replacetags is replacing parseCopyTemplateCmd functionality, does replacetags need 'namespace', 'tagKey' or 'tagVal' flags?
	replaceTagsCmd.PersistentFlags().StringVar(&infile, "infile", "", "Input file to parse, such as helm values file or docker compose file")
	replaceTagsCmd.PersistentFlags().StringVar(&outfile, "outfile", "", "Output file with parsed values (optional, if not supplied - outputs to stdout)")
	replaceTagsCmd.PersistentFlags().StringVar(&inDirectory, "indirectory", "", "Path to directory of input files to parse (either infile or indirectory is required)")
	replaceTagsCmd.PersistentFlags().StringVar(&outDirectory, "outdirectory", "", "Path to directory of output files (required if indirectory is used)")
	replaceTagsCmd.PersistentFlags().StringVar(&tagSourceFile, "tagsource", "", "Source file with tags (optional - specify either source file or instance id and revision)")
	replaceTagsCmd.PersistentFlags().StringVar(&environment, "env", "", "Environment for which to generate tags (optional)")
	replaceTagsCmd.PersistentFlags().StringVar(&instance, "instance", "", "Instance UUID for which to generate tags (optional)")
	replaceTagsCmd.PersistentFlags().StringVar(&instanceURI, "instanceuri", "", "Instance URI for which to generate tags (optional)")
	replaceTagsCmd.PersistentFlags().StringVar(&revision, "revision", "", "Instance revision for which to generate tags (optional)")
	replaceTagsCmd.PersistentFlags().StringVar(&definitionReferenceFile, "defsource", "", "Source file for definitions (optional). For helm, should be output of helm template command")
	replaceTagsCmd.PersistentFlags().StringVar(&typeVal, "type", "cyclonedx", "Type of source tags file: cyclonedx (default) or text")
	replaceTagsCmd.PersistentFlags().StringVar(&version, "version", "", "Bundle version for which to generate tags (optional - required when using bundle)")
	replaceTagsCmd.PersistentFlags().StringVar(&bundle, "bundle", "", "Bundle for which to generate tags (optional)")
	replaceTagsCmd.PersistentFlags().BoolVar(&provenance, "provenance", true, "Set --provenance=[true|false] flag to enable/disable adding provenance (metadata) to beginning of outfile. (optional)")

	exportInstCmd.PersistentFlags().StringVar(&instance, "instance", "", "UUID of instance for which export from (optional)")
	exportInstCmd.PersistentFlags().StringVar(&instanceURI, "instanceuri", "", "URI of instance for which to export from (optional)")
	exportInstCmd.PersistentFlags().StringVar(&revision, "revision", "", "Revision of instance for which to export from (optional)")

	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(printversionCmd)
	rootCmd.AddCommand(addreleaseCmd)
	rootCmd.AddCommand(approveReleaseCmd)
	rootCmd.AddCommand(checkReleaseByHashCmd)
	rootCmd.AddCommand(getLatestReleaseCmd)
	rootCmd.AddCommand(getMyReleaseCmd)
	rootCmd.AddCommand(createProjectCmd)
	rootCmd.AddCommand(getVersionCmd)
	rootCmd.AddCommand(instDataCmd)
	rootCmd.AddCommand(matchBundleCmd)
	rootCmd.AddCommand(parseCopyTemplatesCmd)
	rootCmd.AddCommand(replaceTagsCmd)
	rootCmd.AddCommand(exportInstCmd)
	rootCmd.AddCommand(isApprovalNeededCmd)
}

func printResponse(err error, resp *resty.Response) {
	if debug == "true" {
		// Explore response object
		fmt.Println("Response Info:")
		fmt.Println("Error      :", err)
		fmt.Println("Status Code:", resp.StatusCode())
		fmt.Println("Status     :", resp.Status())
		fmt.Println("Time       :", resp.Time())
		fmt.Println("Received At:", resp.ReceivedAt())
		fmt.Println("Body       :\n", resp)
		fmt.Println()
	} else {
		fmt.Println(resp)
	}

	if resp.StatusCode() != 200 {
		fmt.Println("Error Response Info:")
		fmt.Println("Error      :", err)
		var jsonError ErrorBody
		errJson := json.Unmarshal(resp.Body(), &jsonError)
		if errJson != nil {
			fmt.Println("Error when decoding error json data: ", errJson)
		}
		fmt.Println("Error Message:", jsonError.Message)
		fmt.Println("Status Code:", resp.StatusCode())
		fmt.Println("Status     :", resp.Status())
		fmt.Println("Time       :", resp.Time())
		fmt.Println("Received At:", resp.ReceivedAt())
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig(cmd *cobra.Command) {
	v := viper.New()

	if cfgFile != "" {
		// Use config file from the flag.
		v.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// Search config in home directory with name ".reliza-cli" (without extension).
		v.AddConfigPath(home)
		v.SetConfigName(defaultConfigFilename)
	}
	v.SetEnvPrefix(envPrefix)

	// Attempt to read the config file.
	if err := v.ReadInConfig(); err != nil {
		if debug == "true" {
			fmt.Println(err)
		}
	} else {
		if debug == "true" {
			fmt.Println("Using config file:", v.ConfigFileUsed())
		}
	}

	v.AutomaticEnv() // read in environment variables that match
	bindFlags(cmd, v)

}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
