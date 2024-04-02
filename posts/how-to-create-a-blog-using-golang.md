After getting fed up with React, SPAs, and Javascript around 2021 I decided to re-write my personal webpage in Rust and wrote an [article](/posts/how-to-build-a-simple-blog-using-rust) on how you could build a simple blog, purely using Rust. It ended up becoming one of my most popular articles and for good reason; Rust is exciting, fun to write, and blazingly fast. After a while though, I started to feel frustrated with the development process for adding new features to my site: the feedback loop was simply too long.

I've always been interested in solo entrepreneurship and technology. But, as I'm getting older, I realize that I might have been more interested in trying out new technologies. Great for learning and growing as an engineer, bad for shipping projects, and starting to see that sweet MRR grow. I decided to re-re-write my site once again, this time in Go for multiple reasons: 
- I write Go for a living
- Simple language, with a decent type system and fast performance
- Blazingly fast compile times

In this post, I will show how you can create your own personal blog using Go. I'll assume you're familiar with Go, and know how to configure a router/database/server, and so on. Should you not, feel free to grab a clone of my Go starter template [Grafto](https://github.com/mbv-labs/grafto) that has lots of things configured for you out of the box.

## Foundations

I write all my stuff in markdown; if you're a developer who also wants to start blogging, chances are you also are quite familiar with it. After Go 1.16 where we got the `embed` package included in the standard library most of our work is already done. We basically only need to have a way of storing some filenames, associating them with an ID or a slug, grabbing the file, and serving it to the user. Pretty simple.

Whether you've created your own setup or grabbed a copy of Grafto, create a new directory in the root of your project called `posts` and in there, create a file called `posts.go`. Open it and add the following:

```go
package posts

// imports omitted

//go:embed *.md
var Assets embed.FS
```

Now, any file with the `.md` extension will get included in the binary that ultimately gets built once we run `go run build`. We can simply grab the files from our global `Assets` variable using `Assets.ReadFile(name-of-file)` to handle any error that might occur or return the file as a string, e.g:

```go
file, err := Assets.ReadFile("my-post.md")
if err != nil {
	return err
}

return string(file)
```

It won't be pretty (we'll fix that later), but it gets the basic idea out, which we can build on.

We have a way to include our blog posts in the binary; let's add a way to associate them with a slug. You could use an ID here, but it just looks better to have URLs like: "https://acme.com/my-first-blog-post" compared to "http://acme.com/44a2530b-567d-472f-9495-e2ee64e7ae6d". So, assuming you have a database up and running, add this table:

```sql
create table posts (
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    title varchar(255) not null,
    filename varchar(255) not null,
    slug varchar(255) not null
);
```

Not much exciting going on here. Basically, for the unaware readers, the slug column above will be the title but URL-friendly. So if you have a post with the title "My First Blog Post" the slug equivalent would be "my-first-blog-post" easy.

Lastly, we need to be able to serve this to readers. That flow would typically involve them hitting a landing page showing a list of articles they can choose from, which links to the article.

The implementation of this will depend a bit on your setup, but let us implement the handler to deal with grabbing a specific article using Echo as our router. Assuming you have a route like `/posts/:slug`, create the following:

```go
type ArticleStorage interface {
	GetPostBySlug(slug string) (Post, error)
}

func ArticleHandler(ctx echo.Context, storage ArticleStorage) error {
	postSlug := ctx.Param("postSlug")
	
	postModel, err := storage.GetPostBySlug(postSlug)
	if err != nil {
		return err
	}
	
	postContent, err := posts.Assets.ReadFile(postModel.Filename)
	if err != nil {
		return err
	}

	return ctx.String(http.StatusOK, string(postContent))
}
```

I'm putting some decisions on your here in terms of implementing the ArticleStorage. We just need _something_ that grabs the data on the post from the DB, based on the slug.

This is the foundation of what we need...but it's not pretty let's fix that by letting the server do what it was always supposed to do: return HTML.

## Enter templ

If you've spent time in the Go ecosystem, chances are you've probably heard about [templ](https://templ.guide). It lets you write HTML templates as Go packages and it's just such a pleasant way of building out a UI. Add some HTMX and alpine.js and you've at least 95% of what you get with SPAs, with the added complexity.

It's good practice to have a base template that wraps around your other templates, so we have a single point for adding things like stylesheets, javascript, metadata, etc. Create a directory in root called views and add the following to a file called `base.templ`.

```go
package views

templ base() {
	<!DOCTYPE html>
	<html lang="en">
		<nav>
			<a href="/">MBV Labs</a>
		</nav>
		<body>
			{ children... }
			<footer>
				<aside>
					<p>Copyright ©2024 </p>
					<p>All right reserved by MBV Labs</p>
				</aside>
			</footer>
		</body>
	</html>
}
```

For this to work, you'll need to install templ and run `templ generate` which will produce a file called `base_templ.go` that we can then import into other templates to wrap around them. For the sake of brevity, we'll only create the template to show the actual article. Create a file called `article.templ`, and add the following:

```go
type ArticlePageData struct {
	Title             string
	Content           string
}

templ ArticlePage(data ArticlePageData) {
	@layouts.Base() {
		<div>
			<div>
				<h1>{ data.Title }</h1>
			</div>
			<article>
				@unsafe(data.Content)
			</article>
		</div>
	}
}
```

and run `templ generate` once again.

We can now go back and update our ArticleHandler handler:

```go
func ArticleHandler(ctx echo.Context, storage ArticleStorage) error {
	postSlug := ctx.Param("postSlug")
	
	postModel, err := storage.GetPostBySlug(postSlug)
	if err != nil {
		return err
	}
	
	postContent, err := posts.Assets.ReadFile(postModel.Filename)
	if err != nil {
		return err
	}

	return views.ArticlePage(
		views.ArticlePagedata{
			Title: postModel.Title,
			Content: postContent,
		},
	).Render(
		ctx.Request().Context(), 
		ctx.Response().Writer,
	)
}
```

If you run your application now and visit a valid URL, you should see a (rather ugly) page showing the markdown of your article but this time, with some sweet hypertext markup.

## Making things (slightly) less ugly

In terms of styling things, throwing some tailwind or vanilla CSS at what we have now will get you a long way. But, we still show raw markdown to the user when they visit our articles. Additionally, we might want to show some nicely formatted code snippets in our articles. Let's fix this now.

For this, we need something that can transform the markdown into HTML components e.g

```markdown
## Some sub header
```

into

```html
<h2>Some sub header</h2>
```

Luckily, there already is a create library for this: Goldmark. So let's refactor the `posts/posts.go` file to parse the content we store using embed.

```go
//go:embed *.md
var assets embed.FS // unexport assets

type Manager struct {
	posts embed.FS
	markdownParser goldmark.Markdown
}

func NewManager() Manager {
	md := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)

	return Manager{
		posts:           assets,
		markdownHandler: md,
	}
}

func (m *Manager) Parse(name string) (string, error) {
	source, err := m.posts.ReadFile(name)
	if err != nil {
		return "", err
	}

	// Parse Markdown content
	var htmlOutput bytes.Buffer
	if err := m.markdownHandler.Convert(source, &htmlOutput); err != nil {
		return "", err
	}

	return htmlOutput.String(), nil
}
```

Lastly, update the ArticleHandler to use the Manager:

```go
func ArticleHandler(
	ctx echo.Context, 
	storage ArticleStorage,
	postManager posts.Manager
) error {
	postSlug := ctx.Param("postSlug")
	
	postModel, err := storage.GetPostBySlug(postSlug)
	if err != nil {
		return err
	}
	
	postContent, err := postManager.Parse(postModel.Filename)
	if err != nil {
		return err
	}
	

	return views.ArticlePage(
		views.ArticlePagedata{
			Title: postModel.Title,
			Content: postContent,
		},
	).Render(
		ctx.Request().Context(), 
		ctx.Response().Writer,
	)
}
```

Try and edit your article by adding some code blocks and they will now get nicely formatted. You can add custom themes to the parser so your code snippets will be shown with your favorite theme.
