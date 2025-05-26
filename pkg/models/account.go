package models

// Account represents an Upbank account
type Account struct {
	Type       string         `json:"type"`
	ID         string         `json:"id"`
	Attributes AccountAttr    `json:"attributes"`
	Links      AccountLinks   `json:"links"`
	Relations  AccountRelations `json:"relationships"`
}

// AccountAttr represents the attributes of an account
type AccountAttr struct {
	DisplayName    string      `json:"displayName"`
	AccountType    string      `json:"accountType"`
	OwnershipType  string      `json:"ownershipType"`
	Balance        Balance     `json:"balance"`
	CreatedAt      string      `json:"createdAt"`
}

// Balance represents the balance information
type Balance struct {
	CurrencyCode    string `json:"currencyCode"`
	Value           string `json:"value"`
	ValueInBaseUnits int64  `json:"valueInBaseUnits"`
}

// AccountLinks represents the links associated with an account
type AccountLinks struct {
	Self string `json:"self"`
}

// AccountRelations represents the relationships of an account
type AccountRelations struct {
	Transactions TransactionLinks `json:"transactions"`
}

// AccountsResponse represents the API response for accounts
type AccountsResponse struct {
	Data []Account `json:"data"`
}

// ByTypeAndName implements sort.Interface for []Account based on
// the AccountType and DisplayName fields.
type ByTypeAndName []Account

func (a ByTypeAndName) Len() int      { return len(a) }
func (a ByTypeAndName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByTypeAndName) Less(i, j int) bool {
	if a[i].Attributes.AccountType != a[j].Attributes.AccountType {
		return a[i].Attributes.AccountType < a[j].Attributes.AccountType
	}
	return a[i].Attributes.DisplayName < a[j].Attributes.DisplayName
} 