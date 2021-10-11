package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Balance struct {
	Symbol    string `json:"symbol"`
	Available string `json:"available"`
	InOrder   string `json:"inOrder"`
}

// function to print json data nicely into console
func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

var RestURL = "https://api.bitvavo.com/v2"

func getKeys() string {
	APIKey := ""

	data, err := ioutil.ReadFile("APIKeys.key")
	if err != nil {
		fmt.Println("Failed to get Keys:", err)
	}
	APIKey = string(data)
	return APIKey
}

func createSignature(timestamp string, method string, endpoint string, body map[string]string, ApiSecret string) string {
	result := timestamp + method + "/v2" + endpoint
	if len(body) != 0 {
		bodyString, err := json.Marshal(body)
		if err != nil {
			fmt.Println(err)
		}
		result = result + string(bodyString)
	}
	// Create a new HMAC with hash type and key
	h := hmac.New(sha256.New, []byte(ApiSecret))
	// Write result
	h.Write([]byte(result))
	// Get result and encode as hexadecimal string
	sha := hex.EncodeToString(h.Sum(nil))

	return sha
}

func sendPrivate(endpoint string, method string, body map[string]string) []byte {
	// create timestamp in milliseconds and convert to string

	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	sig := createSignature(timestamp, method, endpoint, body, "04df5c0e6bdc5f3c2d20e6d379e8973d21945183b07fbb21f3438369a8aecc7881ab26c56b0bf9a3a3494a2271310e17d487234b463c277e757a0057dffb2b69")
	url := RestURL + endpoint

	client := &http.Client{}
	byteBody := []byte{}

	// check if body is empty or not
	if len(body) != 0 {
		bodyString, err := json.Marshal(body)
		if err != nil {
			fmt.Println(err)
		}
		// if body is not empty and gives no error, give byte slice of body to request
		byteBody = []byte(bodyString)
	} else {
		// if body is empty give empty byteBody to request
		byteBody = nil
	}

	//create new HTTP request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(byteBody))
	//Add request headers
	req.Header.Set("Bitvavo-Access-Key", getKeys())
	req.Header.Set("Bitvavo-Access-Signature", sig)
	req.Header.Set("Bitvavo-Access-Timestamp", timestamp)
	// req.Header.Set("Bitvavo-Access-Window", strconv.Itoa(10000))
	req.Header.Set("Content-Type", "application/json")

	// get HTTP response
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	// read HTTP response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return respBody
}

//
func getBalance() ([]Balance, error) {
	jsonResponse := sendPrivate("/balance", "GET", map[string]string{})
	t := make([]Balance, 0)
	err := json.Unmarshal(jsonResponse, &t)
	if err != nil {
		return []Balance{Balance{}}, nil
	}
	return t, nil
}

func main() {
	response, err := getBalance()
	if err != nil {
		fmt.Println(err)
	} else {
		for _, balance := range response {
			PrettyPrint(balance)
		}
	}
}
