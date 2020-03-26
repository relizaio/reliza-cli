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
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var action string
var apiKeyId string
var apiKey string
var artBuildId []string
var artCiMeta []string
var artDigests []string
var artId []string
var artType []string
var branch string
var cfgFile string
var commit string
var debug string
var environment string
var hash string
var imageFilePath string
var imageString string
var metadata string
var modifier string
var namespace string
var relizaHubUri string
var project string
var senderId string
var tagKey string
var tagKeyArr []string
var tagVal string
var tagValArr []string
var version string
var versionSchema string
var vcsUri string
var vcsTag string
var vcsType string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "relizaGoClient",
	Short: "CLI client for programmatic operations on Reliza Hub",
	Long:  `This CLI client would allo programmatic actions on Reliza Hub.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// 	},
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
		if commit != "" {
			commitMap := map[string]string{"uri": vcsUri, "type": vcsType, "commit": commit}
			if vcsTag != "" {
				commitMap["vcsTag"] = vcsTag
			}
			body["sourceCodeEntry"] = commitMap
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
			//fmt.Println(artifacts)
		}
		fmt.Println(body)
		client := resty.New()
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("User-Agent", "Reliza Go Client").
			SetHeader("Accept-Encoding", "gzip, deflate").
			SetBody(body).
			SetBasicAuth(apiKeyId, apiKey).
			Post(relizaHubUri + "/api/programmatic/v1/release/create")

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
			fmt.Println("Error      :", err)
			fmt.Println("Status Code:", resp.StatusCode())
			fmt.Println("Status     :", resp.Status())
			fmt.Println("Time       :", resp.Time())
			fmt.Println("Received At:", resp.ReceivedAt())
			fmt.Println("Body       :\n", resp)
			fmt.Println()
			os.Exit(1)
		}
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
			body["images"] = strings.Fields(imageString)
		} else {
			imageBytes, err := ioutil.ReadFile(imageFilePath)
			if err != nil {
				fmt.Println("Error when reading images file")
				fmt.Print(err)
				os.Exit(1)
			}
			body["images"] = strings.Fields(string(imageBytes))
		}
		body["timeSent"] = time.Now().String()
		if len(namespace) > 0 {
			body["namespace"] = namespace
		}
		if len(senderId) > 0 {
			body["senderId"] = senderId
		}
		client := resty.New()
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("User-Agent", "Reliza Go Client").
			SetHeader("Accept-Encoding", "gzip, deflate").
			SetBody(body).
			SetBasicAuth(apiKeyId, apiKey).
			Put(relizaHubUri + "/api/programmatic/v1/instance/sendAgentData")

		// Explore response object
		if debug == "true" {
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
			fmt.Println("Response Info:")
			fmt.Println("Error      :", err)
			fmt.Println("Status Code:", resp.StatusCode())
			fmt.Println("Status     :", resp.Status())
			fmt.Println("Time       :", resp.Time())
			fmt.Println("Received At:", resp.ReceivedAt())
			fmt.Println("Body       :\n", resp)
			fmt.Println()
			os.Exit(1)
		}
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

		body := map[string]string{"branch": branch}
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

		client := resty.New()
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("User-Agent", "Reliza Go Client").
			SetHeader("Accept-Encoding", "gzip, deflate").
			SetBody(body).
			SetBasicAuth(apiKeyId, apiKey).
			Post(relizaHubUri + "/api/programmatic/v1/project/getNewVersion")

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
			os.Exit(1)
		}
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

		body := map[string]string{"hash": hash}

		client := resty.New()
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("User-Agent", "Reliza Go Client").
			SetHeader("Accept-Encoding", "gzip, deflate").
			SetBody(body).
			SetBasicAuth(apiKeyId, apiKey).
			Post(relizaHubUri + "/api/programmatic/v1/release/getByHash")

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
			os.Exit(1)
		}
	},
}

var getLatestReleaseCmd = &cobra.Command{
	Use:   "getlatestrelease",
	Short: "Obtains latest release for Project or Product",
	Long: `This CLI command would connect to Reliza Hub and would obtain latest release for specified Project and Branch
			or specified Product and Feature Set.`,
	Run: func(cmd *cobra.Command, args []string) {
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
		fmt.Println(path)

		client := resty.New()
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("User-Agent", "Reliza Go Client").
			SetHeader("Accept-Encoding", "gzip, deflate").
			SetBasicAuth(apiKeyId, apiKey).
			Get(path)

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
			os.Exit(1)
		}
	},
}

var getMyReleaseCmd = &cobra.Command{
	Use:   "getmyrelease",
	Short: "Get releases to be deployed on this instance",
	Long: `This CLI command is to be used by programmatic access from instance. 
			It would connect to Reliza Hub which would return release and artifacts versions that should be used on this instance.
			Instance would be identified by the API key that is used`,
	Run: func(cmd *cobra.Command, args []string) {
		path := relizaHubUri + "/api/programmatic/v1/instance/getMyFollowReleases"
		if len(namespace) > 0 {
			path += "?namespace=" + namespace
		}

		client := resty.New()
		resp, err := client.R().
			SetHeader("User-Agent", "Reliza Go Client").
			SetHeader("Accept-Encoding", "gzip, deflate").
			SetBasicAuth(apiKeyId, apiKey).
			Get(path)

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
			fmt.Println("Status Code:", resp.StatusCode())
			fmt.Println("Status     :", resp.Status())
			fmt.Println("Time       :", resp.Time())
			fmt.Println("Received At:", resp.ReceivedAt())
			os.Exit(1)
		}
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
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.relizaGoClient.yaml)")
	rootCmd.PersistentFlags().StringVarP(&relizaHubUri, "uri", "u", "https://www.relizahub.com", "FQDN of Reliza Hub server")
	rootCmd.PersistentFlags().StringVarP(&apiKey, "apikey", "k", "", "API Key Secret")
	rootCmd.PersistentFlags().StringVarP(&apiKeyId, "apikeyid", "i", "", "API Key ID")
	rootCmd.PersistentFlags().StringVarP(&debug, "debug", "d", "false", "If set to true, print debug details")

	// flags for addrelease command
	addreleaseCmd.PersistentFlags().StringVarP(&branch, "branch", "b", "", "Name of VCS Branch used")
	addreleaseCmd.PersistentFlags().StringVarP(&version, "version", "v", "", "Release version")
	addreleaseCmd.MarkPersistentFlagRequired("version")
	addreleaseCmd.MarkPersistentFlagRequired("branch")
	addreleaseCmd.PersistentFlags().StringVarP(&vcsUri, "vcsuri", "", "", "URI of VCS repository")
	addreleaseCmd.PersistentFlags().StringVarP(&vcsType, "vcstype", "", "", "Type of VCS repository: git, svn")
	addreleaseCmd.PersistentFlags().StringVarP(&commit, "commit", "", "", "Commit id")
	addreleaseCmd.PersistentFlags().StringVarP(&vcsTag, "vcstag", "", "", "VCS Tag")
	addreleaseCmd.PersistentFlags().StringArrayVar(&artId, "artid", []string{}, "Artifact ID (multiple allowed)")
	addreleaseCmd.PersistentFlags().StringArrayVar(&artBuildId, "artbuildid", []string{}, "Artifact Build ID (multiple allowed)")
	addreleaseCmd.PersistentFlags().StringArrayVar(&artCiMeta, "artcimeta", []string{}, "Artifact CI Meta (multiple allowed)")
	addreleaseCmd.PersistentFlags().StringArrayVar(&artType, "arttype", []string{}, "Artifact Type (multiple allowed)")
	addreleaseCmd.PersistentFlags().StringArrayVar(&artDigests, "artdigests", []string{}, "Artifact Digests (multiple allowed, separate several digests for one artifact with commas)")
	addreleaseCmd.PersistentFlags().StringArrayVar(&tagKeyArr, "tagkey", []string{}, "Artifact Tag Keys (multiple allowed, separate several tag keys for one artifact with commas)")
	addreleaseCmd.PersistentFlags().StringArrayVar(&tagValArr, "tagval", []string{}, "Artifact Tag Values (multiple allowed, separate several tag values for one artifact with commas)")

	// flags for instance data command
	instDataCmd.PersistentFlags().StringVarP(&imageFilePath, "imagefile", "f", "/resources/images.txt", "Path to image file, ignored if --images parameter is supplied")
	instDataCmd.PersistentFlags().StringVar(&imageString, "images", "", "Whitespace separated images with digests or simply digests, if supplied takes precedence over imagefile")
	instDataCmd.PersistentFlags().StringVar(&namespace, "namespace", "default", "Namespace to submit instance data to")
	instDataCmd.PersistentFlags().StringVar(&senderId, "sender", "default", "Namespace to submit instance data to")

	// flags for getmyrelease command
	getMyReleaseCmd.PersistentFlags().StringVar(&namespace, "namespace", "", "Namespace to submit instance data to")

	// flags for get version command
	getVersionCmd.PersistentFlags().StringVarP(&branch, "branch", "b", "", "Name of VCS Branch used")
	getVersionCmd.MarkPersistentFlagRequired("branch")
	getVersionCmd.PersistentFlags().StringVar(&action, "action", "", "Bump action name: bump | bumppatch | bumpminor | bumpmajor | bumpdate")
	getVersionCmd.PersistentFlags().StringVar(&metadata, "metadata", "", "Version metadata")
	getVersionCmd.PersistentFlags().StringVar(&modifier, "modifier", "", "Version modifier")
	getVersionCmd.PersistentFlags().StringVar(&versionSchema, "pin", "", "Version pin if creating new branch")

	// flags for check release by hash command
	checkReleaseByHashCmd.PersistentFlags().StringVar(&hash, "hash", "", "Hash of artifact to check")

	// flags for latest project or product release
	getLatestReleaseCmd.PersistentFlags().StringVar(&project, "project", "", "Project or Product UUID from Reliza Hub of project or product from which to obtain latest release")
	getLatestReleaseCmd.PersistentFlags().StringVar(&branch, "branch", "", "Name of branch or Feature Set from Reliza Hub for which latest release is requested")
	getLatestReleaseCmd.PersistentFlags().StringVar(&environment, "env", "", "Environment to obtain approvals details from (optional)")
	getLatestReleaseCmd.PersistentFlags().StringVar(&tagKey, "tagkey", "", "Tag key to use for picking artifact (optional)")
	getLatestReleaseCmd.PersistentFlags().StringVar(&tagVal, "tagval", "", "Tag value to use for picking artifact (optional)")
	getLatestReleaseCmd.MarkPersistentFlagRequired("project")
	getLatestReleaseCmd.MarkPersistentFlagRequired("branch")

	rootCmd.AddCommand(addreleaseCmd)
	rootCmd.AddCommand(instDataCmd)
	rootCmd.AddCommand(getVersionCmd)
	rootCmd.AddCommand(getMyReleaseCmd)
	rootCmd.AddCommand(checkReleaseByHashCmd)
	rootCmd.AddCommand(getLatestReleaseCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".relizaGoClient" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".relizaGoClient")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
