package bankly

const (
	// BanksPath ...
	BanksPath = "banklist"
)

// BankData ...
type BankData struct {
	ISPB string `json:"ispb,omitempty"`
	Name string `json:"name,omitempty"`
	Code string `json:"compe,omitempty"`
}

// BankDataResponse ...
type BankDataResponse struct {
	Name        string   `json:"name,omitempty"`
	ISPB        string   `json:"ispb,omitempty"`
	Code        string   `json:"code,omitempty"`
	ShortName   string   `json:"shortName,omitempty"`
	IsSPIDirect bool     `json:"isSPIDirect,omitempty"`
	Products    []string `json:"products,omitempty"`
}

// FilterBankListRequest ...
type FilterBankListRequest struct {
	IDs      []string
	Name     *string
	Product  *string
	Page     *int
	PageSize *int
}