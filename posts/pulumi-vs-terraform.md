I've started a new contract, where I'm tasked with consolidating the company's infrastructure under aws. They (the company) are starting to see some serious growth, have a team in-place currently using a PaaS and various other cloud providers. Infrastructure as code tools have been popular for a long time; I've mostly been a bit skeptical about using them given a small team size, often the ui is more than enough. But, when you need to start spinning up multiple environments (dev, staging, prod, demo etc), this start to become a lot of overhead. An issue I've often seen is that it becomes the responsibility of one or two developers to maintain this, as the rest of the team is not familiar with the tool and language syntax. Obviously, this is not ideal, so I wanted to research the space to see if there was a middle way where you can get the benefits of IaC, without the overhead of having to learn a new language and tool.

I've previously used terraform in other companies I've done work for but never from scratch and (full disclosure) always found the experience rather frustrating. So, I began to search for terraform alternatives and stumbled upon pulumi, which lets you provision resources using a number of programming languages (including Go) and pulumi's aws provider support is excellent.

But, as in any other situation where you are considering marrying yourself to a new tool, it's very worthwhile to do some research and preferably build some small proof of concept before you commit. I've had some interesting discoveries through this process, why I choose pulumi over terraform and when I would consider terraform over pulumi.

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

---

<div id="subscribe-form" class="max-w-6xl px-4 sm:px-6 lg:px-8 mx-auto"></div>

---

## Building the POC

All of the applications are containerized being we can take full advantage of aws's elastic container service, so the POC was a simple hello world api packed in a docker image, served using fargate and a load balancer in front. I've been aware of ECS and fargate (plus, vpc, subnets, security groups etc), but never used them, so there was also some learning involved.

The first thing that became clear, terraform has a lot of great documentation and tutorial available. While pulumi does also have documentation and available tutorials, they are not as "complete" as what you find with terraform. However, I would find myself hitting a problem in the infrastructure setup with terraform, switch to pulumi and even with less available information, I was able to reason my way out of the problem because I simply had to read Go code. The overhead of a new language and platform would be frustrating, but staying in something I already knew, removed one leg of complexity and the aws concepts was now the only thing I had to figure out. I suspect that would be one of the main feelings of other engineers that don't do full-time devOps as well.

I managed getting the poc running using both tools but a big help came from staying in an environment I was already familiar with. I'm a heavy (neo)vim user and while there are lots of good plugins for terraform, the Go environment is something I'm in everyday so nothing new were needed.

If I dedicated the time, I don't suspect it would take super long to get familiar with hcl and terraform to the point where the environment would feel close to as familiar as Go. But, given that I wouldn't have to work everyday with terraform after getting the initial setup done, that familiarity would soon start to fade. Meaning, every time a new infrastructure part had to be added, I would need to re-familiarize myself. And, for the rest of the team, this would be same but with the added overhead of getting into a part of the code they didn't help write. This could be mitigated by involving everyone from the get-go, but this is also a business, there are other priorities that needs building too.

Documentation might be the answer some of you are screaming at the screen right now, and I hear you, but documentation quickly goes stale. This is especially true in smaller teams where things move fast, if the majority can be documented by the code itself, you save yourself much trouble.

## Why I choose pulumi over terraform

I ended up going with pulumi. I know from experience that the "correctness" of these choices only show its true self after you're well passed the point of making the switch to the other candidate feasible. But, from my current point of view, this is what will best benefit the team I'm currently working with.

Two questions kept coming up, as I were trying to decide on which tool to use:
1. What are our infrastructure requirements now and in the future? (i.e. are we ever going to need a full-time devOps)
2. What does the current team know? (i.e. programming languages, cloud providers etc.)

The first might seem arbitrary, but the likely hood of finding a full-time devOps engineer that knows terraform is properly much higher, than them knowing Go and pulumi. Given that the needed infrastructure now and in the foreseeable future most likely will not require a full-time devOps role, terraform begins to feel a bit overkill. It's a new language, everybody on the team would need to learn as well as any best practices that comes with it. All the while, getting into the nitty gritty of aws which will also be a new learning experience for the team (myself included).

The second, to me, was very important since I've, from multiple experiences, seen what happens when you have a big terraform setup that maybe one engineer was in charge of creating. Most of the time, it ends up being that one person that has to maintain and fix it, effectively making them a devOps person (which might be against their will -> they switch job -> team is fucked). The rest of the team, despite the person writing documentation, making presentations etc etc, are still very reluctant to touch it. It's something they are unfamiliar with, the infrastructure system can be quite difficult to "have in your head", you can potentially take down everything without knowing how to get back to a running state. It creates a lot of friction. 

You might argue: "well, they will just have to learn it" and I agree, that would be a sensible path forward. But we're not trying to solve a technical problem right now, but rather, a human one. And without allocating the time for people to learn it, it probably will not happen. So, in lieu of a devOps team, you need to optimize for what gets people to actually want to learn infrastructure. Over the two weeks I spend researching and building some POCs, I would often find very good tutorials and materials on setting up infrastructure with terraform. But as soon as I hit my head against wall, not knowing a specific syntax/setup/approach used in the tutorial or what the specific aws services was about, it become very hard to move on. Multiple times, I switched back to pulumi where their docs had examples in Go code, which was much easier for me to wrap my head around since I already knew Go. I "only" had to figure out the aws part. After, I would switch back to terraform and fix the thing with my new understanding from pulumi/go. But it's important to note where that knowledge originated from, by reading code, which most of us are doing pretty much every day.

Taking away the overhead of both having to learn a new language and a cloud provider, to only the cloud provider, helps a tremendous amount. You also get the added benefit of already knowing best practices, when and why to break those practices. You get an api from pulumi so you can start provisioning infrastructure directly in your application(s), which makes spinning up demo environments for sales processes trivial. 

## When I would choose terraform

Up to a certain point, something becomes a standard because of it quality and appeal to a large part of the tools target group. After that point, better alternatives might actually exists but since the standard requirement is to know said tool, everyone else just continues learning that (I'm looking at you, React. Switch those frontends back to html and use htmx for interactivity!). If that's whats going on here with terraform and pulumi is hard to say. I actually had a much better experience exploring terraform, compared to previous experiences, where I only had to tweak it a little bit.

If your team includes members that already know terraform or you know that you have some large scale infrastructure needs in the near future, that will include a dedicated devOps team, terraform still seem like the better choice. I did pick terraform and hcl up rather quickly, but I'm still quite confident that this would fade again after not using it daily, so if you don't have that limitation, terraform would win for me.

The main argument from my side here is to remove friction from your team when they don't deal with infrastructure on a daily basis, learning a new language while your PM is blasting new tickets your way like his life depended on it, is tough. And having to wait for the one team member who truly knows the infrastructure code coming back from holiday to fix things, is very sub-optimal.

## Was this a good decision

It's too early to tell, I will revisit this in a few months once we've had some more production experience with it and the team has started working more with it as well.

