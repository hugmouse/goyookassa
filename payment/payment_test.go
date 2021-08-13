package payment

import (
	"github.com/shopspring/decimal"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestNewKassa(t *testing.T) {
	tests := []struct {
		name string
		want *Kassa
	}{
		{name: "Default", want: &Kassa{}},
		{name: "Default with custom shop id", want: (&Kassa{}).
			SetShopID("1")},
		{name: "Default with custom secret key", want: (&Kassa{}).
			SetSecretKey("key")},
		{name: "Default with custom secret key and shop id", want: (&Kassa{}).
			SetShopID("1").
			SetSecretKey("key")},
	}
	for num, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch num {
			case 0:
				if got := NewKassa(); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewKassa() = %v, want %v", got, tt.want)
				}
			case 1:
				got := NewKassa()
				got.ShopID = "1"
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewKassa().SetShopID(\"1\") = %v, want %v", got, tt.want)
				}
			case 2:
				got := NewKassa()
				got.SecretKey = "key"
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewKassa().SetSecretKey(\"key\") = %v, want %v", got, tt.want)
				}
			case 3:
				got := NewKassa()
				got.ShopID = "1"
				got.SecretKey = "key"
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewKassa().SetShopID(\"1\").SetSecretKey(\"key\") = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestNewPayment(t *testing.T) {
	ShopSecretKey := os.Getenv("SHOP_SECRET_KEY")
	if ShopSecretKey == "" {
		t.Errorf("SHOP_SECRET_KEY environment variable was not set. SHOP_SECRET_KEY = %v", ShopSecretKey)
	}

	ShopID := os.Getenv("SHOP_ID")
	if ShopSecretKey == "" {
		t.Errorf("SHOP_ID environment variable was not set. SHOP_SECRET_KEY = %v", ShopID)
	}

	// Pseudo-random for testing purposes
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789")
	length := 16
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}

	Kassa := NewKassa().SetSecretKey(ShopSecretKey).SetShopID(ShopID)

	tests := []struct {
		name string
		want *Payment
	}{
		// TODO: Add test cases.
		{name: "Default", want: NewPayment().
			SetKassa(Kassa).
			SetIdempotenceKey(b.String()).
			SetAmount(decimal.NewFromInt(666), "RUB").
			SetCapture(true).SetConfirmation(
			Confirmation{
				Type:      "redirect",
				ReturnURL: "https://www.merchant-website.com/return_url",
			}).
			SetDescription("Default GoYooKassa test")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPayment(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPayment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPayment_Do(t *testing.T) {
	ShopSecretKey := os.Getenv("SHOP_SECRET_KEY")
	if ShopSecretKey == "" {
		t.Errorf("SHOP_SECRET_KEY environment variable was not set. SHOP_SECRET_KEY = %v", ShopSecretKey)
	}

	ShopID := os.Getenv("SHOP_ID")
	if ShopSecretKey == "" {
		t.Errorf("SHOP_ID environment variable was not set. SHOP_SECRET_KEY = %v", ShopID)
	}

	// Pseudo-random for testing purposes
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789")
	length := 16
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}

	OurKassa := NewKassa().SetSecretKey(ShopSecretKey).SetShopID(ShopID)
	type fields struct {
		Kassa          *Kassa
		IdempotenceKey string
		Amount         Amount
		Capture        bool
		Confirmation   Confirmation
		Description    string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *YooKassaResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "Default", fields: struct {
			Kassa          *Kassa
			IdempotenceKey string
			Amount         Amount
			Capture        bool
			Confirmation   Confirmation
			Description    string
		}{
			Kassa:          OurKassa,
			IdempotenceKey: b.String(),
			Amount: Amount{
				Value:    decimal.NewFromInt(666),
				Currency: "RUB",
			}, Capture: true, Confirmation: Confirmation{
				Type:      "redirect",
				ReturnURL: "https://www.merchant-website.com/return_url",
			}, Description: "Default test in GoYooKassa package"},
			want: nil, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Payment{
				Kassa:          tt.fields.Kassa,
				IdempotenceKey: tt.fields.IdempotenceKey,
				Amount:         tt.fields.Amount,
				Capture:        tt.fields.Capture,
				Confirmation:   tt.fields.Confirmation,
				Description:    tt.fields.Description,
			}
			_, err := p.Do()
			if (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
