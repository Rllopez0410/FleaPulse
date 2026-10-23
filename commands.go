package main

import (
	"fmt"
	"time"
)

type cliCommand struct {
	name		string
	description	string
	callback	func() error
}

var commandRegistry = map[string]cliCommand {
	//create a command that will change the sample count for min prices
	"kill": {
		name:        	"kill",
		description: 	"Exit the scanner.",
		callback:    	commandKill,
	},

	"mode": {
		name: 			"mode",
		description:	"Checking mode...",
		callback: 		commandMode,
	},

	"rate": {
		name:			"rate",
		description: 	"Changing Get request rate...",
		callback: 		commandRate,
	},

	"items": {
		name: 			"items",
		description:	"Searching for items...",
		callback: 		fetchItemsInfo,
	},

	"flea": {
		name: 			"flea",
		description: 	"Searching for specified item on the flea...",
		callback: 		fetchItemOnFlea,
	},

	"profits": {
		name: 			"profits",
		description: 	"Retrieving a list of potential profitable items...",
		callback: 		displayProfit,
	},

	"pulse": {
    	name:        	"pulse",
    	description: 	"Scannning the flea market. This may take a few minutes...",
    	callback:    	priceSnapShot,
	},

	"man": {
		name:			"man",
		description: 	
		`
			 Welcome to Tarkov FleaPulse. This program is designed to find potentially 
			 profitable items on the flea market that can be flipped by selling them to 
			 traders. The program analyzes recorded flea market prices and compares 
			 them against trader sell prices to identify potential profit opportunities. 
			 Results are based on historical pricing data and do not guarantee current 
			 flea market prices or profits. FleaPulse serves as a guide to help you 
			 identify which items may be worth paying attention to or searching for on 
			 the flea market.
 
			 General Commands:
 
			 Name:
			 "kill"
 
			 Description: 	
			 "Terminate the program."
 
			 -------------------------------------------------------------------------

			 Name:
			 "mode"
 
			 Description: 	
			 "Displays the current game mode and allows you to change it."
 
			 -------------------------------------------------------------------------

			 Name:
			 "rate"

			 Description:
			 "Displays and allows you to change the delay between flea market API requests."

			 -------------------------------------------------------------------------
 
			 Name:
			 "items"
 
			 Description:
			 "Searches for essential information about each item, such as its ID, name, 
			  and type."
 
			 -------------------------------------------------------------------------
 
			 Name:
			 "flea"
 
			 Description:
			 "Searches for a specified item and returns its pricing history."
 
			 -------------------------------------------------------------------------
  
			 Name:
			 "pulse"
			 
			 Description:
			 "Scans the flea market for pricing data, calculates the average minimum 
			 price of each item, and compares it to the item's trader sell price to 
			 identify potential profit opportunities."
			 
			 -------------------------------------------------------------------------
			 
			 Name:
			 "profits"
 
			 Description:
			 "Retrieves a list of potentially profitable items."
 
			 -------------------------------------------------------------------------
			
			 Name:
			 "man"
 
			 Description:
			 "Displays a manual of general commands for navigating and using the CLI."
		`,

		callback: 	commandMan,
	},
}

func welcomeMsg() {
	timer := time.NewTicker(150*time.Millisecond)
	defer timer.Stop()

	fmt.Println(`		
			 ███████╗██╗     ███████╗ █████╗ ██████╗ ██╗   ██╗██╗     ███████╗███████╗`)
	<-timer.C
	fmt.Println(`			 ██╔════╝██║     ██╔════╝██╔══██╗██╔══██╗██║   ██║██║     ██╔════╝██╔════╝`)
	<-timer.C	 
	fmt.Println(`			 █████╗  ██║     █████╗  ███████║██████╔╝██║   ██║██║     ███████╗█████╗`)
	<-timer.C	 
	fmt.Println(`			 ██╔══╝  ██║     ██╔══╝  ██╔══██║██╔═══╝ ██║   ██║██║     ╚════██║██╔══╝`)
	<-timer.C	 
	fmt.Println(`			 ██║     ███████╗███████╗██║  ██║██║     ╚██████╔╝███████╗███████║███████╗`)
	<-timer.C	 
	fmt.Println(`			 ╚═╝     ╚══════╝╚══════╝╚═╝  ╚═╝╚═╝      ╚═════╝ ╚══════╝╚══════╝╚══════╝
	`)
	
	fmt.Println("			 Current version: \033[1mv1.0.0\033[0m\n")
	fmt.Println("			 \033[1mWelcome To FleaPulse\033[0m...\n")
	fmt.Println("			 Type [ \033[1mman\033[0m ] to see a list of general commands.\n")
}