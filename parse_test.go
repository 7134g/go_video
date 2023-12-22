package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	bs, err := os.ReadFile("curl.txt")
	if err != nil {
		t.Fatal(err)
	}
	//reg, err := regexp.Compile("(?s)curl (.*?)compressed")
	reg, err := regexp.Compile("(?s)(.*?)compressed\n+")
	if err != nil {
		t.Fatal(err)
	}
	result := reg.FindAllString(string(bs), -1)
	fmt.Println(result)

	for _, value := range result {
		lines := strings.Split(value, "\n")
		url := lines[0][6 : len(lines[0])-4]
		fmt.Println(url)
	}
}
