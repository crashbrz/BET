![License](https://img.shields.io/badge/license-sushiware-red)
![Issues open](https://img.shields.io/github/issues/crashbrz/BET)
![GitHub pull requests](https://img.shields.io/github/issues-pr-raw/crashbrz/BET)
![GitHub closed issues](https://img.shields.io/github/issues-closed-raw/crashbrz/BET)
![GitHub last commit](https://img.shields.io/github/last-commit/crashbrz/BET)

# BET (Burp Enterprise Toolkit)
This repository intends to have a set of tools to take advantage of (not available or partially available on the web interface) features on the Burp Enterprise.

## BIS (Burp Importer & Scheduler) ##
- Since bulk schedules are not available on Burp Web Interface, this tool automatically imports and schedules all sites from an input(txt) file.
- File: bis.go

Usage example:
```
──(crash㉿Anubis)-[~]
└─$ go run bis.go -u https://burpserver.yourcompany.com:8080 -k BvujYxnHNNKPtNXfULfxhjXuyUjngCQn -i url_list.txt -r "FREQ=MONTHLY;INTERVAL=1" -s ab1c234d-56e7-8efa-9b0a-1b24c56de789 -t "2099-01-15T12:05:00+00:00"

 ```
- The above command will import all sites listed in url_list.txt and schedule each one to start the scan on January 15, 2099, at 12:05. Also, the scan will execute on the same date/time every month after the starting date.

## BSI (Burp Scan ID) ##
- Retrieves all scans names and ID's from Burp Enterprise.
- The Scan ID can be used as -s flag on Burp Importer & Scheduler(bis.go)
- Grepable output
- File: bis.go

Usage example:
```
──(crash㉿Anubis)-[~]
└─$ go run bis.go -u https://burpserver.yourcompany.com:8080 -k BvujYxnHNNKPtNXfULfxhjXuyUjngCQn | grep -i YourScanName | cut -d ":" -f 2
```

## GetFolderId ##
- Retrieves all folders names and ID's from Burp Enterprise.
- The Folder ID can be used as -i flag on BurpDeleteFolder.go
- Grepable output
- File: GetFolderId.go

Usage example:
```
──(crash㉿Anubis)-[~]
└─$ go run GetFolderId.go -u https://burpserver.yourcompany.com:8080 -k BvujYxnHNNKPtNXfULfxhjXuyUjngCQn | grep -i YourFolderName | cut -d ":" -f 2
```

### Usage/Help ###
Please refer to the output of -h and -v for usage information and general help. Also, you can contact me (@crashbrz) on Twitter<br>

### Installation ###
Clone the repository in the desired location.<br>
Remember to install GoLang.<br>

### License ###
BET and the tools are licensed under the SushiWare license. Check [docs/license.txt](docs/license.txt) for more information.
 
### Go Version ###
Tested on:<br>
Go version: go1.17.6 linux/amd64<br>
Kali 2021.4
Ubunto 20
