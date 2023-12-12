So, you have been reading about rust for quite some time, learned that it’s been the most loved programming language multiple years in a row in Stack Overflow’s annual survey and you want in. But, you’re not a system-level programmer and Reddit keeps telling you to ‘use the right tool for the job’ and ‘rust is more aimed at low-level stuff like writing databases and programming drone’ (the first is actually good advice, I just need it to create a dramatic setting). So you get disheartened and choose to use what you’ve always been using. I think that’s a mistake as rust is an amazing language both in terms of speed of the execution but also developer experience. Not many languages give you a handy pair-programmer, i.e. the compiler, that only gives good advice, and shuts up in the meantime (rust’s compiler is amazing). I’ve been toying with Rust for some time now and overall it’s been a great experience (both learning and developer). So to get my hands a bit more dirty I decided to re-write my blog using Rust which was previously built using Next.js (This has also been a thing coming since I’m doing my best to move away from SPAs and Javascript, but that’s for another post).

The goal of the rewrite is to get _something_ out which could serve as a learning example and be built upon to further my understanding of the language. I hope this post does the same for you.

_Side note_: I write Go for a living and have been coding professionally for the last ~5 years, so while I’m confident in my overall skills this tutorial might not always show idiomatic Rust and best practices. There’re a bunch of much more rust-abled developers than me who have written some great material, much of which lay the foundation for what I’m about to show you. At the end of the article, you’ll find a resource list that you can check out after reading (and hopefully coding along) this article. If you’re just interested in seeing the code you can find the repo here.

---

<div class="max-w-6xl py-10 px-4 sm:px-6 lg:px-8 lg:py-16 mx-auto">
  <div class="max-w-xl text-center mx-auto">
    <div class="mb-5">
      <h3 class="text-2xl font-bold md:text-3xl md:leading-tight text-white">Consider subscribing to my newsletter</h2>
    </div>
    <form hx-post="/subscribe" hx-target="this" hx-swap="outerHTML" method="POST" action="/subscribe" >
      <div class="mt-5 lg:mt-8 flex flex-col items-center gap-2 sm:flex-row sm:gap-3">
        <input type="text" id="hero-input" name="hero-input" class="py-3 px-4 block w-full border-gray-200 rounded-lg 
            text-sm focus:border-blue-500 focus:ring-blue-500 disabled:opacity-50 disabled:pointer-events-none 
            bg-slate-900 border-gray-700 text-gray-400 focus:ring-gray-600" placeholder="Enter your email">
        <button type="submit" class="w-full sm:w-auto whitespace-nowrap py-3 px-4 inline-flex justify-center 
            items-center gap-x-2 text-sm font-semibold rounded-lg border bg-slate-600 
            text-white hover:bg-slate-900 disabled:opacity-50 disabled:pointer-events-none focus:outline-none 
            focus:ring-1 focus:ring-gray-600">
          Subscribe
        </button>
      </div>
    </form>
  </div>
</div>

---

## Kicking things off
So, to build your new blog using Rust, you're going to have Rust installed on your system and then run `cargo new awesome-blog`.
Then run `cd awesome-blog` into the directory, open `Cargo.toml` and add the following under `[dependencies]`:
```toml
[dependencies]
tera = "1"
actix-web = "4"
env_logger = "0.9.0"
lazy_static = "1.4.0"
actix-files = "0.6.1"
pulldown-cmark = { version = "0.9.1", default-features = false }
ignore = "0.4"
toml = "0.5"
serde = { version = "1.0", features = ["derive"] }
serde_json = "1.0
```
Feel free to use another version than the ones listed above; it's what I used at the time of writing this piece.

Let me briefly touch upon some of these libraries. The main one here is actix-web which is one of the big players in Rust web frameworks. 
There are multiple other choices (see for example [warp](https://github.com/seanmonstar/warp), [rocket](https://github.com/SergioBenitez/Rocket), [axum](https://github.com/tokio-rs/axum)) but I thoroughly enjoy `actix`'s API and have built multiple projects using it. Furthermore, 
it performs really well and can handle _a lot of_ requests.

Next is `serde` which is a library you most likely are familiar with if you've done any Rust programming. If not, it's basically a set of 
methods for serializing and deserializing data structures efficiently and generically. If you can't find an implementation for your desired data structure, 
you simply implement [the methods](https://serde.rs/custom-serialization.html) for your data structure and you're off to the races.
For this project, `serde` will be used for serializing the metadata of a blog post (also named front matter) from a `.toml` file into a Rust struct.

Lastly, I want to touch upon `tera` which is the template engine I use for this blog. It's inspired by Jinja2 which you might be familiar with if 
you've ever done any web development with `Django`. Back in the day, I started out building web apps using `Django` and really enjoyed the way it does 
templating so wanted to replicate that experience in my personal projects. In a world where everything is components and the slightest duplication is a
death sin, the experience of using "old school" templating has been a relief.

## A simple server
Right out of the box, after running `cargo new awesome-blog`, you should be able to open up `src/main.rs` and be presented with something like this:

```rust
  fn main() {
	println!("Hello, world!");
  }
```

Executing the command: `cargo run` will (are pulling all the dependencies) output the classic `Hello, world!` in your terminal.
We're going to treat `main.rs` as a thin entry point into our application, which basically just means it should do some high-level configs 
and call start-up functions. With that requirement, we end up with something like this:

```rust
use std::net::TcpListener;

use awesome_blog::start;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
	std::env::set_var("RUST_LOG", "actix_web=info");
	env_logger::init();
	
	let listener = TcpListener::bind("0.0.0.0:8080")?;
	start_blog(listener)?.await?;
	Ok(())
}
```

Not much going on here, except a bunch of errors, so let's break it down real quick: 

- We set an env variable `RUST_LOG` that determines what kind of log statements we're getting out
- We initialize a simple logger that is configured using env variables
- We create a listener that listens on port `8080`
- We pass the listener to the `start_blog` and we're off to the races

We will deal with `start_blog` later, which means that the only thing left that hasn’t been addressed is `#[axtic_web::main]`.
This is macro, or more precisely, a `proc_macro` which generates some code for us that makes our `awesome_blog` run asynchronously.
The `rust` team made a conscious decision not to provide an async runtime in the standard library but instead opted to provide essentials and let the community do the rest.
The reason(s) why are manyfold, but one reason could be that there is no consensus on what such runtime should [look like](https://www.reddit.com/r/rust/comments/ui7ayd/comment/i7akj6d/?utm_source=share&utm_medium=web2x&context=3).
That topic is also out of scope for this article; you just have to know that it generates code for us, which makes an async runtime possible.

## Peeling off the layers
We need to tackle the `start_blog` function, so create a new file called `lib.rs` under `src` and open it up.
We will be building it in steps, starting with adding the code needed to run an `HttpServer`:

```rust
use std::net::TcpListener;
use actix_web::{dev::Server, web, App, HttpResponse, HttpServer, middleware};

pub fn start_blog(listener: TcpListener) -> Result<Server, std::io::Error> {
	let srv = HttpServer::new(move || {
	    App::new()
		   .wrap(middleware::Logger::default()) // enable logger
		   .route("/health", web::get().to(HttpResponse::Ok))
	})
	.listen(listener)?
	.run();
	
	Ok(srv)
}
```

So, we did a little more than just create an `HttpServer`. We also added in the logger defined in `main.rs` and created a route that 
pings back a `200 Ok` response. 

Most of the above code should make sense to you (assuming familiarity with `rust`) but let's quickly touch upon the `move` keyword.
According to the docs: a `move` converts any variables captured by reference or mutable reference to variables captured by value.
If you're a bit confused about the above statement, join the club. But, it can be broken down into a more (to me, at least) 
understandable sentence. What it boils down to is `rust`'s ownership model and how it deals with memory (de)allocation and references.
What we see in this line `HttpServer::new( move || { App::new() })` is a closure and closures might [escape](https://huonw.github.io/blog/2015/05/finding-closure-in-rust/).
To (over)simply, `actix_web` will spin up multiple instances of your app, assuming a multi-thread environment, so variables/values passed to `App::new()` might outlive `App`, i.e. the closure might escape.
That would leave dangling references laying around that are not accounted far and that's no bueno in `rust`.
So, to get around that, we tell `App` to take ownership of the values passed to it so that each instance owns what is passed to it.

Even though we're going a bit "olds chool" here and using pure simple html we would still like to have some reusability.
And this is exactly where `tera` comes in. However, since the templates will remain the same once we chuck them into prod, it will probably be a good idea only to load them one time and then pass a reference.
To do this, we're going to be using `lazy_static!` so go ahead and add the following to your `lib.rs`:

```rust
use std::net::TcpListener;
use actix_web::{dev::Server, web, App, HttpResponse, HttpServer};
use tera::Tera;

#[macro_use]
extern crate lazy_static;

lazy_static! {
	pub static ref TEMPLATES: Tera = {
	    let mut tera = match Tera::new("templates/**/*.html") {
			Ok(t) => t,
			Err(e) => {
				println!("Parsing error(s): {}", e);
				::std::process::exit(1);
			}
		};
		tera.autoescape_on(vec![".html", ".sql"]);
		tera
	};
}

pub fn start_blog(listener: TcpListener) -> Result<Server, std::io::Error> {
	let srv = HttpServer::new(move || {
	    App::new()
		   .app_data(web::Data::new(TEMPLATES.clone()))
		   .wrap(middleware::Logger::default()) // enable logger
		   .route("/health", web::get().to(HttpResponse::Ok))
	})
	.listen(listener)?
	.run();
	
	Ok(srv)
}
```

All good, but one thing missing. We need to add a templates directory in the root directory so go 
ahead and `mkdir templates`. We will get back to `main.rs`, `lib.rs` and all that soon but for now,
let's quickly touch upon the templates so we have something to show.

## Dude, where are the components?
I started out with `react` when I got my first job as a front-end developer. Components were all the rage (well, they still 
are but this makes me sound a bit older and wiser) back then and it was how I started thinking about UI elements. I've since 
changed my mind quite a bit on this but that will be for another article. As you know, we're not going to be doing components here 
but `partials` and `blocks` so in your `templates` directory, create a file called `base.html`

```html
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    {% block head %}
    <title>{% block title %}{% endblock title %}</title>
    {% endblock head %}
    <link rel="stylesheet" type="text/css" href="/static/css/index.css">
</head>

<body class="flex flex-col justify-between min-h-screen font-sans leading-normal tracking-normal">
    <div class="h-16">{% block header %}{% endblock header %}</div>

    <main class="container flex-1 w-full md:max-w-3xl mx-auto overflow-x-hidden">
        {% block content %}{% endblock content %}
    </main>
    {% block footer %}
    <div class="flex w-full h-20 justify-center items-center">
        &copy; 2022 by <a class="pl-2" href="https://awesomeblog.com/">Awesome Blog</a>.
    </div>
    {% endblock footer %}
</body>

<script src="/static/js/highlight.min.js"></script>
<script>hljs.highlightAll();</script>
</script>

</html>
```

The above is going to be the base for everything else we do.

The part that's going to let us extend `base.html` (or, inherit..) is this `{% block content %}{% endblock content %}`.
Let's see how that is done, so go ahead and create a new file under `templates` called `home.html`:

```html
{% extends "base.html" %}
{% block title %}Awesome blog | by DeveloperMan{% endblock title %}

{% block content %}
<div class="w-full h-full px-4 md:px-6 text-xl text-gray-800 leading-normal">
    <h1 class="text-center text-3xl">All about tech</h1>
    <h3 class="text-center text-lg mt-2 text-gray-600">More words describing the blog here</h3>
    {% for fm in posts %}
    <a class="mb-3" href="/posts/{{fm.file_name}}">
        <div class="flex flex-col mb-5 px-4 py-6 cursor-pointer">
            <h2 class="text-2xl font-semibold hover:underline">{{fm.title}}</h2>
            <div class="flex items-center mt-2">
                <p class="hidden md:block text-base ml-2 mr-2">{{fm.posted}}</p>
                <p class="hidden md:block text-base mx-2">{{fm.author}}</p>
                <p class="text-base ml-0 md:ml-2">Reading time:
                {{fm.estimated_reading_time}} min</p>
            </div>
        </div>
    </a>
    {% endfor %}
</div>
{% endblock content %}
```

First off, notice the `{% extends "base.html" %}` at the top. This lets us interact with all of the `blocks` we created in the `base.html` template. 
Furthermore, we have a for loop in the markup: `{% for frontmatter in posts %}`. We can pass data as variables when serving the html as part of 
the response to, say, `http://awesomeblog.com/` which in the above case would return a variable named `posts` that is an array of posts. 

If you try and spin up the application now, not much would be different, so let's implement some handlers/controllers/etc to actually serve some html!

## The controller/handler layer
Now we're getting into the actual meat of this thing and starting to write some actual `rust` code!

We need a few things to get started, so go ahead and create a new folder under `src` called `handlers`.
In that holder, create a file called `mod.rs` and `home_handler.rs`. Let's start simple with some test data to show something when we load the home page. 
Open `home_handler.rs` and add the following:

```rust
use actix_web::{get, web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug)]
pub struct Frontmatter {
    title: String,
    file_name: String,
    description: String,
    posted: String,
    tags: Vec<String>,
    author: String,
    estimated_reading_time: u32,
    order: u32,
}

#[get("/")]
pub async fn index(templates: web::Data<tera::Tera>) -> impl Responder {
   let mut context = tera::Context::new();
   
   let frontmatters = vec![Frontmatter{
	tags: vec!["Rusty".to_string(), "Test".to_string()],
	title: "Test posts".to_string(),
	file_name: "test_posts.md".to_string(),
	description: "Just testing out the system".to_string(),
	posted: "2022-08-09".to_string(),
	author: "MBvisti".to_string(),
	estimated_reading_time: 12,
	order: 1,
   }];
   
   context.insert("posts", &frontmatters);
   
   match templates.render("home.html", &context) {
	Ok(s) => HttpResponse::Ok().content_type("text/html").body(s),
	Err(e) => {
		println!("{:?}", e);
		HttpResponse::InternalServerError()
			.content_type("text/html")
			.body("<p>Something went wrong!</p>")
	}
   }
}
```
If you've been through the `rust` book most of this should make sense to you.

If not, let's do a quick rundown. We create a struct and derive (the `#[derive()]` on top of `Frontmatter`) some methods from `serde` 
and the standard library. Next, we create a context through `tera`. Remember how we talked about providing data to the templates through our handlers? 
that's how we do it. We then simply call `.insert()` on the context, provide the variable name and the data and we're good to go.
In the return part of the function, we put `impl Responder` so we just have to return *something* that implements `Responder`. It just so happens 
that `HttpResponse` does implement `Responder`, so all there is left to do is to pull out the correct template from `tera`.
And since `template.render()` can fail we provide a backup that, admittedly, is a bit lazy for now but we can always come back and improve on this.
Lastly, you might be wondering how we actually access the templates we reference in the handler as `templates: web::Data<tera::TerA>`.

Actix has a neat way of sharing data though `App::app_data` and the `struct Data<T: ?Sized>(_)` type.
By passing *something* wrapped in `web::Data::new` to `App.app_data`  (side note: we could improve upon this by passing `web::Data::clone` to `app_data` 
instead of `web::Data::new` and just wrap our data once in `web::Data::new` and simply pass a reference to that in `web::Data::copy`. This is because `Data` 
uses `Arc` internally which makes it very cheap to clone) we can now extract it in your handlers!

Lastly, before we can access our new shiny handler, we need to make it accessible so open up `mod.rs` under `src/handlers`:

```rust
mod home_handler;

pub use home_handler::index;
```

We just need to add a route so we can serve the content to the user. Open up `main.rs` and add the following:
```rust
use std::net::TcpListener;
use actix_web::{dev::Server, web, App, HttpResponse, HttpServer};
use tera::Tera;

pub mod handlers; // new line

#[macro_use]
extern crate lazy_static;

lazy_static! {
	pub static ref TEMPLATES: Tera = {
	    let mut tera = match Tera::new("templates/**/*.html") {
			Ok(t) => t,
			Err(e) => {
				println!("Parsing error(s): {}", e);
				::std::process::exit(1);
			}
		};
		tera.autoescape_on(vec![".html", ".sql"]);
		tera
	};
}

pub fn start_blog(listener: TcpListener) -> Result<Server, std::io::Error> {
	let srv = HttpServer::new(move || {
	    App::new()
		   .app_data(web::Data::new(TEMPLATES.clone()))
		   .wrap(middleware::Logger::default()) // enable logger
		   .route("/health", web::get().to(HttpResponse::Ok))
            .service(handlers::index) // new line
	})
	.listen(listener)?
	.run();
	
	Ok(srv)
}
```

Give it a spin and you should see (a rather ugly) web page with your blog post!

## Static assets and the path of least resistance
While we were setting up our html templates you might have noticed what looks like `tailwindcss` classes and then promptly wondered what the hell kind of end result it had on the home page.
First of all, good on you for choosing great frontend tooling. Secondly, we haven't yet included the css needed for our styles to actually do something.

If you check out `base.html`, you will see this line: `<link rel="stylesheet" type="text/css" href="/static/css/index.css">` which hints at we need a `css` folder within a `static`.
Now, we can go about this in a number of ways depending on what kind of traffic we expect to get and how (over)engineered we are planning to make this.

One option could be to serve the content of `index.css` through an S3 bucket and then point an URL directly at it which then becomes referenced here.
It would certainly keep the resulting binary smaller but it does add some overhead that is not really warranted right now (more complex CI/CD flow, handling various paths in dev and prod etc).
Another option (the one we're going with), would be to just add it as part of the overall binary.

It keeps things simple for now and should our blog become a raving success we can most likely worry about this at that point.
So, create two new folders at the root, namely `static` and under that, `css`. Next up we need to do some tailwind setup that you can find on their own site, so for your convenience, here are the `cmd`s you can run in your terminal:

```bash
mkdir tailwind && cd tailwind

npm init -y

npm install -D tailwindcss

npx tailwindcss init

touch base.css
```

Almost there, just a few more steps so open up `tailwind.config.js` and copypaste the following:

```js
module.exports = {
  content: ["../templates/**/*.{html,js}"],
  theme: {
    extend: {},
  },
  plugins: [],
}
```

Two more things, open `base.css` and copypaste the following:

```css
@tailwind base;
@tailwind components;
@tailwind utilities;
```

Lastly, open up `package.json` and add the following under scripts:

```json
{
	"scripts": {
		"watch-css": "npx tailwindcss -i ./base.css -o ../static/css/index.css --watch",
		"build-css-prod": "npx tailwindcss -i ./base.css -o ../static/css/index.css --minify"
	},
}
```

Those two scripts give you a `cmd` to watch for any changes in development and then build for production when we merge.
That also implies that you would have to remember to run `build-css-prod` before merging to `master`.
This can definitely be improved upon but since you're the lead (and only) developer on this, it will be fine for our needs.

Last step, we need to make `actix_web` aware of our static files. So open up `lib.rs` and update it as follows:

```rust
use actix_files::Files; // new line
use actix_web::{dev::Server, middleware, web, App, HttpResponse, HttpServer};
use std::net::TcpListener;
use tera::Tera;

pub mod handlers;

#[macro_use]
extern crate lazy_static;

lazy_static! {
    pub static ref TEMPLATES: Tera = {
        let mut tera = match Tera::new("templates/**/*.html") {
            Ok(t) => t,
            Err(e) => {
                println!("Parsing error(s): {}", e);
                ::std::process::exit(1);
            }
        };
        tera.autoescape_on(vec![".html", ".sql"]);
        tera
    };
}

pub fn start_blog(listener: TcpListener) -> Result<Server, std::io::Error> {
    let srv = HttpServer::new(move || {
        App::new()
            .app_data(web::Data::new(TEMPLATES.clone()))
            .wrap(middleware::Logger::default())
            .service(Files::new("/static", "static/").use_last_modified(true)) // new line
            .route("/health", web::get().to(HttpResponse::Ok))
            .service(handlers::index) 
    })
    .listen(listener)?
    .run();

    Ok(srv)
}
```

We just need to build the `css` so from the root directory, run: `cd tailwind && npm run build-css-prod && cd ..`.

Now, give `cargo run` a spin, and you should see a much better-looking site.

### A quick interlude

Two things are currently missing: 1. a place to store the _actual_ blog posts and 2. a page to show a post in all its glory.

Let's start with the 2nd item on the list since we pretty much already did that. To not bore you with the same thing twice, 
go ahead and open this [link](https://github.com/MBvisti/awesome-blog/blob/master/templates/post.html) (link to the repo with the complete code), 
copy the content and create a new file under `templates` called `post.html`.

Next, go to [here](https://github.com/MBvisti/awesome-blog/blob/master/tailwind/base.css) copy the content and paste it into your `base.css`. 
Most of this should look familiar to you if you've done any `css` and just applies some rather simple styles to the post page.

## Steal like an artist (or, copy smarter people than yourself)
Everything is _mostly_ done, we just need a way to store our awesome blog posts and a way to retrieve them so we can remove the hard coding 
we did earlier. And the way to do this is heavily inspired by another blog post written by a much smarter guy than myself:
[fasterthanlime](https://fasterthanli.me/) whom, if you haven't already, should check out. This guy does some *seriously* deep dives and really 
knows his craft.

Go ahead and make a new folder called `posts` at the same level as the `static` folder. Inside this one, create another folder, give it a 
name like `my-first-article` and `cd` into it. I write all my stuff using `markdown` so that is what we're going to do here as well.

Create a file called `post.md` and copy-paste some lorem ipsum text (or, an article you've prepared). Next, create a file called 
`post_frontmatter.toml`, open it and add the following:

```toml
title = 'This is my first article'
file_name = 'my-first-article'
description = 'The first article authored by me'
tags = []
posted = '22/08/2022'
estimated_reading_time = 13
author = 'Morten Vistisen' # feel free to swap with your own name
order = 1
```

If this looks suspiciously like our `struct Frontmatter {....}`, then its because it's. 

Next order business, add some logic to extract all the frontmatters we might have and display it on our home page.
Let's start with the logic for getting all the frontmatters of our awesome blogs, so open `home_handler.rs` and add the following:

```rust
use std::{fs, io::Error}; // add these imports
use ignore::WalkBuilder; // add these imports

fn find_all_frontmatters() -> Result<Vec<Frontmatter>, std::io::Error> {
    let mut t = ignore::types::TypesBuilder::new();
    t.add_defaults();
    let toml = match t.select("toml").build() {
	  Ok(t)=> t,
	  Err(e)=> {
		println!("{:}", e);
		return Err(Error::new(std::io::ErrorKind::Other,
		"could not build toml file type matcher"))
	  }
    };
    
    let file_walker = WalkBuilder::new("./posts").types(toml).build();

    let mut frontmatters = Vec::new();
    for frontmatter in file_walker {
        match frontmatter {
	      Ok(fm) => {
		    if fm.path().is_file() {
		        let fm_content = fs::read_to_string(fm.path())?;
			  let frontmatter: Frontmatter = toml::from_str(&fm_content)?;
			  
			  frontmatters.push(frontmatter);
		    }
		}
		Err(e) => {
		    println!("{:}", e); // we're just going to print the error for now
		    return Err(Error::new(std::io::ErrorKind::NotFound, "could not locate frontmatter"))
		}
	  }
    }
    
    Ok(frontmatters)
}
```

Since we're storing our posts as part of the binary we need a way to locate all the frontmatter of our posts and for this, we're going to use `ignore::WalkBuilder`.
This gives us a recursive directory iterator with a large number of configs that we can set, depending on the type of action we want to do.
For this case, we want it to look in the `posts` directory and look for all files with the `.toml` extension, and it does it, *BLAZINGLY* fast.

One thing to note here is that `find_all_frontmatters` returns a `Result` since there are actions that can fail, which means we also have to handle some errors.
As you might be able to tell, I haven't spent the most time on errors in this function and basically, just log whatever error comes from `ignore` and then return
the most similar one from `std::io::ErrorKind`. This can be done better using libraries like `anyerror` or `thiserror`. I encourage you to play around 
with this yourself or open a pull request if you've improvements [here](https://github.com/mbvisti/awesome-blog).

To actually show some dynamic data on the home page, we need to use the above function in our `index` handler so open up `home_handler.rs` and make the following adjustments:

```rust
#[get("/")]
pub async fn index(templates: web::Data<tera::Tera>) -> impl Responder {
    let mut context = tera::Context::new();

    let mut frontmatters = match find_all_frontmatters() {
        Ok(fm) => fm,
        Err(e) => {
            println!("{:?}", e);
            return HttpResponse::InternalServerError()
                .content_type("text/html")
                .body("<p>Something went wrong!</p>");
        }
    };
    frontmatters.sort_by(|a, b| b.order.cmp(&a.order));

    context.insert("posts", &frontmatters);

    match templates.render("home.html", &context) {
        Ok(s) => HttpResponse::Ok().content_type("text/html").body(s),
        Err(e) => {
            println!("{:?}", e);
            HttpResponse::InternalServerError()
                .content_type("text/html")
                .body("<p>Something went wrong!</p>")
        }
    }
}
```

Give `cargo run` a spin and see the results!

The last thing that's needed is to add the handler for the `post.html` page and extract the markdown of the article, 
convert it to `html` and serve it. Since we already have `post.html` ready, go ahead and create `post_handler.rs` under 
`src/handlers` and open it up. We need to add two functions: 1. to extract the frontmatter from a specific post and 2. 
to extract the markdown of the post.

```rust
use std::{io::Error, fs};

use super::home_handler::Frontmatter;

fn extract_markdown(post_name: &str) -> Result<String, Error> {
	let markdown = match fs::read_to_string(format!("./posts/{}/post.md", post_name)) {
		Ok(markdown) => markdown,
		Err(e) => {
			println!("{:?}", e);
			return Err(e)
		}
	};

    Ok(markdown)
}

fn extract_frontmatter(post_name: &str) -> Result<Frontmatter, Error> {
	let frontmatter_input =	match fs::read_to_string(format!("./posts/{}/post_frontmatter.toml", post_name)) {
		Ok(s) => s,
		Err(e) => {
			println!("{:?}", e);
			return Err(e)
		}
	};
	
	let frontmatter = match toml::from_str(&frontmatter_input) {
		Ok(fm) => fm,
		Err(e) => {
			println!("{:?}", e);
			return Err(Error::new(std::io::ErrorKind::Other, "could not find post frontmatter"))
		}
	};

    Ok(frontmatter)
}
```

Not much new here, only the use of `toml::from_str` to deserialize a string into a specific type, in this case: `Frontmatter`.

Last thing we need is to create the handler and add it as a `service` to our `App` in `lib.rs`, so add the following to `post_handler.rs`:
```rust
use actix_web::{web, get, Responder, HttpResponse};
use pulldown_cmark::{Options, Parser, html};

#[get("/posts/{post_name}")]
pub async fn post(
	tmpl: web::Data<tera::Tera>,
	post_name: web::Path<String>,
) -> impl Responder {
	let mut context = tera::Context::new();
	let options = Options::empty(); // used as part of pulldown_cmark for setting flags to enable extra features - we're not going to use any of those, hence the `empty();`
	
	let markdown_input = match extract_markdown(&post_name) {
		Ok(s) => s,
		Err(e) => {
			println!("{:?}", e);
			return HttpResponse::NotFound()
				.content_type("text/html")
				.body("<p>Could not find post - sorry!</p>")
		}
	};
	
	let frontmatter = match extract_frontmatter(&post_name) {
		Ok(s) => s,
		Err(e) => {
			println!("{:?}", e);
			return HttpResponse::NotFound()
				.content_type("text/html")
				.body("<p>Could not find post - sorry!</p>")
		}
	};
	
	let parser = Parser::new_ext(&markdown_input, options);
	
	let mut html_output = String::new();
	html::push_html(&mut html_output, parser);
	
	context.insert("post", &html_output);
	context.insert("meta_data", &frontmatter);
	
	match tmpl.render("post.html", &context) {
		Ok(s) => HttpResponse::Ok().content_type("text/html").body(s),	
		Err(e) => {
			println!("{:?}", e);
			return HttpResponse::NotFound()
				.content_type("text/html")
				.body("<p>Could not find post - sorry!</p>")
		}
	}
}
```

Much of this is familiar, the difference here is that we make two variables available to our template: `post` and `meta_data`.

Last thing, we need to include the new handler in `lib.rs` and our `handlers/mod.rs`. I will leave that up to you!

After adding that, give `cargo run` a spin again and you should be able to click on the article from the home page, be re-directed to the post page and see your article!

## Closing thoughts

I hope you enjoyed the walk-through and are interested in doing more web development with `rust`, it's in a really good state right now, with an exciting future ahead.
Yes, the beginning might be a little bit frustrating and getting your head wrapped around ownership, borrowing and referencing can take some time.
But I promise you, it's well worth it!

The compiler errors you get can be annoying, but will also guide you towards the correct path, and compared to something like `Typescript` it's a godsend.
The more you write the easier it becomes.

It has quickly become my go-to language for writing small services as it plays very nice with `aws lambda`s and due to how fast it is, cold starts are not even an issue anymore.
So, please take what you have here and extend it as much as you want.

## Resources
Here is a collection of all the materials I've used to get started with `rust`. Some of it is free and some of it is paid, but all of them are highly recommendable
so I hope you find some further learning!
  - [Zero 2 Production](https://www.zero2prod.com/)
  - [The rust book](https://doc.rust-lang.org/book/title-page.html)
  - [A new website for 2020](https://fasterthanli.me/articles/a-new-website-for-2020)
  - [Rust by example](https://doc.rust-lang.org/stable/rust-by-example/)

