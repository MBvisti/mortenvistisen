My first interaction with infrastructure as code tools came when I was working for a german start-up, which were heavily invested in Terraform. To this day, I'm pretty the vast majority of the team didn't know the ins and outs of the infrastructure setup. Would this have been different if .

<!--Last year I overcame my (unfounded) disstate of PHP and tried out Laravel; the productivity of Laravel was an eye opener especially after being used to having to write or setup most things from scratch, when developing with Go. But, I love writing Go so I wanted to create something similar with purely in Go which resulted in [Grafto](https://github.com/mbv-labs/grafto). It's still, at the time of writing this, a work in progress but has much of what you need. Authentication, emails, background jobs etc.-->

<!--For one of my clients, I was tasked with consolidating their infrastructure under AWS as they currently use a multiple of cloud providers. IaC tools have long been a thing in the industry with Terraform (in my experience, at least) being the default choice. I often felt that I had to relearn HCL everytime I had to touch infrastructure, which created friction. I wrote a bit about it here [Pulumi vs Terraform](/posts/pulumi-vs-terraform). Since Grafto is focused around Go, and Pulumi has a great SDK for provisioning infrastructure in Go, I wanted to show how you could go about hosting Grafto using Pulumi and AWS.-->

We're going to be using my starter template project, [Grafto](https://github.com/mbv-labs/grafto), as a the basis for this tutorial; since it is containerzed, utilizing services such as ECS and Fargate becomes a breeze. Please note, this will be quite code heavy and will not utilize any design patterns (which, imo, is where pulumi really shines. Will be saved for a future article), so I will only show snippets relevant to the concept being explained. You can find the complete code at the end of the article.

Let's go over how you can setup a production ready infrastructure for Grafto using Pulumi.

## Networking

We'll be creating a somewhat simple network; it should suffice for a long time and follow some best practics. Our application will be running in multiple private networks located in 2 > availability zones for high availability. 

An important note here on the choice of network setup, aws's recommends to keep as much as possible in private networks (i.e. not accessible from the public internet), it increases security. This means that the only entry point will be through an internet gateway which we will create later, but, it also means our application cannot reach out to the internet. This is not always feasible since most app use other apps and need to be able to call their APIs. To allow for this, and still only have one-way communication, we need a nat gateway. This will our applications reach out to the internet but not the other way around. It is also, of course, relatively expensive. Do note that it's possible to place your applications in public subnet, allowing them to reach out to the public internet thus removing the need for the nat gateway. With this approach, you'll need to be very strict with your security group settings and only allowing ingress on ports you trust. At the time of writing, an exploit was just found in OpenSSH so having port 22 public accessible is a security issue.

Alright, let's continue.

We need a few things: a VPC, some subnets to place our resources in, a load balancer to distribute traffic and some gateways. Let's start with creating the VPC and corresponding subnets.

A VPC, virtual private cloud, works like a private network where you can fence in resources. AWS automatically create a default one for you but we will create our own so we can control things like CIDR ranges and how many available ip addresses that are in the network. We'll be using "10.0.0.0/16" which gives us a total of 65.536 addresses.

As for subnets, these can be seen as a postal address, which can be used to increase both security and effiency of communications. As with the VPC, subnets also need a CIDR range. Going back to our earlier discussion on security, public and private networks, we can now limit our potential attack surface by limiting the number of addresses in the public subnet and increase the ones in our private.

Notice how we havea `startingSubnetCidrRange` that ends with `/20`. This gives us a potential of 4.096 number of addresses. And, to ensure high availability, we create 2 private and public subnets in each availability zone. We will also be ignoring the above advice about having more addresses in the private subnets than the public.

Utilizing a rather naive function, we loop over the availability zones we want (in this case only two) and create a private and public subnet in each. Notice how we change the starting cidr range, so that we don't have overlapping addresses.

```go
    // VPC
	vpc, err := ec2.NewVpc(ctx, "grafto-vpc", &ec2.VpcArgs{
		CidrBlock:          pulumi.String("10.0.0.0/16"),
		EnableDnsHostnames: pulumi.Bool(true),
		EnableDnsSupport:   pulumi.Bool(true),
	})
	if err != nil {
		return err
	}

	startingSubnetCidrRange := "10.0.0.0/20"

	availabilityZones := []string{"us-east-1a", "us-east-1b"}

	// SUBNETS
	subnets := make(map[string][]*ec2.Subnet, len(availabilityZones))
	for i, az := range availabilityZones {
		var cidrRangePublic string
		var cidrRangePrivate string
		if i == 0 {
			cidrRangePublic = startingSubnetCidrRange
			cidrRangePrivate = fmt.Sprintf("10.0.%v.0/20", 16)
		} else {
			cidrRangePublic = fmt.Sprintf("10.0.%v.0/20", 16*(i+1))
			cidrRangePrivate = fmt.Sprintf("10.0.%v.0/20", 16*(i+2))
		}

		publicSubnet, err := ec2.NewSubnet(
			ctx,
			fmt.Sprintf("grafto-%s-subnet-%v", "public", i+1),
			&ec2.SubnetArgs{
				VpcId:            vpc.ID(),
				CidrBlock:        pulumi.String(cidrRangePublic),
				AvailabilityZone: pulumi.String(az),
			},
		)
		if err != nil {
			return err
		}

		subnets["public"] = append(subnets["public"], publicSubnet)

		privateSubnet, err := ec2.NewSubnet(
			ctx,
			fmt.Sprintf("grafto-%s-subnet-%v", "private", i+1),
			&ec2.SubnetArgs{
				VpcId:            vpc.ID(),
				CidrBlock:        pulumi.String(cidrRangePrivate),
				AvailabilityZone: pulumi.String(az),
			},
		)
		if err != nil {
			return err
		}

		subnets["private"] = append(subnets["private"], privateSubnet)
	}
```

```go
internetGateway, := err ec2.NewInternetGateway(ctx, payload.resourceName, &ec2.InternetGatewayArgs{
	VpcId: vpc.ID(),
})
	
publicRouteTable, err := ec2.NewRouteTable(ctx, payload.resourceName, &ec2.RouteTableArgs{
	VpcId: vpc.ID(),
})

_, err = ec2.NewRoute(ctx, payload.resourceName, &ec2.RouteArgs{
	DestinationCidrBlock: pulumi.String("0.0.0.0/0"),
	GatewayId: internetGateway.ID(),
	RouteTableId: publicRouteTbl.ID(),
})
	
_, err = ec2.NewRouteTableAssociation(ctx, payload.resourceName, &ec2.RouteTableAssociationArgs{
	RouteTableId: publicRouteTable.ID(),
	SubnetId:     subnets["public"][0].ID(),
})
	
_, err = ec2.NewRouteTableAssociation(ctx, payload.resourceName, &ec2.RouteTableAssociationArgs{
	RouteTableId: publicRouteTable.ID(),
	SubnetId:     subnets["public"][1].ID(),
})
```

TODO: mention to repeat one more time for other private subnet (i.e. -b)
```go
elasticIpA, err := ec2.NewEip(ctx, "elastic-ip-a", &ec2.EipArgs{})
	
natGatewayA, err := ec2.NewNatGateway(ctx, "nat-gateway-a", &ec2.NatGatewayArgs{
	AllocationId: elasticIpA.ID(),
	SubnetId: subnets["public"][0].ID(),
})
	

privateRouteTableA, err := ec2.NewRouteTable(ctx, "private-route-table-a", &ec2.RouteTableArgs{
	VpcId: vpc.ID(),
})

_, err = ec2.NewRoute(ctx, "private-route-a", &ec2.RouteArgs{
	DestinationCidrBlock: pulumi.String("0.0.0.0/0"),
	GatewayId: natGatewayA.ID(),
	RouteTableId: privateRouteTableA.ID(),
})
	
_, err = ec2.NewRouteTableAssociation(ctx, "private-route-association-a", &ec2.RouteTableAssociationArgs{
	RouteTableId: privateRouteTableA.ID(),
	SubnetId:     subnets["public"][0].ID(),
})
```

```go
rdsSubnetGroup, err := rds.NewSubnetGroup(ctx, name, &rds.SubnetGroupArgs{
	SubnetIds: pulumi.StringArray{
		subnets["private"][0].ID(),
		subnets["private"][1].ID(),
	},
})
```

## Security

TODO: explain other sgs needed for ecs and rds

```go
loadBalancerSg, err :=  ec2.NewSecurityGroup(ctx, payload.name, &ec2.SecurityGroupArgs{
	VpcId: vpc.ID(),
	Ingress: ec2.SecurityGroupIngressArray{
		&ec2.SecurityGroupIngressArgs{
			CidrBlocks: []pulumi.StringInput{
				pulumi.String("0.0.0.0/0"),
			},
			FromPort:   pulumi.Int(0),
			ToPort:     pulumi.Int(0),
			Protocol:   pulumi.String("-1"),
		},
	},
	Egress: ec2.SecurityGroupEgressArray{
		&ec2.SecurityGroupEgressArgs{
			CidrBlocks: []pulumi.StringInput{
				pulumi.String("0.0.0.0/0"),
			},
			FromPort:   pulumi.Int(0),
			ToPort:     pulumi.Int(0),
			Protocol:   pulumi.String("-1"),
		},
	},
})
```

## Database

```go
rds, err := rds.NewInstance(ctx, "grafto-rds-psql", &rds.InstanceArgs{
	AllocatedStorage:    pulumi.Int(10),
	DbName:              pulumi.String("grafto"),
	Password:            pulumi.String("password"),
	Username:            pulumi.String("grafto"),
	Engine:              pulumi.String("postgres"),
	EngineVersion:       pulumi.String("15.0"),
	InstanceClass:       pulumi.String("db.t2.micro"),
	ParameterGroupName:  pulumi.String("default.postgres15"),
	DbSubnetGroupName:   rdsSubnetGroup.Name,  
	VpcSecurityGroupIds: pulumi.StringArray{
		rdsSecurityGroup.ID(),
	},
	SkipFinalSnapshot:   pulumi.Bool(true),
	PubliclyAccessible:  pulumi.Bool(false),
})
	
ctx.Export("grafto-db-address", rds.Address)
```

## Load balancing

```go
loadBalancer, err := lb.NewLoadBalancer(ctx, "grafto-load-balancer", &lb.LoadBalancerArgs{
	Internal:                 pulumi.Bool(false),
	LoadBalancerType:         pulumi.String("application"),
	SecurityGroups:           pulumi.StringArray{
			albSecurityGroup.ID(),
	},
	Subnets:                 pulumi.StringArray{
		subnets["public"][0].ID(),
		subnets["public"][1].ID(),
	}, 
	EnableDeletionProtection: pulumi.Bool(false),
})
```

```go
targetGroup, err := lb.NewTargetGroup(ctx, "target-group", &lb.TargetGroupArgs{
	HealthCheck: &lb.TargetGroupHealthCheckArgs{
		Path: pulumi.String("/api/health"),
		Protocol: pulumi.String(payload.protocol),
	},
	Name: pulumi.String(payload.name),
	Port: pulumi.Int(payload.port),
	Protocol:        pulumi.String(payload.protocol),
	TargetType: pulumi.String(payload.targetType),
	VpcId:      payload.vpcID,
})
```

```go
listener, err := lb.NewListener(ctx, payload.resourceName, &lb.ListenerArgs{
	DefaultActions: lb.ListenerDefaultActionArray{
		lb.ListenerDefaultActionArgs{
			TargetGroupArn: targetGroup.Arn,
			Type:           pulumi.String("forward"),
		},
	},
	LoadBalancerArn: loadBalancer.Arn,
	Port:     pulumi.Int(80),
	Protocol: pulumi.String("HTTP"),
})
```

## ECS


