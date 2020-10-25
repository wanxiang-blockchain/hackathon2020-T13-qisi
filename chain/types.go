package chain

type ReqOrder struct {
	EndAt    int64  `json:"end_at"`
	From     string `json:"from"`
	Funds    int64  `json:"funds"`
	Location string `json:"location"`
	StartAt  int64  `json:"start_at"`
	To       string `json:"to"`
}

type ReqRegisterRoom struct {
	Landlord	string `json:"landlord"`
	Property    string `json:"property"`
	Factory     string `json:"factory"`
	Location    string `json:"location"`
	Price int64 `json:"price"`
	Area  string `json:"area"`
	Status uint8 `json:"status"`
	Description string `json:"description"`
}

type ReqConfirm struct {
	Property string `json:"property"`
	OrderId int64 `json:"order_id"`
}