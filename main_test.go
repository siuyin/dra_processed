package main

import (
	"regexp"
	"testing"
	"time"
)

func TestReformat(t *testing.T) {
	if reformat("A,B,C", ",", "|") != "A|B|C" {
		t.Error("reformat")
	}
}
func TestMatch(t *testing.T) {
	re := regexp.MustCompile("WARN|INFO|ERROR")
	dat := []struct {
		s string
		e bool
	}{
		{"INFO", true},
		{"InFO", false},
		{"ERROR", true},
		{"WARNING", true},
	}
	for _, v := range dat {
		if match(v.s, re) != v.e {
			t.Error("match", v.s)
		}
	}
}
func TestExtractDate(t *testing.T) {
	re := regexp.MustCompile(`(\d{8})-\d{6}:`)
	d, err := extractDate("INFO: 20161011-123456:", re)
	if err != nil || !d.Equal(time.Date(2016, 10, 11, 0, 0, 0, 0, time.Local)) {
		t.Error("extractDate", d, err)
	}
}
func TestDateInRange(t *testing.T) {
	s := time.Date(2016, 10, 11, 0, 0, 0, 0, time.Local)
	dat := []struct {
		d time.Time
		i int
		e bool
	}{
		{time.Date(2016, 10, 12, 0, 0, 0, 0, time.Local), 1, false},
		{time.Date(2016, 10, 11, 0, 0, 0, 0, time.Local), 1, false},
		{time.Date(2016, 10, 10, 0, 0, 0, 0, time.Local), 1, true},
		{time.Date(2016, 10, 9, 0, 0, 0, 0, time.Local), 1, false},
	}
	for _, v := range dat {
		if dateInRange(v.d, s, v.i) != v.e {
			t.Error("dateInRange", v.d, s, v.i)
		}
	}
}
func TestYesterday(t *testing.T) {
	now := time.Date(2016, 10, 11, 23, 59, 59, 999999, time.Local)
	yes := yesterday(now)
	y, m, d := yes.Date()
	if y != 2016 || m != 10 || d != 10 {
		t.Error("yesterday", y, m, d)
	}
}
