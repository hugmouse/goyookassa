package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"hugmouse/goyookassa/consts"
	"io"
	"net/http"
	"time"
)

// Kassa struct is used to provide basic auth for YooKassa's endpoint
type Kassa struct {
	// ShopID is your, well, shop id
	ShopID string
	// SecretKey It is required for sending requests to API,
	// and it allows making any transaction in YooMoney using your name.
	//
	// Keep your key in a safe place: if you lose it, it will need to be reissued.
	SecretKey string
}

type Payment struct {
	*Kassa `json:"-"`

	// IdempotenceKey in the context of API, idempotence is the concept of multiple requests having the same effect
	// as a single request.
	//
	// Upon receiving a new request with identical parameters,
	// YooMoney will respond with results of the original request.
	//
	// Such behavior helps prevent unwanted repetition of transactions:
	// for example, if during the payment process the Internet connection was interrupted due to network problems,
	// you’ll be able to safely repeat the request for an unlimited number of times.
	IdempotenceKey string `json:"-"`

	Amount Amount `json:"amount"`

	// Capture that was set to true means you will receive the money immediately after the payment.
	// If the value is false, the required amount will be held on the user’s account,
	// and you’ll be able to capture it whenever convenient for you
	//
	// See more at: https://yookassa.ru/en/developers/payments/payment-process#capture-and-cancel
	Capture bool `json:"capture,omitempty"`

	// Confirmation information required to initiate the selected payment confirmation scenario by the user.
	//
	// More about confirmation scenarios: https://yookassa.ru/en/developers/payments/payment-process#user-confirmation
	Confirmation Confirmation `json:"confirmation"`

	// Description is used if you want to add a payment description that’ll be displayed in the Merchant Profile to you,
	// and during the payment to the user
	//
	// Also description must not exceed 128 characters
	Description string `json:"description,omitempty"`
}

// YooKassaResponse is default YooKassa endpoint response to payment creation request
type YooKassaResponse struct {
	ID           string                   `json:"id"`
	Status       string                   `json:"status"`
	Paid         bool                     `json:"paid"`
	Amount       AmountFromResponse       `json:"amount"`
	Confirmation ConfirmationFromResponse `json:"confirmation"`
	CreatedAt    time.Time                `json:"created_at"`
	Description  string                   `json:"description"`
	Recipient    Recipient                `json:"recipient"`
	Refundable   bool                     `json:"refundable"`
	Test         bool                     `json:"test"`
}

// YooKassaErrorResponse is used for handling error responses from YooKassa's endpoint
type YooKassaErrorResponse struct {
	Type        string `json:"type"`
	ID          string `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Parameter   string `json:"parameter"`
}

func (y *YooKassaErrorResponse) Error() string {
	return fmt.Sprintf("api returned error: %s (in parameter %s)", y.Description, y.Parameter)
}

// Recipient Payment.
//
// Required for separating payment flows within one account or making payments to other accounts.
type Recipient struct {
	AccountID string `json:"account_id"`
	GatewayID string `json:"gateway_id"`
}

// Amount
//
// Sometimes YooMoney's partners charge additional commission from the users that is not included in this amount.
type Amount struct {
	Value    decimal.Decimal `json:"value"`
	Currency string          `json:"currency,omitempty"`
}

type Confirmation struct {
	Type      string `json:"type,omitempty"`
	ReturnURL string `json:"return_url,omitempty"`
}

type AmountFromResponse struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type ConfirmationFromResponse struct {
	Type            string `json:"type"`
	ConfirmationURL string `json:"confirmation_url"`
}

// NewKassa creates and initializes a new Kassa (YooKassa shop id and shop secret key)
func NewKassa() *Kassa {
	return &Kassa{}
}

// SetShopID sets Kassa's shop id
func (c *Kassa) SetShopID(id string) *Kassa {
	c.ShopID = id
	return c
}

// SetSecretKey sets Kassa's secret key
//
// You can find it under API keys section: https://yookassa.ru/my/merchant/integration/api-keys
func (c *Kassa) SetSecretKey(key string) *Kassa {
	c.SecretKey = key
	return c
}

// NewPayment creates and initializes a new Payment
//
// Learn more: https://yookassa.ru/en/developers/api#create_payment
func NewPayment() *Payment {
	return &Payment{}
}

// SetKassa sets payment's YooKassa info (your shop id and shop secret key)
func (p *Payment) SetKassa(kassa *Kassa) *Payment {
	p.Kassa = kassa
	return p
}

// SetIdempotenceKey sets payment's idempotence key
func (p *Payment) SetIdempotenceKey(key string) *Payment {
	p.IdempotenceKey = key
	return p
}

// SetAmount sets payment's amount of money and money's type
//
// Example: payment.NewPayment().SetAmount(decimal.NewFromInt(500), "RUB")
func (p *Payment) SetAmount(value decimal.Decimal, moneyType string) *Payment {
	p.Amount = Amount{
		Value:    value,
		Currency: moneyType,
	}
	return p
}

// SetCapture sets payment's capture bool value
func (p *Payment) SetCapture(cap bool) *Payment {
	p.Capture = cap
	return p
}

// SetConfirmation sets payment's confirmation info
func (p *Payment) SetConfirmation(conf Confirmation) *Payment {
	p.Confirmation = conf
	return p
}

// SetDescription sets payment's description (128 character max)
func (p *Payment) SetDescription(desc string) *Payment {
	p.Description = desc
	return p
}

// Do sends an HTTP request to YooKassa payment endpoint
func (p *Payment) Do() (*YooKassaResponse, error) {
	payloadBytes, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", consts.Endpoint+"payments", body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(p.Kassa.ShopID, p.Kassa.SecretKey)
	req.Header.Set(consts.IdempotentHeader, p.IdempotenceKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	stuff, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	yooKassaError := &YooKassaErrorResponse{}
	err = json.Unmarshal(stuff, yooKassaError)
	if err != nil {
		return nil, err
	}

	if yooKassaError.Type == "error" {
		return nil, yooKassaError
	}

	respKassa := &YooKassaResponse{}
	err = json.Unmarshal(stuff, respKassa)
	if err != nil {
		return nil, err
	}

	return respKassa, nil
}
