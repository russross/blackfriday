//
// Blackfriday Markdown Processor
// Available at http://github.com/russross/blackfriday
//
// Copyright © 2018 Russ Ross <russ@russross.com>.
// Distributed under the Simplified BSD License.
// See README.md for details.
//

// +build commonmark

package commonmark_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"testing"

	blackfriday "github.com/russross/blackfriday"
)

const testsDump = "testdata/commonmark-tests-dump.json"

type commonMarkTestCase struct {
	Section   string `json:"section"`
	Markdown  string `json:"markdown"`
	HTML      string `json:"html"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
	Example   int    `json:"example"`
}

type sectionStats struct {
	numTests int
	numPass  int
}

type stats struct {
	perSection map[string]*sectionStats
	numTests   int
	numPass    int
}

func load(testsuite string) []commonMarkTestCase {
	by, err := ioutil.ReadFile(testsDump)
	if err != nil {
		log.Fatal(err)
	}
	var tests []commonMarkTestCase
	err = json.Unmarshal(by, &tests)
	if err != nil {
		log.Fatal(err)
	}
	return tests
}

func prepStats(tests []commonMarkTestCase) *stats {
	statsPerSection := map[string]*sectionStats{}
	for _, test := range tests {
		s, ok := statsPerSection[test.Section]
		if !ok {
			s = &sectionStats{}
		}
		statsPerSection[test.Section] = s
	}
	return &stats{
		perSection: statsPerSection,
	}
}

func reportStats(s *stats) {
	coveredSections := &stats{
		perSection: map[string]*sectionStats{},
	}
	var coveredNames []string
	var uncoveredNames []string
	for section, stats := range s.perSection {
		if stats.numTests == 0 {
			continue
		}

		if stats.numTests == stats.numPass {
			coveredSections.perSection[section] = stats
			coveredNames = append(coveredNames, section)
			delete(s.perSection, section)
		} else {
			uncoveredNames = append(uncoveredNames, section)
		}
	}
	sort.Strings(coveredNames)
	fmt.Println("Completely covered sections:")
	for _, section := range coveredNames {
		s := coveredSections.perSection[section]
		fmt.Printf("  %q: %d/%d tests passed\n", section, s.numTests, s.numPass)
	}
	sort.Strings(uncoveredNames)
	fmt.Println("\nRemaining sections:")
	for _, section := range uncoveredNames {
		stats := s.perSection[section]
		fmt.Printf("  %q: %d/%d tests passed (%0.f%%)\n", section, stats.numPass, stats.numTests, float64(stats.numPass)/float64(stats.numTests)*100)
	}
	fmt.Printf("\nTotal %d/%d tests passed (%.0f%%)\n\n", s.numPass, s.numTests, float64(s.numPass)/float64(s.numTests)*100)
}

func TestCommonMark(t *testing.T) {
	tests := load(testsDump)
	stats := prepStats(tests)
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s %d", test.Section, test.Example), func(t *testing.T) {
			sect := stats.perSection[test.Section]
			sect.numTests++
			stats.numTests++

			html := string(blackfriday.Run([]byte(test.Markdown)))
			if html == test.HTML {
				stats.numPass++
				sect.numPass++
				return
			}

			t.Errorf("Input: %q\nExpect: %q\nActual: %q", test.Markdown, test.HTML, html)
		})
	}
	reportStats(stats)
}
