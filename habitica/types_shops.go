package habitica

// ShopItem represents an entry in the shop.
type ShopItem struct {
	Key         string  `json:"key"`
	Text        string  `json:"text"`
	Notes       string  `json:"notes"`
	Value       float64 `json:"value"`
	PurchaseType string `json:"purchaseType"`
	Class       string  `json:"class"`
}

