package transfer

type RTransfer struct {
	Id int	`json:"id"`
	Bid string	`json:"bid"`
	ThirdBid string	`json:"third_bid"`
	SourceAccountId int	`json:"source_account_id"`
	DestAccountId int	`json:"dest_account_id"`
	SourceUserId int	`json:"source_user_id"`
	DestUserId int	`json:"dest_user_id"`
	Amount int	`json:"amount"`
	SourceAmount int	`json:"source_amount"`
	DestAmount int	`json:"dest_amount"`
	Action string	`json:"action"`
	OriginAction string	`json:"origin_action"`
	TransferType string	`json:"transfer_type"`
	CreatedAt  string	`json:"created_at"`
}