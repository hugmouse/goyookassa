package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hugmouse/goyookassa/consts"
	"github.com/shopspring/decimal"
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

// Amount (payment amount)
//
// Sometimes YooMoney's partners charge additional commission from the users that is not included in this amount.
type Amount struct {
	// Value is how much money you want to get from someone
	Value decimal.Decimal `json:"value"`
	// Currency is three letter currency code (ex: RUB)
	Currency string `json:"currency,omitempty"`
}

// Confirmation information required to initiate the selected payment confirmation scenario by the user.
//
// More about confirmation scenarios: https://yookassa.ru/en/developers/payments/payment-process#user-confirmation
type Confirmation struct {
	// Type Confirmation scenario code
	Type string `json:"type"`
	// Enforce a request for making a payment with authentication by 3-D Secure.
	//
	// It works if you accept bank card payments without user confirmation by default.
	// In other cases, the 3-D Secure authentication will be handled by YooMoney.
	// If you would like to accept payments without additional confirmation by the user, contact your YooMoney manager.
	//
	// Works only with ConfirmationType == Redirect
	Enforce bool `json:"enforce"`
	// ReturnURL is the URL that the user will return to after confirming or canceling the payment on the webpage
	ReturnURL string `json:"return_url"`
}

type AmountFromResponse struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type ConfirmationFromResponse struct {
	Type            string `json:"type"`
	ConfirmationURL string `json:"confirmation_url"`
}

// ConfirmationType is one of the Confirmation scenarios
//
// More about confirmation scenarios: https://yookassa.ru/en/developers/payments/payment-process#confirmation-scenarios
type ConfirmationType int

const (
	// Embedded confirmation scenario: actions required for payment confirmation will depend on the payment method
	// selected by the user in the YooMoney Checkout Widget.
	// YooMoney will receive the confirmation from the user: all you need to do is embed the widget to your page.
	Embedded ConfirmationType = iota

	// External confirmation scenario:
	// to continue, the user takes action in an external system (for example, responds to a text message).
	// All you need to do is let them know how to proceed.
	External

	// The MobileApplication confirmation scenario: to confirm a payment, the user needs to complete an action
	// in a mobile app (for example, in an online banking app).
	// You need to redirect the user to the ConfirmationURL (ConfirmationFromResponse) received in the payment.
	// After the payment is made successfully (or if something goes wrong), YooMoney will redirect the user back
	// to the return_url that you send in your request for creating the payment.
	// This payment confirmation scenario only works on mobile devices (via mobile app or mobile web version).
	MobileApplication

	// QR confirmation scenario: to confirm the payment, the user scans a QR code.
	// You will need to generate the QR code using any available tools and display it on the payment page.
	QR

	// Redirect confirmation scenario: the user takes action on the YooMoney’s page or its partner’s page
	// (for example, enters bank card details or completes identification process via 3-D Secure).
	// You must redirect the user to ConfirmationURL (ConfirmationFromResponse) received in the payment.
	// If the payment is successful (or if something goes wrong),
	// YooMoney will return the user to return_url that you’ll send in the payment creation request.
	Redirect
)

func (c ConfirmationType) String() string {
	return [...]string{"embedded", "external", "mobile_application", "qr", "redirect"}[c]
}

type List struct {
	Type       string  `json:"type"`
	Items      []Items `json:"items"`
	NextCursor string  `json:"next_cursor"`
}
type Metadata struct {
}
type Card struct {
	First6        string `json:"first6"`
	Last4         string `json:"last4"`
	ExpiryMonth   string `json:"expiry_month"`
	ExpiryYear    string `json:"expiry_year"`
	CardType      string `json:"card_type"`
	IssuerCountry string `json:"issuer_country"`
	IssuerName    string `json:"issuer_name"`
}
type Method struct {
	Type  string `json:"type"`
	ID    string `json:"id"`
	Saved bool   `json:"saved"`
	Card  Card   `json:"card"`
	Title string `json:"title"`
}

type Items struct {
	ID            string    `json:"id"`
	Status        string    `json:"status"`
	Paid          bool      `json:"paid"`
	Amount        Amount    `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
	Description   string    `json:"description"`
	ExpiresAt     time.Time `json:"expires_at"`
	Metadata      Metadata  `json:"metadata"`
	PaymentMethod Method    `json:"payment_method"`
	Recipient     Recipient `json:"recipient"`
	Refundable    bool      `json:"refundable"`
	Test          bool      `json:"test"`
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

func (c *Kassa) ListPayments() *List {
	req, err := http.NewRequest("GET", consts.Endpoint+"payments", nil)
	if err != nil {
		return nil
	}
	req.SetBasicAuth(c.ShopID, c.SecretKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	listFromJSON := new(List)
	err = json.NewDecoder(resp.Body).Decode(&listFromJSON)
	if err != nil {
		return nil
	}

	return listFromJSON
}
