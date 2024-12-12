package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type dao struct {
	RunnerWrapper
}

type DAO interface {
	NewClicksQuery(ctx context.Context) ClicksQuery
}

func NewDAO(db *pgxpool.Pool) DAO {
	return &dao{db}
}

func (d *dao) NewClicksQuery(ctx context.Context) ClicksQuery {
	return &clicksQuery{BaseQuery{
		ctx:    ctx,
		runner: d.RunnerWrapper,
	}}
}
