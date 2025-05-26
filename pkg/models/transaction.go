package models

import "time"

// MoneyObject represents a monetary amount with currency
type MoneyObject struct {
	CurrencyCode    string `json:"currencyCode"`
	Value           string `json:"value"`
	ValueInBaseUnits int64  `json:"valueInBaseUnits"`
}

// Transaction represents an Upbank transaction
type Transaction struct {
	Type       string            `json:"type"`
	ID         string            `json:"id"`
	Attributes TransactionAttr   `json:"attributes"`
	Links      TransactionLinks  `json:"links"`
	Relations  TransactionRel    `json:"relationships"`
}

// TransactionAttr represents the attributes of a transaction
type TransactionAttr struct {
	Status            string            	`json:"status"`
	RawText           *string           	`json:"rawText"`
	Description       string            	`json:"description"`
	Message           string            	`json:"message"`
	IsCategorizable   bool             		`json:"isCategorizable"`
	HoldInfo          *HoldInfo         	`json:"holdInfo"`
	RoundUp           *RoundUp         		`json:"roundUp"`
	Cashback          *Cashback           	`json:"cashback"`
	Amount            MoneyObject           `json:"amount"`
	ForeignAmount     *MoneyObject          `json:"foreignAmount"`
	CardPurchaseMethod *CardPurchaseMethod  `json:"cardPurchaseMethod"`
	SettledAt         time.Time             `json:"settledAt"`
	CreatedAt         time.Time             `json:"createdAt"`
	TransactionType   *string               `json:"transactionType"`
	Note              *Note               	`json:"note"`
	PerformingCustomer PerformingCustomer 	`json:"performingCustomer"`
	DeepLinkURL       string              	`json:"deepLinkURL"`
}

// HoldInfo represents the hold information for a transaction
type HoldInfo struct {
	Amount         MoneyObject  `json:"amount"`
	ForeignAmount  *MoneyObject `json:"foreignAmount"`
}

// RoundUp represents the round up information for a transaction
type RoundUp struct {
	Amount        MoneyObject  `json:"amount"`
	BoostPortion  *MoneyObject `json:"boostPortion"`
}

// Cashback represents cashback information for a transaction
type Cashback struct {
	Description string      `json:"description"`
	Amount      MoneyObject `json:"amount"`
}

// CardPurchaseMethod represents the card purchase method information
type CardPurchaseMethod struct {
	Method           string  `json:"method"`
	CardNumberSuffix *string `json:"cardNumberSuffix"`
}

// Note represents a customer provided note about the transaction
type Note struct {
	Text string `json:"text"`
}

// PerformingCustomer represents the customer who performed the transaction
type PerformingCustomer struct {
	DisplayName string `json:"displayName"`
}

// TransactionRel represents the relationships for transactions
type TransactionRel struct {
	Account struct {
		Data struct {
			Type string `json:"type"`
			ID   string `json:"id"`
		} `json:"data"`
		Links struct {
			Related string `json:"related"`
		} `json:"links"`
	} `json:"account"`
	TransferAccount struct {
		Data *struct {
			Type string `json:"type"`
			ID   string `json:"id"`
		} `json:"data"`
	} `json:"transferAccount"`
	Category struct {
		Data  *struct {
			Type string `json:"type"`
			ID   string `json:"id"`
		} `json:"data"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
	} `json:"category"`
	ParentCategory struct {
		Data *struct {
			Type string `json:"type"`
			ID   string `json:"id"`
		} `json:"data"`
	} `json:"parentCategory"`
	Tags struct {
		Data []struct {
			Type string `json:"type"`
			ID   string `json:"id"`
		} `json:"data"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
	} `json:"tags"`
	Attachment struct {
		Data *struct {
			Type string `json:"type"`
			ID   string `json:"id"`
		} `json:"data"`
	} `json:"attachment"`
}

// TransactionLinks represents the links for transactions
type TransactionLinks struct {
	Self string `json:"self"`
}

// TransactionsResponse represents the API response for transactions
type TransactionsResponse struct {
	Data  []Transaction `json:"data"`
	Links struct {
		Prev *string `json:"prev"`
		Next *string `json:"next"`
	} `json:"links"`
}

// ByDate sorts transactions by date (newest first)
type ByDate []Transaction

func (t ByDate) Len() int           { return len(t) }
func (t ByDate) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByDate) Less(i, j int) bool { return t[i].Attributes.CreatedAt.After(t[j].Attributes.CreatedAt) } 