package renderer

import "context"

type contextKey string

const requestPathContextKey contextKey = "request_path"

func WithRequestPath(ctx context.Context, path string) context.Context {
	return context.WithValue(ctx, requestPathContextKey, path)
}

func RequestPathFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	path, ok := ctx.Value(requestPathContextKey).(string)
	if !ok {
		return ""
	}

	return path
}
