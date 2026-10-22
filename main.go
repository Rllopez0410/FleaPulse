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

func checkMode() error {
	timer := time.NewTicker(50*time.Millisecond)
	defer timer.Stop()

	for {
		fmt.Printf("Current game mode: %s\n", mode)
		<-timer.C
		fmt.Printf("Would like to change the mode: [y/n]")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		choice := scanner.Text()
		scanner.Err()
	
		switch choice {
			case "y":
				<-timer.C
				fmt.Print("Choose mode [pve/regular/pvp-season]: ")
				scanner.Scan()
				newMode := scanner.Text()

				if newMode != "pve" && newMode != "regular" && newMode != "pvp-season" {
					fmt.Printf("Invalid mode... Please try again.\n")
				} else {
					<-timer.C
					fmt.Printf("Mode switched to \033[1m%s\033[0m...\n", newMode)
					return nil
				}
			case "n":
				<-timer.C
				fmt.Printf("Mode unchanged...\n")
				return nil
			default:
				fmt.Printf("Invalid choice... Please try again.\n")	
		}
	}
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
	timer := time.NewTicker(15*time.Millisecond)
	defer timer.Stop()

	fmt.Printf("Insert Item ID => ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	id := scanner.Text()
	scanner.Err()

	res, err := http.Get("https://json.tarkov.dev/" + mode + "/prices/" + id)
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
		<-timer.C
	}

	return nil
}

func fetchItemsInfo() error {
	timer := time.NewTicker(2*time.Millisecond)
	defer timer.Stop()
	res, err := http.Get("https://json.tarkov.dev/" + mode + "/items")
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
			item.Name = strings.ReplaceAll(item.Name, "-", " ")
			fmt.Printf("ID: %s | Name: %s | Type %s\n", item.ID, item.Name, t)
		}

		<-timer.C
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
	welcomeMsg()
	startRepl()
}