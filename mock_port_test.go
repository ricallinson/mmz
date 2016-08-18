package main

import (
	"bytes"
	. "github.com/ricallinson/simplebdd"
	"reflect"
	"testing"
)

func TestMockPort(t *testing.T) {

	Describe("MockPort()", func() {

		It("should return a MockPort object", func() {
			p := NewMockPort()
			AssertEqual(reflect.TypeOf(p).String(), "*main.MockPort")
		})

		It("should get the home screen", func() {
			p := NewMockPort()
			p.Write([]byte{27})
			b := make([]byte, 1000)
			p.Read(b)
			AssertEqual(bytes.Contains(b, []byte("d) Display settings")), true)
		})

		It("should get the settings screen", func() {
			p := NewMockPort()
			p.Write([]byte("d"))
			b := make([]byte, 1000)
			p.Read(b)
			AssertEqual(bytes.Contains(b, []byte("Display only, change with menu")), true)
		})

		It("should get the battery screen", func() {
			p := NewMockPort()
			p.Write([]byte("b"))
			b := make([]byte, 1000)
			p.Read(b)
			AssertEqual(bytes.Contains(b, []byte("a)BA, v)LBV, i)LBVI")), true)
		})

		It("should get the motor screen", func() {
			p := NewMockPort()
			p.Write([]byte("m"))
			b := make([]byte, 1000)
			p.Read(b)
			AssertEqual(bytes.Contains(b, []byte("Motor Settings")), true)
		})

		It("should get the rev screen", func() {
			p := NewMockPort()
			p.Write([]byte("s"))
			b := make([]byte, 1000)
			p.Read(b)
			AssertEqual(bytes.Contains(b, []byte("Rev limits")), true)
		})

		It("should get the options screen", func() {
			p := NewMockPort()
			p.Write([]byte("o"))
			b := make([]byte, 1000)
			p.Read(b)
			AssertEqual(bytes.Contains(b, []byte("Options: Enter letter to change")), true)
		})

		It("should get the special screen", func() {
			p := NewMockPort()
			p.Write([]byte("p"))
			b := make([]byte, 1000)
			p.Read(b)
			AssertEqual(bytes.Contains(b, []byte("Special Menu:")), true)
		})
	})

	Report(t)
}
