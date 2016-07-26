package main

import (
    . "github.com/ricallinson/simplebdd"
    "testing"
    "fmt"
    "io/ioutil"
)

func loadFixtureToBuffer(filename string) []byte {
    dat, err := ioutil.ReadFile("./fixtures/" + filename + ".txt")
    if (err != nil) {
        fmt.Println(err)
        return []byte{}
    }
    return dat
}

func TestZilla(t *testing.T) {

    Describe("Zilla()", func() {

        It("should read settings", func() {
            z := CreateZilla()
            z.LastZillaOutput = loadFixtureToBuffer("settings")
            AssertEqual(len(z.LastZillaOutput), 348)
        })
    })

    Report(t)
}
