package cmd

import "encoding/json"

type Product struct {
	BarCode                string  `json:"bar_code"`
	Label                  string  `json:"label"`
	Tax                    string  `json:"tax"`
	Price                  float32 `json:"price"`
	Items                  float32 `json:"items"`
	SpecificTax            string  `json:"specific_tax"`
	SpecificTaxDesc        string  `json:"specific_tax_desc"`
	OriginalPrice          string  `json:"original_price"`
	PriceChangeExplanation string  `json:"price_change_explanation"`
}

type Payment struct {
	Mode   string  `json:"mode"`
	Amount float32 `json:"amount"`
}

type Bill struct {
	SellerName string    `json:"seller_name"`
	SellerId   string    `json:"seller_id"`
	BuyerIFU   string    `json:"buyer_ifu"`
	BuyerName  string    `json:"buyer_name"`
	AIB        string    `json:"aib"`
	Products   []Product `json:"products"`
	Payments   []Payment `json:"payments"`
	RT         string    `json:"rt"`
	RN         string    `json:"rn"`
	VT         string    `json:"vt"`
}

func newBillFromJson(jsonStr []byte) (*Bill, error) {
	bill := Bill{}
	err := json.Unmarshal(jsonStr, &bill)
	return &bill, err
}
