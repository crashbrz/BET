/*
 * --------------------------------------------------------------------------------
 * "SUSHIWARE LICENSE" (Revision 01):
 * As long as you retain this notice, you can do whatever you want with this code.
 * If we meet someday around the universe, and you think this code was useful,
 * if you want, you can pay me a sushi round in return.
 * Ewerson Guimaraes a.k.a Crash
 * --------------------------------------------------------------------------------
 */
package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type ScanID struct {
	Data struct {
		ScanConfigurations []struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"scan_configurations"`
	}
}

const bgsi = "0.1"

func main() {

	var burp_url string
	flag.StringVar(&burp_url, "u", "", "Your BurpSuite Enterprise URL. Ex. https://burpserver.yourcompany.com:8080")
	var api_key string
	flag.StringVar(&api_key, "k", "", "Your BurpSuite Enterprise API Key. Ex. BvujYxnHNNKPtNXfULfxhjXuyUjngCQn")
	version := flag.Bool("v", false, "Prints the current version and exit.")
	flag.Parse()

	if os.Args == nil || len(os.Args) < 2 {
		fmt.Println("You must provide at least: The BurpSuite Enterprise URL(-u), the API Key(-k)")
		fmt.Println("Version:", bgsi)
		fmt.Println("Under the SushiWare license.")
		//version()
		os.Exit(0)
	}

	if *version {
		fmt.Println("Burp Enterprise Scan ID Version:", bgsi)
		fmt.Println("Under the SushiWare license.")
		os.Exit(0)
	}

	getid(burp_url, api_key) //call to getid
}

func getid(burp_url string, api_key string) {
	url := burp_url + "/graphql/v1" //static  variable
	method := "POST"                //static  variable
	payload := strings.NewReader("{\"query\": \"query GetScanConfigurations{scan_configurations {name id}}\"}")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Accept", "*/*")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36")
	req.Header.Add("Connection", "close")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Authorization", api_key)
	req.Header.Add("Accept-Language", "en-US,en;q=0.9")
	req.Header.Add("Content-Type", "application/json")

	res2, err := client.Do(req)
	res := *res2
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	var result ScanID
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Unmarshal error")
	}
	for _, info := range result.Data.ScanConfigurations {
		fmt.Println(info.Name + ":" + info.ID)
	}
}
