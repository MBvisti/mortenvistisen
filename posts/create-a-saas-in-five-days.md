I love myself a good challenge once in a while. Not the cringy kind you find on MrBeast's channel where they exploit the ever-growing economic inequality. But a simple "Build X thing in Y amount of time" kind of challenge.

I've suffered from wantreprenurship for a long time. I've [written about it](/posts/wantrepreneurship-and-getting-out-of-tutorial-hell) and attempted a similar challenge [before](/posts/one-month-one-dollar-part-one) to cure it. So far, I have not had much success to report on.

So, doubling down on another challenge, I wanted to take my [saas starter template](https://github.com/mbvlabs/grafto) for a spin and see how well it worked out in real life. I also wanted to get better at creating videos on YouTube, so to add even more work to what is already a tight timeline, I decided to create a video about each day of the challenge. 

I decided to try and create a micro-saas in 5 days.

## What can you build in five days?

Given my aforementioned affliction, I've long subscribed to all the newsletters that focused on "business ideas" that I could get my grubby little hands on. 

I have previously also tried to get my creative juices flowing by writing down 10 ideas every day that lasted for a long enough time, to have a decent list of not-so-decent ideas. The whole idea that ideas are important was busted a long time ago. What matters, like really matters, is the execution and _actually_ doing something. Anything.

I've tried to take this to heart. Prioritize getting things out into the wild and see what happens. I'm not going to try and argue that doing some market research beforehand and choosing the most promising one is going to be a waste of time. But when just starting out, you probably shouldn't be spending much time here as you still (like me) lack the gut feel for where to focus your attention. Look around, and see what other people are doing and making money on. Pick something that matches your interest and skillset, try to create a competitor and take it from there.

Anyway, back to my long list of genius ideas. I did spend some time filtering the ideas based on how much they spoke to me and whether I thought they were feasible to do in 5 working days. The time constraint did help me cut a lot of ideas, many of which would involve acquiring lots of data for the AI overloads to gobble up. So what could I, realistically, build and ship a scrappy version of in five days?

A no-code blog.

## Game Plan

I've written and rewritten my own blog many times now. Started out in next.js back when I was huge into JavaScript/TypeScript and had my "portfolio" website with some light blogging content, to learning and rewriting it in Rust to finally just settling on Go. The language I've been working in for the last many years. I think it's one of the best ways for developers to get a sense of how to build and ship software, especially early on, as it takes you through all of the steps that happen in building web apps.

Anyway, I have a fair bit of experience writing blogging software.

I had recently begun to use Traefik in favour of Caddy for reverse proxying my applications. Combined with a docker plugin called `rollout` you get zero-downtime deployments basically for free. This sparked an idea for how I easily could spin up blogs for customers, keeping them updated and restarted when needed.

I would simply have a template file that would parse some variables, create the file and store it in a specific place and call docker-compose pointing to this file. Easy. This ended up looking something like this:

```docker
services:
  blog--id-{{.BlogID}}:
    image: mbvofdocker/the-bloggening:{{.Version}}
    environment:
      - ENV_VAR_1={{.EnvVarOne}}
      - ENV_VAR_2={{.EnvVarTwo}}
    labels:
      - "traefik.enable=true"
      - "traefik.http.services.blog-{{.BlogID}}.loadbalancer.server.port=3333"
```

I've left labels and fields out for brevity but this should get the idea across. 

Traefik would take care of discovery and routing traffic to the individual blogs, generating SSL certs and taking the health status of each container.

Next, there needs to be some way to spin up blogs quickly with various designs, color themes, layouts etc. You might have guessed that this is what happens in this part of the above docker-compose file: `mbvofdocker/the-bloggeing:{{.Version}}` and you're right. 

I went with only having one design and two colour themes, dark and light, to begin with. Basically, I'm just re-creating this blog you're reading this article on, but on an automatic basis.

Then, I needed a way to take in some information about the blog like a name, title, description etc. Maybe even some socials and have it be dynamically updated whenever the user chooses to do so. 

If you've spent any time in the tech space lately you probably heard about the hype of SQLite and Turso. Turso has open-sourced an upgraded version of SQLite called libsql that allows you to talk to your database using HTTP requests and have embedded read replicas. The combination of these two was exactly what I was looking for.

So the components and plan were in place. I would use my starter template, create a docker-compose file for each blog a customer created, route traffic to it using traefik and sync/provide data through Turso. But, as a wise man once said:
> Everybody has a plan until they get punched in the mouth

## Where are my databases?

I was putting the above together fast, like fast fast. I could definitely make this in those 5 days, no problem. I had a working demo and wanted to have my wife try and create her own blog, so I enthusiastically forced her to drop whatever she was doing to sit down and be an equally enthusiastic test user.

She went through the signup flow, with no issues, same for the onboarding. I had taken a very liberal approach to error handling at this point, which basically meant I only logged when stuff errored. But, nothing happened. It was reporting overall success and I could see the container image was running and healthy. So we went onto her blog that was being served locally on `mortens-wifes-blog.localhost`, aaaaand nothing. 

Not wanting her to think that the financial success I had practically guaranteed from this project was a pipe dream, I frantically looked for why. Turns on, that since I was on the free version of Turso they would sometimes put my "cluster" to sleep. However, their API didn't return any errors or notice that the database hadn't been created. A bit of a bummer, but again, free software so I can't really complain much. Also, I realized that as long as I used my own database to hold the information I could just as easily sync the data with each blog through an endpoint staying with regular sqlite.

I made the changes, had her sign up again and everything worked. Great. Then, my next (and much bigger headache): exit code 18.

## Calling docker from within docker

I love Docker. Having your application containerised makes a lot of sense. In this case, it ended up causing a lot of frustration.

Everything was working locally, but I develop outside of docker on my machine and packaged it up and ship when I push to production. So calling docker from my application was not an issue until I had to do it from with the docker image that doesn't automatically have access to resources on the host system. It was a fairly easy fix, just add a volume that contained the docker socket and you're good to go.

```docker
volume:
  - /var/run/docker.sock:/var/run/docker.sock
```

However, I kept getting an exit error from the system that I couldn't for the life of me figure out what meant. Nothing came from googling and asking chatgippity. I even tried Stackoverflow to no avail. I then replicated the entire setup on my local machine but with everything running in docker, where it worked, which left me with more questions than answers.

Turned out, that I hadn't pulled the correct version of the docker image on the server so it couldn't start. Normally, when an image is not available, docker will try and pull it from some hub. I, rather naively, assumed the same would happen here which it of course didn't since permissions are different.

This would cause so much delay that I was super tight on time to create the onboarding experience for which I would spend all of day 5.

In the end, I had a version live that allowed users to sign up, create a blog and see something live, e.g. mortens-blog.mbvlabs.com. But, since I was so strapped for time on day 5 I didn't get to implement the functionality that would allow users to upload or write blog posts. Saying I finished without that feature is probably pushing it but I got super close and have since added it to the setup. The product still needs work but pushing through those 5 days gave it a quick foundation on which to stand. Retrospectively, I probably wouldn't do this and have to create a video about it for each day. Just mentally having to build a thing and also make it work in a story that shows progression is quite tough. You can't just speed through the building process. You need to get relevant clips and snippets of videos throughout the entire process. But it was super fun.

But, what's the takeaway?

## Urgency is good

Having publically, even to a rather small audience, announced an end date was incredibly beneficial in cutting into the core of what the product is and what was needed. I realized multiple times during that week that I was spending time on things or adding stuff that wasn't really necessary at the core of what I wanted to do: let users quickly spin up a good-looking blog and add content.

I think that, in general, urgency is a good thing in software. It's super easy to get into a long drawn-out analysis of what is needed right now vs what is needed in the future. Spend time only focusing on the now and you might have a hard time evolving the thing in the future. Spend time preparing for the future and it will be easier to deal with future requirements. If you're right about those future requirements. And that's the core of it, you need to be correct in your assumptions of the future which most people rarely are. I think optimizing for the now will let you get something out in the real world that users can test and give feedback on. As long as you follow good practices such as the single responsibility principle, separation of concern and YAGNI, dealing with future requirements most likely wouldn't be a huge problem. And, should it turn out to be so and you need to rewrite a lot of code, isn't that a good thing? This means you learned something about the domain in which you are working and that is ultimately what should drive the software.

Stop overthinking, ship some things and see what happens.
