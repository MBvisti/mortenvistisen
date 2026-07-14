package seeds

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"mortenvistisen/internal/storage"
	"mortenvistisen/models"
	"mortenvistisen/models/factories"
)

// Development is intended to run once after resetting the local database.
func Development(ctx context.Context, exec storage.Executor) error {
	now := time.Now().UTC().Truncate(time.Second)

	admin, err := factories.CreateUser(ctx, exec,
		factories.WithEmail("admin@example.com"),
		factories.WithIsAdmin(true),
		factories.WithValidatedEmail(),
	)
	if err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	user, err := factories.CreateUser(ctx, exec,
		factories.WithEmail("user@example.com"),
		factories.WithValidatedEmail(),
	)
	if err != nil {
		return fmt.Errorf("failed to create regular user: %w", err)
	}

	if _, err := factories.CreateSubscriber(
		ctx,
		exec,
		factories.WithSubscribersEmail(
			sql.NullString{String: "alex.morgan@example.com", Valid: true},
		),
		factories.WithSubscribersSubscribedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -62), Valid: true},
		),
		factories.WithSubscribersReferer(
			sql.NullString{String: "/articles/boring-software-scales", Valid: true},
		),
		factories.WithSubscribersIsVerified(sql.NullBool{Bool: true, Valid: true}),
	); err != nil {
		return fmt.Errorf("failed to create subscriber alex.morgan@example.com: %w", err)
	}
	if _, err := factories.CreateSubscriber(
		ctx,
		exec,
		factories.WithSubscribersEmail(
			sql.NullString{String: "jamie.chen@example.com", Valid: true},
		),
		factories.WithSubscribersSubscribedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -55), Valid: true},
		),
		factories.WithSubscribersReferer(sql.NullString{String: "/newsletter", Valid: true}),
		factories.WithSubscribersIsVerified(sql.NullBool{Bool: true, Valid: true}),
	); err != nil {
		return fmt.Errorf("failed to create subscriber jamie.chen@example.com: %w", err)
	}
	if _, err := factories.CreateSubscriber(
		ctx,
		exec,
		factories.WithSubscribersEmail(
			sql.NullString{String: "sam.rivera@example.com", Valid: true},
		),
		factories.WithSubscribersSubscribedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -48), Valid: true},
		),
		factories.WithSubscribersReferer(
			sql.NullString{String: "/articles/operationally-kind-go-services", Valid: true},
		),
		factories.WithSubscribersIsVerified(sql.NullBool{Bool: true, Valid: true}),
	); err != nil {
		return fmt.Errorf("failed to create subscriber sam.rivera@example.com: %w", err)
	}
	if _, err := factories.CreateSubscriber(
		ctx,
		exec,
		factories.WithSubscribersEmail(
			sql.NullString{String: "priya.shah@example.com", Valid: true},
		),
		factories.WithSubscribersSubscribedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -41), Valid: true},
		),
		factories.WithSubscribersReferer(sql.NullString{String: "https://github.com", Valid: true}),
		factories.WithSubscribersIsVerified(sql.NullBool{Bool: true, Valid: true}),
	); err != nil {
		return fmt.Errorf("failed to create subscriber priya.shah@example.com: %w", err)
	}
	if _, err := factories.CreateSubscriber(
		ctx,
		exec,
		factories.WithSubscribersEmail(
			sql.NullString{String: "noah.williams@example.com", Valid: true},
		),
		factories.WithSubscribersSubscribedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -36), Valid: true},
		),
		factories.WithSubscribersReferer(sql.NullString{String: "/", Valid: true}),
		factories.WithSubscribersIsVerified(sql.NullBool{Bool: false, Valid: true}),
	); err != nil {
		return fmt.Errorf("failed to create subscriber noah.williams@example.com: %w", err)
	}
	if _, err := factories.CreateSubscriber(
		ctx,
		exec,
		factories.WithSubscribersEmail(
			sql.NullString{String: "maya.patel@example.com", Valid: true},
		),
		factories.WithSubscribersSubscribedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -33), Valid: true},
		),
		factories.WithSubscribersReferer(
			sql.NullString{String: "/articles/designing-admin-tools-for-speed", Valid: true},
		),
		factories.WithSubscribersIsVerified(sql.NullBool{Bool: true, Valid: true}),
	); err != nil {
		return fmt.Errorf("failed to create subscriber maya.patel@example.com: %w", err)
	}
	if _, err := factories.CreateSubscriber(
		ctx,
		exec,
		factories.WithSubscribersEmail(
			sql.NullString{String: "lucas.martin@example.com", Valid: true},
		),
		factories.WithSubscribersSubscribedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -29), Valid: true},
		),
		factories.WithSubscribersReferer(sql.NullString{String: "/newsletter", Valid: true}),
		factories.WithSubscribersIsVerified(sql.NullBool{Bool: true, Valid: true}),
	); err != nil {
		return fmt.Errorf("failed to create subscriber lucas.martin@example.com: %w", err)
	}
	if _, err := factories.CreateSubscriber(
		ctx,
		exec,
		factories.WithSubscribersEmail(
			sql.NullString{String: "sofia.garcia@example.com", Valid: true},
		),
		factories.WithSubscribersSubscribedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -26), Valid: true},
		),
		factories.WithSubscribersReferer(
			sql.NullString{String: "/articles/small-interfaces-better-go", Valid: true},
		),
		factories.WithSubscribersIsVerified(sql.NullBool{Bool: true, Valid: true}),
	); err != nil {
		return fmt.Errorf("failed to create subscriber sofia.garcia@example.com: %w", err)
	}
	if _, err := factories.CreateSubscriber(
		ctx,
		exec,
		factories.WithSubscribersEmail(
			sql.NullString{String: "ethan.brooks@example.com", Valid: true},
		),
		factories.WithSubscribersSubscribedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -23), Valid: true},
		),
		factories.WithSubscribersReferer(sql.NullString{String: "/", Valid: true}),
		factories.WithSubscribersIsVerified(sql.NullBool{Bool: false, Valid: true}),
	); err != nil {
		return fmt.Errorf("failed to create subscriber ethan.brooks@example.com: %w", err)
	}
	if _, err := factories.CreateSubscriber(
		ctx,
		exec,
		factories.WithSubscribersEmail(
			sql.NullString{String: "amina.yusuf@example.com", Valid: true},
		),
		factories.WithSubscribersSubscribedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -19), Valid: true},
		),
		factories.WithSubscribersReferer(sql.NullString{String: "/newsletter", Valid: true}),
		factories.WithSubscribersIsVerified(sql.NullBool{Bool: true, Valid: true}),
	); err != nil {
		return fmt.Errorf("failed to create subscriber amina.yusuf@example.com: %w", err)
	}
	if _, err := factories.CreateSubscriber(
		ctx,
		exec,
		factories.WithSubscribersEmail(
			sql.NullString{String: "liam.kelly@example.com", Valid: true},
		),
		factories.WithSubscribersSubscribedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -16), Valid: true},
		),
		factories.WithSubscribersReferer(
			sql.NullString{String: "/articles/measuring-what-matters", Valid: true},
		),
		factories.WithSubscribersIsVerified(sql.NullBool{Bool: true, Valid: true}),
	); err != nil {
		return fmt.Errorf("failed to create subscriber liam.kelly@example.com: %w", err)
	}
	if _, err := factories.CreateSubscriber(
		ctx,
		exec,
		factories.WithSubscribersEmail(
			sql.NullString{String: "zoe.thompson@example.com", Valid: true},
		),
		factories.WithSubscribersSubscribedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -13), Valid: true},
		),
		factories.WithSubscribersReferer(
			sql.NullString{String: "https://www.linkedin.com", Valid: true},
		),
		factories.WithSubscribersIsVerified(sql.NullBool{Bool: true, Valid: true}),
	); err != nil {
		return fmt.Errorf("failed to create subscriber zoe.thompson@example.com: %w", err)
	}
	if _, err := factories.CreateSubscriber(
		ctx,
		exec,
		factories.WithSubscribersEmail(
			sql.NullString{String: "daniel.kim@example.com", Valid: true},
		),
		factories.WithSubscribersSubscribedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -11), Valid: true},
		),
		factories.WithSubscribersReferer(
			sql.NullString{String: "/articles/the-maintenance-mindset", Valid: true},
		),
		factories.WithSubscribersIsVerified(sql.NullBool{Bool: true, Valid: true}),
	); err != nil {
		return fmt.Errorf("failed to create subscriber daniel.kim@example.com: %w", err)
	}
	if _, err := factories.CreateSubscriber(
		ctx,
		exec,
		factories.WithSubscribersEmail(
			sql.NullString{String: "ines.santos@example.com", Valid: true},
		),
		factories.WithSubscribersSubscribedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -9), Valid: true},
		),
		factories.WithSubscribersReferer(sql.NullString{String: "/newsletter", Valid: true}),
		factories.WithSubscribersIsVerified(sql.NullBool{Bool: false, Valid: true}),
	); err != nil {
		return fmt.Errorf("failed to create subscriber ines.santos@example.com: %w", err)
	}
	if _, err := factories.CreateSubscriber(
		ctx,
		exec,
		factories.WithSubscribersEmail(
			sql.NullString{String: "oliver.jones@example.com", Valid: true},
		),
		factories.WithSubscribersSubscribedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -7), Valid: true},
		),
		factories.WithSubscribersReferer(sql.NullString{String: "/", Valid: true}),
		factories.WithSubscribersIsVerified(sql.NullBool{Bool: true, Valid: true}),
	); err != nil {
		return fmt.Errorf("failed to create subscriber oliver.jones@example.com: %w", err)
	}
	if _, err := factories.CreateSubscriber(
		ctx,
		exec,
		factories.WithSubscribersEmail(
			sql.NullString{String: "fatima.hassan@example.com", Valid: true},
		),
		factories.WithSubscribersSubscribedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -5), Valid: true},
		),
		factories.WithSubscribersReferer(
			sql.NullString{String: "/articles/boring-software-scales", Valid: true},
		),
		factories.WithSubscribersIsVerified(sql.NullBool{Bool: true, Valid: true}),
	); err != nil {
		return fmt.Errorf("failed to create subscriber fatima.hassan@example.com: %w", err)
	}
	if _, err := factories.CreateSubscriber(
		ctx,
		exec,
		factories.WithSubscribersEmail(
			sql.NullString{String: "theo.andersen@example.com", Valid: true},
		),
		factories.WithSubscribersSubscribedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -3), Valid: true},
		),
		factories.WithSubscribersReferer(sql.NullString{String: "/newsletter", Valid: true}),
		factories.WithSubscribersIsVerified(sql.NullBool{Bool: true, Valid: true}),
	); err != nil {
		return fmt.Errorf("failed to create subscriber theo.andersen@example.com: %w", err)
	}
	if _, err := factories.CreateSubscriber(
		ctx,
		exec,
		factories.WithSubscribersEmail(
			sql.NullString{String: "clara.rossi@example.com", Valid: true},
		),
		factories.WithSubscribersSubscribedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -1), Valid: true},
		),
		factories.WithSubscribersReferer(
			sql.NullString{String: "/articles/operationally-kind-go-services", Valid: true},
		),
		factories.WithSubscribersIsVerified(sql.NullBool{Bool: false, Valid: true}),
	); err != nil {
		return fmt.Errorf("failed to create subscriber clara.rossi@example.com: %w", err)
	}

	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -52), Valid: true},
		),
		factories.WithArticlesPublished(true),
		factories.WithArticlesTitle("Boring Software Scales"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "Why predictable systems often outperform clever ones as teams and products grow.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{String: "Boring Software Scales", Valid: true},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "A practical argument for predictable architecture, explicit code, and calm operations.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("boring-software-scales"),
		factories.WithArticlesImageLink(sql.NullString{}),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 7, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# Boring Software Scales\n\nClear boundaries, familiar tools, and explicit control flow make software easier to operate. Choose the smallest design that explains the problem honestly.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article boring-software-scales: %w", err)
	}
	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -41), Valid: true},
		),
		factories.WithArticlesPublished(true),
		factories.WithArticlesTitle("Operationally Kind Go Services"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "Designing Go services that explain what they are doing when production gets difficult.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{String: "Operationally Kind Go Services", Valid: true},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "Patterns for building observable, diagnosable, and operationally friendly Go services.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("operationally-kind-go-services"),
		factories.WithArticlesImageLink(sql.NullString{}),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 9, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# Operationally Kind Go Services\n\nGood logs, useful metrics, and clear failure modes are part of the product. Reliability grows from small choices that reduce uncertainty.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article operationally-kind-go-services: %w", err)
	}
	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -33), Valid: true},
		),
		factories.WithArticlesPublished(true),
		factories.WithArticlesTitle("Designing Admin Tools for Speed"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "Admin interfaces should reduce hesitation and make common work feel obvious.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{String: "Designing Admin Tools for Speed", Valid: true},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "Practical interface principles for fast, legible, and trustworthy administration tools.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("designing-admin-tools-for-speed"),
		factories.WithArticlesImageLink(sql.NullString{}),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 6, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# Designing Admin Tools for Speed\n\nShow the numbers people ask for every morning and keep recent records close to those numbers.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article designing-admin-tools-for-speed: %w", err)
	}
	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -25), Valid: true},
		),
		factories.WithArticlesPublished(true),
		factories.WithArticlesTitle("Small Interfaces, Better Go"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "Keeping Go interfaces close to consumers produces simpler and more flexible code.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{String: "Small Interfaces, Better Go", Valid: true},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "How consumer-owned interfaces improve testing, reduce coupling, and clarify Go programs.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("small-interfaces-better-go"),
		factories.WithArticlesImageLink(sql.NullString{}),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 8, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# Small Interfaces, Better Go\n\nStart with concrete types. Introduce an interface where a caller needs substitution, then keep it beside the caller.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article small-interfaces-better-go: %w", err)
	}
	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -16), Valid: true},
		),
		factories.WithArticlesPublished(true),
		factories.WithArticlesTitle("Measuring What Matters"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "A useful metric should help someone decide what to do next.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{String: "Measuring What Matters", Valid: true},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "A grounded approach to product and engineering metrics that support real decisions.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("measuring-what-matters"),
		factories.WithArticlesImageLink(sql.NullString{}),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 5, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# Measuring What Matters\n\nWrite down the decision first, then identify the smallest measurement that reduces uncertainty around it.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article measuring-what-matters: %w", err)
	}
	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -8), Valid: true},
		),
		factories.WithArticlesPublished(true),
		factories.WithArticlesTitle("The Maintenance Mindset"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "Treating future changes as a design constraint leads to calmer software.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{String: "The Maintenance Mindset", Valid: true},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "Why maintainability starts with today's naming, boundaries, tests, and operational choices.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("the-maintenance-mindset"),
		factories.WithArticlesImageLink(sql.NullString{}),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 7, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# The Maintenance Mindset\n\nMaintenance begins with the first decision another person will need to understand or reverse.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article the-maintenance-mindset: %w", err)
	}
	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(sql.NullTime{}),
		factories.WithArticlesPublished(false),
		factories.WithArticlesTitle("Notes on Sustainable Delivery"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "Early notes on building a delivery rhythm that stays useful under pressure.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{String: "Notes on Sustainable Delivery", Valid: true},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "Draft notes about sustainable engineering delivery and healthy feedback loops.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("notes-on-sustainable-delivery"),
		factories.WithArticlesImageLink(sql.NullString{}),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 4, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# Notes on Sustainable Delivery\n\nSmall batches, fast feedback, and quiet time create a sustainable delivery rhythm.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article notes-on-sustainable-delivery: %w", err)
	}
	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(sql.NullTime{}),
		factories.WithArticlesPublished(false),
		factories.WithArticlesTitle("A Field Guide to Useful Constraints"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "A working draft about constraints that sharpen decisions without limiting learning.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{String: "A Field Guide to Useful Constraints", Valid: true},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "Draft field guide to technical constraints that improve focus and system clarity.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("field-guide-useful-constraints"),
		factories.WithArticlesImageLink(sql.NullString{}),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 6, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# A Field Guide to Useful Constraints\n\nGood constraints narrow the decision space while leaving teams free to solve the actual problem.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article field-guide-useful-constraints: %w", err)
	}

	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -70), Valid: true},
		),
		factories.WithArticlesPublished(true),
		factories.WithArticlesTitle("Workshop Notes 01: Keeping Systems Legible"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "Practical notes on keeping systems legible for teams building software that should remain calm and understandable.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{String: "Workshop Notes 01: Keeping Systems Legible", Valid: true},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "A practical guide to keeping systems legible for software teams that value clear decisions, reliable delivery, and maintainable systems.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("workshop-notes-01"),
		factories.WithArticlesImageLink(
			sql.NullString{
				String: "https://images.example.com/articles/workshop-notes-01.jpg",
				Valid:  true,
			},
		),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 5, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# Keeping Systems Legible\n\nClear constraints, visible trade-offs, and small feedback loops help teams apply these ideas in everyday product work.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article workshop-notes-01: %w", err)
	}

	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -71), Valid: true},
		),
		factories.WithArticlesPublished(true),
		factories.WithArticlesTitle("Workshop Notes 02: Designing for Reversibility"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "Practical notes on designing for reversibility for teams building software that should remain calm and understandable.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{String: "Workshop Notes 02: Designing for Reversibility", Valid: true},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "A practical guide to designing for reversibility for software teams that value clear decisions, reliable delivery, and maintainable systems.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("workshop-notes-02"),
		factories.WithArticlesImageLink(
			sql.NullString{
				String: "https://images.example.com/articles/workshop-notes-02.jpg",
				Valid:  true,
			},
		),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 6, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# Designing for Reversibility\n\nClear constraints, visible trade-offs, and small feedback loops help teams apply these ideas in everyday product work.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article workshop-notes-02: %w", err)
	}

	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -72), Valid: true},
		),
		factories.WithArticlesPublished(true),
		factories.WithArticlesTitle("Workshop Notes 03: Writing Practical Error Messages"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "Practical notes on writing practical error messages for teams building software that should remain calm and understandable.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{
				String: "Workshop Notes 03: Writing Practical Error Messages",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "A practical guide to writing practical error messages for software teams that value clear decisions, reliable delivery, and maintainable systems.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("workshop-notes-03"),
		factories.WithArticlesImageLink(
			sql.NullString{
				String: "https://images.example.com/articles/workshop-notes-03.jpg",
				Valid:  true,
			},
		),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 7, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# Writing Practical Error Messages\n\nClear constraints, visible trade-offs, and small feedback loops help teams apply these ideas in everyday product work.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article workshop-notes-03: %w", err)
	}

	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -73), Valid: true},
		),
		factories.WithArticlesPublished(true),
		factories.WithArticlesTitle("Workshop Notes 04: Simple Queues and Clear Work"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "Practical notes on simple queues and clear work for teams building software that should remain calm and understandable.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{String: "Workshop Notes 04: Simple Queues and Clear Work", Valid: true},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "A practical guide to simple queues and clear work for software teams that value clear decisions, reliable delivery, and maintainable systems.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("workshop-notes-04"),
		factories.WithArticlesImageLink(
			sql.NullString{
				String: "https://images.example.com/articles/workshop-notes-04.jpg",
				Valid:  true,
			},
		),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 8, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# Simple Queues and Clear Work\n\nClear constraints, visible trade-offs, and small feedback loops help teams apply these ideas in everyday product work.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article workshop-notes-04: %w", err)
	}

	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -74), Valid: true},
		),
		factories.WithArticlesPublished(true),
		factories.WithArticlesTitle("Workshop Notes 05: Naming the Important Things"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "Practical notes on naming the important things for teams building software that should remain calm and understandable.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{String: "Workshop Notes 05: Naming the Important Things", Valid: true},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "A practical guide to naming the important things for software teams that value clear decisions, reliable delivery, and maintainable systems.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("workshop-notes-05"),
		factories.WithArticlesImageLink(
			sql.NullString{
				String: "https://images.example.com/articles/workshop-notes-05.jpg",
				Valid:  true,
			},
		),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 9, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# Naming the Important Things\n\nClear constraints, visible trade-offs, and small feedback loops help teams apply these ideas in everyday product work.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article workshop-notes-05: %w", err)
	}

	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -75), Valid: true},
		),
		factories.WithArticlesPublished(true),
		factories.WithArticlesTitle("Workshop Notes 06: Interfaces That Explain Themselves"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "Practical notes on interfaces that explain themselves for teams building software that should remain calm and understandable.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{
				String: "Workshop Notes 06: Interfaces That Explain Themselves",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "A practical guide to interfaces that explain themselves for software teams that value clear decisions, reliable delivery, and maintainable systems.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("workshop-notes-06"),
		factories.WithArticlesImageLink(
			sql.NullString{
				String: "https://images.example.com/articles/workshop-notes-06.jpg",
				Valid:  true,
			},
		),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 10, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# Interfaces That Explain Themselves\n\nClear constraints, visible trade-offs, and small feedback loops help teams apply these ideas in everyday product work.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article workshop-notes-06: %w", err)
	}

	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -76), Valid: true},
		),
		factories.WithArticlesPublished(true),
		factories.WithArticlesTitle("Workshop Notes 07: Reliable Work Without Heroics"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "Practical notes on reliable work without heroics for teams building software that should remain calm and understandable.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{String: "Workshop Notes 07: Reliable Work Without Heroics", Valid: true},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "A practical guide to reliable work without heroics for software teams that value clear decisions, reliable delivery, and maintainable systems.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("workshop-notes-07"),
		factories.WithArticlesImageLink(
			sql.NullString{
				String: "https://images.example.com/articles/workshop-notes-07.jpg",
				Valid:  true,
			},
		),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 5, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# Reliable Work Without Heroics\n\nClear constraints, visible trade-offs, and small feedback loops help teams apply these ideas in everyday product work.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article workshop-notes-07: %w", err)
	}

	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -77), Valid: true},
		),
		factories.WithArticlesPublished(true),
		factories.WithArticlesTitle("Workshop Notes 08: Useful Boundaries in Go"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "Practical notes on useful boundaries in go for teams building software that should remain calm and understandable.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{String: "Workshop Notes 08: Useful Boundaries in Go", Valid: true},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "A practical guide to useful boundaries in go for software teams that value clear decisions, reliable delivery, and maintainable systems.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("workshop-notes-08"),
		factories.WithArticlesImageLink(
			sql.NullString{
				String: "https://images.example.com/articles/workshop-notes-08.jpg",
				Valid:  true,
			},
		),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 6, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# Useful Boundaries in Go\n\nClear constraints, visible trade-offs, and small feedback loops help teams apply these ideas in everyday product work.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article workshop-notes-08: %w", err)
	}

	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -78), Valid: true},
		),
		factories.WithArticlesPublished(true),
		factories.WithArticlesTitle("Workshop Notes 09: Shipping Smaller Changes"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "Practical notes on shipping smaller changes for teams building software that should remain calm and understandable.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{String: "Workshop Notes 09: Shipping Smaller Changes", Valid: true},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "A practical guide to shipping smaller changes for software teams that value clear decisions, reliable delivery, and maintainable systems.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("workshop-notes-09"),
		factories.WithArticlesImageLink(
			sql.NullString{
				String: "https://images.example.com/articles/workshop-notes-09.jpg",
				Valid:  true,
			},
		),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 7, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# Shipping Smaller Changes\n\nClear constraints, visible trade-offs, and small feedback loops help teams apply these ideas in everyday product work.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article workshop-notes-09: %w", err)
	}

	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -79), Valid: true},
		),
		factories.WithArticlesPublished(true),
		factories.WithArticlesTitle("Workshop Notes 10: Feedback Loops That Work"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "Practical notes on feedback loops that work for teams building software that should remain calm and understandable.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{String: "Workshop Notes 10: Feedback Loops That Work", Valid: true},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "A practical guide to feedback loops that work for software teams that value clear decisions, reliable delivery, and maintainable systems.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("workshop-notes-10"),
		factories.WithArticlesImageLink(
			sql.NullString{
				String: "https://images.example.com/articles/workshop-notes-10.jpg",
				Valid:  true,
			},
		),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 8, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# Feedback Loops That Work\n\nClear constraints, visible trade-offs, and small feedback loops help teams apply these ideas in everyday product work.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article workshop-notes-10: %w", err)
	}

	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -80), Valid: true},
		),
		factories.WithArticlesPublished(true),
		factories.WithArticlesTitle("Workshop Notes 11: Making Ownership Visible"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "Practical notes on making ownership visible for teams building software that should remain calm and understandable.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{String: "Workshop Notes 11: Making Ownership Visible", Valid: true},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "A practical guide to making ownership visible for software teams that value clear decisions, reliable delivery, and maintainable systems.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("workshop-notes-11"),
		factories.WithArticlesImageLink(
			sql.NullString{
				String: "https://images.example.com/articles/workshop-notes-11.jpg",
				Valid:  true,
			},
		),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 9, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# Making Ownership Visible\n\nClear constraints, visible trade-offs, and small feedback loops help teams apply these ideas in everyday product work.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article workshop-notes-11: %w", err)
	}

	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -81), Valid: true},
		),
		factories.WithArticlesPublished(true),
		factories.WithArticlesTitle("Workshop Notes 12: Calm Incident Response"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "Practical notes on calm incident response for teams building software that should remain calm and understandable.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{String: "Workshop Notes 12: Calm Incident Response", Valid: true},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "A practical guide to calm incident response for software teams that value clear decisions, reliable delivery, and maintainable systems.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("workshop-notes-12"),
		factories.WithArticlesImageLink(
			sql.NullString{
				String: "https://images.example.com/articles/workshop-notes-12.jpg",
				Valid:  true,
			},
		),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 10, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# Calm Incident Response\n\nClear constraints, visible trade-offs, and small feedback loops help teams apply these ideas in everyday product work.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article workshop-notes-12: %w", err)
	}

	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -82), Valid: true},
		),
		factories.WithArticlesPublished(true),
		factories.WithArticlesTitle("Workshop Notes 13: Documentation as Product Design"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "Practical notes on documentation as product design for teams building software that should remain calm and understandable.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{
				String: "Workshop Notes 13: Documentation as Product Design",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "A practical guide to documentation as product design for software teams that value clear decisions, reliable delivery, and maintainable systems.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("workshop-notes-13"),
		factories.WithArticlesImageLink(
			sql.NullString{
				String: "https://images.example.com/articles/workshop-notes-13.jpg",
				Valid:  true,
			},
		),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 5, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# Documentation as Product Design\n\nClear constraints, visible trade-offs, and small feedback loops help teams apply these ideas in everyday product work.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article workshop-notes-13: %w", err)
	}

	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -83), Valid: true},
		),
		factories.WithArticlesPublished(true),
		factories.WithArticlesTitle("Workshop Notes 14: Testing the Decisions That Matter"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "Practical notes on testing the decisions that matter for teams building software that should remain calm and understandable.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{
				String: "Workshop Notes 14: Testing the Decisions That Matter",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "A practical guide to testing the decisions that matter for software teams that value clear decisions, reliable delivery, and maintainable systems.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("workshop-notes-14"),
		factories.WithArticlesImageLink(
			sql.NullString{
				String: "https://images.example.com/articles/workshop-notes-14.jpg",
				Valid:  true,
			},
		),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 6, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# Testing the Decisions That Matter\n\nClear constraints, visible trade-offs, and small feedback loops help teams apply these ideas in everyday product work.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article workshop-notes-14: %w", err)
	}

	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -84), Valid: true},
		),
		factories.WithArticlesPublished(true),
		factories.WithArticlesTitle("Workshop Notes 15: Choosing Boring Dependencies"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "Practical notes on choosing boring dependencies for teams building software that should remain calm and understandable.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{String: "Workshop Notes 15: Choosing Boring Dependencies", Valid: true},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "A practical guide to choosing boring dependencies for software teams that value clear decisions, reliable delivery, and maintainable systems.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("workshop-notes-15"),
		factories.WithArticlesImageLink(
			sql.NullString{
				String: "https://images.example.com/articles/workshop-notes-15.jpg",
				Valid:  true,
			},
		),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 7, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# Choosing Boring Dependencies\n\nClear constraints, visible trade-offs, and small feedback loops help teams apply these ideas in everyday product work.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article workshop-notes-15: %w", err)
	}

	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -85), Valid: true},
		),
		factories.WithArticlesPublished(true),
		factories.WithArticlesTitle("Workshop Notes 16: Good Defaults for Internal Tools"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "Practical notes on good defaults for internal tools for teams building software that should remain calm and understandable.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{
				String: "Workshop Notes 16: Good Defaults for Internal Tools",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "A practical guide to good defaults for internal tools for software teams that value clear decisions, reliable delivery, and maintainable systems.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("workshop-notes-16"),
		factories.WithArticlesImageLink(
			sql.NullString{
				String: "https://images.example.com/articles/workshop-notes-16.jpg",
				Valid:  true,
			},
		),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 8, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# Good Defaults for Internal Tools\n\nClear constraints, visible trade-offs, and small feedback loops help teams apply these ideas in everyday product work.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article workshop-notes-16: %w", err)
	}

	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -86), Valid: true},
		),
		factories.WithArticlesPublished(true),
		factories.WithArticlesTitle("Workshop Notes 17: Observability Without Noise"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "Practical notes on observability without noise for teams building software that should remain calm and understandable.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{String: "Workshop Notes 17: Observability Without Noise", Valid: true},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "A practical guide to observability without noise for software teams that value clear decisions, reliable delivery, and maintainable systems.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("workshop-notes-17"),
		factories.WithArticlesImageLink(
			sql.NullString{
				String: "https://images.example.com/articles/workshop-notes-17.jpg",
				Valid:  true,
			},
		),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 9, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# Observability Without Noise\n\nClear constraints, visible trade-offs, and small feedback loops help teams apply these ideas in everyday product work.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article workshop-notes-17: %w", err)
	}

	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -87), Valid: true},
		),
		factories.WithArticlesPublished(true),
		factories.WithArticlesTitle("Workshop Notes 18: Planning for Routine Maintenance"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "Practical notes on planning for routine maintenance for teams building software that should remain calm and understandable.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{
				String: "Workshop Notes 18: Planning for Routine Maintenance",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "A practical guide to planning for routine maintenance for software teams that value clear decisions, reliable delivery, and maintainable systems.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("workshop-notes-18"),
		factories.WithArticlesImageLink(
			sql.NullString{
				String: "https://images.example.com/articles/workshop-notes-18.jpg",
				Valid:  true,
			},
		),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 10, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# Planning for Routine Maintenance\n\nClear constraints, visible trade-offs, and small feedback loops help teams apply these ideas in everyday product work.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article workshop-notes-18: %w", err)
	}

	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(sql.NullTime{}),
		factories.WithArticlesPublished(false),
		factories.WithArticlesTitle("Workshop Notes 19: Reducing Coordination Costs"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "Practical notes on reducing coordination costs for teams building software that should remain calm and understandable.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{String: "Workshop Notes 19: Reducing Coordination Costs", Valid: true},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "A practical guide to reducing coordination costs for software teams that value clear decisions, reliable delivery, and maintainable systems.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("workshop-notes-19"),
		factories.WithArticlesImageLink(
			sql.NullString{
				String: "https://images.example.com/articles/workshop-notes-19.jpg",
				Valid:  true,
			},
		),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 5, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# Reducing Coordination Costs\n\nClear constraints, visible trade-offs, and small feedback loops help teams apply these ideas in everyday product work.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article workshop-notes-19: %w", err)
	}

	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(sql.NullTime{}),
		factories.WithArticlesPublished(false),
		factories.WithArticlesTitle("Workshop Notes 20: Learning From Production"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "Practical notes on learning from production for teams building software that should remain calm and understandable.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{String: "Workshop Notes 20: Learning From Production", Valid: true},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "A practical guide to learning from production for software teams that value clear decisions, reliable delivery, and maintainable systems.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("workshop-notes-20"),
		factories.WithArticlesImageLink(
			sql.NullString{
				String: "https://images.example.com/articles/workshop-notes-20.jpg",
				Valid:  true,
			},
		),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 6, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# Learning From Production\n\nClear constraints, visible trade-offs, and small feedback loops help teams apply these ideas in everyday product work.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article workshop-notes-20: %w", err)
	}

	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(sql.NullTime{}),
		factories.WithArticlesPublished(false),
		factories.WithArticlesTitle("Workshop Notes 21: Building Trustworthy Admin Flows"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "Practical notes on building trustworthy admin flows for teams building software that should remain calm and understandable.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{
				String: "Workshop Notes 21: Building Trustworthy Admin Flows",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "A practical guide to building trustworthy admin flows for software teams that value clear decisions, reliable delivery, and maintainable systems.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("workshop-notes-21"),
		factories.WithArticlesImageLink(
			sql.NullString{
				String: "https://images.example.com/articles/workshop-notes-21.jpg",
				Valid:  true,
			},
		),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 7, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# Building Trustworthy Admin Flows\n\nClear constraints, visible trade-offs, and small feedback loops help teams apply these ideas in everyday product work.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article workshop-notes-21: %w", err)
	}

	if _, err := factories.CreateArticle(
		ctx,
		exec,
		factories.WithArticlesFirstPublishedAt(sql.NullTime{}),
		factories.WithArticlesPublished(false),
		factories.WithArticlesTitle("Workshop Notes 22: Designing for the Next Change"),
		factories.WithArticlesExcerpt(
			sql.NullString{
				String: "Practical notes on designing for the next change for teams building software that should remain calm and understandable.",
				Valid:  true,
			},
		),
		factories.WithArticlesMetaTitle(
			sql.NullString{String: "Workshop Notes 22: Designing for the Next Change", Valid: true},
		),
		factories.WithArticlesMetaDescription(
			sql.NullString{
				String: "A practical guide to designing for the next change for software teams that value clear decisions, reliable delivery, and maintainable systems.",
				Valid:  true,
			},
		),
		factories.WithArticlesSlug("workshop-notes-22"),
		factories.WithArticlesImageLink(
			sql.NullString{
				String: "https://images.example.com/articles/workshop-notes-22.jpg",
				Valid:  true,
			},
		),
		factories.WithArticlesReadTime(sql.NullInt32{Int32: 8, Valid: true}),
		factories.WithArticlesContent(
			sql.NullString{
				String: "# Designing for the Next Change\n\nClear constraints, visible trade-offs, and small feedback loops help teams apply these ideas in everyday product work.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create article workshop-notes-22: %w", err)
	}

	tagTitles := []string{
		"Architecture",
		"Backend",
		"Career",
		"Delivery",
		"Design",
		"Developer Experience",
		"Engineering Management",
		"Go",
		"Infrastructure",
		"Leadership",
		"Maintainability",
		"Metrics",
		"Observability",
		"Operations",
		"Performance",
		"Product",
		"Reliability",
		"Software Design",
		"Testing",
		"Tooling",
	}
	tags := make([]models.TagEntity, 0, len(tagTitles))
	for _, title := range tagTitles {
		tag, err := factories.CreateTag(ctx, exec, factories.WithTagsTitle(title))
		if err != nil {
			return fmt.Errorf("failed to create tag %s: %w", title, err)
		}
		tags = append(tags, tag)
	}

	articles, err := models.Article.All(ctx, exec)
	if err != nil {
		return fmt.Errorf("failed to load articles for tag connections: %w", err)
	}
	for articleIndex, article := range articles {
		tagCount := 2 + articleIndex%3
		for offset := range tagCount {
			tag := tags[(articleIndex*3+offset*7)%len(tags)]
			if _, err := factories.CreateArticleTagConnection(
				ctx,
				exec,
				article.ID,
				tag.ID,
			); err != nil {
				return fmt.Errorf(
					"failed to attach tag %s to article %s: %w",
					tag.Title,
					article.Slug,
					err,
				)
			}
		}
	}

	if _, err := factories.CreateNewsletter(
		ctx,
		exec,
		factories.WithNewslettersTitle("Field Notes 01: Choosing Calm Technology"),
		factories.WithNewslettersSlug(
			sql.NullString{String: "field-notes-01-choosing-calm-technology", Valid: true},
		),
		factories.WithNewslettersMetaTitle("Field Notes 01: Choosing Calm Technology"),
		factories.WithNewslettersMetaDescription(
			"Notes on choosing tools that reduce operational load and support steady product work.",
		),
		factories.WithNewslettersIsPublished(sql.NullBool{Bool: true, Valid: true}),
		factories.WithNewslettersReleasedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -49), Valid: true},
		),
		factories.WithNewslettersContent(
			sql.NullString{
				String: "# Choosing Calm Technology\n\nChoose tools based on fit, ownership, and operational clarity.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create newsletter field-notes-01: %w", err)
	}
	if _, err := factories.CreateNewsletter(
		ctx,
		exec,
		factories.WithNewslettersTitle("Field Notes 02: Better Defaults"),
		factories.WithNewslettersSlug(
			sql.NullString{String: "field-notes-02-better-defaults", Valid: true},
		),
		factories.WithNewslettersMetaTitle("Field Notes 02: Better Defaults"),
		factories.WithNewslettersMetaDescription(
			"How thoughtful defaults remove repeated decisions and improve product consistency.",
		),
		factories.WithNewslettersIsPublished(sql.NullBool{Bool: true, Valid: true}),
		factories.WithNewslettersReleasedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -37), Valid: true},
		),
		factories.WithNewslettersContent(
			sql.NullString{
				String: "# Better Defaults\n\nGood defaults make the safe and common path the easiest one to follow.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create newsletter field-notes-02: %w", err)
	}
	if _, err := factories.CreateNewsletter(
		ctx,
		exec,
		factories.WithNewslettersTitle("Field Notes 03: Make State Visible"),
		factories.WithNewslettersSlug(
			sql.NullString{String: "field-notes-03-make-state-visible", Valid: true},
		),
		factories.WithNewslettersMetaTitle("Field Notes 03: Make State Visible"),
		factories.WithNewslettersMetaDescription(
			"A short guide to making important state legible in interfaces, logs, and workflows.",
		),
		factories.WithNewslettersIsPublished(sql.NullBool{Bool: true, Valid: true}),
		factories.WithNewslettersReleasedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -24), Valid: true},
		),
		factories.WithNewslettersContent(
			sql.NullString{
				String: "# Make State Visible\n\nVisibility reduces coordination cost and helps people make confident decisions.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create newsletter field-notes-03: %w", err)
	}
	if _, err := factories.CreateNewsletter(
		ctx,
		exec,
		factories.WithNewslettersTitle("Field Notes 04: The Value of Small Batches"),
		factories.WithNewslettersSlug(
			sql.NullString{String: "field-notes-04-small-batches", Valid: true},
		),
		factories.WithNewslettersMetaTitle("Field Notes 04: The Value of Small Batches"),
		factories.WithNewslettersMetaDescription(
			"Why smaller changes improve feedback, reduce risk, and keep work understandable.",
		),
		factories.WithNewslettersIsPublished(sql.NullBool{Bool: true, Valid: true}),
		factories.WithNewslettersReleasedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -12), Valid: true},
		),
		factories.WithNewslettersContent(
			sql.NullString{
				String: "# The Value of Small Batches\n\nSmall changes are easier to review, release, observe, and reverse.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create newsletter field-notes-04: %w", err)
	}
	if _, err := factories.CreateNewsletter(
		ctx,
		exec,
		factories.WithNewslettersTitle("Field Notes 05: Working With Uncertainty"),
		factories.WithNewslettersSlug(
			sql.NullString{String: "field-notes-05-working-with-uncertainty", Valid: true},
		),
		factories.WithNewslettersMetaTitle("Field Notes 05: Working With Uncertainty"),
		factories.WithNewslettersMetaDescription(
			"Draft notes on progress when requirements and technical constraints are still moving.",
		),
		factories.WithNewslettersIsPublished(sql.NullBool{Bool: false, Valid: true}),
		factories.WithNewslettersReleasedAt(sql.NullTime{}),
		factories.WithNewslettersContent(
			sql.NullString{
				String: "# Working With Uncertainty\n\nPrototypes and reversible decisions help teams learn without pretending uncertainty is gone.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create newsletter field-notes-05: %w", err)
	}
	if _, err := factories.CreateNewsletter(
		ctx,
		exec,
		factories.WithNewslettersTitle("Field Notes 06: A Useful Review Culture"),
		factories.WithNewslettersSlug(
			sql.NullString{String: "field-notes-06-useful-review-culture", Valid: true},
		),
		factories.WithNewslettersMetaTitle("Field Notes 06: A Useful Review Culture"),
		factories.WithNewslettersMetaDescription(
			"Draft ideas for reviews that improve work while sharing context across a team.",
		),
		factories.WithNewslettersIsPublished(sql.NullBool{Bool: false, Valid: true}),
		factories.WithNewslettersReleasedAt(sql.NullTime{}),
		factories.WithNewslettersContent(
			sql.NullString{
				String: "# A Useful Review Culture\n\nClear intent and curious questions make review a place for shared understanding.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create newsletter field-notes-06: %w", err)
	}

	if _, err := factories.CreateNewsletter(
		ctx,
		exec,
		factories.WithNewslettersTitle("Field Notes 07: Reading the System"),
		factories.WithNewslettersSlug(
			sql.NullString{String: "field-notes-07-reading-the-system", Valid: true},
		),
		factories.WithNewslettersMetaTitle("Field Notes 07: Reading the System"),
		factories.WithNewslettersMetaDescription(
			"Field notes on reading the system and the practical choices that help software teams work with greater clarity, confidence, and consistency.",
		),
		factories.WithNewslettersIsPublished(sql.NullBool{Bool: true, Valid: true}),
		factories.WithNewslettersReleasedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -68), Valid: true},
		),
		factories.WithNewslettersContent(
			sql.NullString{
				String: "# Reading the System\n\nSmall, explicit choices make this practice easier to understand, repeat, and improve across a team.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create newsletter field-notes-07: %w", err)
	}

	if _, err := factories.CreateNewsletter(
		ctx,
		exec,
		factories.WithNewslettersTitle("Field Notes 08: Defaults Shape Behavior"),
		factories.WithNewslettersSlug(
			sql.NullString{String: "field-notes-08-defaults-shape-behavior", Valid: true},
		),
		factories.WithNewslettersMetaTitle("Field Notes 08: Defaults Shape Behavior"),
		factories.WithNewslettersMetaDescription(
			"Field notes on defaults shape behavior and the practical choices that help software teams work with greater clarity, confidence, and consistency.",
		),
		factories.WithNewslettersIsPublished(sql.NullBool{Bool: true, Valid: true}),
		factories.WithNewslettersReleasedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -62), Valid: true},
		),
		factories.WithNewslettersContent(
			sql.NullString{
				String: "# Defaults Shape Behavior\n\nSmall, explicit choices make this practice easier to understand, repeat, and improve across a team.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create newsletter field-notes-08: %w", err)
	}

	if _, err := factories.CreateNewsletter(
		ctx,
		exec,
		factories.WithNewslettersTitle("Field Notes 09: Small Releases"),
		factories.WithNewslettersSlug(
			sql.NullString{String: "field-notes-09-small-releases", Valid: true},
		),
		factories.WithNewslettersMetaTitle("Field Notes 09: Small Releases"),
		factories.WithNewslettersMetaDescription(
			"Field notes on small releases and the practical choices that help software teams work with greater clarity, confidence, and consistency.",
		),
		factories.WithNewslettersIsPublished(sql.NullBool{Bool: true, Valid: true}),
		factories.WithNewslettersReleasedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -56), Valid: true},
		),
		factories.WithNewslettersContent(
			sql.NullString{
				String: "# Small Releases\n\nSmall, explicit choices make this practice easier to understand, repeat, and improve across a team.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create newsletter field-notes-09: %w", err)
	}

	if _, err := factories.CreateNewsletter(
		ctx,
		exec,
		factories.WithNewslettersTitle("Field Notes 10: Healthy Constraints"),
		factories.WithNewslettersSlug(
			sql.NullString{String: "field-notes-10-healthy-constraints", Valid: true},
		),
		factories.WithNewslettersMetaTitle("Field Notes 10: Healthy Constraints"),
		factories.WithNewslettersMetaDescription(
			"Field notes on healthy constraints and the practical choices that help software teams work with greater clarity, confidence, and consistency.",
		),
		factories.WithNewslettersIsPublished(sql.NullBool{Bool: true, Valid: true}),
		factories.WithNewslettersReleasedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -50), Valid: true},
		),
		factories.WithNewslettersContent(
			sql.NullString{
				String: "# Healthy Constraints\n\nSmall, explicit choices make this practice easier to understand, repeat, and improve across a team.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create newsletter field-notes-10: %w", err)
	}

	if _, err := factories.CreateNewsletter(
		ctx,
		exec,
		factories.WithNewslettersTitle("Field Notes 11: Operational Feedback"),
		factories.WithNewslettersSlug(
			sql.NullString{String: "field-notes-11-operational-feedback", Valid: true},
		),
		factories.WithNewslettersMetaTitle("Field Notes 11: Operational Feedback"),
		factories.WithNewslettersMetaDescription(
			"Field notes on operational feedback and the practical choices that help software teams work with greater clarity, confidence, and consistency.",
		),
		factories.WithNewslettersIsPublished(sql.NullBool{Bool: true, Valid: true}),
		factories.WithNewslettersReleasedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -44), Valid: true},
		),
		factories.WithNewslettersContent(
			sql.NullString{
				String: "# Operational Feedback\n\nSmall, explicit choices make this practice easier to understand, repeat, and improve across a team.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create newsletter field-notes-11: %w", err)
	}

	if _, err := factories.CreateNewsletter(
		ctx,
		exec,
		factories.WithNewslettersTitle("Field Notes 12: Clear Ownership"),
		factories.WithNewslettersSlug(
			sql.NullString{String: "field-notes-12-clear-ownership", Valid: true},
		),
		factories.WithNewslettersMetaTitle("Field Notes 12: Clear Ownership"),
		factories.WithNewslettersMetaDescription(
			"Field notes on clear ownership and the practical choices that help software teams work with greater clarity, confidence, and consistency.",
		),
		factories.WithNewslettersIsPublished(sql.NullBool{Bool: true, Valid: true}),
		factories.WithNewslettersReleasedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -38), Valid: true},
		),
		factories.WithNewslettersContent(
			sql.NullString{
				String: "# Clear Ownership\n\nSmall, explicit choices make this practice easier to understand, repeat, and improve across a team.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create newsletter field-notes-12: %w", err)
	}

	if _, err := factories.CreateNewsletter(
		ctx,
		exec,
		factories.WithNewslettersTitle("Field Notes 13: Maintenance Windows"),
		factories.WithNewslettersSlug(
			sql.NullString{String: "field-notes-13-maintenance-windows", Valid: true},
		),
		factories.WithNewslettersMetaTitle("Field Notes 13: Maintenance Windows"),
		factories.WithNewslettersMetaDescription(
			"Field notes on maintenance windows and the practical choices that help software teams work with greater clarity, confidence, and consistency.",
		),
		factories.WithNewslettersIsPublished(sql.NullBool{Bool: true, Valid: true}),
		factories.WithNewslettersReleasedAt(
			sql.NullTime{Time: now.AddDate(0, 0, -32), Valid: true},
		),
		factories.WithNewslettersContent(
			sql.NullString{
				String: "# Maintenance Windows\n\nSmall, explicit choices make this practice easier to understand, repeat, and improve across a team.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create newsletter field-notes-13: %w", err)
	}

	if _, err := factories.CreateNewsletter(
		ctx,
		exec,
		factories.WithNewslettersTitle("Field Notes 14: Useful Documentation"),
		factories.WithNewslettersSlug(
			sql.NullString{String: "field-notes-14-useful-documentation", Valid: true},
		),
		factories.WithNewslettersMetaTitle("Field Notes 14: Useful Documentation"),
		factories.WithNewslettersMetaDescription(
			"Field notes on useful documentation and the practical choices that help software teams work with greater clarity, confidence, and consistency.",
		),
		factories.WithNewslettersIsPublished(sql.NullBool{Bool: false, Valid: true}),
		factories.WithNewslettersReleasedAt(sql.NullTime{}),
		factories.WithNewslettersContent(
			sql.NullString{
				String: "# Useful Documentation\n\nSmall, explicit choices make this practice easier to understand, repeat, and improve across a team.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create newsletter field-notes-14: %w", err)
	}

	if _, err := factories.CreateNewsletter(
		ctx,
		exec,
		factories.WithNewslettersTitle("Field Notes 15: Making Work Visible"),
		factories.WithNewslettersSlug(
			sql.NullString{String: "field-notes-15-making-work-visible", Valid: true},
		),
		factories.WithNewslettersMetaTitle("Field Notes 15: Making Work Visible"),
		factories.WithNewslettersMetaDescription(
			"Field notes on making work visible and the practical choices that help software teams work with greater clarity, confidence, and consistency.",
		),
		factories.WithNewslettersIsPublished(sql.NullBool{Bool: false, Valid: true}),
		factories.WithNewslettersReleasedAt(sql.NullTime{}),
		factories.WithNewslettersContent(
			sql.NullString{
				String: "# Making Work Visible\n\nSmall, explicit choices make this practice easier to understand, repeat, and improve across a team.",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create newsletter field-notes-15: %w", err)
	}

	if _, err := factories.CreateProject(
		ctx,
		exec,
		factories.WithProjectsPublished(true),
		factories.WithProjectsTitle("Andurel Framework"),
		factories.WithProjectsSlug("andurel-framework"),
		factories.WithProjectsStartedAt(sql.NullTime{Time: now.AddDate(0, 0, -820), Valid: true}),
		factories.WithProjectsStatus("Active"),
		factories.WithProjectsDescription(
			"A productive Go framework for building explicit web applications with clear conventions, strong tooling, and calm operational defaults.",
		),
		factories.WithProjectsMetaTitle("Andurel Framework for Productive Go Web Development"),
		factories.WithProjectsMetaDescription(
			"Explore Andurel, a productive Go web framework with explicit conventions, integrated tooling, and practical defaults for maintainable applications.",
		),
		factories.WithProjectsContent(
			"# Andurel Framework\n\nAndurel brings routing, persistence, admin interfaces, jobs, email, and deployment conventions into one coherent Go workflow.",
		),
		factories.WithProjectsProjectUrl(
			sql.NullString{String: "https://example.com/projects/andurel", Valid: true},
		),
	); err != nil {
		return fmt.Errorf("failed to create project andurel-framework: %w", err)
	}

	if _, err := factories.CreateProject(
		ctx,
		exec,
		factories.WithProjectsPublished(true),
		factories.WithProjectsTitle("DeployCrate"),
		factories.WithProjectsSlug("deploycrate"),
		factories.WithProjectsStartedAt(sql.NullTime{Time: now.AddDate(0, 0, -640), Valid: true}),
		factories.WithProjectsStatus("Active"),
		factories.WithProjectsDescription(
			"A deployment platform focused on predictable releases, visible infrastructure state, and a smaller operational surface for product teams.",
		),
		factories.WithProjectsMetaTitle("DeployCrate for Predictable Application Deployments"),
		factories.WithProjectsMetaDescription(
			"Deploy applications with predictable releases, visible infrastructure state, and a focused operational workflow designed for small product teams.",
		),
		factories.WithProjectsContent(
			"# DeployCrate\n\nDeployCrate turns application delivery into a visible, repeatable workflow with explicit configuration and useful operational feedback.",
		),
		factories.WithProjectsProjectUrl(
			sql.NullString{String: "https://example.com/projects/deploycrate", Valid: true},
		),
	); err != nil {
		return fmt.Errorf("failed to create project deploycrate: %w", err)
	}

	if _, err := factories.CreateProject(
		ctx,
		exec,
		factories.WithProjectsPublished(true),
		factories.WithProjectsTitle("Shadowfax"),
		factories.WithProjectsSlug("shadowfax"),
		factories.WithProjectsStartedAt(sql.NullTime{Time: now.AddDate(0, 0, -510), Valid: true}),
		factories.WithProjectsStatus("Maintained"),
		factories.WithProjectsDescription(
			"A compact development tool that keeps common project workflows fast, discoverable, and consistent across local environments.",
		),
		factories.WithProjectsMetaTitle("Shadowfax Developer Workflows and Project Automation"),
		factories.WithProjectsMetaDescription(
			"Shadowfax keeps everyday development workflows fast and discoverable with focused project automation and consistent local tooling conventions.",
		),
		factories.WithProjectsContent(
			"# Shadowfax\n\nShadowfax packages the repetitive parts of local development into small commands that explain what they do and fail clearly.",
		),
		factories.WithProjectsProjectUrl(
			sql.NullString{String: "https://example.com/projects/shadowfax", Valid: true},
		),
	); err != nil {
		return fmt.Errorf("failed to create project shadowfax: %w", err)
	}

	if _, err := factories.CreateProject(
		ctx,
		exec,
		factories.WithProjectsPublished(true),
		factories.WithProjectsTitle("Admin Content Studio"),
		factories.WithProjectsSlug("admin-content-studio"),
		factories.WithProjectsStartedAt(sql.NullTime{Time: now.AddDate(0, 0, -370), Valid: true}),
		factories.WithProjectsStatus("Completed"),
		factories.WithProjectsDescription(
			"An editorial workspace for drafting, reviewing, publishing, and maintaining structured long-form content without unnecessary complexity.",
		),
		factories.WithProjectsMetaTitle("Admin Content Studio for Calm Editorial Workflows"),
		factories.WithProjectsMetaDescription(
			"A focused editorial workspace for drafting, reviewing, publishing, and maintaining structured content with clear validation and status visibility.",
		),
		factories.WithProjectsContent(
			"# Admin Content Studio\n\nThe studio keeps content, metadata, publication readiness, and destructive actions visible in one consistent editorial interface.",
		),
		factories.WithProjectsProjectUrl(
			sql.NullString{
				String: "https://example.com/projects/admin-content-studio",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create project admin-content-studio: %w", err)
	}

	if _, err := factories.CreateProject(
		ctx,
		exec,
		factories.WithProjectsPublished(true),
		factories.WithProjectsTitle("Go Service Blueprint"),
		factories.WithProjectsSlug("go-service-blueprint"),
		factories.WithProjectsStartedAt(sql.NullTime{Time: now.AddDate(0, 0, -280), Valid: true}),
		factories.WithProjectsStatus("Completed"),
		factories.WithProjectsDescription(
			"A reference architecture for observable Go services with explicit boundaries, useful failure modes, and production-friendly defaults.",
		),
		factories.WithProjectsMetaTitle("Go Service Blueprint for Observable Production Systems"),
		factories.WithProjectsMetaDescription(
			"A practical reference architecture for observable Go services with explicit boundaries, useful failure modes, and production-friendly defaults.",
		),
		factories.WithProjectsContent(
			"# Go Service Blueprint\n\nThis blueprint demonstrates service boundaries, structured errors, health checks, telemetry, and graceful lifecycle management.",
		),
		factories.WithProjectsProjectUrl(
			sql.NullString{
				String: "https://example.com/projects/go-service-blueprint",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create project go-service-blueprint: %w", err)
	}

	if _, err := factories.CreateProject(
		ctx,
		exec,
		factories.WithProjectsPublished(true),
		factories.WithProjectsTitle("Observability Toolkit"),
		factories.WithProjectsSlug("observability-toolkit"),
		factories.WithProjectsStartedAt(sql.NullTime{Time: now.AddDate(0, 0, -190), Valid: true}),
		factories.WithProjectsStatus("Active"),
		factories.WithProjectsDescription(
			"A practical toolkit for adding useful traces, metrics, and logs while keeping instrumentation focused on real operational decisions.",
		),
		factories.WithProjectsMetaTitle("Observability Toolkit for Useful Operational Signals"),
		factories.WithProjectsMetaDescription(
			"Add useful traces, metrics, and logs with an observability toolkit focused on real operational decisions instead of collecting unnecessary noise.",
		),
		factories.WithProjectsContent(
			"# Observability Toolkit\n\nThe toolkit favors a small set of actionable signals and consistent instrumentation patterns that teams can understand under pressure.",
		),
		factories.WithProjectsProjectUrl(
			sql.NullString{
				String: "https://example.com/projects/observability-toolkit",
				Valid:  true,
			},
		),
	); err != nil {
		return fmt.Errorf("failed to create project observability-toolkit: %w", err)
	}

	if _, err := factories.CreateProject(
		ctx,
		exec,
		factories.WithProjectsPublished(false),
		factories.WithProjectsTitle("Research Notebook"),
		factories.WithProjectsSlug("research-notebook"),
		factories.WithProjectsStartedAt(sql.NullTime{Time: now.AddDate(0, 0, -90), Valid: true}),
		factories.WithProjectsStatus("Planned"),
		factories.WithProjectsDescription(
			"An early concept for connecting technical research notes, experiments, and decisions in a searchable project record.",
		),
		factories.WithProjectsMetaTitle("Research Notebook for Connected Technical Decisions"),
		factories.WithProjectsMetaDescription(
			"An early project exploring how technical research notes, experiments, and decisions can become a searchable and connected engineering record.",
		),
		factories.WithProjectsContent(
			"# Research Notebook\n\nThis draft explores lightweight ways to connect questions, experiments, evidence, and architectural decisions.",
		),
		factories.WithProjectsProjectUrl(sql.NullString{}),
	); err != nil {
		return fmt.Errorf("failed to create project research-notebook: %w", err)
	}

	if _, err := factories.CreateProject(
		ctx,
		exec,
		factories.WithProjectsPublished(false),
		factories.WithProjectsTitle("Developer Portal"),
		factories.WithProjectsSlug("developer-portal"),
		factories.WithProjectsStartedAt(sql.NullTime{Time: now.AddDate(0, 0, -45), Valid: true}),
		factories.WithProjectsStatus("Planned"),
		factories.WithProjectsDescription(
			"A draft developer portal that brings service ownership, documentation, operational links, and common actions into one place.",
		),
		factories.WithProjectsMetaTitle("Developer Portal for Clear Service Ownership"),
		factories.WithProjectsMetaDescription(
			"A draft developer portal bringing service ownership, documentation, operational links, and common engineering actions into one coherent workspace.",
		),
		factories.WithProjectsContent(
			"# Developer Portal\n\nThe portal is being shaped around fast discovery, visible ownership, and direct access to the actions engineers use every day.",
		),
		factories.WithProjectsProjectUrl(sql.NullString{}),
	); err != nil {
		return fmt.Errorf("failed to create project developer-portal: %w", err)
	}

	fmt.Printf("Created admin user: %s\n", admin.Email)
	fmt.Printf("Created regular user: %s\n", user.Email)
	fmt.Println("Created 18 subscribers")
	fmt.Println("Created 30 articles")
	fmt.Println("Created 15 newsletters")
	fmt.Println("Created 8 projects")

	return nil
}
