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

// Get API-key from key-file
func getKeys() string {
	data, err := ioutil.ReadFile("APIKeys.key")
	if err != nil {
		fmt.Println("Failed to get Keys:", err)
	}
	var APIKey = string(data)
	return APIKey
}

// Create signature for sending with API request in the header
func createSignature(timestamp string, method string, endpoint string, body map[string]string, ApiSecret string) string {
	// create string to convert to signature
	result := timestamp + method + "/v2" + endpoint
	// check if body is not empty
	if len(body) != 0 {
		bodyString, err := json.Marshal(body)
		if err != nil {
			fmt.Println(err)
		}
		// if body is not empty add body to string
		result = result + string(bodyString)
	}
	// Create a new HMAC with hash type and secret key
	h := hmac.New(sha256.New, []byte(ApiSecret))
	// Write result
	h.Write([]byte(result))
	// Get result and encode as hexadecimal string
	sha := hex.EncodeToString(h.Sum(nil))

	return sha
}

//
func sendPrivate(endpoint string, method string, body map[string]string) []byte {
	// create timestamp in milliseconds and convert to string
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	//create signature
	sig := createSignature(timestamp, method, endpoint, body, "04df5c0e6bdc5f3c2d20e6d379e8973d21945183b07fbb21f3438369a8aecc7881ab26c56b0bf9a3a3494a2271310e17d487234b463c277e757a0057dffb2b69")
	// create url
	url := RestURL + endpoint

	// create new HTTP client
	client := &http.Client{}
	byteBody := []byte{}

	// check if body is empty or not
	if len(body) != 0 {
		bodyString, err := json.Marshal(body)
		if err != nil {
			fmt.Println(err)
		}
		// if body is not empty and gives no error, give byte slice of body to HTTP Request
		byteBody = []byte(bodyString)
	} else {
		// if body is empty give empty byteBody to HTTP Request
		byteBody = nil
	}

	//create new HTTP request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(byteBody))
	//Add request headers
	req.Header.Set("Bitvavo-Access-Key", getKeys())       //API-key
	req.Header.Set("Bitvavo-Access-Signature", sig)       //signature
	req.Header.Set("Bitvavo-Access-Timestamp", timestamp) //current timestamp in milliseconds
	// req.Header.Set("Bitvavo-Access-Window", strconv.Itoa(10000)) // Optional: Setting the update limit
	req.Header.Set("Content-Type", "application/json") //conrent type

	// send HTTP request and get HTTP response
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	// close response if an error is given
	defer resp.Body.Close()

	// read HTTP response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return respBody
}

// Get balance from the account
func getBalance() ([]Balance, error) {
	// get HTTP response from balance endpoint
	jsonResponse := sendPrivate("/balance", "GET", map[string]string{})
	// use stuct to put in the HTTP response
	t := make([]Balance, 0)
	// unmarshel json into the Balance struct with using a pointer
	err := json.Unmarshal(jsonResponse, &t)
	if err != nil {
		// if no error is given return the converted HTTP response and a nil error
		return []Balance{Balance{}}, nil
	}
	// To do: Print error when given (json data)
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
