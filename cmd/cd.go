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

var (
	artDigest  string
	sealedCert string
)

func init() {
	rootCmd.AddCommand(cdCmd)

	artifactGetSecrets.PersistentFlags().StringVar(&instance, "instance", "", "UUID of instance for which to generate (either this, or instanceuri must be provided)")
	artifactGetSecrets.PersistentFlags().StringVar(&instanceURI, "instanceuri", "", "URI of instance for which to generate (either this, or instanceuri must be provided)")
	artifactGetSecrets.PersistentFlags().StringVar(&artDigest, "artdigest", "", "Digest or hash of the artifact to resolve secrets for")
	artifactGetSecrets.PersistentFlags().StringVar(&namespace, "namespace", "", "Namespace to use for secrets (optional, defaults to default namespace)")
	artifactGetSecrets.MarkPersistentFlagRequired("artdigest")

	cdCmd.AddCommand(artifactGetSecrets)

	isInstHasSecretCertCmd.PersistentFlags().StringVar(&instance, "instance", "", "UUID of instance for which export from (optional)")
	isInstHasSecretCertCmd.PersistentFlags().StringVar(&instanceURI, "instanceuri", "", "URI of instance for which to export from (optional)")
	cdCmd.AddCommand(isInstHasSecretCertCmd)

	setInstSecretCertCmd.PersistentFlags().StringVar(&instance, "instance", "", "UUID of instance for which export from (optional)")
	setInstSecretCertCmd.PersistentFlags().StringVar(&instanceURI, "instanceuri", "", "URI of instance for which to export from (optional)")
	setInstSecretCertCmd.PersistentFlags().StringVar(&sealedCert, "cert", "", "Sealed certificate used by the instance (required)")
	cdCmd.AddCommand(setInstSecretCertCmd)
}

var cdCmd = &cobra.Command{
	Use:   "cd",
	Short: "Umbrella for commands specific to CD",
	Long:  `Set of commands for continuous delivery.`,
}

var artifactGetSecrets = &cobra.Command{
	Use:   "artsecrets",
	Short: "Get secrets to download specific artifact",
	Long: `Command to get secrets for specific. Artifact must belong to the organizaiton.
			Secret names are returned`,
	Run: func(cmd *cobra.Command, args []string) {
		var respData ProjectAuthResp

		if len(instance) <= 0 && len(instanceURI) <= 0 && !strings.HasPrefix(apiKeyId, "INSTANCE__") {
			fmt.Println("instance or instanceURI not specified!")
			os.Exit(1)
		}

		if len(namespace) <= 1 {
			namespace = "default"
		}

		client := graphql.NewClient(relizaHubUri + "/graphql")
		req := graphql.NewRequest(`
			query ($instanceUuid: ID, $instanceUri: String, $artDigest: String!, $namespace: String) {
				artifactDownloadSecrets(instanceUuid: $instanceUuid, instanceUri: $instanceUri, artDigest: $artDigest, namespace: $namespace) {
					login
					password
					type
				}
			}
		`)
		req.Var("instanceUuid", instance)
		req.Var("instanceUri", instanceURI)
		req.Var("artDigest", artDigest)
		req.Var("namespace", namespace)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Reliza CLI")
		req.Header.Set("Accept-Encoding", "gzip, deflate")

		if len(apiKeyId) > 0 && len(apiKey) > 0 {
			auth := base64.StdEncoding.EncodeToString([]byte(apiKeyId + ":" + apiKey))
			req.Header.Add("Authorization", "Basic "+auth)
		}

		session, _ := getSession()
		if session != nil {
			req.Header.Set("X-CSRF-Token", session.Token)
			req.Header.Set("Cookie", "JSESSIONID="+session.JSessionId)
		}
		if err := client.Run(context.Background(), req, &respData); err != nil {
			printGqlError(err)
			os.Exit(1)
		}

		respJson, err := json.Marshal(respData)
		if err != nil {
			panic(err)
		}

		fmt.Print(string(respJson))
	},
}

var isInstHasSecretCertCmd = &cobra.Command{
	Use:   "iscertinit",
	Short: "Use to check whether instance has sealed cert property configured",
	Long: `Bitnami Sealed Certificate property is used to encrypt secrets for instance.
	This command checks whether this property is configured for the particular instance.`,
	Run: func(cmd *cobra.Command, args []string) {
		var respData IsHasCertRHResp
		client := graphql.NewClient(relizaHubUri + "/graphql")
		req := graphql.NewRequest(`
			query ($instanceUuid: ID, $instanceUri: String) {
				isInstanceHasSealedSecretCert(instanceUuid: $instanceUuid, instanceUri: $instanceUri)
			}
		`)
		req.Var("instanceUuid", instance)
		req.Var("instanceUri", instanceURI)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Reliza CLI")
		req.Header.Set("Accept-Encoding", "gzip, deflate")

		if len(apiKeyId) > 0 && len(apiKey) > 0 {
			auth := base64.StdEncoding.EncodeToString([]byte(apiKeyId + ":" + apiKey))
			req.Header.Add("Authorization", "Basic "+auth)
		}
		session, _ := getSession()
		if session != nil {
			req.Header.Set("X-CSRF-Token", session.Token)
			req.Header.Set("Cookie", "JSESSIONID="+session.JSessionId)
		}

		if err := client.Run(context.Background(), req, &respData); err != nil {
			printGqlError(err)
			os.Exit(1)
		}

		jsonResp, _ := json.Marshal(respData.Responsewrapper)
		fmt.Println(string(jsonResp))
	},
}

var setInstSecretCertCmd = &cobra.Command{
	Use:   "setsecretcert",
	Short: "Use to to set sealed cert property on the instance",
	Long: `Bitnami Sealed Certificate property is used to encrypt secrets for instance.
	This command sets this certificate for the particular instance.
	Only supports instance own API Key.`,
	Run: func(cmd *cobra.Command, args []string) {
		var respData SetCertRHResp
		client := graphql.NewClient(relizaHubUri + "/graphql")
		req := graphql.NewRequest(`
			mutation ($instanceUuid: ID, $instanceUri: String, $sealedCert: String!) {
				setInstanceSealedSecretCert(instanceUuid: $instanceUuid, instanceUri: $instanceUri,
					sealedCert: $sealedCert)
			}
		`)
		req.Var("instanceUuid", instance)
		req.Var("instanceUri", instanceURI)
		req.Var("sealedCert", sealedCert)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Reliza CLI")
		req.Header.Set("Accept-Encoding", "gzip, deflate")

		if len(apiKeyId) > 0 && len(apiKey) > 0 {
			auth := base64.StdEncoding.EncodeToString([]byte(apiKeyId + ":" + apiKey))
			req.Header.Add("Authorization", "Basic "+auth)
		}

		session, _ := getSession()
		if session != nil {
			req.Header.Set("X-CSRF-Token", session.Token)
			req.Header.Set("Cookie", "JSESSIONID="+session.JSessionId)
		}
		if err := client.Run(context.Background(), req, &respData); err != nil {
			printGqlError(err)
			os.Exit(1)
		}

		jsonResp, _ := json.Marshal(respData.Responsewrapper)
		fmt.Println(string(jsonResp))
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
