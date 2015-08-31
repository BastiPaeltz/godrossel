package main

import (
	"os"
	"log"
	"fmt"
	"github.com/BastiPaeltz/godrossel/utils"
)

// usage from cmd line:
//   godrossel SERVER_ADDR GOOGLE_SEARCH_API_KEY GOOGLE_SEARCH_CX_ID
// more info about cxid and api key:
// https://console.developers.google.com/project/ (api key)
// https://cse.google.com/all (cx id)
func main() {
	startLogger("godrossel.log")
	utils.StartWebserver(string(os.Args[1]))
}

// starts the logger.
// Creates new log file, if none exists.
// Else write to already existing one.
func startLogger(fileName string) {
	logfile, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil && os.IsNotExist(err) {
		logfile, err = os.Create(fileName)
		if err != nil {
			fmt.Println("Unable to generate log file.")
			return
		}
	}
	log.SetOutput(logfile)
}
