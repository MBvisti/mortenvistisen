# Public Pages Implementation Plan

Status: **COMPLETED**

## Goal

Implement the remaining public surface using the approved retro-futuristic design:

- `GET /posts`
- `GET /posts/:slug`
- `GET /newsletters`
- `GET /newsletters/:slug`
- `GET /projects`
- `GET /projects/:slug`
- `GET /about`

The `/about` page must link to [MBV Labs](https://mbvlabs.com) and explain that visitors who want to work with Morten can do so through MBV Labs.

## Current Design State

The landing page design is approved.

Preserve these decisions:

- Dark background using the existing warm ivory, orange, and green palette.
- Sharp rectangular geometry.
- No rounded public cards.
- No structural shadows or gradients.
- Large uppercase display type with safe wrapping.
- Small uppercase monospace labels.
- Header sits directly on the page canvas with no border, separate background, or sticky behavior.
- Header contains only the name and navigation. No logo boxes and no online indicator.
- Hero has no portrait, no box logo, and no “Available online” label.
- Footer headline uses the smaller card-heading scale.
- Long headings must not overflow their panels.

## Required Breadcrumbs

Every index and show page must start with a visible breadcrumb inspired by the provided MBV Labs references.

Examples:

```text
HOME / POSTS
HOME / POSTS / THE MAINTENANCE MINDSET
HOME / NEWSLETTERS
HOME / NEWSLETTERS / FIELD NOTES 04
HOME / PROJECTS
HOME / PROJECTS / ANDUREL
HOME / ABOUT
```

Requirements:

- Semantic `<nav aria-label="Breadcrumb">`.
- Uppercase monospace typography.
- Linked ancestor items.
- Final item uses `aria-current="page"`.
- Long breadcrumb labels wrap without horizontal overflow.
- Breadcrumb appears before the page title.

## Mandatory Andurel CLI Workflow

The missing controllers, actions, routes, and views must be generated with the Andurel CLI before customization.

### Discovery

Run before mutation:

```bash
andurel --agent --help
andurel project info --json
andurel commands --json
andurel generate controller --help
```

### Controller generation

Run each resource sequentially. Do not run generator previews in parallel because all controller generators touch `controllers/controller.go`.

```bash
andurel generate controller Post index show --model-name=Article --dry-run --diff --json
andurel generate controller Post index show --model-name=Article --json

andurel generate controller Newsletter index show --dry-run --diff --json
andurel generate controller Newsletter index show --json

andurel generate controller Project index show --dry-run --diff --json
andurel generate controller Project index show --json

andurel generate controller Page about --dry-run --diff --json
andurel generate controller Page about --json
```

Inspect the structured output after every preview and application:

- `files_created`
- `files_updated`
- `routes_added`
- `warnings`
- Any skipped or command artifact arrays

### Known CLI dry-run defect

Observed during planning with the currently installed CLI:

- `andurel generate controller ... --dry-run` removes the previewed resource files as expected.
- It incorrectly leaves constructor and invoke registrations in `controllers/controller.go`.
- A following generator then fails because those registered controller types do not exist.

Safe sequence for each resource:

1. Run one dry-run.
2. Inspect its output.
3. Remove only the leaked preview registration from `controllers/controller.go` using the edit tool.
4. Confirm no preview files remain.
5. Run the real generator.
6. Continue to the next resource.

The planning pass restored `controllers/controller.go`. It currently has no leaked public controller registration and no preview-generated resource files.

### About generator warning

The `Page about` dry-run showed that the generator wants to rewrite the existing root route from `HomePage` at `/` to a generated `/pages/home` route.

After applying the generator:

- Preserve the existing `HomePage` route at `/`.
- Add an `AboutPage` route at `/about`.
- Keep the current landing action and route intact.
- Customize the generated About action and Templ view only after generation.

## Generated Baseline

Expected resource artifacts:

```text
controllers/posts.go
router/routes/posts.go
views/posts_resource.templ

controllers/newsletters.go
router/routes/newsletters.go
views/newsletters_resource.templ

controllers/projects.go
router/routes/projects.go
views/projects_resource.templ

views/pages_resource.templ
controllers/pages.go             # updated by Page about generator
router/routes/pages.go           # updated by Page about generator
controllers/controller.go        # public controller registration
```

The generator previews confirmed that the default output:

- Uses integer `:id` show routes.
- Calls the existing admin-oriented `Paginate` methods, which include drafts.
- Calls `Find` by integer ID.
- Generates admin-like table and field-detail markup.

These generated files are the required starting point, not the final public implementation.

## Final Route Definitions

Customize generated routes to this contract:

```go
const PostPrefix = "/posts"

var PostIndex = routing.NewSimpleRoute(
    "",
    "posts.index",
    PostPrefix,
)

var PostShow = routing.NewRouteWithSlug(
    "/:slug",
    "posts.show",
    PostPrefix,
)
```

Apply the same pattern to newsletters and projects.

Final route constants:

- `routes.PostIndex`
- `routes.PostShow`
- `routes.NewsletterIndex`
- `routes.NewsletterShow`
- `routes.ProjectIndex`
- `routes.ProjectShow`
- `routes.HomePage`
- `routes.AboutPage`

After route changes, run:

```bash
andurel routes --json
andurel generate routes --json
```

This project has an Inertia route manifest, so `resources/js/routes.ts` must remain synchronized even though the public pages use Templ.

## Public Model Queries

Do not change the existing `Find`, `All`, or `Paginate` behavior because admin screens need drafts.

Add explicit public methods to each model.

### Articles

```go
func (a article) PaginatePublished(
    ctx context.Context,
    db storage.Executor,
    page, pageSize int64,
) (PaginatedArticles, error)

func (a article) FindPublishedBySlug(
    ctx context.Context,
    db storage.Executor,
    slug string,
) (ArticleEntity, error)
```

Rules:

- `published = TRUE`
- `first_published_at IS NOT NULL`
- Order by `first_published_at DESC`
- Missing or unpublished slug returns `models.ErrNotFound`

### Newsletters

```go
func (n newsletter) PaginatePublished(
    ctx context.Context,
    db storage.Executor,
    page, pageSize int64,
) (PaginatedNewsletters, error)

func (n newsletter) FindPublishedBySlug(
    ctx context.Context,
    db storage.Executor,
    slug string,
) (NewsletterEntity, error)
```

Rules:

- `is_published = TRUE`
- `released_at IS NOT NULL`
- Order by `released_at DESC`
- Missing or unpublished slug returns `models.ErrNotFound`

### Projects

```go
func (p project) PaginatePublished(
    ctx context.Context,
    db storage.Executor,
    page, pageSize int64,
) (PaginatedProjects, error)

func (p project) FindPublishedBySlug(
    ctx context.Context,
    db storage.Executor,
    slug string,
) (ProjectEntity, error)
```

Rules:

- `published = TRUE`
- `started_at IS NOT NULL`
- Order by `started_at DESC`
- Missing or unpublished slug returns `models.ErrNotFound`

### Pagination behavior

- Fixed page size: 12.
- Accept only the `page` query parameter.
- Invalid or non-positive page values fall back to page 1.
- Do not expose generated `per_page` behavior.
- Clamp pages using the same safe behavior as existing pagination methods.
- Return total count and total pages to the views.

## Newsletter Slug Constraint

Articles and projects already have unique database slugs. Newsletters do not.

Before relying on `/newsletters/:slug`, add a unique newsletter slug index.

Command discovery confirmed that this Andurel version has no migration generator. Follow the required manual fallback:

1. Run `andurel migrations --json`.
2. State that no migration generator exists.
3. Add the next migration following repository conventions.

Expected file:

```text
database/migrations/00019_make_newsletter_slug_unique.sql
```

Expected SQL:

```sql
-- +goose Up
CREATE UNIQUE INDEX newsletters_slug_unique
ON newsletters (slug)
WHERE slug IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS newsletters_slug_unique;
```

Check for duplicate non-null newsletter slugs before applying the migration.

## Controller Customization

Keep controllers as HTTP adapters.

### Index actions

Each generated index action should:

1. Parse the `page` query parameter.
2. Call the resource’s `PaginatePublished` method with page size 12.
3. Propagate unexpected database errors to existing error handling.
4. Render the generated resource index view with items and pagination metadata.

### Show actions

Each generated show action should:

1. Read `etx.Param("slug")`.
2. Call `FindPublishedBySlug`.
3. Render the public not-found page with HTTP 404 for `models.ErrNotFound`.
4. Propagate unexpected errors instead of presenting them as not found.
5. Render the generated resource show view.

Draft content must return 404 even when its slug exists.

## Shared View Primitives

Add one shared Templ source for reused public-page presentation:

```text
views/public_components.templ
```

Include:

- Breadcrumb component.
- Pagination component.
- Any small public date or metadata formatting helpers that are genuinely reused.

Suggested breadcrumb shape:

```go
type PublicBreadcrumb struct {
    Label string
    URL   string
}
```

The last breadcrumb is rendered as current text rather than a link.

Pagination should:

- Render only when `TotalPages > 1`.
- Include previous and next links when valid.
- Show the current page and total page count.
- Build query strings with `net/url`, not string concatenation.

## Markdown Rendering

Article, newsletter, and project content is stored as Markdown. The public Templ pages need server-rendered HTML.

Use `github.com/yuin/goldmark` as the single new Go dependency.

Add a presentation helper in:

```text
views/markdown.go
```

Rules:

- Goldmark raw HTML remains disabled.
- Convert Markdown to HTML on the server.
- Only pass Goldmark’s safe renderer output to `templ.Raw`.
- Never pass database content directly to `templ.Raw`.
- Conversion errors must propagate.

Add a `.public-prose` style to `css/base.css` covering:

- Paragraphs.
- `h2`, `h3`, and `h4`.
- Ordered and unordered lists.
- Links and visible focus states.
- Blockquotes.
- Inline code.
- Fenced code blocks with internal horizontal scrolling.
- Images constrained to the content width.

Keep the prose column narrow enough for comfortable reading.

## Page Designs

### Index page structure

All three index pages follow this pattern:

```text
breadcrumb
large editorial title
short introduction
responsive card/list grid
pagination when needed
```

The introduction should sit directly on the page canvas rather than inside a bordered panel.

#### `/posts`

Each item shows:

- Publication date.
- Read time when available.
- Title.
- Excerpt.
- Link to `routes.PostShow.URL(article.Slug)`.

Suggested heading: `Writing from the workshop.`

#### `/newsletters`

Each item shows:

- Release date.
- Title.
- Meta description.
- Link to `routes.NewsletterShow.URL(newsletter.Slug.String)`.

Suggested heading: `Field notes and dispatches.`

#### `/projects`

Each item shows:

- Start date.
- Status.
- Title.
- Description.
- Link to `routes.ProjectShow.URL(project.Slug)`.

Suggested heading: `Products and experiments.`

### Show page structure

All show pages follow this pattern:

```text
breadcrumb with current title
large title
resource metadata
optional image
narrow Markdown content column
resource-specific CTA when applicable
```

#### `/posts/:slug`

Use:

- `MetaTitle` for page title when available.
- `MetaDescription` for description.
- `ImageLink` and title-derived alt text when available.
- `FirstPublishedAt`.
- `ReadTime`.
- Rendered `Content`.
- `MetaType` set to `article`.

#### `/newsletters/:slug`

Use:

- `MetaTitle`.
- `MetaDescription`.
- `ImageLink` and title-derived alt text.
- `ReleasedAt`.
- Rendered `Content`.
- `MetaType` set to `article`.

#### `/projects/:slug`

Use:

- `MetaTitle`.
- `MetaDescription`.
- `StartedAt`.
- `Status`.
- `Description`.
- Rendered `Content`.
- A clear external CTA when `ProjectUrl.Valid`.

### `/about`

The generated About page should be customized into a personal, concise page rather than a full resume.

Include:

- Breadcrumb: `HOME / ABOUT`.
- Heading: `Engineer, freelancer, bootstrapper.`
- Short introduction as a Danish software engineer living in Spain.
- Current focus on Go, practical systems, products, and writing.
- Relevant social/profile links.
- A strong work-with-me panel linking to `https://mbvlabs.com`.

Required work CTA copy direction:

> I work with companies through MBV Labs. Visit mbvlabs.com to learn how we can work together.

## SEO and Head Data

Extend `views/head.templ` only as needed.

Add a `SetMetaType` option for article-like detail pages.

Each page should set:

- Title.
- Description.
- Canonical slug/path.
- Image and alt text when available.
- Open Graph type.

Index paths:

- `/posts`
- `/newsletters`
- `/projects`
- `/about`

Show canonical paths use the public slug route.

## Navigation and Landing Integration

After route constants exist:

### `views/layout.templ`

Replace temporary anchor links with route helpers:

- Writing -> `routes.PostIndex.URL()`
- Newsletters -> `routes.NewsletterIndex.URL()`
- Projects -> `routes.ProjectIndex.URL()`
- About -> `routes.AboutPage.URL()`

Apply the same changes to footer navigation.

### `views/landing.templ`

Make landing content cards link to internal show routes:

- Article -> `routes.PostShow.URL(article.Slug)`
- Newsletter -> `routes.NewsletterShow.URL(newsletter.Slug.String)`
- Project -> `routes.ProjectShow.URL(project.Slug)`

The project’s external URL belongs on its show page, not as the landing card destination.

### Skip link

The layout’s skip link currently targets the landing container. Change all public pages to use a shared `id="main-content"` target so the skip link works everywhere.

## Error and Accessibility Requirements

- Unknown slug returns HTTP 404.
- Draft slug returns HTTP 404.
- Unexpected database errors are not converted to 404.
- Breadcrumb uses semantic navigation.
- Current breadcrumb has `aria-current="page"`.
- All links have visible keyboard focus.
- External links using a new tab include `rel="noopener noreferrer"`.
- Images have meaningful alt text.
- Heading order remains valid.
- Long titles use safe wrapping.
- Code blocks scroll internally instead of overflowing the page.
- Layout works at 360, 768, 1280, and 1440 pixel widths.

## Implementation Phases

### Phase 1: Generate all missing surfaces (COMPLETED)

- Repeat CLI discovery.
- Preview and apply Post index/show generation.
- Preview and apply Newsletter index/show generation.
- Preview and apply Project index/show generation.
- Preview and apply Page about generation.
- Handle the dry-run registration leak after each preview.
- Preserve the root route after Page generation.

### Phase 2: Add public data contracts (COMPLETED)

- Add published pagination methods.
- Add published slug finders.
- Add newsletter slug uniqueness migration.
- Keep admin methods unchanged.

### Phase 3: Add shared public rendering (COMPLETED)

- Add breadcrumbs.
- Add pagination.
- Add Goldmark rendering helper.
- Add public prose styling.

### Phase 4: Build index pages (COMPLETED)

- Replace generated tables with editorial index layouts.
- Include breadcrumbs and pagination.
- Add empty states.

### Phase 5: Build show pages and About (COMPLETED)

- Replace generated field-detail layouts.
- Add metadata and Markdown content.
- Add project URL CTA.
- Add MBV Labs work CTA to About.

### Phase 6: Connect and verify (COMPLETED)

- Update shared navigation.
- Link landing cards.
- Regenerate Templ.
- Regenerate TypeScript route helpers.
- Compile CSS.
- Run static checks.

## Expected File Map

Generated and customized:

```text
controllers/posts.go
controllers/newsletters.go
controllers/projects.go
controllers/pages.go
controllers/controller.go

router/routes/posts.go
router/routes/newsletters.go
router/routes/projects.go
router/routes/pages.go

views/posts_resource.templ
views/newsletters_resource.templ
views/projects_resource.templ
views/pages_resource.templ
```

Manually added or extended after generation:

```text
views/public_components.templ
views/markdown.go
views/layout.templ
views/landing.templ
views/head.templ

models/article.go
models/newsletter.go
models/project.go

database/migrations/00019_make_newsletter_slug_unique.sql

css/base.css
assets/css/style.css
resources/js/routes.ts
```

Generated output must also be refreshed for all changed `.templ` files.

## Verification

Do not run the project runtime unless explicitly requested.

Run:

```bash
andurel generate view --json
andurel routes --json
andurel generate routes --json
./bin/tailwindcli -i ./css/base.css -o ./assets/css/style.css --minify
gofmt -w <changed Go and generated Go files>
go vet ./...
andurel views --json
andurel doctor --json
```

Inspect that:

- All seven required public routes exist.
- Show routes use `:slug`, never `:id`.
- TypeScript route helpers match the Go route manifest.
- Every `.templ` source has current generated Go output.
- No preview-generated files or leaked controller registrations remain.
- No draft query is used by a public controller.
- No headings or prose overflow.

The repository currently has an existing non-blocking Andurel version mismatch warning between `andurel.lock` and the installed CLI.

## Intentionally Deferred

Do not add during this implementation:

- Tags on public pages.
- Search.
- RSS.
- Sitemap expansion.
- Related content.
- Theme switching.
- Newsletter signup.
- New JavaScript navigation.

Add these only through separate requests.
