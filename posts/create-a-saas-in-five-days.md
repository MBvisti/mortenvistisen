## Intro

I love myself a good challenge once in a while. Not the cringy kind you find on MrBeast's channel where they exploit the ever growing economic inequality. No no, the ones that test your skills either mentally, physically or even better, both at the same time.

I've suffered from wantreprenurship for a long time. I've [written about it](/posts/wantrepreneurship-and-getting-out-of-tutorial-hell) and attempted a similar challenge [before](/posts/one-month-one-dollar-part-one) in order to cure it. So far, without much success to report on.

So, doubling down on another challenge, I wanted to take my [saas starter template](https://github.com/mbvlabs/grafto) for a spin and see how well it worked out in real life. I also wanted to get better at creating videos on YouTube, so to add even more work to what is already a tight timeline, I decided to create a video about each day of the challenge. I decided to try and create a micro-saas in 5 days.

## What can you build in five days?

Given my aforementioned affliction, I've long subscribed to all the newsletters that focused on "business ideas" that I could get my hands on. I have previously also tried to get my creative juices flowing by writing down 10 ideas everyday that lasted for long enough time, to have a decent list of not so decent ideas. The whole idea that ideas are important has been busted long time ago. What matters, like really matters, is the execution and _actually_ doing something. Anything.

I've tried to take this to heart. Prioritize getting things out into the wild and see what happens. I'm not going to try and argue for market research and chosing a promising one is going to be a waste of time. But if you spend more than an afternoon (as a developer, at least), you're overdoing it imo.

Anyways, back to my long list of genius ideas. I did spend some time filtering the ideas based on how much they spoke to me and whether I thought they were feasble to do in 5 working days. 

I landed on trying to create a no-code blog.

## Game Plan

I had recently begun to use Traefik in favor of Caddy for reverse proxying my applications. Combined with a docker plugin called `rollout` you get zero-downtime deployments almost for free. This sparked an idea for how I easily could spin up blogs for customers, keeping them updated and restarted when needed.

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

Next, there need to be some way to spin up blogs quickly with various designs, color themes, layouts etc. You might have guessed that this is what happens with the image: mbvofdocker/the-bloggeing:some-version-here and you're right. I choose to only have one design and two color themes, dark and light to begin with. Basically, I'm just re-creating this blog on an automatic basis.

I then needed a way to take in some information about the blog like a name, title, description. Maybe even some socials and have it be dynamically updated whenever the user choose to do so. If you've spend any time in the tech space lately you probably heard about the hype of sqlite and turso. Turso has open-sourced an upgraded version of sqlite called libsql that allows you to talk to your database using http requests and have embedded read replicas. The combination of these two was exactly what I was looking for.

So the components were in place. I would use my startr template, create a docker-compose file for each blog a customer created, route traffic to it using traefik and sync/provide data through turso.

## Calling docker from within docker

I love docker. Having your application containerised makes a lot of sense. But this would end up cause a lot of headaches on day 3.

Everything was working locally, but I develop outside of docker in my machine and only package it up and ship when I push to prod. So calling docker from my application was not an isue until I had to do it from with the docker image that doesn't automatically have access to resources on the host system. It was a fairy easy fix, just add a volume that contained the docker socket and you're good to do.
