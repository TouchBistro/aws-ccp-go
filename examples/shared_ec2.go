//go:build example

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/TouchBistro/aws-ccp-go/clients/_ec2"
	"github.com/TouchBistro/aws-ccp-go/providers"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func main() {

	fmt.Println("loading provider")
	sh1, err := providers.NewSharedConfigCredsProvider(context.Background(), "sh1", providers.WithConfigProfile("main"))
	if err != nil {
		fmt.Println("error initializing provider")
		os.Exit(1)
	}

	fmt.Println("initializing ec2 client")
	client, err := _ec2.Client((sh1))
	if err != nil {
		fmt.Println("error initializing ec2 client")
		os.Exit(1)
	}

	fmt.Println("listing 5 ec2 instances")
	out, err := client.DescribeInstances(context.Background(), &ec2.DescribeInstancesInput{
		MaxResults: aws.Int32(5),
	})
	if err != nil {
		fmt.Println("error initializing ec2 client")
		os.Exit(1)
	}

	for _, r := range out.Reservations {
		for _, i := range r.Instances {
			fmt.Println(*i.InstanceId)
		}
	}
}
