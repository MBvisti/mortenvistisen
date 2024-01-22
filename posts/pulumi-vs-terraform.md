---
title: "Pulumi vs Terraform yada yada"
focus: "compare pulumi and terraform from a developer perspective, include current experiences and a short guide for setting up pulumi and go. Section about optimizing for humans/team, not necessarily purely technical reason. Size of team, infrastructure needs, team experience level."
keywords: ['terraform alternative', 'terraform vs pulumi', 'pulumi tutorial', 'pulumi aws', 'infrastructure as code tools']
---

I've started a new contract, where I'm tasked with consolidating the company's infrastructure under AWS. They (the company) are starting to see some serious growth, both in terms of reveneue but also in terms of the team, so my natural reaction was to begin researching infrastructure as code tools (IaC tools). I've previously used terraform in other companies I've done work for, but never from scratch and always found the experience rather frustrating. So, I began to search for terraform alternatives and stumbled upon pulumi, which lets you provision resources using a number of programming languages (including Go) and has excellent support for aws providers.

But, as in any other situation where you are considering marrying yourself to a new tool, it's very worthwhile to do some research and preferably build some small proof of concept before you commit. Since this is what I've been doing for the past many weeks, I wanted to share my experience, and provide you with a small tutorial for deploying a small Go app, using the following AWS service:
- ECR
- ECS
- Fargate

both in terraform and pulumi.

## Short intro to terraform

I think it's fair to say terraform has been the big dog, in the IaC space for a long time. They've support for the majority of cloud providers, been battle-tested at scale and has a bunch of great tutorials. They also used to have a open source license, using Mozilla Public License v2, but decided recently to switch to what is known as the Business License, BSL, which can be considered as "source-available". I will not go into a longer debate about why this is not ideal, if it's fair or not, but it does pose a threat to the development of third-party tools and your options in terms of making tweeks if needed.

For the unfamiliar, terraform is written in their own "DSL" known as HCL. It's a declarative language, so you tell it the desired state you want your infrastructure to be in, and terraform figures out how to get your infra to that state.

A short example of how you would provision a S3 bucket would look like this:

```hcl
resource "aws_s3_bucket" "my_bucket" {
  bucket = "my-bucket"
  acl    = "private"
}
```

Which would be written in a file with the `.tf` extension. As we'll see later on, you have lots of options in terms of structuring your code into modules, different types of options for code-reuse: `data` blocks, `variables`, `outputs` etc.

## Short intro to pulumi

Pulumi is rather new in the IaC space but have made some significant advances. Some years ago, I worked at a place who utilized AWS's CDK, and I imagine people would feel pulumi to be quite similar.

You've a range of choices in terms of the language you want to use, JS/TS, python, Go etc. All having the same API, so switching __shouldn't__ be the biggest issue. I write all my stuff in Go these days, so the same example as above written using Go and pulumi, would look like this:

```go
func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := s3.NewBucket(ctx, "my-bucket", &s3.BucketArgs{
			Bucket: pulumi.String("my-bucket"),
			Acl:    pulumi.String("private"),
		})
		if err != nil {
			return err
		}
		return nil
	})
}
```

To anyone that has ever written Go, this should immediately make sense, even without knowing the pulumi api or the aws ecosystem. And this is one of the main benefits we will get back to, again and again, over the course of this post.

Pulumi is also "true open source", in the sense that they operate with an Apache 2.0 license, which is a lot more permissive than the BSL license terraform uses. Which is great for the community and adoption rate, but, so was terraform at one point. For now, it's a huge plus.

However, it isn't all sunshine and rainbows. Pulumi is still young, so documentation is still not great when actually writing code. Their general purpose docs do make up for this, but it can be frustrating having to jump in and out of the editor to look up what you're currently working on.


## What we're building

Containerization is pretty much the standard these days, and when working with Go, this becomes even easier (try to deploy a ml model using python and docker; it's a totally different experience..). So we'll build a simply api that just returns "word" whenever you ping the endpoint `/hello`.

As eluded to in the introduction, we'll deploy this using ECS and Fargate, and we'll use ECR to store our docker image. On the Go side, we will use `chi/v5` as our router.

If you're coding along, go ahead and create a new folder, initialize a go module and create the folder structure `cmd/app/main.go`. In `main.go`, add the following:

```go
package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	
	r.Use(middleware.Logger)
	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("world"))
	})

	http.ListenAndServe(":8000", r)
}
```

Quite straightforward, but we also need to create a docker image to run our app. So, create a `Dockerfile` in the root of your project, and add:

```docker


```































