package domain

import "time"

// Click - структура для хранения данных о клике
type Click struct {
	BannerID  int       `json:"banner_id"`
	ClickTime time.Time `json:"click_time"`
}

// BannerClicksFilter структура для получения выборки нужных кликов
type BannerClicksFilter struct {
	BannerID int       `json:"banner_id"`
	TimeFrom time.Time `json:"time_from"`
	TimeTo   time.Time `json:"time_to"`
}

// ClickStatistics - структура для хранения данных о кол-ве кликов за определенный момент времени
type ClickStatistics struct {
	Count     int64     `json:"count"`
	ClickTime time.Time `json:"click_time"`
}
