package clearance

type RRecord struct {
	Id int `json:"id"`
	Bid string `json:"bid"`
	SourceUserId int `json:"source_user_id"`
	DestUserId int `json:"dest_user_id"`
	Role string `json:"role"`
	Amount int `json:"amount"`
	Ratio float64 `json:"ratio"`
	SettledAt string `json:"cleared_at"`
}