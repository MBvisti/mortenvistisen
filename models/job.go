package models

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"mortenvistisen/internal/storage"

	"github.com/riverqueue/river/rivertype"
	"github.com/uptrace/bun"
)

type JobSummary struct {
	bun.BaseModel `bun:"table:river_job,alias:river_job"`
	ID            int64        `bun:"id,pk"`
	State         string       `bun:"state"`
	Attempt       int16        `bun:"attempt"`
	MaxAttempts   int16        `bun:"max_attempts"`
	AttemptedAt   sql.NullTime `bun:"attempted_at"`
	CreatedAt     time.Time    `bun:"created_at"`
	FinalizedAt   sql.NullTime `bun:"finalized_at"`
	ScheduledAt   time.Time    `bun:"scheduled_at"`
	Priority      int16        `bun:"priority"`
	Kind          string       `bun:"kind"`
	Queue         string       `bun:"queue"`
}

type JobFilter struct {
	State  string
	Search string
}

type PaginatedJobs struct {
	Jobs       []JobSummary
	TotalCount int64
	Page       int64
	PageSize   int64
	TotalPages int64
}

type JobStats struct {
	Total   int64
	ByState map[string]int64
}

func applyJobFilters(query *bun.SelectQuery, filter JobFilter) *bun.SelectQuery {
	if filter.State != "" {
		query = query.Where("state = ?", filter.State)
	}

	if search := strings.TrimSpace(filter.Search); search != "" {
		pattern := "%" + search + "%"
		query = query.Where(
			"(kind ILIKE ? OR queue ILIKE ? OR CAST(id AS TEXT) = ?)",
			pattern,
			pattern,
			search,
		)
	}

	return query
}

func (j job) Paginate(
	ctx context.Context,
	db storage.Executor,
	filter JobFilter,
	page int64,
	pageSize int64,
) (PaginatedJobs, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 25
	}
	if pageSize > 100 {
		pageSize = 100
	}

	totalCount, err := applyJobFilters(
		db.NewSelect().Model((*JobSummary)(nil)),
		filter,
	).Count(ctx)
	if err != nil {
		return PaginatedJobs{}, err
	}

	jobs := make([]JobSummary, 0, pageSize)
	if err := applyJobFilters(
		db.NewSelect().Model(&jobs),
		filter,
	).
		OrderExpr("created_at DESC, id DESC").
		Limit(int(pageSize)).
		Offset(int((page - 1) * pageSize)).
		Scan(ctx); err != nil {
		return PaginatedJobs{}, err
	}

	return PaginatedJobs{
		Jobs:       jobs,
		TotalCount: int64(totalCount),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: (int64(totalCount) + pageSize - 1) / pageSize,
	}, nil
}

func (j job) Stats(ctx context.Context, db storage.Executor) (JobStats, error) {
	stateCounts := make([]struct {
		State string `bun:"state"`
		Count int64  `bun:"count"`
	}, 0, len(rivertype.JobStates()))

	if err := db.NewSelect().
		TableExpr("river_job").
		ColumnExpr("state::text AS state").
		ColumnExpr("COUNT(*) AS count").
		GroupExpr("state").
		Scan(ctx, &stateCounts); err != nil {
		return JobStats{}, err
	}

	stats := JobStats{ByState: make(map[string]int64, len(rivertype.JobStates()))}
	for _, state := range rivertype.JobStates() {
		stats.ByState[string(state)] = 0
	}
	for _, stateCount := range stateCounts {
		stats.ByState[stateCount.State] = stateCount.Count
		stats.Total += stateCount.Count
	}

	return stats, nil
}
