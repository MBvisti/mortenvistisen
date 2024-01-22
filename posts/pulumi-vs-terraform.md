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

Pulumi is just a very strange name. Growing up, I had a friend we for some reason called pummi (or maybe it was bummy, can't recall for sure) and all I can see when reading the name, is his face. They're also rather new in the IaC space but have made some significant advances. Some years ago, I worked at a place who utilized AWS's CDK, and I imagine people would feel pulumi to be quite similar.

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

```dockerfile
FROM golang:1.21 AS build

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -mod=readonly -v -o /main cmd/app/main.go

FROM scratch

WORKDIR /app

COPY --from=build /main /

CMD ["/main"]
EXPOSE 8000
```

I'm not going to dive too much into details about docker, we're simply utilizing a multi-step build process, where we copy the files into the container, install the dependencies and build our binary. Then, we copy this into a layer that is based on "scratch", which creates a very small image, roughly ~4mb.

## Getting started with Terraform

To make this run, in a more realistic example, we're going to need quite a few things so buckle up, no half-baked "hello world" tutorials here. After all my years in business school and finance, I learned one thing, a visual representation of a topic can really help with understanding (or, in the finance world, make you seem knowledable about an area you just googled yesterday). Anyway, this is what we'll need:

![app infrastructure](/static/images/pulumi-vs-terraform-infrastructure.svg)

### First steps

Before we dive in, we need to setup some things first, so, in the root of the project, add a `main.tf` file and add:

```hcl
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.10.0"
    }
  }
}

provider "aws" {
  region = "eu-west-1"
}
```

We need to initialize terraform, so, in the root of the project go ahead and run(assuming you already installedit, ofc, if not, please go and do that now):

```bash
terraform init
```

This will download an install the aws provider, and we're ready to go.

### Networking & Security

I'm by no means an expert in those subjects, by I'm rather good at google and reading documentation, so we'll try our best here to follow best practices from aws when setting up our app. And for that, we need a VPC, some subnets, security groups, gateways and load balancers or in other words, quite a few things.

In `main.tf`, add:

```hcl
resoruce "aws_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support = true
  tags = {
    Name = "pulumi-vs-terraform-vpc"
  }

resource "aws_subnet" "subnet_public_a" {
  vpc_id                  = aws_vpc.vpc.id
  cidr_block              = "10.0.16.0/20"
  availability_zone       = "eu-west-1a"

  tags = {
    name = "subnet_public_a"
  }
}

resource "aws_subnet" "subnet_public_b" {
  vpc_id                  = aws_vpc.vpc.id
  cidr_block              = "10.0.0.0/20"
  availability_zone       = "eu-west-1b"

  tags = {
    name = "subnet_public_b"
  }
}

resource "aws_subnet" "subnet_private_a" {
  vpc_id            = aws_vpc.vpc.id
  cidr_block        = "10.0.144.0/20"
  availability_zone = "eu-west-1a"

  tags = {
    name = "subnet_private_a"
  }
}

resource "aws_subnet" "subnet_private_b" {
  vpc_id            = aws_vpc.vpc.id
  cidr_block           =  "10.0.32.0/20"
  availability_zone = "eu-west-1b"

  tags = {
    name = "subnet_private_b"
  }
}
```

We add a VPC, which if you check the diagram is the red box an isolates everything else within. Nothing comes in or out unless we give it permission. Next, we follow some best practices for high availability and create 2 public and 2 private subnets,  Next, we follow some best practices for high availability and create 2 public and 2 private subnet, both in different availability zones (the green and blue lines in the diagram).

Next, we add the internet gateway, route table and their associations (not depicted but necessary):

```hcl
resource "aws_internet_gateway" "internet_gateway" {
  vpc_id = aws_vpc.vpc.id

  tags = {
    name = "internet_gateway"
  }
}

resource "aws_route_table" "public_route_table" {
  vpc_id = aws_vpc.vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.internet_gateway.id
  }

  tags = {
    name = "public_route_table"
  }
}

resource "aws_route_table_association" "public_route_table_association_a" {
  subnet_id      = aws_subnet.subnet_public_a.id
  route_table_id = aws_route_table.public_route_table.id
}

resource "aws_route_table_association" "public_route_table_association_b" {
  subnet_id      = aws_subnet.subnet_public_b.id
  route_table_id = aws_route_table.public_route_table.id
}

resource "aws_route_table_association" "private_route_table_association_a" {
  subnet_id      = aws_subnet.subnet_private_a.id
  route_table_id = aws_route_table.private_route_table.id
}

resource "aws_route_table_association" "private_route_table_association_b" {
  subnet_id      = aws_subnet.subnet_private_b.id
  route_table_id = aws_route_table.private_route_table.id
}

resource "aws_eip" "elastic_ip" {
  tags = {
    name = "nat_gateway"
  }
}

resource "aws_nat_gateway" "nat_gateway" {
  allocation_id = aws_eip.elastic_ip.id
  subnet_id     = aws_subnet.subnet_public_a.id

  tags = {
    name = "nat_gateway"
  }
}

resource "aws_route_table" "private_route_table" {
  vpc_id = aws_vpc.vpc.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.nat_gateway.id
  }

  tags = {
    Name = "private_route_table"
  }
}
```

Lastly, we add some security groups, which will be used to control traffic in and out of our app. We need one for our upcoming load balancer and one for ECS:

```hcl
resource "aws_security_group" "ecs_sg" {
  name   = "ecs_sg"
  vpc_id = aws_vpc.vpc.id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    name = "ecs_sg"
  }
}

resource "aws_security_group" "alb_sg" {
  name   = "alb_sg"
  vpc_id = aws_vpc.vpc.id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = -1
    self        = "false"
    cidr_blocks = ["0.0.0.0/0"]
    description = "any"
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
```

If you're currently thinking "but, they have the exact same properties but with different names", you're absolutely right. We're cutting corners a bit, but in a real world app, you would most likely need different security groups with different ingress and egress paramters, so that's what we're simulating.

Thats, quite a mouthful, especially if you're not familiar with HCL, but all we've done is specifying what resources we need and how they should be connected.

### Load balancer

We're almost at the point where we can begin setting up ECS and deploy our app, we just need to setup our load balancer, so we can actually access this thing once it's live.

```hcl
resource "aws_alb" "app_alb" {
  name_prefix        = "app_alb"
  internal           = false
  load_balancer_type = "application"
  idle_timeout       = 60
  security_groups    = [aws_security_group.alb_sg.id]
  subnets            = [aws_subnet.subnet_public_a.id, aws_subnet.subnet_public_b.id]

  tags = {
    Name = app-alb-public"
  }
}

resource "aws_alb_listener" "ecs_alb_listener" {
  load_balancer_arn = aws_alb.app_alb.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_alb_target_group.public_tg.arn
  }
}

resource "aws_alb_target_group" "public_tg" {
  name        = "ecs-target-group"
  port        = 8000
  protocol    = "HTTP"
  target_type = "ip"
  vpc_id      = aws_vpc.vpc.id

  health_check {
    path = "/"
  }
}

output "lb_dns" {
  description = "The DNS name of the load balancer"
  value       = "http://${aws_alb.app_alb.dns_name}"
}
```

### ECX

With that in place, we need some of the elastic something something services setup and first on the agenda is ECR: elastic container registry, so we have a place to host our image(s).

```hcl
resource "aws_ecr_repository" "repo" {
  name                 = "repo"
  image_tag_mutability = "MUTABLE"
  force_delete         = true

  tags = {
    "Name" = "repo"
  }
}

output "repo_url" {
  value = aws_ecr_repository.repo.repository_url
}
```

And we've arrived at a chicken or egg situation since we need to build and push our image to ECR, before we can successfully setup ECS and Fargate. So, assuming you've an AWS account that can provision new resources open up a terminal and run:

```bash
terrform apply
```

This will create a plan for what to provision, show you a summary and if you agree, simply type 'yes' and enter.

Assuming everything went through without any issues, you should now be able to get the url of the repo so you can push a docker image to it. A little bash scripting never hurts, so if you want, you can utilize this so build, tag and push to the repo. Just remember to update the image name, tag and registry url:

```bash
#!/bin/bash

# Set your Docker image name and tag
IMAGE_NAME="app"
TAG_LATEST="latest"
DOCKERFILE="Dockerfile"
REGISTRY_URL="1234567890.dkr.ecr.eu-west-1.amazonaws.com/repo"

# Build the Docker image
docker build -t "${REGISTRY_URL}" -f "${DOCKERFILE}" .

# Tag the Docker image
docker tag "${REGISTRY_URL}"  "${REGISTRY_URL}:${TAG_LATEST}"

aws ecr get-login-password | docker login --username AWS --password-stdin ${REGISTRY_URL}

# Push the Docker image
docker push "${REGISTRY_URL}:${TAG_LATEST}"

# Clean up local images (optional)
docker image prune -f
```

Verify that your image is in the repo and lets continue.


### ECS
