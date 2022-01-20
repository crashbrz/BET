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
	"bufio"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Response struct {
	Data struct {
		CreateSite struct {
			Site struct {
				ID       string `json:"id"`
				ParentID string `json:"parent_id"`
				Scope    struct {
					IncludedUrls []string `json:"included_urls"`
				} `json:"scope"`
				ApplicationLogins struct {
					LoginCredentials []interface{} `json:"login_credentials"`
					RecordedLogins   []interface{} `json:"recorded_logins"`
				} `json:"application_logins"`
			} `json:"site"`
		} `json:"create_site"`
	} `json:"data"`
}

const bemi = "0.1"

func main() {

	var burp_url string
	flag.StringVar(&burp_url, "u", "", "Your BurpSuite Enterprise URL. Ex. https://burpserver.yourcompany.com:8080")
	var api_key string
	flag.StringVar(&api_key, "k", "", "Your BurpSuite Enterprise API Key. Ex. BvujYxnHNNKPtNXfULfxhjXuyUjngCQn")
	var input_file string
	flag.StringVar(&input_file, "i", "", "A file containing the list of URLs to be imported, one by line.")
	var rrule string
	flag.StringVar(&rrule, "r", "", "RRULE to schedule the scan. Ex. \"FREQ=MONTHLY;INTERVAL=1\"") // FREQ=MONTHLY;BYDAY=TU;INTERVAL=1 -> BREAKS the SCAN
	var scan_id string
	flag.StringVar(&scan_id, "s", "", "BurpSuite Enterprise Scan ID. Ex. ab1c234d-56e7-8efa-9b0a-1b24c56de789")
	var time string
	flag.StringVar(&time, "t", "", "Date and time when the scan must start. Ex. \"2099-01-18T12:05:00+00:00\"")
	version := flag.Bool("v", false, "Prints the current version and exit.")
	flag.Parse()

	if os.Args == nil || len(os.Args) < 3 {
		fmt.Println("You must provide at least: The BurpSuite Enterprise URL(-u), the API Key(-k), and file contaning the URL list(-i)")
		fmt.Println("If the schedule scan configurations are not provided, the Importer will create the sites without any scan configuration/schedule")
		fmt.Println("Version:", bemi)
		fmt.Println("Under the SushiWare license.")
		//version()
		os.Exit(0)

	}

	if *version {
		fmt.Println("Burp Enterprise Mass Import Version:", bemi)
		fmt.Println("Under the SushiWare license.")
		os.Exit(0)
	}

	add(burp_url, api_key, input_file, rrule, scan_id, time) //call to add sites.
}

func add(burp_url string, api_key string, input_file string, rrule string, scan_id string, time string) {
	file, err := os.Open(input_file) // passar o arquivo pra parametro
	if err != nil {
		log.Fatalf("Failed to open the URL input list file")
		os.Exit(0)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string

	for scanner.Scan() {
		text = append(text, scanner.Text())

	}
	file.Close()

	url := burp_url + "/graphql/v1" //static  variable
	method := "POST"                //static  variable

	// Loop to read the url input file
	for _, each_ln := range text {
		fmt.Println("Site to add: " + each_ln) // Prints the current site to be added
		payload_aux := ("{\"query\": \"mutation CreateSite {\\n create_site(input: {name: \\\"")
		payload_aux1 := each_ln
		payload_aux2 := ("\\\", parent_id: 0, application_logins: {login_credentials: [], recorded_logins: []} , scope: {included_urls: \\\"")
		payload_aux3 := ("\\\"},}) {\\n   site {\\nid\\nparent_id\\nscope {\\nincluded_urls\\n} application_logins {login_credentials {id\\nlabel\\nusername\\n}\\nrecorded_logins {\\nid\\nlabel\\n}\\n}\\n}\\n }\\n}\"}")
		payload := strings.NewReader(payload_aux + payload_aux1 + payload_aux2 + payload_aux1 + payload_aux3)
		//		fmt.Println(payload) // remove the comment for debugging
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

		//		fmt.Println(string(body)) // remove the comment for debugging

		//Json Parsing
		response1 := Response{}
		jsonErr := json.Unmarshal(body, &response1)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		//		fmt.Println(response1.Data.CreateSite.Site.ID) // Get the Scan ID from Json - remove the comment for debugging
		site_id := response1.Data.CreateSite.Site.ID
		schedule(site_id, api_key, burp_url, rrule, scan_id, time) // call the scan schedule function

	}
}

func schedule(site_id string, api_key string, burp_url string, rrule string, scan_id string, time string) {
	/*	Debugging
		    fmt.Println("Schedule function called") //debug
			fmt.Println("API KEY")//debug
			fmt.Println(api_key)//debug
			fmt.Println("burp url")//debug
			fmt.Println(burp_url)//debug
	*/
	url := burp_url + "/graphql/v1" //static  variable
	method := "POST"                //static  variable
	// fmt.Println("Fullurl")//debug
	// fmt.Println(url)//debug
	payload_sch := "{{\"operationName\":\"CreateScheduleItem\",\"variables\":{\"input\":{\"site_id\":\""
	payload_sch_aux1 := site_id
	payload_sch_aux2 := "\",\"schedule\":{\"initial_run_time\":\"" + time + "\",\"rrule\":\"" + rrule + "\"},\"scan_configuration_ids\":[\"" + scan_id + "\"]}},\"query\":\"mutation CreateScheduleItem($input: CreateScheduleItemInput!) {create_schedule_item(input: $input) {schedule_item {id}}}\"}"
	payload_sch_full := payload_sch + payload_sch_aux1 + payload_sch_aux2
	sch_advjust_var := payload_sch_full[1:]
	final_sch_payload := strings.NewReader(sch_advjust_var)

	/*	Debugging
		fmt.Println("Payload: Schedule") //debug
		fmt.Println(payload_sch_full) //debug
	*/
	tr1 := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client1 := &http.Client{Transport: tr1}
	req1, err1 := http.NewRequest(method, url, final_sch_payload)

	if err1 != nil {
		fmt.Println(err1)
	}
	req1.Header.Add("Accept", "*/*")
	req1.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36")
	req1.Header.Add("Connection", "close")
	req1.Header.Add("Accept-Encoding", "gzip, deflate")
	req1.Header.Add("Authorization", api_key)
	req1.Header.Add("Upgrade-Insecure-Requests", "1")
	req1.Header.Add("Accept-Language", "en-US,en;q=0.9")
	req1.Header.Add("Content-Type", "application/json")

	res21, err1 := client1.Do(req1)
	res1 := *res21
	defer res1.Body.Close()
	/*	Debugging
			//body1, err1 := ioutil.ReadAll(res1.Body)

		fmt.Println("Print schedule JSON result") //debug
		fmt.Println(string(body1)) //debug
	*/
}

