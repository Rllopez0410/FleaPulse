package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type cliCommand struct {
	name		string
	description	string
	callback	func() error
}

var commandRegistry = map[string]cliCommand {
	"kill": {
		name:        	"kill",
		description: 	"Exit the scanner",
		callback:    	commandKill,
	},

	"items": {
		name: 			"items",
		description:	"Searches for items",
		callback: 		fetchItemsInfo,
	},

	"flea": {
		name: 			"flea",
		description: 	"Searching for specified item on the flea.",
		callback: 		fetchItemOnFlea,
	},

	"profits": {
		name: 			"profits",
		description: 	"Retrieves a list of potential profitable items.",
		callback: 		profitResults,
	},

	"pulse": {
    	name:        "pulse",
    	description: "Scannning the flea market. This may take a few minutes...",
    	callback:    priceSnapShot,
},

	"man": {
		name:		"help",
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
			 "profits"
 
			 Description:
			 "Retrieves a list of potentially profitable items."
 
			 -------------------------------------------------------------------------
 
			 Name:
			 "pulse"
 
			 Description:
			 "Scans the flea market for pricing data, calculates the average minimum 
			  price of each item, and compares it to the item's trader sell price to 
			  identify potential profit opportunities."
 
			 -------------------------------------------------------------------------
 
			 Name:
			 "man"
 
			 Description:
			 "Displays a manual of general commands for navigating and using the CLI."
		`,

		callback: commandMan,
	},
}

//create a function that will clean the names of the items and remove special characters
//create a gamemode variable so that users can change between them

func cleanInput(text string) []string {
	wordArr := strings.Fields(strings.ToLower(text))
	return wordArr
}

func commandKill() error {
	fmt.Println("Closing the FleaScanner... Goodluck!")
	os.Exit(0)
	return nil
}

func commandMan() error {
	return nil
}

func priceSnapShot() error {
	if len(item_cache.cachedItems) == 0 {
		collectItemData()
	}

	if len(trader_cache.cachedTraders) == 0 {
		collectTraderData()
	}

	collectFleaPrices()
	calculateProfit()

	return nil
}

func profitResults() error {
	timer := time.NewTicker(50*time.Millisecond)
	defer timer.Stop()

	if len(profitibleItems) == 0 {
		fmt.Println(len(item_cache.cachedItems))
		fmt.Println("There are no results. Try again later!")
		return nil
	}

	for i := 0; i < len(profitibleItems); i++ {
		fmt.Printf("Name: %s | Profit: %d\n", profitibleItems[i].name, profitibleItems[i].profit)
		<- timer.C
	}

	return nil
}

func fetchItemOnFlea() error {
	fmt.Printf("Insert Item ID =>")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	id := scanner.Text()
	scanner.Err()

	res, err := http.Get("https://json.tarkov.dev/pve/prices/" + id)
	if err != nil {
		fmt.Println("Error: cannot connect to api")
		return err
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)

	if res.StatusCode > 299 {
		log.Fatalf("Status Code %v", res.StatusCode)
	}

	if err != nil {
		log.Fatal(err)
	}

	var price PriceResponse
	if err := json.Unmarshal(data, &price); err != nil {
		return err
	}

	for _, price := range price.Data[len(price.Data)-5:] {
		recordedTime := int64(price.Timestamp)
		t := time.UnixMilli(recordedTime).String()
		fmt.Printf("Listed Price Value: [%d]	|	Current Lowest Price: [%d]	|	Timestamp: %s\n", price.Price, price.PriceMin, t)
	}

	return nil
}

func fetchItemsInfo() error {
	res, err := http.Get("https://json.tarkov.dev/pve/items")
	if err != nil {
		fmt.Println("Error: Cannot connect to api")
		return err
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)

	if res.StatusCode > 299 {
		log.Fatalf("Status Code: %v", res.StatusCode)
	}

	if err != nil {
		log.Fatal(err)
	}

	var item ItemResponse
	if err := json.Unmarshal(data, &item); err != nil {
		return err
	}

	itemCount := 0

	for _, item := range item.Data.Items {
		var t string
		if len(item.Type) == 0 || item.Type[0] == "noFlea" {
			continue
		} else {
			t = item.Type[0]
			itemCount += 1
			fmt.Printf("ID: %s | Name: %s | Type %s\n", item.ID, item.Name, t)
		}
	}
	
	fmt.Printf("Items found: %d\n", itemCount)

	return nil
}

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Tarkov => ")
		scanner.Scan()
		text := scanner.Text()
		scanner.Err()
		textSlice := cleanInput(text)

		matchKey := false
		for key := range commandRegistry {
			if textSlice[0] == key {
				fmt.Println(commandRegistry[key].description)
				commandRegistry[key].callback()
				matchKey = true
				break
			}
		}
		
		if !matchKey {
			fmt.Println("Unknown command")
		}
	}
}

func main() {
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
	
	fmt.Println("			 \033[1mWelcome To FleaPulse\033[0m...\n")
	fmt.Println("			 Type [ \033[1mman\033[0m ] to see a list of general commands.\n")
	
	startRepl()
}