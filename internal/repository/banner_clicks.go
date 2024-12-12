package repository

import (
	"bkh-ecom/internal/domain"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

const bannerClicksTableName = "banner_clicks"

type clicksQuery struct {
	BaseQuery
}

type ClicksQuery interface {
	InsertClicks(clicks []domain.Click) []error
	ListClicks(filter domain.BannerClicksFilter) ([]domain.ClickStatistics, error)
}

func (q *clicksQuery) InsertClicks(clicks []domain.Click) []error {
	batch := &pgx.Batch{}

	for _, click := range clicks {
		query := fmt.Sprintf("INSERT INTO %v (banner_id, click_time) VALUES ($1, $2)", bannerClicksTableName)
		batch.Queue(query, click.BannerID, click.ClickTime)
	}

	results := q.runner.SendBatch(q.Context(), batch)
	defer results.Close()

	var errs []error
	for i := 0; i < len(clicks); i++ {
		_, err := results.Exec()
		if err != nil {
			err = errors.Wrap(err, "Ошибка отправки запроса")
			errs = append(errs, err)
		}
	}

	return errs
}

func (q *clicksQuery) ListClicks(filter domain.BannerClicksFilter) ([]domain.ClickStatistics, error) {
	query := fmt.Sprintf("SELECT click_time, count(banner_id) as count FROM %v WHERE banner_id = $1 AND click_time BETWEEN $2 AND $3 GROUP BY click_time;", bannerClicksTableName)
	rows, err := q.runner.Query(q.Context(), query, filter.BannerID, filter.TimeFrom, filter.TimeTo)
	if err != nil {
		return nil, errors.Wrap(err, "Ошибка получения статистики по банеру")
	}
	defer rows.Close()

	resp := make([]domain.ClickStatistics, 0)
	for rows.Next() {
		item := domain.ClickStatistics{}
		err = rows.Scan(&item.ClickTime, &item.Count)
		if err != nil {
			return nil, errors.Wrap(err, "Ошибка преобразования данных в структуру")
		}
		resp = append(resp, item)
	}

	return resp, nil
}
