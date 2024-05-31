package usecases

import (
	"fmt"
	"os"
)

func FormatArticleSlug(slug string) string {
	return fmt.Sprintf("posts/%s", slug)
}

func BuildURLFromSlug(slug string) string {
	return fmt.Sprintf("%s://%s/%s", os.Getenv("APP_SCHEME"), os.Getenv("APP_HOST"), slug)
}
