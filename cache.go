package main

import (
	"time"
)

type itemCache struct {
	timeCached		time.Time
	cachedItems		[]Item
}

type traderCache struct {
	cachedTraders	[]Trader
}

var item_cache itemCache
var trader_cache traderCache

func cacheItem(id, name string) {
	item_cache.timeCached = time.Now()
	var newItem Item
	newItem.ID = id
	newItem.Name = name
	item_cache.cachedItems = append(item_cache.cachedItems, newItem)
}

func cachedTraders(id, name string) {
	var newTrader Trader
	newTrader.ID = id
	newTrader.Name = name
	trader_cache.cachedTraders = append(trader_cache.cachedTraders, newTrader)
}