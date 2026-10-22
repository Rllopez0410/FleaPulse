package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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
	profit		int
	sellPrice	int
	name		string
}

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
	res, err := http.Get("https://json.tarkov.dev/" + mode + "/items")
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
	            skipItem = true
	            break
	        }
	    }

	    if skipItem {
	        continue
	    }

		item.Name = strings.ReplaceAll(item.Name, "-", " ")

	    cacheItem(item)
	}

	return nil
}

func progressBar(current, total int) {
	const barWidth = 50

	percent := float64(current) / float64(total)
	filled := int(percent * float64(barWidth))

	fmt.Printf("\r[")

	for i := 0; i < barWidth; i++ {
		if i < filled {
			fmt.Print("#")
		} else {
			fmt.Print(" ")
		}
	}

	fmt.Printf("] %d/%d (%.0f%%) ", current, total, percent*100)
}

func collectFleaPrices() error {
	timer := time.NewTicker(rate * time.Millisecond)
	defer timer.Stop()

	var requestMade int

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
			res, err := http.Get("https://json.tarkov.dev/" + mode + "/prices/" + itemId)
			if err != nil {
				time.Sleep(30 * time.Second)

				res, err = http.Get("https://json.tarkov.dev/" + mode + "/prices/" + itemId)
				if err != nil {
					return err
				}
			}

			data, err := io.ReadAll(res.Body)
			res.Body.Close()

			//but this data in the progress bar function
			/*if res.StatusCode > 299 {
    			fmt.Printf("FAILED: %s | %s | status: %d\n", itemName, itemId, res.StatusCode)
				failCount++
    			continue
			} else {
				fmt.Printf("SUCCESS: %s | %s | status: %d\n", itemName, itemId, res.StatusCode)
			}*/

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

		requestMade++
		progressBar(requestMade, len(item_cache.cachedItems))
		<- timer.C
	}

	requestMade = 0
	<-timer.C
	fmt.Println("Data fetching complete...")
	<-timer.C
	fmt.Println("Averaging prices...")
	<-timer.C
	return nil
}

func collectTraderData() error {
	res, err := http.Get("https://json.tarkov.dev/" + mode + "/traders")
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
	timer := time.NewTicker(3*time.Second)
	defer timer.Stop()

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

			if priceCount == 0 || sellPrice == 0 {
    			continue
			}
			
			avgMinPrice = sumOfPrices / priceCount
			avgProfit = sellPrice - avgMinPrice

			if avgProfit > 1000 {

				var newGoodItem goodItem
				newGoodItem.name = curItemName
				newGoodItem.profit = avgProfit
				newGoodItem.sellPrice = sellPrice
				profitibleItems = append(profitibleItems, newGoodItem)
			} else {
				continue
			}
		}
	} 

	<-timer.C
	fmt.Println("Calculation completed!")
}