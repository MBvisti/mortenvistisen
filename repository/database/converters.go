package database

import (
	"encoding/json"

	"github.com/MBvisti/mortenvistisen/domain"
)

func (n *Newsletter) ConvertNewsletterToDomain(
	associatedArticleSlug string,
) (domain.Newsletter, error) {
	var paragraphs []string
	if err := json.Unmarshal(n.Body, &paragraphs); err != nil {
		return domain.Newsletter{}, err
	}

	return domain.Newsletter{
		ID:          n.ID,
		Title:       n.Title,
		Edition:     n.Edition.Int32,
		ReleasedAt:  n.ReleasedAt.Time,
		Released:    n.Released.Bool,
		Paragraphs:  paragraphs,
		ArticleSlug: associatedArticleSlug,
	}, nil
}
