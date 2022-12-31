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

package tests

import (
	"os"
	"testing"

	"github.com/relizaio/reliza-cli/cmd"
)

func TestReplaceTags(t *testing.T) {
	tagSourceFile := "mafia_tag_source_cdx.json"
	typeVal := "cyclonedx"
	infile := "values_mafia.yaml"
	replacedTags := cmd.ReplaceTags(tagSourceFile, typeVal, "", "", "", "", "", "", "", "", "", infile, "")
	expectedReplacement, err := os.ReadFile("expected_values_mafia.yaml")
	if err != nil {
		t.Fatalf("failed reading expected values file")
	}
	if replacedTags != string(expectedReplacement) {
		t.Fatalf("replaced tags do not equal expected, actual = " + replacedTags)
	}
}
