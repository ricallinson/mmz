package main

import (
	"bytes"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strings"
)

func copyIntoArray(s []byte, d []byte) int {
	for i, _ := range d {
		if i >= len(s) || i >= len(d) {
			return i
		}
		d[i] = s[i]
	}
	return len(d)
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

// Reads a file into a byte array or exits.
func readFileToByteArray(p string) []byte {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		log.Fatalf("Error reading file: #%v ", err)
	}
	return b
}

func readYamlFileToExecutorCommands(p string) *ExecutorCommands {
	r := &ExecutorCommands{}
	if p != "" {
		err := yaml.Unmarshal(readFileToByteArray(p), r)
		if err != nil {
			log.Fatalf("YAML Parse Error: %v", err)
		}
	}
	return r
}

func interfaceToJson(i interface{}) []byte {
	d, err := json.Marshal(i)
	if err != nil {
		log.Fatalf("YAML Write Error: %v", err)
	}
	return d
}

func interfaceToYaml(i interface{}) []byte {
	d, err := yaml.Marshal(i)
	if err != nil {
		log.Fatalf("YAML Write Error: %v", err)
	}
	return d
}
