package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

type Event struct {
	UserName  string `json:"userName"`
	UserGroup string `json:"userGroup"`
}

// Function to create an IAM group
func createIAMGroup(ctx context.Context, iamClient *iam.Client, groupName string) (*iam.CreateGroupOutput, error) {
	log.Printf("Creating IAM group: %s", groupName)
	createGroupOutput, err := iamClient.CreateGroup(ctx, &iam.CreateGroupInput{
		GroupName: aws.String(groupName),
	})
	if err != nil {
		log.Fatalf("Error creating IAM group: %s", err)
		return nil, fmt.Errorf("error creating IAM group: %s", err)
	}
	log.Printf("IAM group created: %s", groupName)
	return createGroupOutput, nil
}

// Function to assign a policy to an IAM group
func attachPolicyToGroup(ctx context.Context, iamClient *iam.Client, groupName string, policyArn string) error {
	log.Printf("Assigning policy %s to IAM group: %s", policyArn, groupName)
	_, err := iamClient.AttachGroupPolicy(ctx, &iam.AttachGroupPolicyInput{
		GroupName: aws.String(groupName),
		PolicyArn: aws.String(policyArn),
	})
	if err != nil {
		log.Fatalf("Error attaching policy to IAM group: %s", err)
		return fmt.Errorf("error attaching policy to IAM group: %s", err)
	}
	log.Printf("Policy assigned to IAM group: %s", groupName)
	return nil
}

// Function to add a user to an IAM group
func addUserToGroup(ctx context.Context, iamClient *iam.Client, userName string, groupName string) error {
	log.Printf("Adding user '%s' to IAM group: %s", userName, groupName)
	_, err := iamClient.AddUserToGroup(ctx, &iam.AddUserToGroupInput{
		GroupName: aws.String(groupName),
		UserName:  aws.String(userName),
	})
	if err != nil {
		log.Fatalf("Error adding user to IAM group: %s", err)
		return fmt.Errorf("error adding user to IAM group: %s", err)
	}
	log.Printf("User '%s' added to IAM group: %s", userName, groupName)
	return nil
}

func HandleRequest(ctx context.Context, event Event) (string, error) {
	log.Printf("Starting the creation of IAM group: %s", event.UserGroup)
	log.Printf("With IAM user: %s", event.UserName)

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Error loading AWS configuration: %s", err)
		return "", fmt.Errorf("error loading AWS configuration: %s", err)
	}

	groupName := event.UserGroup
	userName := event.UserName
	iamClient := iam.NewFromConfig(cfg)

	// Create IAM group
	_, err = createIAMGroup(ctx, iamClient, groupName)
	if err != nil {
		return "", err
	}

	// Create IAM user
	log.Printf("Creating IAM user: %s", userName)
	createUserOutput, err := iamClient.CreateUser(ctx, &iam.CreateUserInput{
		UserName: aws.String(userName),
	})
	if err != nil {
		log.Fatalf("Error creating IAM user: %s", err)
		return "", fmt.Errorf("error creating IAM user: %s", err)
	}
	log.Printf("IAM user created: %s", *createUserOutput.User.UserName)

	// Create access key
	log.Printf("Creating access key for user: %s", userName)
	createAccessKeyOutput, err := iamClient.CreateAccessKey(ctx, &iam.CreateAccessKeyInput{
		UserName: aws.String(userName),
	})
	if err != nil {
		log.Fatalf("Error creating access key: %s", err)
		return "", fmt.Errorf("error creating access key: %s", err)
	}
	accessKeyID := *createAccessKeyOutput.AccessKey.AccessKeyId
	secretAccessKey := *createAccessKeyOutput.AccessKey.SecretAccessKey
	log.Printf("Access key created: %s", accessKeyID)
	log.Printf("Secret access key created: %s", secretAccessKey)

	// Assign policy to IAM group
	policyArn := "arn:aws:iam::aws:policy/AmazonS3FullAccess" // Example policy ARN
	err = attachPolicyToGroup(ctx, iamClient, groupName, policyArn)
	if err != nil {
		return "", err
	}

	// Add user to IAM group
	err = addUserToGroup(ctx, iamClient, userName, groupName)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("User '%s' created and added to group '%s' with S3 access.", userName, groupName), nil
}

func main() {
	lambda.Start(HandleRequest)
}
