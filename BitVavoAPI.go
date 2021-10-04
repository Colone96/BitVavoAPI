package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func getKeys() {
	APIKey := ""

	data, err := ioutil.ReadFile("APIKeys.key")
	if err != nil {
		fmt.Println("Failed to get Keys:", err)
	}
	APIKey = string(data)
	fmt.Println(APIKey)
}

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

func main() {
	getKeys()
}
