package business

type IUser interface {
	GetId() int
}

type ICorp interface {
	GetId() int
	GetPlatformId() int
	IsPlatform() bool
	IsValid() bool
}

type IOrder interface {
	GetId() int
	GetBid() string
	GetDeductableMoney() float64
}

type IBusiness interface {
	GetBid() string
	GetBusinessType() string
	GetActionType() string
}