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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/machinebox/graphql"
	"github.com/spf13/cobra"
)

var rebomUri string

type BomInput struct {
	Meta string                 `json:"meta"`
	Bom  map[string]interface{} `json:"bom"`
	Tags map[string]interface{} `json:"tags"`
}

func init() {
	putBomCmd.PersistentFlags().StringVar(&infile, "infile", "", "Input file with bom json")
	putBomCmd.PersistentFlags().StringVar(&rebomUri, "rebomuri", "http://localhost:4000", "Rebom URI")
	putBomCmd.MarkPersistentFlagRequired("infile")

	rootCmd.AddCommand(rebomCmd)
	rebomCmd.AddCommand(putBomCmd)
}

var rebomCmd = &cobra.Command{
	Use:   "rebom",
	Short: "Set of commands to interact with rebom tool",
	Long:  `Set of commands to interact with rebom tool`,
	Run: func(cmd *cobra.Command, args []string) {
		addBomToRebomFunc()
	},
}

var putBomCmd = &cobra.Command{
	Use:   "put",
	Short: "Send bom file to rebom",
	Long:  `Send bom file to rebom`,
	Run: func(cmd *cobra.Command, args []string) {
		addBomToRebomFunc()
	},
}

func addBomToRebomFunc() {
	// open infile
	// Make sure infile is a file and not a directory
	fileInfo, err := os.Stat(infile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else if fileInfo.IsDir() {
		fmt.Println("Error: infile must be a path to a file, not a directory!")
		os.Exit(1)
	}
	// Read infile if not directory:
	fileContentByteSlice, _ := ioutil.ReadFile(infile)
	// fileContent := string(fileContentByteSlice)

	// Parse file content into json
	var bomJSON map[string]interface{}
	parseError := json.Unmarshal(fileContentByteSlice, &bomJSON)
	if parseError != nil {
		fmt.Println("Error unmarshalling json bom file")
		fmt.Println(parseError)
		os.Exit(1)
	}
	//var bomJSON map[string]interface{}
	// json.Unmarshal(jsonizedTResult, &bomJSON)
	// cleanJson, _ := json.Marshal(bomJSON)

	var bomInput BomInput
	bomInput.Meta = "sent from reliza cli"
	bomInput.Bom = bomJSON

	fmt.Println(bomInput)

	req := graphql.NewRequest(`
		mutation addBom ($bomInput: BomInput) {
			addBom(bomInput: $bomInput) {
				uuid
				meta
			}
		}
	`)
	req.Var("bomInput", bomInput)
	fmt.Println("adding bom...")
	fmt.Println(sendRequestWithUri(req, "addBom", rebomUri+"/graphql"))
}
