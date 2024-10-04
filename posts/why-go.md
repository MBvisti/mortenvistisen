I've been writing Go for some time now, close to 7 years. But it wasn't where I originally started out in my tech career. 

My first job was at a place where their core competency where backend and databases-bootstrap, the defacto fronted library for multiple years at the time, was something completely new to them. UI/UX was mostly an afterthought but since the clients were large accounting companies, it didn't matter. They were happy with the colour blue and lots of boxes.

I don't have a traditional background in computer science and have been self-taught all the way which probably was also why I ended up in frontend to begin with. The thinking at my first job seemed to be: 

> How could you trust a guy who hasn't studied every algorithm under the sun to make changes to a server? A ridiculous notion. 

So I started learning how I did everything else, typing "how to learn frontend" into Google.

React was just breaking away, in terms of popularity, from Angular at the time and was still in the good 'ol pre hook days. I eventually became proficient with React and would mainly write it (and JS/TS) for the next couple of years.

But, something kept happening. Everything seemed consistently out of date as new tools/libraries/frameworks were being shipped every day that you just _had to use_. But even more important/discouraging was that things seemed to become super complex, for simple things. Do you want to submit a form? You need event listeners, and different hooks to efficiently update the state, validate the input AND then do it once more on the backend. It turned out to be a rabbit hole of complexity where you know need a PhD in hooks to be able to efficiently use them.

And it makes sense, React has grown immensely over the last ~10 years. But they seem to have been trying to cover every requirement known to man.

None of the above is new in any sense, it has been written about and discussed extensively in the tech community for some time now. 

I, thankfully, don't write much React anymore and have gone back to the roots: plain old HTML which, thanks to Go, its ecosystem, and some guy from Montana still gets like 90% of the interactivity you get with a SPA app.

So, let me explain why I think you should use Go and why I fell in love with it. 

If you just want the tl;dr: simplicity and limited choice.

## Limited choice

When writing Go, you'd often find that a solution to a problem tends to converge to a common conclusion across different teams/people/etc.

Some time ago I was frustrated that a validation library I was using didn't make it easy to do what I wanted it to. So I did what every other good engineer does and wasted a couple of days writing my own solution. 

This was on my own time, after all, so fuck it. 

About halfway through, I started to look into what other people had done and across multiple different projects, the solution to the problem was implemented in a very similar way.

Of course, this might be due to every library that comes after the first one copies the code. But that's not my general experience. Solutions to similar problems seem to converge due to Golang's strict feature set.

One criticism I often see is that Go doesn't allow developers to "express" themselves in code. Which, honestly, is bullshit. You're not Picasso, you get paid tons of money to sit in a cosy office and drink kombucha directly from tab. Express yourself in system designs and elegant solutions if you'd like, but leave the code mate.

It's very difficult to spot an individual's contribution to a Go codebase without looking into the commit history. There is not much room to be eccentric or have your own style. Again, sounds rather boring but it ends up working extremely well when you care about building and shipping digital products.

The solution is what you focus on.

### Simple syntax & easy onboarding

The syntax is very minimalistic with only 25 keywords to learn from which you probably need like, what, two-thirds? 

It's so fast to pick up which has also been my experience whenever I have introduced Go to teams where I've worked as a consultant. Doesn't really matter if they're completely fresh out of school, work in a different way (I'm looking at you here, data scientist) or have a lot of experience in another language.

People generally tend to pick up Go fast which I think comes from its readability. 

To me, it's not beautiful as you might look at a piece of Rust code you spend 8 story points writing, only to find out that it didn't match stakeholder expectations once it finishes compiling at the end of the sprint. 

No, it's practical and boring and allows me to focus on what should be the more interesting part of my job, the problem. There is no need to argue over which linter style (or which linter) to use, it's decided for you and will be the same for everybody on the team.

## Everything you need in one place

The standard library is fucking amazing. It brings a certain amount of calm to your projects as you can be certain that it will be maintained and won't have to make tons of rewrites to update.

Most things you want to build related to the web you can get done using just the standard library. There are other libraries that make certain parts better but they are not strictly necessary.

But, maybe just as important, the Golang toolchain ships with a linter and a standard for how to apply linting rules. This is generally something most developers have an opinion on that doesn't really add value, only take away value, if different styles are used. Just use one and get on with it.

## Go is ripe for Web Development

In recent years, projects in the Go ecosystem have made it even easier to build full-stack apps only using Go. 

The main thing, in my opinion, that has come out of the ecosystem that will help propel full-stack development in Go is Templ. 

The creators behind it talk about (paraphrasing) that one of the "issues" we see is that lots of frontend/full-stack developers have only seen the SPA way of doing things. Hypermedia, anchor tags and forms are, in their native form, strange concepts. So, they wanted to create something that would allow them to think in components but still be completely valid HTML. 

If you've ever done any Django, rails or Laravel development this way of doing frontends is most likely very familiar to you. 

Use a base template to wrap the application, have some views for specific pages and abstract elements like buttons, tables etc into components for re-usability so that your API is no longer JSON-driven but hypermedia-driven. 

You're creating a contract as you do with JSON-based APIs but now the consumer of the contract, the client, is much more tightly coupled to the API. This makes breaking changes much less likely, maintenance requirements fall and the consequences of changes become much easier to reason about.

## A straightforward database layer

When I first started out in Go, the database layer was one of the most repetitive parts. Lots of "=+" for setting up query statements, scanning the rows into structures etc. You were writing sql but also not really.

Sqlc gets mentioned on forums quite often, and for good reasons, it lets you stay in the domain of SQL and solves the majority of queries I have.

It removes a lot of the boilerplate by generating type-safe code and takes care of scanning the results into what matches your database structure, validates that the fields are correct.

You'll have to reach for something else when doing dynamic queries, but adding something like squirrel to the mix is quite straightfoward. And since you can quite easily just add the generated methods from sqlc as a dependency to a database struct, having both is trivial.

## Deployment is a joy

Once you've tasted the joy of single binary deployments it's hard to Go back.

You can go super simple: rent a VPS, add a service to Systemd and you're up and running.

But as you might have experienced if you've written some Go, it works very well with Docker. Given that it builds and outputs a single binary, you can get super slim docker images. Like 10mbs slim. Just getting started doing that with Python and you quickly approach 1gb. But storage is cheap I guess.

I tend to build a docker image that contains the app and some tooling, for things like running migrations, and I often end up with a docker image around 100mb.

## Concurrency

It's rather hard to make a Go appreciation post without talking about concurrency. It's very good, but I honestly don't use it that often. 

I usually take advantage of it whenever I have to do some data processing but in most of my work, it's not really that common. It's built into the packages you rely on and given its, comparatively, easy approach it becomes much more approachable.

It's there if I do need it.

## Performant and efficient

This is always nice to throw in the mix when trying to shill Golang to other people. 

Even though it's not the main reason for choosing Go it is very nice to know that your APIs are performant almost straight out of the box and don't eat up a lot of RAM. Depending on what you're doing this can have a financial benefit as you'd be using up fewer cloud resources. 

I'm hosting quite a few apps on a small VM right now where the memory overhead is negligible.
