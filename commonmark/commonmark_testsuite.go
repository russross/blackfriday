package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"

	blackfriday "github.com/russross/blackfriday.v2"
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
		s.numTests++
		statsPerSection[test.Section] = s
	}
	return &stats{
		perSection: statsPerSection,
		numTests:   len(tests),
	}
}

func reportStats(s *stats) {
	coveredSections := &stats{
		perSection: map[string]*sectionStats{},
	}
	var coveredNames []string
	var uncoveredNames []string
	for section, stats := range s.perSection {
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
		fmt.Printf("%q: %d/%d tests passed\n", section, s.numTests, s.numPass)
	}
	sort.Strings(uncoveredNames)
	fmt.Println("\nRemaining sections:")
	for _, section := range uncoveredNames {
		stats := s.perSection[section]
		fmt.Printf("%q: %d/%d tests passed\n", section, stats.numPass, stats.numTests)
	}
	fmt.Printf("\nTotal %d/%d tests passed\n", s.numPass, s.numTests)
}

func main() {
	tests := load(testsDump)
	stats := prepStats(tests)
	for _, test := range tests {
		s := stats.perSection[test.Section]
		html := string(blackfriday.Run([]byte(test.Markdown)))
		if html == test.HTML {
			stats.numPass++
			s.numPass++
		}
	}
	reportStats(stats)
}
