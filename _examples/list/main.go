package main

import (
	"encoding/json"
	"fmt"
	"hugmouse/goyookassa/payment"
)

func main() {
	kassa := payment.NewKassa().SetShopID("YOUR_SHOP_ID").SetSecretKey("YOUR_API_SECRET_KEY")
	s, _ := json.MarshalIndent(kassa.ListPayments(), "", "\t")
	fmt.Printf("%v\n", string(s))
}
