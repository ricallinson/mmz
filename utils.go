package main

import (
	"bytes"
	"strings"
)

func copyIntoArray(s []byte, d []byte) {
	for i, _ := range d {
		if i >= len(s) || i >= len(d) {
			return
		}
		d[i] = s[i]
	}
}

func findValueLine(b []byte, end string) []byte {
	start := bytes.Index(b, []byte(end))
	b = b[start:]
	b = b[bytes.Index(b, []byte("\n"))+1:]
	return b[:bytes.Index(b, []byte("\n"))]
}

// Return a boolean of "On" or "Off".
func truthy(s string) bool {
	return strings.Contains(s, "On")
}

// Returns a string array of values with white space removed.
func split(s string, sep string) []string {
	values := []string{}
	tokens := strings.Split(s, sep)
	for _, token := range tokens {
		token = strings.TrimSpace(token)
		if len(token) > 0 {
			values = append(values, token)
		}
	}
	return values
}

// Convert a boolean to an integer of 0 or 1.
func Btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}
