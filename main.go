package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func match(s string, re *regexp.Regexp) bool {
	return re.MatchString(s)
}
func reformat(s, sep, join string) string {
	return strings.Join(strings.Split(s, sep), join)
}
func extractDate(s string, re *regexp.Regexp) (time.Time, error) {
	ms := re.FindStringSubmatch(s)
	if len(ms) < 2 {
		return time.Now(), errors.New("No matching date found")
	}
	tm, err := time.ParseInLocation("20060102", ms[1], time.Local)
	return tm, err
}
func dateInRange(d, s time.Time, i int) bool {
	dur := s.Sub(d)
	switch {
	case dur <= 0:
		return false
	case dur <= time.Hour*time.Duration(i*24):
		return true
	}
	return false
}
func yesterday(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local).Add(-24 * time.Hour)
}
func main() {
	// Get files to process.
	fmt.Println("Daily Processed Report")
	files, err := filepath.Glob(os.Getenv("FILE_GLOB") + "**")
	if err != nil {
		log.Fatal(err)
	}

	// Get number of files to limit processing to.
	fmt.Printf("Analysing: %q\n", files)
	limit, err := strconv.Atoi(os.Getenv("NUM_FILES"))
	if err != nil {
		log.Fatal("bad NUM_FILES:", err)
	}

	// Report Levels
	rptLvls := strings.Split(os.Getenv("REPORT_LEVELS"), ",")
	fmt.Printf("%q\n", rptLvls)

	// Text to scan.
	scnTxt := strings.Split(os.Getenv("SCAN_TEXT"), ",")
	fmt.Printf("%q\n", scnTxt)

loop:
	for i, v := range files {
		fmt.Printf("%v: %v\n", i, v)
		file, err := os.Open(v)
		if err != nil {
			log.Fatal("file open:", err)
		}

		// Create a new scanner.
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			for _, lvl := range rptLvls {
				for _, scn := range scnTxt {
					if strings.Contains(line, lvl) && strings.Contains(line, scn) {
						fmt.Println(scanner.Text())
					}
				}
			}
		}
		err = scanner.Err()
		if err != nil {
			log.Fatal("scanner:", err)
		}

		// Break if limit is read.
		if i >= limit-1 {
			break loop
		}
	}
}
