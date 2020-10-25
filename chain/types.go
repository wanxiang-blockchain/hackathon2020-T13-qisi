package chain

type ReqOrder struct {
	EndAt    int64  `fmt:"end_at"`
	From     string `fmt:"from"`
	Funds    int64  `fmt:"funds"`
	Location string `fmt:"location"`
	StartAt  int64  `fmt:"start_at"`
	To       string `fmt:"to"`
}
