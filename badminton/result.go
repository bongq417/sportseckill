package badminton

import "time"

type SaveResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type AddResult struct {
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
	Pay      int    `json:"pay"`
	OrderNum string `json:"OrderNum"`
}

type BuyInfo struct {
	Username  string        `json:"username"`
	Mobile    string        `json:"mobile"`
	IdCard    string        `json:"id_card"`
	ShowId    string        `json:"showId"`
	HallTime  string        `json:"hallTime"`
	OffsetDay time.Duration `json:"offsetDay"`
}
