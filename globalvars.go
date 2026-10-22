package main

import "time"

var dataForDay = make([]itemHistory, 0, 5000)

var profitibleItems = make([]goodItem, 0, 500)

var priceMap = make(map[string][]itemHistory)

var item_cache itemCache

var trader_cache traderCache

var mode string = "pve"

var rate time.Duration = 50