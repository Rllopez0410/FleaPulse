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

func cacheItem(item Item) {
	item_cache.timeCached = time.Now()
	item_cache.cachedItems = append(item_cache.cachedItems, item)
}

func cachedTraders(trader Trader) {
	trader_cache.cachedTraders = append(trader_cache.cachedTraders, trader)
}