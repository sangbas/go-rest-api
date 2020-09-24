package entity

type MovieReq struct {
	Name     string `json:"name"`
	Duration int    `json:"duration"`
	Genre    string `json:"genre"`
}

type MovieResp struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Duration int    `json:"duration"`
	Genre    string `json:"genre"`
}
