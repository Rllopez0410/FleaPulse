package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type priceRecord struct {
	minPrice	int
	time		string
}

type itemHistory struct {
	name		string
	captured	[]priceRecord
}

type goodItem struct {
	name		string
	profit		int
}

var requestMade int

func incrementCount() {
	requestMade += 1
}

func showCount() error {
	fmt.Printf("Request made: %d / %d\n", requestMade, len(item_cache.cachedItems))
	return nil
}

var dataForDay = make([]itemHistory, 0, 5000)
var profitibleItems = make([]goodItem, 0, 2000)
//create a function that'll sort most to least profitiable

var priceMap = make(map[string][]itemHistory)

func addPrice(min int, name, t string) {
	var newData priceRecord
	newData.minPrice = min
	newData.time = t

	 for i := 0; i < len(dataForDay); i++ {
        if dataForDay[i].name == name {
			dataForDay[i].captured = append(dataForDay[i].captured, newData)
            return
        }
    }

	var newHistory itemHistory
	newHistory.name = name
	newHistory.captured = append(newHistory.captured, newData)
	dataForDay = append(dataForDay, newHistory)
}

func addToMap() {
	day := time.Now().Weekday().String()
	priceMap[day] = append(priceMap[day], dataForDay...)
}

func collectItemData() error {
	res, err := http.Get("https://json.tarkov.dev/pve/items")
	if err != nil {
		fmt.Println("Error: Cannot connect to api")
		return err
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)

	if res.StatusCode > 299 {
		log.Fatalf("Status Code for items: %v", res.StatusCode)
	}
	if err != nil {
		log.Fatal(err)
	}

	var item ItemResponse
	if err := json.Unmarshal(data, &item); err != nil {
		return err
	}

	for _, item := range item.Data.Items {
	    skipItem := false

	    for _, itemType := range item.Type {
	        if itemType == "noFlea" ||
	            itemType == "gun" ||
	            itemType == "mods" ||
	            itemType == "pistolGrip" ||
	            itemType == "suppressor" ||
				itemType == "preset" || 
				itemType == "wearable" ||
				itemType == "ammoBox" {
					//possibly remove keys from the equation as well
	            skipItem = true
	            break
	        }
	    }

	    if skipItem {
	        continue
	    }
		
	    cacheItem(item)
	}

	return nil
}

func collectFleaPrices() error {
	timer := time.NewTicker(1*time.Second)

	for i := 0; i < len(item_cache.cachedItems); i++ {
		itemName := item_cache.cachedItems[i].Name
		itemId := item_cache.cachedItems[i].ID
		itemType := item_cache.cachedItems[i].Type

		skipItem := false

		for i := 0; i < len(itemType); i++ {
			if itemType[i] == "noFlea" || itemType[i] == "unknown" {
				skipItem = true
			}
		}

		if skipItem {
			continue
		} else {
			res, err := http.Get("https://json.tarkov.dev/pve/prices/" + itemId)
			if err != nil {
				fmt.Printf("Request failed for %s (%s): %v\n", itemName, itemId, err)
				
				time.Sleep(30 * time.Second)

				res, err = http.Get("https://json.tarkov.dev/pve/prices/" + itemId)
				if err != nil {
					return err
				}
			}

			data, err := io.ReadAll(res.Body)
			res.Body.Close()

			if res.StatusCode > 299 {
    			fmt.Printf("FAILED: %s | %s | status: %d\n", itemName, itemId, res.StatusCode)
    			continue
			} else {
				fmt.Printf("SUCCESS: %s | %s | status: %d\n", itemName, itemId, res.StatusCode)
			}

			if err != nil {
				log.Fatal(err)
			}

			var price PriceResponse
			if err := json.Unmarshal(data, &price); err != nil {
				return err
			}

			if len(price.Data) < 5 {
				for _, price := range price.Data {
					recordedTime := int64(price.Timestamp)
					t := time.UnixMilli(recordedTime).String()
					addPrice(price.PriceMin, itemName, t)
				}
			} else {
				for _, price := range price.Data[len(price.Data)-5:] {
					recordedTime := int64(price.Timestamp)
					t := time.UnixMilli(recordedTime).String()
					addPrice(price.PriceMin, itemName, t)
				}
			}

			addToMap()
			dataForDay = dataForDay[:0]
		}

		incrementCount()
		showCount()
		<- timer.C
	}

	fmt.Println("Data fetching complete...")

	return nil
}

func collectTraderData() error {
	res, err := http.Get("https://json.tarkov.dev/pve/traders")
	if err != nil {
		fmt.Println("Error: cannot connect to api")
		return err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)	
	if res.StatusCode > 299 {
		log.Fatalf("Status Code for traders %v", res.StatusCode)
	}
	if err != nil {
		log.Fatal(err)
	}
	
	var traders TraderResponse
	if err := json.Unmarshal(data, &traders); err != nil {
		return err
	}

	for _, trader := range traders.Data {
    	
    	cachedTraders(trader)
	}

	return nil
}

func calculateProfit() {
	fmt.Println("Days in priceMap:", len(priceMap))
    fmt.Println("Traders cached:", len(trader_cache.cachedTraders))
    fmt.Println("Items cached:", len(item_cache.cachedItems))

	for _, trader := range trader_cache.cachedTraders {
    	fmt.Printf("Trader: %s | ID: %s\n", trader.Name, trader.ID)
	}

    for day, value := range priceMap {
        fmt.Printf("%s has %d item histories\n", day, len(value))
	}

	for _, value := range priceMap {
		
		for i := 0; i < len(value); i++ {
			var avgProfit int
			var avgMinPrice int
			var sellPrice int
			var sumOfPrices int 
			var curItemName string

			priceCount := 0
			curItemName = value[i].name

			for _, priceStruct := range value[i].captured {
				sumOfPrices += priceStruct.minPrice
				priceCount++
			}

			for _, v := range item_cache.cachedItems {
				if v.Name == curItemName {
					for i := 0; i < len(v.SellToTrader); i++ {
						for _, j := range trader_cache.cachedTraders {
							if v.SellToTrader[i].Trader == j.ID && j.Name == "therapist" {
    							sellPrice = v.SellToTrader[i].PriceRub
							}
						}
					}
				}
			}

			fmt.Printf(
    		"%s | prices: %d | sellPrice: %d\n", curItemName, priceCount, sellPrice)

			if priceCount == 0 || sellPrice == 0 {
    			continue
			}
			
			avgMinPrice = sumOfPrices / priceCount
			avgProfit = sellPrice - avgMinPrice

			if avgProfit > 1000 {
				var newGoodItem goodItem
				newGoodItem.name = curItemName
				newGoodItem.profit = avgProfit
				profitibleItems = append(profitibleItems, newGoodItem)
			} else {
				continue
			}
		}
	} 
}