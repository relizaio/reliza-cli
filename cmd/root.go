/*
The MIT License (MIT)

Copyright (c) 2020 - 2022 Reliza Incorporated (Reliza (tm), https://reliza.io)

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
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/machinebox/graphql"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/spf13/viper"
)

var action string
var aggregated bool
var apiKeyId string
var apiKey string
var artBuildId []string
var artBuildUri []string
var artBomFilePaths []string
var artCiMeta []string
var artDigests []string
var artId []string
var artType []string
var artifactType string
var artVersion []string
var artPublisher []string
var artGroup []string
var artPackage []string
var approvalType string
var branch string
var bundle string
var cfgFile string
var closedDate string
var commit string
var commit2 string
var commitMessage string
var commits string // base64-encoded list of commits obtained with: git log $LATEST_COMMIT..$CURRENT_COMMIT --date=iso-strict --pretty='%H|||%ad|||%s' | base64 -w 0
var createdDate string
var dateActual string
var dateStart []string
var dateEnd []string
var debug string
var defaultBranch string
var disapprove bool // approve (default) or disapprove
var endpoint string
var environment string
var featureBranchVersioning string
var filePath string
var hash string
var imageFilePath string
var imageString string
var imageStyle string
var includeApi bool
var instance string
var instanceURI string
var manual bool
var mergedDate string
var metadata string
var modifier string
var namespace string
var number string
var onlyVersion bool
var outDirectory string
var parseDirectory string
var inDirectory string
var infile string
var outfile string
var tagSourceFile string
var definitionReferenceFile string
var provenance bool  // add provenance (default), or do not add provenance
var parseMode string // "simple" || "extended" || "strict" mode
var releaseId string
var releaseVersion string
var releaseNs string
var relizaHubUri string
var revision string
var product string
var project string
var projectName string
var projectType string
var senderId string
var state string
var status string
var tagKey string
var tagKeyArr []string
var tagVal string
var tagValArr []string
var targetBranch string
var title string
var typeVal string
var version string
var version2 string
var versionSchema string
var vcsName string
var vcsTag string
var vcsType string
var vcsUri string
var vcsUuid string
var valueFiles []string

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
		tagRecords{
            key
            value
        }
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

type TagRecord struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "reliza-cli",
	Short: "CLI client for programmatic actions on Reliza Hub",
	Long:  `CLI client for programmatic actions on Reliza Hub.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initConfig(cmd)
	},
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

			if len(artBomFilePaths) > 0 && len(artBomFilePaths) != len(artId) {
				fmt.Println("number of --artboms flags must be either zero or match number of --artid flags")
				os.Exit(2)
			} else if len(artBomFilePaths) > 0 {
				for i, bomPath := range artBomFilePaths {

					artifacts[i]["bom"] = ReadBomJsonFromFile(bomPath)
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

					k := make([]TagRecord, 0)
					for j := range tagKeys {
						tr := TagRecord{
							Key:   tagKeys[j],
							Value: tagVals[j],
						}
						k = append(k, tr)
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
					if len(commitParts) > 3 {
						singleCommitEl["commitAuthor"] = commitParts[3]
						singleCommitEl["commitEmail"] = commitParts[4]
					}
					commitsInBody[i] = singleCommitEl

					// if commit is not present but we are here, use first line as commit
					if len(commit) < 1 && i == 0 {
						commitMap := map[string]string{}
						if len(commitParts) > 3 {
							commitMap = map[string]string{"commit": commitParts[0], "dateActual": commitParts[1], "commitMessage": commitParts[2], "commitAuthor": commitParts[3], "commitEmail": commitParts[4]}
						} else {
							commitMap = map[string]string{"commit": commitParts[0], "dateActual": commitParts[1], "commitMessage": commitParts[2]}
						}
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

		if debug == "true" {
			fmt.Println(string(jsonBody))
		}

		req := graphql.NewRequest(`
			mutation ($releaseInputProg: ReleaseInputProg) {
				addReleaseProg(release:$releaseInputProg) {` + RELEASE_GQL_DATA + `}
			}`,
		)
		req.Var("releaseInputProg", body)
		fmt.Println(sendRequest(req, "addReleaseProg"))
	},
}

var addArtifactCmd = &cobra.Command{
	Use:   "addartifact",
	Short: "Add artifacts to a release",
	Long:  `This CLI command would connect to Reliza Hub and add artifacts to a release using a valid API key.`,
	Run: func(cmd *cobra.Command, args []string) {
		if debug == "true" {
			fmt.Println("Using Reliza Hub at", relizaHubUri)
		}

		body := map[string]interface{}{}
		if len(releaseId) > 0 {
			body["release"] = releaseId
		}
		if len(project) > 0 {
			body["project"] = project
		}
		if len(version) > 0 {
			body["version"] = version
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

					k := make([]TagRecord, 0)
					for j := range tagKeys {
						tr := TagRecord{
							Key:   tagKeys[j],
							Value: tagVals[j],
						}
						k = append(k, tr)
					}
					artifacts[i]["tags"] = k
				}
			}

			body["artifacts"] = artifacts
		}

		req := graphql.NewRequest(`
			mutation ($AddArtifactInput: AddArtifactInput) {
				addArtifact(release: $AddArtifactInput) {` + RELEASE_GQL_DATA + `}
			}
		`)
		req.Var("AddArtifactInput", body)
		fmt.Println(sendRequest(req, "addArtifact"))
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

		req := graphql.NewRequest(`
			mutation ($ApproveReleaseInput: ApproveReleaseInput) {
				approveReleaseProg(release:$ApproveReleaseInput) {` + RELEASE_GQL_DATA + `}
			}
		`)
		req.Var("ApproveReleaseInput", body)
		fmt.Println(sendRequest(req, "approveReleaseProg"))
	},
}

var isApprovalNeededCmd = &cobra.Command{
	Use:   "isapprovalneeded",
	Short: "Check if a release needs to be approved using valid API key",
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

		req := graphql.NewRequest(`
			query ($IsApprovalNeededInput: IsApprovalNeededInput) {
				isApprovalNeeded(release:$IsApprovalNeededInput)
			}
		`)
		req.Var("IsApprovalNeededInput", body)
		fmt.Println(sendRequest(req, "isApprovalNeeded"))
	},
}

var downloadableArtifactCmd = &cobra.Command{
	Use:   "addDownloadableArtifact",
	Short: "Add a downloadable artifact to a release using valid API key",
	Long:  `This CLI command would connect to Reliza Hub add downloadable artifact to a release.`,
	Run: func(cmd *cobra.Command, args []string) {
		if debug == "true" {
			fmt.Println("Using Reliza Hub at", relizaHubUri)
		}

		body := map[string]string{}
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
		if len(artifactType) > 0 {
			body["artifactType"] = artifactType
		}
		client := resty.New()
		session, _ := getSession()
		if session != nil {
			client.SetHeader("X-CSRF-Token", session.Token)
			client.SetHeader("Cookie", "JSESSIONID="+session.JSessionId)
		}
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("User-Agent", "Reliza Go Client").
			SetHeader("Accept-Encoding", "gzip, deflate").
			SetHeader("Accept-Encoding", "gzip, deflate").
			SetFile("file", filePath).
			SetFormData(body).
			SetBasicAuth(apiKeyId, apiKey).
			Post(relizaHubUri + "/api/programmatic/v1/artifact/upload")

		printResponse(err, resp)

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
		body["timeSent"] = time.Now().UTC().Format(time.RFC3339)
		if len(namespace) > 0 {
			body["namespace"] = namespace
		}
		if len(senderId) > 0 {
			body["senderId"] = senderId
		}

		if debug == "true" {
			fmt.Println(body)
		}

		req := graphql.NewRequest(`
			mutation ($InstanceDataInput: InstanceDataInput) {
				instData(instance:$InstanceDataInput)
			}
		`)
		req.Var("InstanceDataInput", body)
		fmt.Println(sendRequest(req, "instData"))
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
		body["timeSent"] = time.Now().UTC().Format(time.RFC3339)
		if len(namespace) > 0 {
			body["namespace"] = namespace
		}

		if debug == "true" {
			fmt.Println(body)
		}

		req := graphql.NewRequest(`
			mutation ($InstanceDataInput: InstanceDataInput) {
				matchToProductRelease(release:$InstanceDataInput)
			}
		`)
		req.Var("InstanceDataInput", body)
		fmt.Println(sendRequest(req, "matchToProductRelease"))
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
		if len(defaultBranch) > 0 {
			body["defaultBranch"] = strings.ToUpper(defaultBranch)
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

		req := graphql.NewRequest(`
			mutation ($CreateProjectInput: CreateProjectInput) {
				createProjectProg(project:$CreateProjectInput) {` + PROJECT_GQL_DATA + `}
			}
		`)
		req.Var("CreateProjectInput", body)
		fmt.Println(sendRequest(req, "createProjectProg"))
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

		if commit != "" || commitMessage != "" {
			commitMap := map[string]string{"uri": vcsUri, "type": vcsType, "commit": commit, "commitMessage": commitMessage}
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

		if len(commits) > 0 {
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
					if len(commitParts) > 3 {
						singleCommitEl["commitAuthor"] = commitParts[3]
						singleCommitEl["commitEmail"] = commitParts[4]
					}
					commitsInBody[i] = singleCommitEl

					// if commit is not present but we are here, use first line as commit
					if len(commit) < 1 && i == 0 {
						commitMap := map[string]string{}
						if len(commitParts) > 3 {
							commitMap = map[string]string{"commit": commitParts[0], "dateActual": commitParts[1], "commitMessage": commitParts[2], "commitAuthor": commitParts[3], "commitEmail": commitParts[4]}
						} else {
							commitMap = map[string]string{"commit": commitParts[0], "dateActual": commitParts[1], "commitMessage": commitParts[2]}
						}
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

		body["onlyVersion"] = onlyVersion

		req := graphql.NewRequest(`
			mutation ($GetNewVersionInput: GetNewVersionInput) {
				getNewVersion(project:$GetNewVersionInput)
			}
		`)
		req.Var("GetNewVersionInput", body)
		fmt.Println(sendRequest(req, "getNewVersion"))
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

		req := graphql.NewRequest(`
			query ($hash: String!) {
				getReleaseByHash(hash: $hash) {` + RELEASE_GQL_DATA + `}
			}
		`)
		req.Var("hash", hash)
		resp := sendRequest(req, "getReleaseByHash")
		if resp == "null" {
			resp = "{}"
		}
		fmt.Println(resp)
	},
}

var getLatestReleaseCmd = &cobra.Command{
	Use:   "getlatestrelease",
	Short: "Obtains latest release for Project or Product",
	Long: `This CLI command would connect to Reliza Hub and would obtain latest release for specified Project and Branch
			or specified Product and Feature Set.`,
	Run: func(cmd *cobra.Command, args []string) {
		getLatestReleaseFunc(debug, relizaHubUri, project, product, branch, environment, tagKey, tagVal, apiKeyId, apiKey, instance, namespace, status)
	},
}

var getMyReleaseCmd = &cobra.Command{
	Use:   "getmyrelease",
	Short: "Get releases to be deployed on this instance",
	Long: `This CLI command is to be used by programmatic access from instance. 
			It would connect to Reliza Hub which would return release and artifacts versions that should be used on this instance.
			Instance would be identified by the API key that is used`,
	Run: func(cmd *cobra.Command, args []string) {
		if debug == "true" {
			fmt.Println("Using Reliza Hub at", relizaHubUri)
		}

		req := graphql.NewRequest(`
			query ($namespace: String) {
				getMyRelease(namespace: $namespace) {` + FULL_RELEASE_GQL_DATA + `}
			}
		`)
		req.Var("namespace", namespace)
		fmt.Println(sendRequest(req, "getMyRelease"))
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

var exportInstCmd = &cobra.Command{
	Use:   "exportinst",
	Short: "Outputs the Cyclone DX spec of your instance",
	Long:  `Outputs the Cyclone DX spec of your instance`,
	Run: func(cmd *cobra.Command, args []string) {
		cycloneBytes := getInstanceRevisionCycloneDxExportV1(apiKeyId, apiKey, instance, revision, instanceURI, namespace)
		fmt.Println(string(cycloneBytes))
	},
}

var exportBundleCmd = &cobra.Command{
	Use:   "exportbundle",
	Short: "Outputs the Cyclone DX spec of your bundle",
	Long:  `Outputs the Cyclone DX spec of your bundle`,
	Run: func(cmd *cobra.Command, args []string) {
		cycloneBytes := getBundleVersionCycloneDxExportV1(apiKeyId, apiKey, bundle, environment, version)
		fmt.Println(string(cycloneBytes))
	},
}

var getChangelogCmd = &cobra.Command{
	Use:   "getchangelog",
	Short: "Outputs changelog information of your project",
	Long:  `Outputs changelog information of your project`,
	Run: func(cmd *cobra.Command, args []string) {
		if debug == "true" {
			fmt.Println("Using Reliza Hub at", relizaHubUri)
		}

		if len(commit) > 0 && len(commit2) > 0 {
			req := graphql.NewRequest(`
				query ($commit1: String!, $commit2: String!, $aggregated: Boolean) {
					getChangelogBetweenCommits(commit1: $commit1, commit2: $commit2, aggregated: $aggregated)
				}
			`)
			req.Var("commit1", commit)
			req.Var("commit2", commit2)
			req.Var("aggregated", aggregated)
			fmt.Println(sendRequest(req, "getChangelogBetweenCommits"))
		} else if len(version) > 0 && len(version2) > 0 {
			req := graphql.NewRequest(`
				query ($version1: String!, $version2: String!, $projectUuid: ID, $aggregated: Boolean) {
					getChangelogBetweenVersions(version1: $version1, version2: $version2, projectUuid: $projectUuid, aggregated: $aggregated)
				}
			`)
			req.Var("version1", version)
			req.Var("version2", version2)
			req.Var("projectUuid", project)
			req.Var("aggregated", aggregated)
			fmt.Println(sendRequest(req, "getChangelogBetweenVersions"))
		} else {
			fmt.Println("Error: Either commit and commit2, or version and version2 must be set")
			os.Exit(1)
		}
	},
}

var prDataCmd = &cobra.Command{
	Use:   "prdata",
	Short: "Sends pull request data to Reliza Hub",
	Long:  `This CLI command would stream pull request data from ci to Reliza Hub`,
	Run: func(cmd *cobra.Command, args []string) {
		if debug == "true" {
			fmt.Println("Using Reliza Hub at", relizaHubUri)
		}

		body := map[string]interface{}{"branch": branch}

		if len(state) > 0 {
			body["state"] = state
		}
		if len(project) > 0 {
			body["project"] = project
		}

		if len(targetBranch) > 0 {
			body["targetBranch"] = targetBranch
		}
		if len(endpoint) > 0 {
			body["endpoint"] = endpoint
		}
		if len(title) > 0 {
			body["title"] = title
		}
		if len(createdDate) > 0 {
			body["createdDate"] = createdDate
		}
		if len(closedDate) > 0 {
			body["closedDate"] = closedDate
		}
		if len(mergedDate) > 0 {
			body["mergedDate"] = mergedDate
		}
		if len(number) > 0 {
			body["number"] = number
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
					if len(commitParts) > 3 {
						singleCommitEl["commitAuthor"] = commitParts[3]
						singleCommitEl["commitEmail"] = commitParts[4]
					}
					commitsInBody[i] = singleCommitEl

					// if commit is not present but we are here, use first line as commit
					if len(commit) < 1 && i == 0 {
						commitMap := map[string]string{}
						if len(commitParts) > 3 {
							commitMap = map[string]string{"commit": commitParts[0], "dateActual": commitParts[1], "commitMessage": commitParts[2], "commitAuthor": commitParts[3], "commitEmail": commitParts[4]}
						} else {
							commitMap = map[string]string{"commit": commitParts[0], "dateActual": commitParts[1], "commitMessage": commitParts[2]}
						}
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

		if debug == "true" {
			fmt.Println(body)
		}
		req := graphql.NewRequest(`
			mutation ($PullRequestInput: PullRequestInput) {
				setPRData(pullRequest:$PullRequestInput)
			}
		`)
		req.Var("PullRequestInput", body)
		fmt.Println(sendRequest(req, "prdata"))
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
	addreleaseCmd.PersistentFlags().StringArrayVar(&artBomFilePaths, "artboms", []string{}, "Artifact Sbom file paths (multiple allowed)")

	addreleaseCmd.PersistentFlags().StringVar(&status, "status", "", "Status of release - set to 'rejected' for failed releases, otherwise 'completed' is used (optional).")

	addArtifactCmd.PersistentFlags().StringVar(&releaseId, "releaseid", "", "UUID of release to add artifact to (either releaseid or project, branch, and version must be set)")
	addArtifactCmd.PersistentFlags().StringVar(&project, "project", "", "Project UUID for this release if org-wide key is used")
	addArtifactCmd.PersistentFlags().StringVar(&version, "version", "", "Release version")
	addArtifactCmd.PersistentFlags().StringArrayVar(&artId, "artid", []string{}, "Artifact ID (multiple allowed)")
	addArtifactCmd.PersistentFlags().StringArrayVar(&artBuildId, "artbuildid", []string{}, "Artifact Build ID (multiple allowed)")
	addArtifactCmd.PersistentFlags().StringArrayVar(&artBuildUri, "artbuilduri", []string{}, "Artifact Build URI (multiple allowed)")
	addArtifactCmd.PersistentFlags().StringArrayVar(&artCiMeta, "artcimeta", []string{}, "Artifact CI Meta (multiple allowed)")
	addArtifactCmd.PersistentFlags().StringArrayVar(&artType, "arttype", []string{}, "Artifact Type (multiple allowed)")
	addArtifactCmd.PersistentFlags().StringArrayVar(&artDigests, "artdigests", []string{}, "Artifact Digests (multiple allowed, separate several digests for one artifact with commas)")
	addArtifactCmd.PersistentFlags().StringArrayVar(&tagKeyArr, "tagkey", []string{}, "Artifact Tag Keys (multiple allowed, separate several tag keys for one artifact with commas)")
	addArtifactCmd.PersistentFlags().StringArrayVar(&tagValArr, "tagval", []string{}, "Artifact Tag Values (multiple allowed, separate several tag values for one artifact with commas)")
	addArtifactCmd.PersistentFlags().StringArrayVar(&dateStart, "datestart", []string{}, "Artifact Build Start date and time (optional, multiple allowed)")
	addArtifactCmd.PersistentFlags().StringArrayVar(&dateEnd, "dateend", []string{}, "Artifact Build End date and time (optional, multiple allowed)")
	addArtifactCmd.PersistentFlags().StringArrayVar(&artVersion, "artversion", []string{}, "Artifact version, if different from release (multiple allowed)")
	addArtifactCmd.PersistentFlags().StringArrayVar(&artPackage, "artpackage", []string{}, "Artifact package type (multiple allowed)")
	addArtifactCmd.PersistentFlags().StringArrayVar(&artPublisher, "artpublisher", []string{}, "Artifact publisher (multiple allowed)")
	addArtifactCmd.PersistentFlags().StringArrayVar(&artGroup, "artgroup", []string{}, "Artifact group (multiple allowed)")

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

	// flags for is approval needed check command
	downloadableArtifactCmd.PersistentFlags().StringVar(&releaseId, "releaseid", "", "UUID of release (either releaseid or releaseversion and project must be set)")
	downloadableArtifactCmd.PersistentFlags().StringVar(&releaseVersion, "releaseversion", "", "Version of release (either releaseid or releaseversion and project must be set)")
	downloadableArtifactCmd.PersistentFlags().StringVar(&project, "project", "", "UUID of project or product for release (either instance and project or releaseid or releaseversion and project must be set)")
	downloadableArtifactCmd.PersistentFlags().StringVar(&instance, "instance", "", "UUID or URI of instance for release (either instance and project or releaseid or releaseversion and project must be set)")
	downloadableArtifactCmd.PersistentFlags().StringVar(&namespace, "namespace", "", "Namespace of the instance for release (optional, only considered if instance is specified")
	downloadableArtifactCmd.PersistentFlags().StringVarP(&filePath, "file", "f", "", "Path to the artifact")
	downloadableArtifactCmd.PersistentFlags().StringVar(&artifactType, "artifactType", "GENERIC", "Type of artifact - can be (TEST_REPORT, SECURITY_SCAN, DOCUMENTATION, GENERIC) or some user defined value")

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
	createProjectCmd.PersistentFlags().StringVar(&defaultBranch, "defaultbranch", "main", "Default branch name of project, default set to main. Available names are either main or master.")
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
	getVersionCmd.PersistentFlags().StringVar(&commit, "commit", "", "Commit id (required to create Source Code Entry for new release)")
	getVersionCmd.PersistentFlags().StringVar(&commitMessage, "commitmessage", "", "Commit message or subject (optional)")
	getVersionCmd.PersistentFlags().StringVar(&commits, "commits", "", "Base64-encoded list of commits associated with this release, can be obtained with 'git log --date=iso-strict --pretty='%H|||%ad|||%s' | base64 -w 0' command (optional)")
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
	getLatestReleaseCmd.PersistentFlags().StringVar(&status, "status", "", "Status of the release, default is completed (optional)")

	// flags for parse template command
	parseCopyTemplatesCmd.PersistentFlags().StringVar(&environment, "env", "", "Environment to obtain approvals details from (optional)")
	parseCopyTemplatesCmd.PersistentFlags().StringVar(&tagKey, "tagkey", "", "Tag key to use for picking artifact (optional)")
	parseCopyTemplatesCmd.PersistentFlags().StringVar(&tagVal, "tagval", "", "Tag value to use for picking artifact (optional)")
	parseCopyTemplatesCmd.PersistentFlags().StringVar(&parseDirectory, "indirectory", "/indir", "Input directory to parse template files from")
	parseCopyTemplatesCmd.PersistentFlags().StringVar(&outDirectory, "outdirectory", "/outdir", "Output directory to output resulting files with substitutions")
	parseCopyTemplatesCmd.PersistentFlags().StringVar(&instance, "instance", "", "Instance ID for which to check release (optional)")
	parseCopyTemplatesCmd.PersistentFlags().StringVar(&namespace, "namespace", "", "Namespace within instance for which to check release (optional)")

	exportInstCmd.PersistentFlags().StringVar(&instance, "instance", "", "UUID of instance for which export from (optional)")
	exportInstCmd.PersistentFlags().StringVar(&instanceURI, "instanceuri", "", "URI of instance for which to export from (optional)")
	exportInstCmd.PersistentFlags().StringVar(&revision, "revision", "", "Revision of instance for which to export from (optional, default is -1)")
	exportInstCmd.PersistentFlags().StringVar(&namespace, "namespace", "", "Use to define specific namespace for instance export (optional)")

	exportBundleCmd.PersistentFlags().StringVar(&bundle, "bundle", "", "Bundle for which to export from")
	exportBundleCmd.PersistentFlags().StringVar(&version, "version", "", "Bundle version for which to export from, either version or environment must be set")
	exportBundleCmd.PersistentFlags().StringVar(&environment, "env", "", "Bundle environment for which to export from last approved bundle, either version or environment must be set")
	exportBundleCmd.MarkPersistentFlagRequired("bundle")

	getChangelogCmd.PersistentFlags().StringVar(&project, "project", "", "Project UUID if org-wide key is used and attaining changelog using versions")
	getChangelogCmd.PersistentFlags().StringVar(&commit, "commit", "", "Commit id, either this and commit2 or version and version2 must be supplied")
	getChangelogCmd.PersistentFlags().StringVar(&commit2, "commit2", "", "Second commit id to construct changelog from")
	getChangelogCmd.PersistentFlags().StringVar(&version, "version", "", "Release version, either this and version2 or commit and commit2 must be supplied")
	getChangelogCmd.PersistentFlags().StringVar(&version2, "version2", "", "Second release version to construct changelog from")
	getChangelogCmd.PersistentFlags().BoolVar(&aggregated, "aggregated", false, "something")

	prDataCmd.PersistentFlags().StringVarP(&branch, "branch", "b", "", "Name of VCS Branch used")
	prDataCmd.PersistentFlags().StringVarP(&state, "state", "s", "", "State of the Pull Request")
	prDataCmd.PersistentFlags().StringVarP(&targetBranch, "targetBranch", "t", "", "Name of target branch")
	prDataCmd.PersistentFlags().StringVar(&project, "project", "", "Project UUID if org-wide key is used")
	prDataCmd.PersistentFlags().StringVar(&endpoint, "endpoint", "", "HTML endpoint of the Pull Request")
	prDataCmd.PersistentFlags().StringVar(&title, "title", "", "Title of the Pull Request")
	prDataCmd.PersistentFlags().StringVar(&createdDate, "createdDate", "", "Datetime when the Pull Request was created")
	prDataCmd.PersistentFlags().StringVar(&closedDate, "closedDate", "", "Datetime when the Pull Request was closed")
	prDataCmd.PersistentFlags().StringVar(&mergedDate, "mergedDate", "", "Datetime when the Pull Request was merged")
	prDataCmd.PersistentFlags().StringVar(&number, "number", "", "Number of the Pull Request")
	prDataCmd.PersistentFlags().StringVar(&commit, "commit", "", "SHA of current commit on the Pull Request (will be merged with existing list)")
	prDataCmd.PersistentFlags().StringVar(&commitMessage, "commitmessage", "", "Commit message or subject (optional)")
	prDataCmd.PersistentFlags().StringVar(&vcsUri, "vcsuri", "", "URI of VCS repository")
	prDataCmd.PersistentFlags().StringVar(&vcsType, "vcstype", "", "Type of VCS repository: git, svn, mercurial")
	prDataCmd.PersistentFlags().StringVar(&commits, "commits", "", "Base64-encoded list of commits associated with this pull request event, can be obtained with 'git log --date=iso-strict --pretty='%H|||%ad|||%s' | base64 -w 0' command (optional)")
	prDataCmd.PersistentFlags().StringVar(&vcsTag, "vcstag", "", "VCS Tag")
	prDataCmd.PersistentFlags().StringVar(&dateActual, "commitdate", "", "Commit date and time in iso strict format, use git log --date=iso-strict (optional).")

	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(printversionCmd)
	rootCmd.AddCommand(addreleaseCmd)
	rootCmd.AddCommand(addArtifactCmd)
	rootCmd.AddCommand(approveReleaseCmd)
	rootCmd.AddCommand(checkReleaseByHashCmd)
	rootCmd.AddCommand(getLatestReleaseCmd)
	rootCmd.AddCommand(getMyReleaseCmd)
	rootCmd.AddCommand(createProjectCmd)
	rootCmd.AddCommand(getVersionCmd)
	rootCmd.AddCommand(instDataCmd)
	rootCmd.AddCommand(matchBundleCmd)
	rootCmd.AddCommand(parseCopyTemplatesCmd)
	rootCmd.AddCommand(exportInstCmd)
	rootCmd.AddCommand(exportBundleCmd)
	rootCmd.AddCommand(getChangelogCmd)
	rootCmd.AddCommand(isApprovalNeededCmd)
	rootCmd.AddCommand(downloadableArtifactCmd)
	rootCmd.AddCommand(prDataCmd)
}

func sendRequest(req *graphql.Request, endpoint string) string {
	return sendRequestWithUri(req, endpoint, relizaHubUri+"/graphql")
}

func sendRequestWithUri(req *graphql.Request, endpoint string, uri string) string {
	session, _ := getSession()
	// if err != nil {
	// 	fmt.Printf("Error making API request: %s\n", err)
	// 	os.Exit(1)
	// }

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Reliza Go Client")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	if session != nil {
		req.Header.Set("X-CSRF-Token", session.Token)
		req.Header.Set("Cookie", "JSESSIONID="+session.JSessionId)
	}
	if len(apiKeyId) > 0 && len(apiKey) > 0 {
		auth := base64.StdEncoding.EncodeToString([]byte(apiKeyId + ":" + apiKey))
		req.Header.Add("Authorization", "Basic "+auth)
	}

	var respData map[string]interface{}
	client := graphql.NewClient(uri)
	if err := client.Run(context.Background(), req, &respData); err != nil {
		printGqlError(err)
		os.Exit(1)
	}

	jsonResponse, _ := json.Marshal(respData[endpoint])
	return string(jsonResponse)
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

func printGqlError(err error) {
	splitError := strings.Split(err.Error(), ":")
	fmt.Println("Error: ", splitError[len(splitError)-1])
}

func getSession() (*RequestSession, error) {
	client := resty.New()
	var result map[string]string
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "Reliza Go Client").
		SetHeader("Accept-Encoding", "gzip, deflate").
		SetResult(&result).
		Get(relizaHubUri + "/api/manual/v1/fetchCsrf")

	if err != nil {
		return nil, err
	}
	// Extract cookies
	session, err := getJSessionIDCookieAndToken(resp)

	if err != nil {
		return nil, err
	}

	return session, err
}

func getJSessionIDCookieAndToken(resp *resty.Response) (*RequestSession, error) {
	// Extract cookies
	cookies := resp.Cookies()
	var jsessionid string
	for _, cookie := range cookies {
		if cookie.Name == "JSESSIONID" {
			jsessionid = cookie.Value
			break
		}
	}

	if jsessionid == "" {
		return nil, fmt.Errorf("JSESSIONID cookie not found")
	}

	// Assume the token is returned in the response body as a JSON object
	var result map[string]interface{}
	err := json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %s", err)
	}

	token, ok := result["token"].(string)
	if !ok {
		return nil, fmt.Errorf("token not found in the response body")
	}

	return &RequestSession{JSessionId: jsessionid, Token: token}, nil
}

type RequestSession struct {
	JSessionId string
	Token      string
}
