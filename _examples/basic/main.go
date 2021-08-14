package main

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"hugmouse/goyookassa/payment"
)

func main() {
	kassa := payment.NewKassa().SetShopID("YOUR_SHOP_ID").SetSecretKey("YOUR_API_SECRET_KEY")
	resp, err := payment.NewPayment().
		SetKassa(kassa).
		SetIdempotenceKey("RANDOM_STRING"). // Just use UUID v4
		SetAmount(decimal.NewFromInt(666), "RUB").
		SetCapture(true).
		SetConfirmation(
			payment.Confirmation{
				Type:      payment.Redirect,
				ReturnURL: "https://www.merchant-website.com/return_url",
			}).
		SetDescription("Default GoYooKassa test").Do()
	if err != nil {
		panic(err)
	}

	s, _ := json.MarshalIndent(resp, "", "\t")
	fmt.Printf("%v\n", string(s))

	// Example output:
	//	{
	//		"id": "0",
	//		"status": "pending",
	//		"paid": false,
	//		"amount": {
	//		"value": "666.00",
	//			"currency": "RUB"
	//	},
	//		"confirmation": {
	//		"type": "redirect",
	//		"confirmation_url": "https://yoomoney.ru/checkout/payments/v2/contract?orderId=0"
	//	},
	//		"created_at": "2021-08-13T14:13:46.45Z",
	//		"description": "Default test in GoYooKassa package",
	//		"recipient": {
	//		"account_id": "666999",
	//		"gateway_id": "1869069"
	//	},
	//		"refundable": false,
	//		"test": true
	//	}
}
