package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"mortenvistisen/internal/storage"
	"mortenvistisen/models"
	"mortenvistisen/models/factories"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	godotenv.Load()

	ctx := context.Background()

	db, err := storage.NewConnection(ctx, buildDatabaseURL())
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	fmt.Println("Seeding database...")

	admin, err := ensureAdminUser(ctx, db.Conn())
	if err != nil {
		return fmt.Errorf("failed to ensure admin user: %w", err)
	}
	fmt.Printf("Admin user ready: %s (password: password123)\n", admin.Email)

	targetTags := 12
	targetArticles := int64(36)
	targetNewsletters := int64(18)
	targetProjects := int64(12)
	targetSubscribers := int64(120)

	tags, createdTags, err := seedTags(ctx, db.Conn(), targetTags)
	if err != nil {
		return fmt.Errorf("failed to seed tags: %w", err)
	}
	fmt.Printf("Tags: +%d (total available: %d)\n", createdTags, len(tags))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	createdArticles, err := seedArticles(ctx, db.Conn(), tags, targetArticles, rng)
	if err != nil {
		return fmt.Errorf("failed to seed articles: %w", err)
	}
	fmt.Printf("Articles: +%d\n", createdArticles)

	createdNewsletters, err := seedNewsletters(ctx, db.Conn(), targetNewsletters)
	if err != nil {
		return fmt.Errorf("failed to seed newsletters: %w", err)
	}
	fmt.Printf("Newsletters: +%d\n", createdNewsletters)

	createdProjects, err := seedProjects(ctx, db.Conn(), targetProjects)
	if err != nil {
		return fmt.Errorf("failed to seed projects: %w", err)
	}
	fmt.Printf("Projects: +%d\n", createdProjects)

	createdSubscribers, err := seedSubscribers(ctx, db.Conn(), targetSubscribers)
	if err != nil {
		return fmt.Errorf("failed to seed subscribers: %w", err)
	}
	fmt.Printf("Subscribers: +%d\n", createdSubscribers)

	fmt.Println("Seeding complete!")
	return nil
}

func ensureAdminUser(
	ctx context.Context,
	exec storage.Executor,
) (models.User, error) {
	const adminEmail = "admin@example.com"

	existing, err := models.FindUserByEmail(ctx, exec, adminEmail)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		created, createErr := factories.CreateUser(
			ctx,
			exec,
			factories.WithEmail(adminEmail),
			factories.WithIsAdmin(true),
			factories.WithValidatedEmail(),
		)
		if createErr != nil {
			return models.User{}, createErr
		}
		return created, nil
	}

	validatedAt := existing.EmailValidatedAt
	if validatedAt.IsZero() {
		validatedAt = time.Now().UTC()
	}

	return models.UpdateUser(ctx, exec, models.UpdateUserData{
		ID:    existing.ID,
		Email: adminEmail,
		EmailValidatedAt: sql.NullTime{
			Time:  validatedAt,
			Valid: true,
		},
		Password: existing.Password,
		IsAdmin:  true,
	})
}

func seedTags(
	ctx context.Context,
	exec storage.Executor,
	targetCount int,
) ([]models.Tag, int, error) {
	baseTitles := []string{
		"Engineering",
		"Product",
		"Design",
		"AI",
		"Startups",
		"Marketing",
		"Operations",
		"Culture",
		"Security",
		"Data",
		"Growth",
		"Leadership",
	}

	tags, err := models.AllTags(ctx, exec)
	if err != nil {
		return nil, 0, err
	}

	created := 0
	for i := len(tags); i < targetCount; i++ {
		title := baseTitles[i%len(baseTitles)]
		if i >= len(baseTitles) {
			title = fmt.Sprintf("%s %d", title, i+1)
		}

		tag, createErr := models.CreateTag(ctx, exec, models.CreateTagData{
			Title: title,
		})
		if createErr != nil {
			return nil, created, createErr
		}
		tags = append(tags, tag)
		created++
	}

	return tags, created, nil
}

func seedArticles(
	ctx context.Context,
	exec storage.Executor,
	tags []models.Tag,
	targetCount int64,
	rng *rand.Rand,
) (int, error) {
	existingCount, err := models.CountArticles(ctx, exec)
	if err != nil {
		return 0, err
	}
	if existingCount >= targetCount {
		return 0, nil
	}

	topics := []string{
		"Shipping Faster With Better Feedback Loops",
		"Practical AI Workflows for Small Teams",
		"Writing Better Product Specs",
		"How to Run Useful Postmortems",
		"Scaling a Content Engine Without Burnout",
		"SEO Basics That Still Matter",
		"Roadmapping With Real Constraints",
		"Onboarding Improvements That Actually Work",
		"Designing for Clarity in Complex Flows",
		"Making Reliability Visible to Customers",
		"Building an Effective Editorial Process",
		"Metrics to Track in Early-Stage Products",
	}

	now := time.Now().UTC()
	created := 0
	for i := existingCount; i < targetCount; i++ {
		articleNumber := i + 1
		topic := topics[i%int64(len(topics))]
		slug := fmt.Sprintf("seed-article-%03d-%d", articleNumber, now.Unix())
		releasedAt := now.AddDate(0, 0, -int(articleNumber))

		article, createErr := models.CreateArticle(ctx, exec, models.CreateArticleData{
			FirstPublishedAt: releasedAt,
			Published:        articleNumber%4 != 0,
			Title:            fmt.Sprintf("Article %d: %s", articleNumber, topic),
			Excerpt:          fmt.Sprintf("A practical take on %s.", topic),
			MetaTitle:        fmt.Sprintf("Article %d", articleNumber),
			MetaDescription:  fmt.Sprintf("Seeded article %d about %s.", articleNumber, topic),
			Slug:             slug,
			ImageLink:        fmt.Sprintf("https://picsum.photos/seed/article-%d/1200/675", articleNumber),
			ReadTime:         int32(3 + (articleNumber % 10)),
			Content: fmt.Sprintf(
				"## %s\n\nThis seeded article exists to populate the admin UI with realistic content. It includes enough text for previews, table listings, and detail pages.\n\n### Key points\n\n- Focus on clear writing.\n- Keep scope practical.\n- Publish consistently.",
				topic,
			),
		})
		if createErr != nil {
			return created, createErr
		}
		created++

		if len(tags) == 0 {
			continue
		}

		linkCount := 1 + rng.Intn(3)
		for j := 0; j < linkCount; j++ {
			tag := tags[(int(articleNumber)+j)%len(tags)]
			if _, execErr := exec.Exec(
				ctx,
				"insert into article_tag_connections (article_id, tag_id) values ($1, $2)",
				article.ID,
				tag.ID,
			); execErr != nil {
				return created, execErr
			}
		}
	}

	return created, nil
}

func seedNewsletters(
	ctx context.Context,
	exec storage.Executor,
	targetCount int64,
) (int, error) {
	existingCount, err := models.CountNewsletters(ctx, exec)
	if err != nil {
		return 0, err
	}
	if existingCount >= targetCount {
		return 0, nil
	}

	now := time.Now().UTC()
	created := 0
	for i := existingCount; i < targetCount; i++ {
		number := i + 1
		releasedAt := now.AddDate(0, 0, -int(number*3))

		_, createErr := models.CreateNewsletter(ctx, exec, models.CreateNewsletterData{
			Title:           fmt.Sprintf("Weekly Update #%02d", number),
			Slug:            fmt.Sprintf("weekly-update-%02d-%d", number, now.Unix()),
			MetaTitle:       fmt.Sprintf("Weekly Update %02d", number),
			MetaDescription: fmt.Sprintf("Highlights and updates for week %02d.", number),
			IsPublished:     number%3 != 0,
			ReleasedAt:      releasedAt,
			Content:         fmt.Sprintf("This is seeded newsletter #%02d with curated updates, articles, and announcements.", number),
		})
		if createErr != nil {
			return created, createErr
		}
		created++
	}

	return created, nil
}

func seedSubscribers(
	ctx context.Context,
	exec storage.Executor,
	targetCount int64,
) (int, error) {
	existingCount, err := models.CountSubscribers(ctx, exec)
	if err != nil {
		return 0, err
	}
	if existingCount >= targetCount {
		return 0, nil
	}

	now := time.Now().UTC()
	referrers := []string{
		"organic",
		"twitter",
		"linkedin",
		"newsletter",
		"friend",
		"google",
	}

	created := 0
	for i := existingCount; i < targetCount; i++ {
		number := i + 1
		_, createErr := models.CreateSubscriber(ctx, exec, models.CreateSubscriberData{
			Email:        fmt.Sprintf("subscriber%03d+%d@example.com", number, now.Unix()),
			SubscribedAt: now.AddDate(0, 0, -int(number%60)),
			Referer:      referrers[i%int64(len(referrers))],
			IsVerified:   number%5 != 0,
		})
		if createErr != nil {
			return created, createErr
		}
		created++
	}

	return created, nil
}

func seedProjects(
	ctx context.Context,
	exec storage.Executor,
	targetCount int64,
) (int, error) {
	existingCount, err := models.CountProjects(ctx, exec)
	if err != nil {
		return 0, err
	}
	if existingCount >= targetCount {
		return 0, nil
	}

	now := time.Now().UTC()
	names := []string{
		"Planetaria",
		"Animaginary",
		"HelioStream",
		"CosmOS",
		"OpenShuttle",
		"Orbit Pulse",
		"Signal Dock",
		"Crewboard",
		"Beacon Ops",
		"Nebula Notes",
		"Mission Forge",
		"Lighthouse UI",
	}

	created := 0
	for i := existingCount; i < targetCount; i++ {
		number := i + 1
		name := names[i%int64(len(names))]
		startedAt := now.AddDate(0, -int(number), -int(number*2))
		baseSlug := fmt.Sprintf("%s-%02d-%d", name, number, now.Unix())

		_, createErr := models.CreateProject(ctx, exec, models.CreateProjectData{
			Published:   number%4 != 0,
			Title:       fmt.Sprintf("%s %02d", name, number),
			Slug:        baseSlug,
			StartedAt:   startedAt,
			Status:      []string{"Planned", "In Progress", "Completed"}[i%3],
			Description: fmt.Sprintf("%s is a seeded project focused on shipping a polished user experience with reliable backend performance.", name),
			ProjectURL:  fmt.Sprintf("https://example.com/projects/%s", baseSlug),
			Content: fmt.Sprintf(
				"## %s\n\nThis seeded project entry exists to populate both admin and public project views.\n\n### What was shipped\n\n- Responsive marketing page\n- Authenticated dashboard\n- Analytics instrumentation\n\n### Learnings\n\nShipping smaller increments made it easier to keep quality high and avoid hidden complexity.",
				name,
			),
		})
		if createErr != nil {
			return created, createErr
		}
		created++
	}

	return created, nil
}

func buildDatabaseURL() string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_KIND"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"),
	)
}
