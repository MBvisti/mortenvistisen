I've started a new contract, where I'm tasked with consolidating the company's infrastructure under AWS. They (the company) are beginning to see some serious growth and have a team currently using PaaS and various other cloud providers. Infrastructure as code tools have been popular for a long time; I've mostly been skeptical about using these. Given a small team size, often, the UI of the cloud provider is more than enough. But, when you need to start spinning up multiple environments (dev, staging, prod, demo, etc), this can become a lot of overhead. An issue I've often seen is that it becomes the responsibility of one or two developers to maintain this, as the rest of the team is not familiar with the tool and language syntax. This is not ideal, so I wanted to research the space to see if there was a middle way where you can get the benefits of IaC, without the overhead of having to learn a new language and tool.

I've previously used Terraform in other companies I've done work for but never from scratch and (full disclosure) always found the experience rather frustrating. So, I began to search for terraform alternatives and stumbled upon Pulumi, which lets you provision resources using several programming languages (including Go) and Pulumi's AWS provider support is excellent.

But, as in any other situation where you are considering marrying yourself to a new tool, it's very worthwhile to do some research and preferably build some small proof of concept before you commit. I've had some interesting discoveries through this process, why I choose Pulumi over Terraform and when I would consider Terraform over Pulumi.

## Short intro to Terraform

I think it's fair to say Terraform has been the big dog, in the IaC space for a long time. They've support for the majority of cloud providers, have been battle-tested at scale, and have a bunch of great tutorials. They also used to have an open source license, using Mozilla Public License v2, but decided recently to switch to what is known as the Business License, BSL, which can be considered as "source-available". I will not go into a longer debate about why this is not ideal if it's fair or not, but it does pose a threat to the development of third-party tools and your options in terms of making tweaks if needed.

For the unfamiliar, terraform is written in its own "DSL" known as HCL. It's a declarative language, so you tell it the desired state you want your infrastructure to be in and Terraform figures out how to get your infra to that state.

A short example of how you would provision a S3 bucket would look like this:

```hcl
resource "aws_s3_bucket" "my_bucket" {
  bucket = "my-bucket"
  acl    = "private"
}
```

Which would be written in a file with the `.tf` extension. You have lots of options in terms of structuring your code into modules, and different types of options for code-reuse: `data` blocks, `variables`, `outputs` etc.

## Short intro to pulumi

Pulumi is just a strange name. Growing up, I had a friend that we for some reason, called Pummi (or maybe it was bummy, can't recall for sure), and all I can see when reading the name is his face. They're also rather new in the IaC space but have made some significant advances. Some years ago, I worked at a place that utilized AWS's CDK, and I imagine people would feel Pulumi to be quite similar.

You have a range of choices in terms of the language you want to use, JS/TS, python, Go, etc. All have the same API, so switching __shouldn't__ be the biggest issue. I write all my stuff in Go these days, so the same example as above written using Go and Pulumi, would look like this:

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

To anyone who has ever written Go, this should immediately make sense, even without knowing the Pulumi API or the AWS ecosystem. And this is one of the main benefits we will get back to, again and again, throughout this post.

Pulumi is also "true open source", in the sense that it operates with an Apache 2.0 license, which is a lot more permissive than the BSL license terraform uses. Which is great for the community and adoption rate, but, so was Terraform at one point. For now, it's a huge plus.

However, it isn't all sunshine and rainbows. Pulumi is still young, so documentation is still not great when writing actual code. Their general-purpose docs do make up for this, but it can be frustrating having to jump in and out of the editor to look up what you're currently working on.

## Building the POC

All of the applications are containerized being we can take full advantage of AWS's elastic container service, so the POC was a simple hello world API packed in a docker image, served using fargate and a load balancer in front. I've been aware of ECS and Fargate (plus, vpc, subnets, security groups, etc), but never used them, so there was also some learning involved.

The first thing that became clear is that Terraform has a lot of great documentation and tutorials available. While Pulumi does also have documentation and available tutorials, they are not as "complete" as what you find with Terraform. However, I would find myself hitting a problem in the infrastructure setup with Terraform, switching to Pulumi, and even with less available information, I was able to reason my way out of the problem because I simply had to read Go code. The overhead of a new language and platform would be frustrating, but staying in something I already knew, removed one leg of complexity and the AWS concepts were now the only thing I had to figure out. I suspect that would be one of the main feelings of other engineers who don't do full-time DevOps as well.

I managed to get the POC running using both tools but a big help came from staying in an environment I was already familiar with. I'm a heavy (neo) Vim user and while there are lots of good plugins for Terraform, the Go environment is something I'm in every day so nothing new was needed.

If I dedicated the time, I don't suspect it would take super long to get familiar with HCL and Terraform to the point where the environment would feel close to as familiar as Go. But, given that I wouldn't have to work every day with Terraform after getting the initial setup done, that familiarity would soon start to fade. Meaning, that every time a new infrastructure part had to be added, I would need to re-familiarize myself. And, for the rest of the team, this would be the same but with the added overhead of getting into a part of the code they didn't help write. This could be mitigated by involving everyone from the get-go, but this is also a business, and other features need building too.

Documentation might be the answer some of you are screaming at the screen right now, and I hear you, but documentation quickly goes stale. This is especially true in smaller teams where things move fast, if the majority can be documented by the code itself, you save yourself much trouble.

## Why I choose Pulumi over Terraform

I ended up going with Pulumi. I know from experience that the "correctness" of these choices only shows its true self after you're well past the point of making the switch to the other candidate feasible. But, from my current point of view, this is what will best benefit the team I'm currently working with.

Two questions kept coming up, as I was trying to decide on which tool to use:
1. What are our infrastructure requirements now and in the future? (i.e. are we ever going to need a full-time DevOps)
2. What does the current team know? (i.e. programming languages, cloud providers, etc.)

The first might seem arbitrary, but the likelihood of finding a full-time DevOps engineer who knows Terraform is properly much higher than knowing Go and Pulumi. Given that the needed infrastructure now and in the foreseeable future most likely will not require a full-time DevOps role, terraform begins to feel a bit overkill. It's a new language, and everybody on the team would need to learn it as well as any best practices that come with it. All the while, getting into the nitty gritty of AWS which will also be a new learning experience for the team (myself included).

The second, to me, was very important since I've, from multiple experiences, seen what happens when you have a big terraform setup that maybe one engineer was in charge of creating. Most of the time, it ends up being that one person that has to maintain and fix it, effectively making them a devOps person (which might be against their will -> they switch job -> team is fucked). The rest of the team, despite the person writing documentation, making presentations, etc etc, are still very reluctant to touch it. It's something they are unfamiliar with, the infrastructure system can be quite difficult to "have in your head", you can potentially take down everything without knowing how to get back to a running state. It creates a lot of friction. 

You might argue: "Well, they will just have to learn it" and I agree, that would be a sensible path forward. But we're not trying to solve a technical problem right now, but rather, a human one. And without allocating the time for people to learn it, it probably will not happen. So, instead of a DevOps team, you need to optimize for what gets people to want to learn infrastructure. Over the two weeks, I spent researching and building some POCs, I would often find very good tutorials and materials on setting up infrastructure with Terraform. But as soon as I hit my head against the wall, not knowing a specific syntax/setup/approach used in the tutorial or what the specific AWS services were about, it became very hard to move on. Multiple times, I switched back to Pulumi where their docs had examples in Go code, which was much easier for me to wrap my head around since I already knew Go. I "only" had to figure out the AWS part. After, I would switch back to Terraform and fix the thing with my new understanding from Pulumi/Go. But it's important to note where that knowledge originated from, by reading code, which most of us are doing pretty much every day.

Taking away the overhead of having to both learn a new language and a cloud provider, to only needing to learn the cloud provider, helps a tremendous amount. You also get the added benefit of already knowing best practices, and when and why to break those practices. You get an API from Pulumi so you can start provisioning infrastructure directly in your application(s), which makes spinning up demo environments for sales processes trivial. 

## When I would choose Terraform

Up to a certain point, something becomes a standard because of its quality and appeal to a large part of the tool's target group. After that point, better alternatives might exist but since the standard requirement is to know said tool, everyone else just continues learning that (I'm looking at you, React. Switch those frontends back to HTML and use htmx for interactivity!). If that's what's going on here with Terraform and Pulumi is hard to say. I had a much better experience exploring terraform, compared to previous experiences, where I only had to tweak it a little bit.

If your team includes members who already know Terraform or you know that you have some large-scale infrastructure needs shortly, that will include a dedicated DevOps team, terraform still seems like the better choice. I did pick Terraform and HCL up rather quickly, but I'm still quite confident that this would fade again after not using it daily, so if you don't have that limitation, terraform would win for me.

The main argument from my side here is to remove friction from your team when they don't deal with infrastructure daily, learning a new language while your PM is blasting new tickets your way like his life depended on it, is tough. And having to wait for the one team member who truly knows the infrastructure code coming back from holiday to fix things, is very sub-optimal.

## Was this a good decision

It's too early to tell, I will revisit this in a few months once we've had some more production experience with it and the team has started working more with it as well.

