package models

type PaginationOptions struct {
	Limit  int32
	Offset int32
}

type PaginationOption func(*PaginationOptions)

func WithOffset(val int32) PaginationOption {
	return func(pgo *PaginationOptions) {
		pgo.Offset = val
	}
}

func WithLimit(val int32) PaginationOption {
	return func(pgo *PaginationOptions) {
		pgo.Limit = val
	}
}

func WithPagination(limit, offset int32) PaginationOption {
	return func(pgo *PaginationOptions) {
		pgo.Offset = offset
		pgo.Limit = limit
	}
}
