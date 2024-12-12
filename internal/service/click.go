package service

import (
	"bkh-ecom/internal/domain"
	"bkh-ecom/internal/dto"
	"bkh-ecom/internal/logger"
	"bkh-ecom/internal/repository"
	"context"
	"fmt"
	"go.uber.org/atomic"
	"sync"
	"time"
)

// flushClicksBatchInterval интервал сброса батча в БД
const flushClicksBatchInterval = time.Minute

type clickService struct {
	dao   repository.DAO
	t     atomic.Time
	m     sync.Mutex
	batch []domain.Click
}

func NewClickService(ctx context.Context, dao repository.DAO) ClickService {
	resp := clickService{
		dao: dao,
		t:   atomic.Time{},
		m:   sync.Mutex{},
	}
	resp.t.Store(time.Now().Truncate(flushClicksBatchInterval))
	resp.start(ctx)
	return &resp
}

func (c *clickService) start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(flushClicksBatchInterval)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				isTimeToFlushBatch := false
				var batchForFlush []domain.Click
				c.m.Lock()
				if time.Since(c.t.Load()) > flushClicksBatchInterval {
					c.t.Store(time.Now().Truncate(flushClicksBatchInterval))
					isTimeToFlushBatch = true
					batchForFlush = make([]domain.Click, len(c.batch))
					copy(batchForFlush, c.batch)
					c.batch = nil
				}
				c.m.Unlock()

				if isTimeToFlushBatch {
					errs := c.dao.NewClicksQuery(ctx).InsertClicks(batchForFlush)
					if len(errs) > 0 {
						logger.ErrorKV(ctx, logger.Data{
							Msg:   "Failed to batch insert clicks",
							Error: errs[0],
						})
					}
				}
			}
		}
	}()
}

type ClickService interface {
	SaveClick(req domain.Click)
	ListClicks(ctx context.Context, filter dto.ClicksStatRequest) ([]domain.ClickStatistics, error)
}

func (c *clickService) SaveClick(req domain.Click) {
	req.ClickTime = req.ClickTime.Truncate(time.Minute)
	c.m.Lock()
	c.batch = append(c.batch, req)
	c.m.Unlock()
	return
}

func (c *clickService) ListClicks(ctx context.Context, filter dto.ClicksStatRequest) ([]domain.ClickStatistics, error) {

	repoFilter := domain.BannerClicksFilter{
		BannerID: filter.BannerID,
		TimeFrom: filter.TsFrom,
		TimeTo:   filter.TsTo,
	}

	clicksStats, err := c.dao.NewClicksQuery(ctx).ListClicks(repoFilter)
	if err != nil {
		logger.ErrorKV(ctx, logger.Data{
			Msg:    "Ошибка получения статистики кликов",
			Error:  err,
			Detail: fmt.Sprintf("%+v", repoFilter),
		})
		return nil, err
	}

	return clicksStats, nil
}
