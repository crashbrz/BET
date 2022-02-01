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

type FolderID struct {
	Data struct {
		SiteTree struct {
			Folders []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"folders"`
		} `json:"site_tree"`
	} `json:"data"`
}

const df = "0.1"

func main() {
	fmt.Println("Burp Delete Folder")
	var burp_url string
	flag.StringVar(&burp_url, "u", "", "Your BurpSuite Enterprise URL. Ex. https://burpserver.yourcompany.com:8080")
	var folder_id string
	flag.StringVar(&folder_id, "i", "", "Your BurpSuite Enterprise URL. Ex. https://burpserver.yourcompany.com:8080")
	var api_key string
	flag.StringVar(&api_key, "k", "", "Your BurpSuite Enterprise API Key. Ex. BvujYxnHNNKPtNXfULfxhjXuyUjngCQn")
	version := flag.Bool("v", false, "Prints the current version and exit.")
	flag.Parse()

	if os.Args == nil || len(os.Args) < 3 {
		fmt.Println("You must provide at least: The BurpSuite Enterprise URL(-u), API Key(-k) and the Folder ID(-i)")
		fmt.Println("Version:", df)
		fmt.Println("Under the SushiWare license.")
		os.Exit(0)
	}

	if *version {
		fmt.Println("Burp Enterprise Delete Folder:", df)
		fmt.Println("Under the SushiWare license.")
		os.Exit(0)
	}

	getid(burp_url, api_key, folder_id) //call to getid
}

func getid(burp_url string, api_key string, folder_id string) {
	url := burp_url + "/graphql/v1" //static  variable
	method := "POST"                //static  variable
		payload := strings.NewReader("{\"variables\":{\"input\":{\"id\":" + folder_id + "}},\"query\":\"mutation DeleteFolder($input: DeleteFolderInput!) { delete_folder(input: $input) {id}}\"}")
	
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
	bodyString := string(body)
	var result FolderID
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Unmarshal error")
	}
	if strings.Contains(bodyString, folder_id) {
		fmt.Println("Folder " + folder_id + " deleted")
	} else {
		fmt.Println(bodyString)
	}
}
