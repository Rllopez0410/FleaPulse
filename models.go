package main

//=====Items structs=====

type SellToTrader struct {
	Trader 		string 	`json:"trader"`
	Currency  	string	`json:"currency"`
	PriceRub	int		`json:"priceRUB"`
	Price		int		`json:"price"`
}

type Item struct {
	ID   			string 			`json:"id"`
	Name 			string 			`json:"normalizedName"`
	Type 			[]string		`json:"types"`	
	SellToTrader	[]SellToTrader	`json:"sellToTraders"`
}

type Data struct {
	Items map[string]Item `json:"items"`
}

type ItemResponse struct {
	Data Data `json:"data"`
}

//=====Trader struct=======

type Trader struct {
	ID		string		`json:"id"`
	Name	string		`json:"normalizedName"`
}

type TraderResponse struct {
	Data map[string]Trader `json:"data"`
}

//=====Flea items struct=====

type ItemPrice struct {
	PriceMin	int		`json:"priceMin"`
	Price		int		`json:"price"`
	Timestamp	int64	`json:"timestamp"`
}

type PriceResponse struct {
	Data	[]ItemPrice	`json:"data"`
}
