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

var dataForDay = make([]itemHistory, 0, 6000)

var priceMap = make(map[string][]itemHistory)

func addPrice(min int, name, t string) {
	var newData priceRecord
	var newHistory itemHistory
	day := time.Now().Weekday().String()
	newHistory.name = name
	newData.minPrice = min
	newData.time = t
	newHistory.captured = append(newHistory.captured, newData)
	dataForDay = append(dataForDay, newHistory)
	priceMap[day] = dataForDay
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
		log.Fatalf("Status Code: %v", res.StatusCode)
	}
	if err != nil {
		log.Fatal(err)
	}

	var item ItemResponse
	if err := json.Unmarshal(data, &item); err != nil {
		return err
	}

	for _, item := range item.Data.Items {
		if len(item.Type) == 0 || item.Type[0] == "noFlea" {
			continue
		} else {
			cacheItem(item.ID, item.Name)
		}
	}

	return nil
}

func collectFleaPrices() error {
	timer := time.NewTicker(2*time.Second)

	for i := 0; i < len(item_cache.cachedItems); i++ {
		<- timer.C
		itemName := item_cache.cachedItems[i].Name
		itemId := item_cache.cachedItems[i].ID
		itemType := item_cache.cachedItems[i].Type

		skipItem := false

		for i := 0; i < len(itemType); i++ {
			if itemType[i] == "noFlea" || itemType[i] == "unknown" {
				skipItem = true
			}
		}

		if skipItem == true {
			res, err := http.Get("https://json.tarkov.dev/pve/prices/" + itemId)
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

			for _, price := range price.Data[len(price.Data)-20:] {
				recordedTime := int64(price.Timestamp)
				t := time.UnixMilli(recordedTime).String()
				addPrice(price.PriceMin, itemName, t)
			}

			dataForDay = dataForDay[:0]

			fmt.Printf("%v\n", priceMap)
		} else {
			continue
		}
	}

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
			log.Fatalf("Status Code %v", res.StatusCode)
		}
		if err != nil {
			log.Fatal(err)
		}
		
		var traders TraderResponse
		if err := json.Unmarshal(data, &traders); err != nil {
			return err
		}

		for _, trader := range traders.Data {
			cachedTraders(trader.ID, trader.Name)
		}

	return nil
}