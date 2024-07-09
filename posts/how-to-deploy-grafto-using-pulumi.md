Prompted by client work where I had to consolidate their infrastructure on AWS I was left with a question if I should use an IaC tool, andif yes, which one. I wrote a bit about that decision process [here](https://mortenvistisen.com/posts/pulumi-vs-terraform). In an continuous attempt to throw the real world at my solutions and decisions, I thought it would be interesting, to see how one could go about hosting 
[Grafto](https://github.com/mbv-labs/grafto), using Pulumi on AWS. Grafto is a small starter template project of mine, that is containerzed, so utilizing services such as ECS and Fargate becomes a breeze.

This is a rather naive approach that doesn't utilize a lot of Pulumi's strengths, like allowing us to use design patterns, when building out infrastructure. But should, hopefully, illustrate the benefits of having your infrastructure as code in code you actually read and write everyday.

## Wall Of Code

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ecs"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/lb"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/rds"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		availabilityZones := []string{"us-east-1a", "us-east-1b"}

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

		// INTERNET GATEWAY
		internetGateway, err := ec2.NewInternetGateway(
			ctx,
			"grafto-internet-gateway",
			&ec2.InternetGatewayArgs{
				VpcId: vpc.ID(),
			},
		)
		if err != nil {
			return err
		}

		publicRouteTable, err := ec2.NewRouteTable(
			ctx,
			"grafto-public-route-table",
			&ec2.RouteTableArgs{
				VpcId: vpc.ID(),
			},
		)
		if err != nil {
			return err
		}

		_, err = ec2.NewRoute(ctx, "grafto-public-route", &ec2.RouteArgs{
			DestinationCidrBlock: pulumi.String("0.0.0.0/0"),
			GatewayId:            internetGateway.ID(),
			RouteTableId:         publicRouteTable.ID(),
		})
		if err != nil {
			return err
		}

		_, err = ec2.NewRouteTableAssociation(
			ctx,
			"grafto-public-route-ass-1",
			&ec2.RouteTableAssociationArgs{
				RouteTableId: publicRouteTable.ID(),
				SubnetId:     subnets["public"][0].ID(),
			},
		)
		if err != nil {
			return err
		}

		_, err = ec2.NewRouteTableAssociation(
			ctx,
			"grafto-public-route-ass-2",
			&ec2.RouteTableAssociationArgs{
				RouteTableId: publicRouteTable.ID(),
				SubnetId:     subnets["public"][1].ID(),
			},
		)
		if err != nil {
			return err
		}

		// NATGATEWAY
		elasticIP, err := ec2.NewEip(ctx, "grafto-elastic-ip", &ec2.EipArgs{})
		if err != nil {
			return err
		}

		natGateway, err := ec2.NewNatGateway(ctx, "grafto-nat-gateway", &ec2.NatGatewayArgs{
			AllocationId: elasticIP.ID(),
			SubnetId:     subnets["public"][0].ID(),
		})
		if err != nil {
			return err
		}

		privateRouteTable, err := ec2.NewRouteTable(
			ctx,
			"grafto-private-route-table",
			&ec2.RouteTableArgs{
				VpcId: vpc.ID(),
			},
		)
		if err != nil {
			return err
		}

		_, err = ec2.NewRoute(ctx, "grafto-private-route", &ec2.RouteArgs{
			DestinationCidrBlock: pulumi.String("0.0.0.0/0"),
			NatGatewayId:         natGateway.ID(),
			RouteTableId:         privateRouteTable.ID(),
		})
		if err != nil {
			return err
		}

		_, err = ec2.NewRouteTableAssociation(
			ctx,
			"grafto-private-route-ass-1",
			&ec2.RouteTableAssociationArgs{
				RouteTableId: privateRouteTable.ID(),
				SubnetId:     subnets["private"][0].ID(),
			},
		)
		if err != nil {
			return err
		}

		_, err = ec2.NewRouteTableAssociation(
			ctx,
			"grafto-private-route-ass-2",
			&ec2.RouteTableAssociationArgs{
				RouteTableId: privateRouteTable.ID(),
				SubnetId:     subnets["private"][1].ID(),
			},
		)
		if err != nil {
			return err
		}

		// SECURITY GROUP
		applicationLoadBalancer, err := ec2.NewSecurityGroup(
			ctx,
			"grafto-alb-sg",
			&ec2.SecurityGroupArgs{
				VpcId: vpc.ID(),
				Ingress: ec2.SecurityGroupIngressArray{
					&ec2.SecurityGroupIngressArgs{
						CidrBlocks: pulumi.StringArray{
							pulumi.String("0.0.0.0/0"),
						},
						FromPort: pulumi.Int(80),
						ToPort:   pulumi.Int(80),
						Protocol: pulumi.String("tcp"),
					},
				},
				Egress: ec2.SecurityGroupEgressArray{
					&ec2.SecurityGroupEgressArgs{
						CidrBlocks: pulumi.StringArray{
							pulumi.String("0.0.0.0/0"),
						},
						FromPort: pulumi.Int(0),
						ToPort:   pulumi.Int(0),
						Protocol: pulumi.String("-1"),
					},
				},
			},
		)
		if err != nil {
			return err
		}

		ecsSG, err := ec2.NewSecurityGroup(
			ctx,
			"grafto-ecs-sg",
			&ec2.SecurityGroupArgs{
				VpcId: vpc.ID(),
				Ingress: ec2.SecurityGroupIngressArray{
					&ec2.SecurityGroupIngressArgs{
						CidrBlocks: pulumi.StringArray{
							pulumi.String("0.0.0.0/0"),
						},
						FromPort: pulumi.Int(0),
						ToPort:   pulumi.Int(0),
						Protocol: pulumi.String("-1"),
					},
				},
				Egress: ec2.SecurityGroupEgressArray{
					&ec2.SecurityGroupEgressArgs{
						CidrBlocks: pulumi.StringArray{
							pulumi.String("0.0.0.0/0"),
						},
						FromPort: pulumi.Int(0),
						ToPort:   pulumi.Int(0),
						Protocol: pulumi.String("-1"),
					},
				},
			},
		)
		if err != nil {
			return err
		}

		rdsSGG, err := ec2.NewSecurityGroup(
			ctx,
			"grafto-rds-sgg",
			&ec2.SecurityGroupArgs{
				VpcId: vpc.ID(),
				Ingress: ec2.SecurityGroupIngressArray{
					&ec2.SecurityGroupIngressArgs{
						CidrBlocks: pulumi.StringArray{
							pulumi.String("0.0.0.0/0"),
						},
						FromPort: pulumi.Int(0),
						ToPort:   pulumi.Int(0),
						Protocol: pulumi.String("-1"),
					},
				},
				Egress: ec2.SecurityGroupEgressArray{
					&ec2.SecurityGroupEgressArgs{
						CidrBlocks: pulumi.StringArray{
							pulumi.String("0.0.0.0/0"),
						},
						FromPort: pulumi.Int(0),
						ToPort:   pulumi.Int(0),
						Protocol: pulumi.String("-1"),
					},
				},
			},
		)
		if err != nil {
			return err
		}

		rdsSg, err := rds.NewSubnetGroup(ctx, "grafto-rds-sg", &rds.SubnetGroupArgs{
			SubnetIds: pulumi.StringArray{
				subnets["private"][0].ID(),
				subnets["private"][1].ID(),
			},
		})
		if err != nil {
			return err
		}

		database, err := rds.NewInstance(ctx, "grafto-rds-psql", &rds.InstanceArgs{
			AllocatedStorage:   pulumi.Int(10),
			DbName:             pulumi.String("grafto"),
			Password:           pulumi.String("password"),
			Username:           pulumi.String("grafto"),
			Engine:             pulumi.String("postgres"),
			EngineVersion:      pulumi.String("16.3"),
			InstanceClass:      pulumi.String("db.t3.micro"),
			ParameterGroupName: pulumi.String("default.postgres16"),
			DbSubnetGroupName:  rdsSg.Name,
			VpcSecurityGroupIds: pulumi.StringArray{
				rdsSGG.ID(),
			},
			SkipFinalSnapshot:  pulumi.Bool(true),
			PubliclyAccessible: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}

		loadBalancer, err := lb.NewLoadBalancer(ctx, "grafto-load-balancer", &lb.LoadBalancerArgs{
			Internal:         pulumi.Bool(false),
			LoadBalancerType: pulumi.String("application"),
			SecurityGroups: pulumi.StringArray{
				applicationLoadBalancer.ID(),
			},
			Subnets: pulumi.StringArray{
				subnets["public"][0].ID(),
				subnets["public"][1].ID(),
			},
			EnableDeletionProtection: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}
		ctx.Export("url", pulumi.Sprintf("http://%s", loadBalancer.DnsName))

		targetGroup, err := lb.NewTargetGroup(ctx, "grafto-alb-target-group", &lb.TargetGroupArgs{
			HealthCheck: &lb.TargetGroupHealthCheckArgs{
				Path:     pulumi.String("/api/health"),
				Protocol: pulumi.String("HTTP"),
			},
			Name:       pulumi.String("grafto-app-tg"),
			Port:       pulumi.Int(80),
			Protocol:   pulumi.String("HTTP"),
			TargetType: pulumi.String("ip"),
			VpcId:      vpc.ID(),
		})
		if err != nil {
			return err
		}

		_, err = lb.NewListener(ctx, "grafto-alb-listener", &lb.ListenerArgs{
			DefaultActions: lb.ListenerDefaultActionArray{
				lb.ListenerDefaultActionArgs{
					TargetGroupArn: targetGroup.Arn,
					Type:           pulumi.String("forward"),
				},
			},
			LoadBalancerArn: loadBalancer.Arn,
			Port:            pulumi.Int(80),
			Protocol:        pulumi.String("HTTP"),
		})
		if err != nil {
			return err
		}

		// IAM RELATED STUFF
		_, err = iam.NewServiceLinkedRole(
			ctx,
			"elastic-container-service",
			&iam.ServiceLinkedRoleArgs{
				AwsServiceName: pulumi.String("ecs.amazonaws.com"),
				Description:    pulumi.String("Role to enable Amazon ECS to manage your cluster."),
			},
		)
		if err != nil {
			return err
		}

		_, err = iam.NewServiceLinkedRole(ctx, "rds", &iam.ServiceLinkedRoleArgs{
			AwsServiceName: pulumi.String("rds.amazonaws.com"),
			Description:    pulumi.String("Role to enable Amazon RDS to manage your cluster."),
		})
		if err != nil {
			return err
		}

		_, err = iam.NewServiceLinkedRole(ctx, "elastic-load-balancer", &iam.ServiceLinkedRoleArgs{
			AwsServiceName: pulumi.String("elasticloadbalancing.amazonaws.com"),
			Description:    pulumi.String("Allows ELB to call AWS services on your behalf"),
		})
		if err != nil {
			return err
		}

		_, err = iam.NewServiceLinkedRole(
			ctx,
			"application-autoscaling",
			&iam.ServiceLinkedRoleArgs{
				AwsServiceName: pulumi.String("ecs.application-autoscaling.amazonaws.com"),
				Description: pulumi.String(
					"Allows application autoscaling to call AWS services on your behalf",
				),
			},
		)
		if err != nil {
			return err
		}

		roleJson, err := json.Marshal(map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []map[string]interface{}{
				{
					"Action": []string{
						"sts:AssumeRole",
					},
					"Principal": map[string]string{"Service": "ecs-tasks.amazonaws.com"},
					"Effect":    "Allow",
				},
			},
		})
		if err != nil {
			return err
		}
		role, err := iam.NewRole(ctx, "grafto-iam-role", &iam.RoleArgs{
			Name:             pulumi.String("grafto-iam-role"),
			AssumeRolePolicy: pulumi.String(string(roleJson)),
		})
		if err != nil {
			return err
		}

		rolePolicyJson, err := json.Marshal(map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []map[string]interface{}{
				{
					"Action": []string{
						"ecr:*",
					},
					"Effect":   "Allow",
					"Resource": "*",
				},
			},
		})
		if err != nil {
			return err
		}
		_, err = iam.NewRolePolicy(ctx, "grafto-iam-role-policy", &iam.RolePolicyArgs{
			Name:   pulumi.String("grafto-iam-role"),
			Role:   role.Name,
			Policy: pulumi.String(string(rolePolicyJson)),
		})
		if err != nil {
			return err
		}

		// ELASTIC CONTAINER SERVICE
		cluster, err := ecs.NewCluster(ctx, "grafto-ecs-cluster", &ecs.ClusterArgs{
			Name: pulumi.String("grafto"),
		})
		if err != nil {
			return err
		}

		taskContainerDefinition := pulumi.JSONMarshal([]map[string]interface{}{
			{
				"name":  "grafto-task",
				"image": "docker.io/mbvofdocker/grafto:pulumi-blog",
				"portMappings": []map[string]interface{}{
					{
						"containerPort": 8080,
						"hostPort":      8080,
						"protocol":      "HTTP",
					},
				},
				"essential": true,
				"command":   []string{"./app"},
				"environment": []map[string]interface{}{
					{
						"name":  "ENVIRONMENT",
						"value": "production",
					},
					{
						"name":  "SERVER_HOST",
						"value": "0.0.0.0",
					},
					{
						"name":  "SERVER_PORT",
						"value": "8080",
					},
					{
						"name":  "DEFAULT_SENDER_SIGNATURE",
						"value": "noreply@mortenvistisen.com",
					},
					{
						"name":  "POSTMARK_API_TOKEN",
						"value": "insert-valid-token-here",
					},
					{
						"name":  "DB_KIND",
						"value": "postgres",
					},
					{
						"name":  "DB_PORT",
						"value": "5432",
					},
					{
						"name": "DB_HOST",
						"value": database.Address.ApplyT(
							func(addr string) string {
								return addr
							},
						).(pulumi.StringOutput),
					},
					{
						"name": "DB_NAME",
						"value": database.DbName.ApplyT(
							func(name string) string {
								return name
							},
						).(pulumi.StringOutput),
					},
					{
						"name": "DB_USER",
						"value": database.Username.ApplyT(
							func(name string) string {
								return name
							},
						).(pulumi.StringOutput),
					},
					{
						"name": "DB_PASSWORD",
						"value": database.Password.ApplyT(
							func(pass *string) string {
								return *pass
							},
						).(pulumi.StringOutput),
					},
					{
						"name":  "DB_SSL_MODE",
						"value": "require",
					},
					{
						"name":  "PASSWORD_PEPPER",
						"value": "lotsandlotsofrandomcharshere",
					},
					{
						"name":  "PROJECT_NAME",
						"value": "Pulumi Grafto BLog Post",
					},
					{
						"name": "APP_HOST",
						"value": loadBalancer.DnsName.ApplyT(func(url string) string {
							return url
						}),
					},
					{
						"name":  "APP_SCHEME",
						"value": "http",
					},
					{
						"name":  "CSRF_TOKEN",
						"value": "lotsandlotsofrandomcharshere",
					},
					{
						"name":  "SESSION_KEY",
						"value": "lotsandlotsofrandomcharshere",
					},
					{
						"name":  "SESSION_ENCRYPTION_KEY",
						"value": "lotsandlotsofrandomcharshere",
					},
					{
						"name":  "TOKEN_SIGNING_KEY",
						"value": "lotsandlotsofrandomcharshere",
					},
				},
			},
		})
		taskDefinition, err := ecs.NewTaskDefinition(ctx, "grafto-task", &ecs.TaskDefinitionArgs{
			ContainerDefinitions: taskContainerDefinition,
			Cpu:                  pulumi.String("256"),
			ExecutionRoleArn:     role.Arn,
			Family:               pulumi.String("grafto"),
			Memory:               pulumi.String("512"),
			NetworkMode:          pulumi.String("awsvpc"),
			TaskRoleArn:          role.Arn,
		})
		if err != nil {
			return err
		}

		_, err = ecs.NewService(ctx, "grafto-service", &ecs.ServiceArgs{
			Cluster:                         cluster.Arn,
			DeploymentMaximumPercent:        pulumi.IntPtr(200),
			DeploymentMinimumHealthyPercent: pulumi.IntPtr(50),
			DesiredCount:                    pulumi.IntPtr(1),
			ForceNewDeployment:              pulumi.Bool(true),
			LoadBalancers: ecs.ServiceLoadBalancerArray{
				&ecs.ServiceLoadBalancerArgs{
					TargetGroupArn: targetGroup.Arn,
					ContainerName:  pulumi.String("grafto-task"),
					ContainerPort:  pulumi.Int(8080),
				},
			},
			NetworkConfiguration: ecs.ServiceNetworkConfigurationArgs{
				Subnets: pulumi.StringArray{
					subnets["private"][0].ID(),
					subnets["private"][1].ID(),
				},
				SecurityGroups: pulumi.StringArray{
					ecsSG.ID(),
				},
			},
			Name:            pulumi.String("grafto-ecs-service"),
			LaunchType:      pulumi.String("FARGATE"),
			PlatformVersion: pulumi.String("1.4.0"),
			TaskDefinition:  taskDefinition.Arn,
		})
		if err != nil {
			return err
		}

		return nil
	})
}
```

## Improvements

An obvious improvement to the above would be to enable HTTPS; if you check the load balancer's security group you can see that we allow ingress traffic on port 80. This is the only entrypoint since our Fargate tasks are all in private networks, so adding a certicate would limiting it to port 443 could go a long way.

Take a look at the calculations of the cidr ranges. If we add too many availability zones this will fail which is something to handle as well. Could be a simple check on how many AZs is requsted and limited it to a certain level, but should still be fixed.

It would also be beneficial to store the environmental variables somewhere like AWS's parameter store, and not directly in the code.

You'll probably also have noticed multiple opportunities for re-using code, through setup functions or, my personal favorite in this case, builders. Builders can simply the code a lot, especially if the amount of tasks you've in your ecs service increase. In a future article we'll improve upon this so we can easily expand upon our intrastructure.

An interesting comparison would be to do the same for Terraform and see how much they differ, and if the effort in making the infrastructure code resuable with different design patterns make sense in the end.

But for now, that's all. Happy hacking!
