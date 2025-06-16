package seeds

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/google/uuid"
	"github.com/mbvisti/mortenvistisen/models"
)

type articleTagSeedData struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string
}

type articleTagSeedOption func(*articleTagSeedData)

func WithArticleTagID(id uuid.UUID) articleTagSeedOption {
	return func(atsd *articleTagSeedData) {
		atsd.ID = id
	}
}

func WithArticleTagCreatedAt(createdAt time.Time) articleTagSeedOption {
	return func(atsd *articleTagSeedData) {
		atsd.CreatedAt = createdAt
	}
}

func WithArticleTagUpdatedAt(updatedAt time.Time) articleTagSeedOption {
	return func(atsd *articleTagSeedData) {
		atsd.UpdatedAt = updatedAt
	}
}

func WithArticleTagTitle(title string) articleTagSeedOption {
	return func(atsd *articleTagSeedData) {
		atsd.Title = title
	}
}

func (s Seeder) PlantArticleTag(
	ctx context.Context,
	opts ...articleTagSeedOption,
) (models.ArticleTag, error) {
	data := &articleTagSeedData{
		ID:        uuid.New(),
		CreatedAt: time.Now().Add(-time.Duration(rand.IntN(365)) * 24 * time.Hour), //nolint:gosec // G404: Weak random for test data is acceptable
		UpdatedAt: time.Now(),
		Title:     getRandomTagTitle(),
	}

	for _, opt := range opts {
		opt(data)
	}

	tag, err := models.NewArticleTag(ctx, s.dbtx, models.NewArticleTagPayload{
		Title: data.Title,
	})
	if err != nil {
		return models.ArticleTag{}, err
	}

	return tag, nil
}

func (s Seeder) PlantArticleTags(
	ctx context.Context,
	amount int,
) ([]models.ArticleTag, error) {
	tags := make([]models.ArticleTag, amount)

	for i := range amount {
		tag, err := s.PlantArticleTag(ctx)
		if err != nil {
			return nil, err
		}
		tags[i] = tag
	}

	return tags, nil
}

func (s Seeder) PlantPredefinedArticleTags(
	ctx context.Context,
) ([]models.ArticleTag, error) {
	predefinedTags := getPredefinedTags()
	tags := make([]models.ArticleTag, len(predefinedTags))

	for i, tagTitle := range predefinedTags {
		tag, err := s.PlantArticleTag(
			ctx,
			WithArticleTagTitle(tagTitle),
		)
		if err != nil {
			return nil, err
		}
		tags[i] = tag
	}

	return tags, nil
}

// getPredefinedTags returns a curated list of common tech tags
func getPredefinedTags() []string {
	return []string{
		// Programming Languages
		"Go",
		"TypeScript",
		"JavaScript",
		"Python",
		"Rust",
		"Java",
		"C++",
		"PHP",

		// Databases
		"PostgreSQL",
		"MySQL",
		"Redis",
		"MongoDB",
		"SQLite",
		"Elasticsearch",
		"CouchDB",

		// Tools & Platforms
		"Docker",
		"Kubernetes",
		"Git",
		"Linux",
		"AWS",
		"Azure",
		"Google Cloud",
		"Terraform",
		"Ansible",
		"Jenkins",
		"GitHub Actions",
		"GitLab CI",

		// Frameworks & Libraries
		"React",
		"Vue.js",
		"Angular",
		"Express.js",
		"Django",
		"Flask",
		"Spring Boot",
		"Laravel",
		"Echo",
		"Gin",

		// Concepts & Practices
		"API",
		"REST",
		"GraphQL",
		"Microservices",
		"Testing",
		"Unit Testing",
		"Integration Testing",
		"TDD",
		"Performance",
		"Security",
		"Authentication",
		"Authorization",
		"CI/CD",
		"DevOps",
		"Monitoring",
		"Logging",
		"Caching",
		"Load Balancing",
		"Scaling",
		"Architecture",
		"Design Patterns",
		"Clean Code",
		"Refactoring",
		"Code Review",
		"Documentation",

		// Frontend
		"HTML",
		"CSS",
		"Sass",
		"Tailwind CSS",
		"Bootstrap",
		"Webpack",
		"Vite",
		"ESLint",
		"Prettier",

		// Backend
		"Node.js",
		"Deno",
		"Express",
		"FastAPI",
		"gRPC",
		"Message Queues",
		"RabbitMQ",
		"Apache Kafka",

		// Mobile
		"React Native",
		"Flutter",
		"Swift",
		"Kotlin",

		// Other
		"Machine Learning",
		"AI",
		"Data Science",
		"Blockchain",
		"WebAssembly",
		"Progressive Web Apps",
		"Serverless",
		"Edge Computing",
	}
}

// getRandomTagTitle returns a random tag from the predefined list
func getRandomTagTitle() string {
	tags := getPredefinedTags()
	return tags[rand.IntN(len(tags))] //nolint:gosec // G404: Weak random for test data is acceptable
}
