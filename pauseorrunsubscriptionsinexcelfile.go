package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/invoiced/invoiced-go"
	"os"
	"strconv"
	"strings"
)

//This program will pause or un-pause subscriptions in the excel sheet

func main() {
	sandBoxEnv := true
	key := flag.String("key", "", "api key in Settings > Developer")
	environment := flag.String("env", "", "your environment production or sandbox")
	fileLocation := flag.String("file", "", "specify your excel file")

	flag.Parse()

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("This program will pause or un-pause subscriptions in the excel sheet")

	if *key == "" {
		fmt.Print("Please enter your API Key: ")
		*key, _ = reader.ReadString('\n')
		*key = strings.TrimSpace(*key)
	}

	*environment = strings.ToUpper(strings.TrimSpace(*environment))

	if *environment == "" {
		fmt.Println("Enter P for Production, S for Sandbox: ")
		*environment, _ = reader.ReadString('\n')
		*environment = strings.TrimSpace(*environment)
	}

	if *environment == "P" {
		sandBoxEnv = false
		fmt.Println("Using Production for the environment")
	} else if *environment == "S" {
		fmt.Println("Using Sandbox for the environment")
	} else {
		fmt.Println("Unrecognized value ", *environment, ", enter P or S only")
		return
	}

	if *fileLocation == "" {
		fmt.Println("Please specify your excel file: ")
		*fileLocation, _ = reader.ReadString('\n')
		*fileLocation = strings.TrimSpace(*fileLocation)
	}

	conn := invdapi.NewConnection(*key, sandBoxEnv)

	f, err := excelize.OpenFile(*fileLocation)

	if err != nil {
		panic(err)
	}

	fmt.Println("Read in excel file ", *fileLocation, ", successfully")

	subscriptionIdIndex := 0
	subscriptionPauseStatusIndex := 3

	rows, err := f.GetRows("Sheet1")

	if err != nil {
		panic("Error trying to get rows for the sheet" + err.Error())
	}

	if len(rows) == 0 {
		fmt.Println("No subscription data to update")
	}

	rows = rows[1:len(rows)]

	for i, row := range rows {

		subscriptionIdStr := row[subscriptionIdIndex]
		subscriptionPauseStatusStr := row[subscriptionPauseStatusIndex]

		subscriptionPauseStatus, err := strconv.ParseBool(subscriptionPauseStatusStr)

		if err != nil {
			fmt.Println("error at row ",i, ",",err)
			continue
		}

		subscriptionId, err := strconv.ParseInt(subscriptionIdStr,10,64)

		if err != nil {
			fmt.Println("error at row ",i, ",",err)
			continue
		}

		subscriptionToUpdate := conn.NewSubscription()
		subscriptionToUpdate.Id = subscriptionId
		subscriptionToUpdate.Paused = subscriptionPauseStatus
		err = subscriptionToUpdate.Save()
		if err != nil {
			fmt.Println("error at row ",i, ",",err)
			continue
		}

		if subscriptionPauseStatus {
			fmt.Println("Updated subscription with id = ",subscriptionId,", subscription is now paused")
		} else {
			fmt.Println("Updated subscription with id = ",subscriptionId,", subscription is now running")
		}

	}

}
