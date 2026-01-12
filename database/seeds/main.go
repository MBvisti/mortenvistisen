package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"mortenvistisen/internal/storage"
	"mortenvistisen/models/factories"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
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

	if err := seedUsers(ctx, db.Conn()); err != nil {
		return err
	}

	if err := seedArticles(ctx, db.Conn()); err != nil {
		return err
	}

	if err := seedSubscribers(ctx, db.Conn()); err != nil {
		return err
	}

	if err := seedNewsletters(ctx, db.Conn()); err != nil {
		return err
	}

	fmt.Println("Seeding complete!")
	return nil
}

func seedUsers(ctx context.Context, db storage.Executor) error {
	fmt.Println("\n--- Seeding Users ---")

	admin, err := factories.CreateUser(ctx, db,
		factories.WithEmail("admin@example.com"),
		factories.WithIsAdmin(true),
		factories.WithValidatedEmail(),
	)
	if err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}
	fmt.Printf("Created admin user: %s\n", admin.Email)

	user, err := factories.CreateUser(ctx, db,
		factories.WithEmail("user@example.com"),
		factories.WithValidatedEmail(),
	)
	if err != nil {
		return fmt.Errorf("failed to create regular user: %w", err)
	}
	fmt.Printf("Created regular user: %s\n", user.Email)

	users, err := factories.CreateUsers(ctx, db, 8, factories.WithValidatedEmail())
	if err != nil {
		return fmt.Errorf("failed to create additional users: %w", err)
	}
	fmt.Printf("Created %d additional users\n", len(users))

	return nil
}

func seedArticles(ctx context.Context, db storage.Executor) error {
	fmt.Println("\n--- Seeding Articles ---")

	now := time.Now()

	tags, err := factories.CreateTags(ctx, db, 5)
	if err != nil {
		return fmt.Errorf("failed to create tags: %w", err)
	}

	tagIDs := make([]uuid.UUID, len(tags))
	for i, tag := range tags {
		tagIDs[i] = tag.ID
	}

	publishedArticles, err := factories.CreateArticles(ctx, db, 15,
		factories.WithArticlesFirstPublishedAt(now.AddDate(0, 0, -30)),
		factories.WithArticlesTags(tagIDs),
	)
	if err != nil {
		return fmt.Errorf("failed to create published articles: %w", err)
	}
	fmt.Printf("Created %d published articles with tags\n", len(publishedArticles))

	draftArticles, err := factories.CreateArticles(ctx, db, 8,
		factories.WithArticlesTags(tagIDs),
	)
	if err != nil {
		return fmt.Errorf("failed to create draft articles: %w", err)
	}
	fmt.Printf("Created %d draft articles with tags\n", len(draftArticles))

	return nil
}

func seedSubscribers(ctx context.Context, db storage.Executor) error {
	fmt.Println("\n--- Seeding Subscribers ---")

	now := time.Now()

	verifiedSubscribers, err := factories.CreateSubscribers(ctx, db, 20,
		factories.WithSubscribersSubscribedAt(now.AddDate(0, 0, -60)),
		factories.WithSubscribersIsVerified(true),
	)
	if err != nil {
		return fmt.Errorf("failed to create verified subscribers: %w", err)
	}
	fmt.Printf("Created %d verified subscribers\n", len(verifiedSubscribers))

	unverifiedSubscribers, err := factories.CreateSubscribers(ctx, db, 5)
	if err != nil {
		return fmt.Errorf("failed to create unverified subscribers: %w", err)
	}
	fmt.Printf("Created %d unverified subscribers\n", len(unverifiedSubscribers))

	return nil
}

func seedNewsletters(ctx context.Context, db storage.Executor) error {
	fmt.Println("\n--- Seeding Newsletters ---")

	now := time.Now()

	publishedNewsletters, err := factories.CreateNewsletters(ctx, db, 10,
		factories.WithNewslettersIsPublished(true),
		factories.WithNewslettersReleasedAt(now.AddDate(0, 0, -14)),
	)
	if err != nil {
		return fmt.Errorf("failed to create published newsletters: %w", err)
	}
	fmt.Printf("Created %d published newsletters\n", len(publishedNewsletters))

	draftNewsletters, err := factories.CreateNewsletters(ctx, db, 5)
	if err != nil {
		return fmt.Errorf("failed to create draft newsletters: %w", err)
	}
	fmt.Printf("Created %d draft newsletters\n", len(draftNewsletters))

	return nil
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
