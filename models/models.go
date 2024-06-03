package models

import "database/sql"

type listOptions struct {
	limit  sql.NullInt32
	offset sql.NullInt32
	// orderBy string
}

type listOpt func(*listOptions)

func WithOffset(val int32) listOpt {
	return func(lso *listOptions) {
		lso.offset = sql.NullInt32{Int32: val, Valid: true}
	}
}

func WithLimit(val int32) listOpt {
	return func(lso *listOptions) {
		lso.limit = sql.NullInt32{Int32: val, Valid: true}
	}
}

func WithPagination(limit, offset int32) listOpt {
	return func(lso *listOptions) {
		lso.offset = sql.NullInt32{Int32: offset, Valid: true}
		lso.limit = sql.NullInt32{Int32: limit, Valid: true}
	}
}
