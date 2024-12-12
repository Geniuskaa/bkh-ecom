package dto

import "time"

type ClicksStatRequest struct {
	BannerID int       `json:"banner_id"`
	TsFrom   time.Time `json:"tsFrom"`
	TsTo     time.Time `json:"tsTo"`
}

type ClicksStatResponse struct {
	Ts          time.Time `json:"ts"`
	ClicksCount int       `json:"clicksCount"`
}
