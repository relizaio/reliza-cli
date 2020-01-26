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

var cfgFile string
var relizaHubUri string
var branch string
var version string
var apiKeyId string
var apiKey string
var vcsUri string
var vcsType string
var commit string
var vcsTag string
var artId []string
var artBuildId []string
var artCiMeta []string
var artType []string
var artDigests []string
var imageFilePath string

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
		fmt.Println("Using Reliza Hub at", relizaHubUri)

		body := map[string]interface{}{"branch": branch, "version": version}
		if commit != "" && vcsType != "" && vcsUri != "" {
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
			body["artifacts"] = artifacts
			//fmt.Println(artifacts)
		}
		client := resty.New()
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("User-Agent", "Reliza Go Client").
			SetHeader("Accept-Encoding", "gzip, deflate").
			SetBody(body).
			SetBasicAuth(apiKeyId, apiKey).
			Post(relizaHubUri + "/api/programmatic/v1/release/create")

		// // Explore response object
		fmt.Println("Response Info:")
		fmt.Println("Error      :", err)
		fmt.Println("Status Code:", resp.StatusCode())
		fmt.Println("Status     :", resp.Status())
		fmt.Println("Time       :", resp.Time())
		fmt.Println("Received At:", resp.ReceivedAt())
		fmt.Println("Body       :\n", resp)
		fmt.Println()

		if resp.StatusCode() != 200 {
			os.Exit(1)
		}
	},
}

var instDataCmd = &cobra.Command{
	Use:   "instdata",
	Short: "Sends instance data to Reliza Hub",
	Long:  `This CLI command would stream agent data from instance to Reliza Hub`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Using Reliza Hub at", relizaHubUri)

		body := map[string]interface{}{"branch": branch, "version": version}
		imageBytes, err := ioutil.ReadFile(imageFilePath)
		if err != nil {
			fmt.Println("Error when reading images file")
			fmt.Print(err)
			os.Exit(1)
		}
		body["images"] = strings.Fields(string(imageBytes))
		body["timeSent"] = time.Now().String()
		client := resty.New()
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("User-Agent", "Reliza Go Client").
			SetHeader("Accept-Encoding", "gzip, deflate").
			SetBody(body).
			SetBasicAuth(apiKeyId, apiKey).
			Put(relizaHubUri + "/api/programmatic/v1/instance/sendAgentData")

		// Explore response object
		fmt.Println("Response Info:")
		fmt.Println("Error      :", err)
		fmt.Println("Status Code:", resp.StatusCode())
		fmt.Println("Status     :", resp.Status())
		fmt.Println("Time       :", resp.Time())
		fmt.Println("Received At:", resp.ReceivedAt())
		fmt.Println("Body       :\n", resp)
		fmt.Println()

		if resp.StatusCode() != 200 {
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
	rootCmd.PersistentFlags().StringVarP(&relizaHubUri, "uri", "u", "https://relizahub.com", "FQDN of Reliza Hub server")
	rootCmd.PersistentFlags().StringVarP(&apiKey, "apikey", "k", "", "API Key Secret")
	rootCmd.PersistentFlags().StringVarP(&apiKeyId, "apikeyid", "i", "", "API Key ID")

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

	// flags for instance data command
	instDataCmd.PersistentFlags().StringVarP(&imageFilePath, "imagefile", "f", "/resources/images.txt", "Path to image file")

	rootCmd.AddCommand(addreleaseCmd)
	rootCmd.AddCommand(instDataCmd)
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
