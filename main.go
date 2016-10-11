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
	ytd := yesterday(time.Now())

	// Report Levels
	rptLvls := regexp.MustCompile(reformat(os.Getenv("REPORT_LEVELS"), ",", "|"))

	// Text to scan.
	scnTxt := regexp.MustCompile(reformat(os.Getenv("SCAN_TEXT"), ",", "|"))
	dateRe := regexp.MustCompile(os.Getenv("DATE_REGEXP"))
	numDays, err := strconv.Atoi(os.Getenv("NUM_DAYS"))
	if err != nil {
		log.Fatal("bad NUM_DAYS:", err)
	}

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
			if rptLvls.MatchString(line) && scnTxt.MatchString(line) {
				dt, _ := extractDate(line, dateRe)
				if dateInRange(dt, ytd, numDays) {
					fmt.Println(line)
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
