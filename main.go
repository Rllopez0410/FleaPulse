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
		description: 	"Searches for specified item and returns pricing data",
		callback: 		fetchItemOnFlea,
	},
	"test": {
		name: 			"test",
		description: 	"Test command",
		callback: 		showCount,
	},
}

func cleanInput(text string) []string {
	wordArr := strings.Fields(strings.ToLower(text))
	return wordArr
}

func commandKill() error {
	fmt.Println("Closing the FleaScanner... Goodluck!")
	os.Exit(0)
	return nil
}

func priceSnapShot() {
	timer := time.NewTicker(2*time.Hour + 10*time.Minute)

	for {
		collectItemData()
		if len(trader_cache.cachedTraders) == 0 {
			collectTraderData()
		}
		collectFleaPrices()
		calculateProfit()
		//Create a function that will compute the differences of the FLEA and SELL price.
		<- timer.C
	}
}

func profitResults() error {
	if len(profitibleItems) == 0 {
		fmt.Printf("Traders cached: %d\n", len(trader_cache.cachedTraders))
		fmt.Println(len(item_cache.cachedItems))

		fmt.Println("There are no results. Try again later!")
		return nil
	}

	for i := 0; i < len(profitibleItems); i++ {
		fmt.Printf("Name: %s -> Profit: %d\n", profitibleItems[i].name, profitibleItems[i].profit)
	}

	return nil
}

func fetchItemOnFlea() error {
	fmt.Printf("Insert Item ID ->")
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

	for _, price := range price.Data[len(price.Data)-10:] {
		recordedTime := int64(price.Timestamp)
		t := time.UnixMilli(recordedTime).String()
		fmt.Printf("Listed Price Value: [%d]	Current Lowest Price: [%d]	Timestamp: %s\n", price.Price, price.PriceMin, t)
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
			fmt.Printf("ID: %s -> Name: %s -> Type %s\n", item.ID, item.Name, t)
		}
	}
	
	fmt.Printf("Items found: %d\n", itemCount)

	return nil
}

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Tarkov >")
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
	fmt.Println("Scanning flea...")
	go priceSnapShot()
	startRepl()
}