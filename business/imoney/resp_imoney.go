package imoney

type RImoney struct {
	Id int `json:"id"`
	DisplayName string `json:"display_name"`
	Code string `json:"code"`
	ExchangeRate float64 `json:"exchange_rate"`
	EnableFraction bool `json:"enable_fraction"`
	IsPayable bool `json:"is_payable"`
	IsDebtable bool `json:"is_debtable"`
}