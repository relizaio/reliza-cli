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
	"fmt"
	"os"
	"strings"

	"github.com/machinebox/graphql"
	"github.com/spf13/cobra"
)

func init() {
	argoGenCmd.PersistentFlags().StringVar(&instance, "instance", "", "UUID of instance for which to generate (either this, or instanceuri must be provided)")
	argoGenCmd.PersistentFlags().StringVar(&instanceURI, "instanceuri", "", "URI of instance for which to generate (either this, or instanceuri must be provided)")
	argoGenCmd.PersistentFlags().StringVar(&revision, "revision", "", "Revision of instance for which to generate (optional, default is -1)")
	argoGenCmd.PersistentFlags().StringVar(&infile, "infile", "", "Input file to parse, such as helm values file or docker compose file")
	argoGenCmd.PersistentFlags().StringVar(&outfile, "outfile", "", "Output file with parsed values (optional, if not supplied - outputs to stdout)")

	argoGenCmd.MarkPersistentFlagRequired("infile")

	rootCmd.AddCommand(argoGenCmd)
}

var argoGenCmd = &cobra.Command{
	Use:   "argogen",
	Short: "Generate ArgoCD specific yaml values",
	Long: `Generates ArgoCD yaml values based on CycloneDX definitions obtained from Reliza Hub.
	This is similar to replacetags, but specific to Argo and outputs files with Argo definitions.
	Secrets are only returned if allowed to be read by the instance, if the instance has sealed certificate set and in the encrypted form.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Obtain CycloneDX definitions
		var respData ProjectAuthResp

		if len(instance) <= 0 && len(instanceURI) <= 0 && !strings.HasPrefix(apiKeyId, "INSTANCE__") {
			//throw error and exit
			fmt.Println("instance or instanceURI not specified!")
			os.Exit(1)
		}

		if len(namespace) <= 1 {
			// TODO: pass argocd namespace name here
			namespace = "default"
		}

		client := graphql.NewClient(relizaHubUri + "/graphql")
		req := graphql.NewRequest(`
			query ($instanceUuid: ID, $instanceUri: String, $artDigest: String!) {
				artifactDownloadSecrets(instanceUuid: $instanceUuid, instanceUri: $instanceUri, artDigest: $artDigest) {
					login
					password
					type
				}
			}
		`)
		req.Var("instanceUuid", instance)
		req.Var("instanceUri", instanceURI)
		req.Var("artDigest", "sha256:dd43ce80c439e1a0d10fcf49df763c7327132c61d4d3552f288d78786be6bb9a")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Reliza CLI")
		req.Header.Set("Accept-Encoding", "gzip, deflate")

		if len(apiKeyId) > 0 && len(apiKey) > 0 {
			auth := base64.StdEncoding.EncodeToString([]byte(apiKeyId + ":" + apiKey))
			req.Header.Add("Authorization", "Basic "+auth)
		}

		if err := client.Run(context.Background(), req, &respData); err != nil {
			printGqlError(err)
			os.Exit(1)
		}

		fmt.Println(respData)
	},
}

type ProjectAuthResp struct {
	Responsewrapper ProjectAuthRespMaps `json:"artifactDownloadSecrets"`
}

type ProjectAuthRespMaps struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Type     string `json:"type"`
}
