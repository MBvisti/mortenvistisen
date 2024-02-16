
## Networking

```go
var availabilityZones = []string{"us-west-2a", "us-west-2b"}
```

```go
vpc, err := ec2.NewVpc(ctx, "grafto-vpc", &ec2.VpcArgs{
	CidrBlock:          pulumi.String("10.0.0.0/16"),
	EnableDnsHostnames: pulumi.Bool(true),
	EnableDnsSupport:   pulumi.Bool(true),
})
```

```go
networkAccessibility := "private"
if isPublic {
	networkAccessibility = "public"
}

name := fmt.Sprintf("backend-%s-subnet-%v", networkAccessibility, number)

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

	publicSubnet, err := ec2.NewSubnet(ctx, name, &ec2.SubnetArgs{
		VpcId:            vpc.ID(),
		CidrBlock:        pulumi.String(cidrRangePublic),
		AvailabilityZone: pulumi.String(az),
	})

	subnets["public"] = append(subnets["public"], publicSubnet)
	
	publicSubnet, err := ec2.NewSubnet(ctx, name, &ec2.SubnetArgs{
		VpcId:            vpc.ID(),
		CidrBlock:        pulumi.String(cidrRangePublic),
		AvailabilityZone: pulumi.String(az),
	})

	subnets["public"] = append(subnets["public"], publicSubnet)
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


