//
// Blackfriday Markdown Processor
// Available at http://github.com/russross/blackfriday
//
// Copyright © 2011 Russ Ross <russ@russross.com>.
// Distributed under the Simplified BSD License.
// See README.md for details.
//

//
// Unit tests for html rendering
//

package blackfriday

import (
	"testing"
)

func doTestsMakeAnchor(t *testing.T, tests []string) {
	// catch and report panics
	var candidate string

	for i := 0; i+1 < len(tests); i += 2 {
		input := tests[i]
		candidate = input
		expected := tests[i+1]
		actual := makeAnchorText([]byte(candidate))
		if actual != expected {
			t.Errorf("\nInput   [%#v]\nExpected[%#v]\nActual  [%#v]",
				candidate, expected, actual)
		}

	}
}

func TestMakeAnchor(t *testing.T) {
	var tests = []string{
		"httpunit",
		"httpunit",

		"Architecture overview:",
		"architecture-overview",

		"http unit \"by hand\":",
		"http-unit-by-hand",

		"http unit \"by file\":",
		"http-unit-by-file",

		"TOML file format:",
		"toml-file-format",

		"Basic test parameters",
		"basic-test-parameters",

		"IP variables",
		"ip-variables",

		"Testing with regular expressions",
		"testing-with-regular-expressions",

		"Additional command line flags",
		"additional-command-line-flags",

		"`-filter string`",
		"-filter-string",

		"`-header string`",
		"-header-string",

		"`-ipmap string`",
		"-ipmap-string",

		"`-no10`",
		"-no10",

		"`-timeout duration`",
		"-timeout-duration",

		// This is being bug compatible with Gitlab.
		// Some would say the result should be: "-v-and--vv"
		"`-v` and `-vv`",
		"-v-and-vv",

		"Common Tasks",
		"common-tasks",

		"Add a rule for a new haproxy port.",
		"add-a-rule-for-a-new-haproxy-port",

		"Oncall Tasks",
		"oncall-tasks",

		"Developer Info",
		"developer-info",

		"Package Creation",
		"package-creation",
	}

	doTestsMakeAnchor(t, tests)
}
