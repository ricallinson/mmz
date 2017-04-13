package main

import (
	. "github.com/ricallinson/simplebdd"
	"testing"
)

func TestUtils(t *testing.T) {

	Describe("utils", func() {

		It("should copy the whole buffer", func() {
			s := []byte{1, 2, 3, 4, 5}
			d := []byte{0, 0, 0, 0, 0}
			i := copyIntoArray(s, d)
			AssertEqual(len(s), 5)
			AssertEqual(len(d), 5)
			AssertEqual(i, 5)
		})

	})

	Report(t)
}
