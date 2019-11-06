package user

type RUser struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Sex    string `json:"sex"`
	Code   string `json:"code"`
}
