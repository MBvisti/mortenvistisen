package seeds

import (
	"context"
	"math/rand/v2"

	"github.com/google/uuid"
	"github.com/mbvisti/mortenvistisen/models"
)

func (s Seeder) PlantArticleTagConnection(
	ctx context.Context,
	articleID, tagID uuid.UUID,
) (models.ArticleTagConnection, error) {
	connection, err := models.NewArticleTagConnection(ctx, s.dbtx, articleID, tagID)
	if err != nil {
		return models.ArticleTagConnection{}, err
	}

	return connection, nil
}

func (s Seeder) PlantArticleTagConnections(
	ctx context.Context,
	articleID uuid.UUID,
	tagIDs []uuid.UUID,
) ([]models.ArticleTagConnection, error) {
	connections := make([]models.ArticleTagConnection, len(tagIDs))

	for i, tagID := range tagIDs {
		connection, err := s.PlantArticleTagConnection(ctx, articleID, tagID)
		if err != nil {
			return nil, err
		}
		connections[i] = connection
	}

	return connections, nil
}

// PlantRandomArticleTagConnections assigns random tags to an article
func (s Seeder) PlantRandomArticleTagConnections(
	ctx context.Context,
	articleID uuid.UUID,
	availableTags []models.ArticleTag,
	minTags, maxTags int,
) ([]models.ArticleTagConnection, error) {
	if len(availableTags) == 0 {
		return []models.ArticleTagConnection{}, nil
	}

	// Ensure min/max are within bounds
	if minTags < 0 {
		minTags = 0
	}
	if maxTags > len(availableTags) {
		maxTags = len(availableTags)
	}
	if minTags > maxTags {
		minTags = maxTags
	}

	// Determine number of tags to assign
	numTags := minTags
	if maxTags > minTags {
		numTags = minTags + rand.IntN(maxTags-minTags+1) //nolint:gosec // G404: Weak random for test data is acceptable
	}

	if numTags == 0 {
		return []models.ArticleTagConnection{}, nil
	}

	// Shuffle available tags and take the first numTags
	shuffledTags := make([]models.ArticleTag, len(availableTags))
	copy(shuffledTags, availableTags)
	rand.Shuffle(len(shuffledTags), func(i, j int) { //nolint:gosec // G404: Weak random for test data is acceptable
		shuffledTags[i], shuffledTags[j] = shuffledTags[j], shuffledTags[i]
	})

	selectedTagIDs := make([]uuid.UUID, numTags)
	for i := 0; i < numTags; i++ {
		selectedTagIDs[i] = shuffledTags[i].ID
	}

	return s.PlantArticleTagConnections(ctx, articleID, selectedTagIDs)
}

// PlantBulkArticleTagConnections assigns tags to multiple articles efficiently
func (s Seeder) PlantBulkArticleTagConnections(
	ctx context.Context,
	articles []models.Article,
	availableTags []models.ArticleTag,
	minTags, maxTags int,
) error {
	for _, article := range articles {
		_, err := s.PlantRandomArticleTagConnections(
			ctx,
			article.ID,
			availableTags,
			minTags,
			maxTags,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
