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
	"os"
	"strings"

	"github.com/machinebox/graphql"
	"github.com/spf13/cobra"
)

var properties []string
var secrets []string

func init() {
	instPropsSecretsCmd.PersistentFlags().StringVar(&instance, "instance", "", "UUID of instance for which export from (optional)")
	instPropsSecretsCmd.PersistentFlags().StringVar(&instanceURI, "instanceuri", "", "URI of instance for which to export from (optional)")
	instPropsSecretsCmd.PersistentFlags().StringVar(&revision, "revision", "", "Revision of instance for which to export from (optional, default is -1)")
	instPropsSecretsCmd.PersistentFlags().StringVar(&namespace, "namespace", "", "Use to define specific namespace for instance export (optional, default is 'default')")
	instPropsSecretsCmd.PersistentFlags().StringArrayVar(&properties, "property", []string{}, "Property to resolve (multiple allowed)")
	instPropsSecretsCmd.PersistentFlags().StringArrayVar(&secrets, "secret", []string{}, "Secret to resolve (multiple allowed)")
	rootCmd.AddCommand(instPropsSecretsCmd)
}

var instPropsSecretsCmd = &cobra.Command{
	Use:   "instprops",
	Short: "Used to retrieve specific properties and secrets per instance",
	Long: `Retrieves a list of properties and secrets for specific instance from Reliza Hub.
	Secrets are only returned if allowed to be read by the instance, if the instance has sealed certificate set and in the encrypted form.`,
	Run: func(cmd *cobra.Command, args []string) {
		props := properties
		secrs := secrets
		retrieveInstancePropsSecretsVerbose(props, secrs)
	},
}

func retrieveInstancePropsSecrets(props []string, secrs []string) SecretPropsRHResp {
	var respData SecretPropsRHResp

	if resolveProps {
		if len(instance) <= 0 && len(instanceURI) <= 0 && !strings.HasPrefix(apiKeyId, "INSTANCE__") {
			//throw error and exit
			fmt.Println("instance or instanceURI not specified!")
			os.Exit(1)
		}

		if "" == revision {
			revision = "-1"
		}

		if len(namespace) <= 1 {
			namespace = "default"
		}

		client := graphql.NewClient(relizaHubUri + "/graphql")
		req := graphql.NewRequest(`
			query ($instanceUuid: ID, $instanceUri: String, $revision: Int!, $namespace: String!, $properties: [String], $secrets: [String]) {
				getInstancePropSecrets(instanceUuid: $instanceUuid, instanceUri: $instanceUri, revision: $revision, namespace: $namespace, properties: $properties, secrets: $secrets) {
					properties {
						key
						value
					}
					secrets {
						key
						value
						lastUpdated
					}
				}
			}
		`)
		req.Var("instanceUuid", instance)
		req.Var("instanceUri", instanceURI)
		req.Var("revision", revision)
		req.Var("namespace", namespace)
		req.Var("properties", props)
		req.Var("secrets", secrs)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Reliza CLI")
		req.Header.Set("Accept-Encoding", "gzip, deflate")

		if len(apiKeyId) > 0 && len(apiKey) > 0 {
			auth := base64.StdEncoding.EncodeToString([]byte(apiKeyId + ":" + apiKey))
			req.Header.Add("Authorization", "Basic "+auth)
		}

		if err := client.Run(context.Background(), req, &respData); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	}

	return respData

	// jsonResp, _ := json.Marshal(respData["getInstancePropSecrets"])
	// fmt.Println(string(jsonResp))
	// return respData["getInstancePropSecrets"].(map[string]interface{})
}

func retrieveInstancePropsSecretsVerbose(props []string, secrs []string) {
	respData := retrieveInstancePropsSecrets(props, secrs)
	jsonResp, _ := json.Marshal(respData.Responsewrapper)
	fmt.Println(string(jsonResp[:]))
}

type SecretPropsRHResp struct {
	Responsewrapper SecretPropsRHRespMaps `json:"getInstancePropSecrets"`
}

type SecretPropsRHRespMaps struct {
	Secrets    []ResolvedSecret   `json:"secrets"`
	Properties []ResolvedProperty `json:"properties"`
}
