package fcrypt

import (
	"fmt"
	"testing"
)

var tests = []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

var testString = []string{"3737xxxx", "e93393xx", "37373783x", "ue83838x", "ASDFAX&993"}

func TestCreateToken(t *testing.T) {

	fmt.Println("Test now TestCreateToken")
	for i, v := range tests {

		val := CreateToken()
		if val == "" {

			t.Fatalf("at index %d, expected %d, go val %s", i, v, val)

		} else {

			fmt.Println("CreateToken: ", val)
		}
	}
}

func TestMd5(t *testing.T) {

	fmt.Println("")
	fmt.Println("Test now TestMd5")
	for i, v := range testString {

		val := Md5(v)

		if val == "" {

			t.Fatalf("at index %d, expected %d, go val %s", i, v, val)

		} else {

			fmt.Println("Md5: ", val)
		}
	}
}
