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

func HandleRequest(ctx context.Context, event Event) (string, error) {
	log.Printf("Starting the deletion process for IAM user: %s and group: %s", event.UserName, event.UserGroup)

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Error loading AWS configuration: %s", err)
		return "", fmt.Errorf("error loading AWS configuration: %s", err)
	}

	iamClient := iam.NewFromConfig(cfg)

	// Delete access keys for the user
	log.Printf("Deleting access keys for the user: %s", event.UserName)
	accessKeys, err := iamClient.ListAccessKeys(ctx, &iam.ListAccessKeysInput{
		UserName: aws.String(event.UserName),
	})
	if err != nil {
		log.Fatalf("Error listing access keys: %s", err)
		return "", fmt.Errorf("error listing access keys: %s", err)
	}

	for _, accessKey := range accessKeys.AccessKeyMetadata {
		_, err = iamClient.DeleteAccessKey(ctx, &iam.DeleteAccessKeyInput{
			AccessKeyId: accessKey.AccessKeyId,
			UserName:    aws.String(event.UserName),
		})
		if err != nil {
			log.Fatalf("Error deleting access key: %s", err)
			return "", fmt.Errorf("error deleting access key: %s", err)
		}
	}

	// Remove user from all groups
	log.Printf("Removing user from all groups: %s", event.UserName)
	groups, err := iamClient.ListGroupsForUser(ctx, &iam.ListGroupsForUserInput{
		UserName: aws.String(event.UserName),
	})
	if err != nil {
		log.Fatalf("Error listing groups for user: %s", err)
		return "", fmt.Errorf("error listing groups for user: %s", err)
	}

	for _, group := range groups.Groups {
		_, err = iamClient.RemoveUserFromGroup(ctx, &iam.RemoveUserFromGroupInput{
			GroupName: group.GroupName,
			UserName:  aws.String(event.UserName),
		})
		if err != nil {
			log.Fatalf("Error removing user from group: %s", err)
			return "", fmt.Errorf("error removing user from group: %s", err)
		}
	}

	// Delete the IAM user
	log.Printf("Deleting IAM user: %s", event.UserName)
	_, err = iamClient.DeleteUser(ctx, &iam.DeleteUserInput{
		UserName: aws.String(event.UserName),
	})
	if err != nil {
		log.Fatalf("Error deleting IAM user: %s", err)
		return "", fmt.Errorf("error deleting IAM user: %s", err)
	}

	// Detach policies from the group
	log.Printf("Detaching policies from the group: %s", event.UserGroup)
	attachedPolicies, err := iamClient.ListAttachedGroupPolicies(ctx, &iam.ListAttachedGroupPoliciesInput{
		GroupName: aws.String(event.UserGroup),
	})
	if err != nil {
		log.Fatalf("Error listing attached group policies: %s", err)
		return "", fmt.Errorf("error listing attached group policies: %s", err)
	}

	for _, policy := range attachedPolicies.AttachedPolicies {
		_, err = iamClient.DetachGroupPolicy(ctx, &iam.DetachGroupPolicyInput{
			PolicyArn: policy.PolicyArn,
			GroupName: aws.String(event.UserGroup),
		})
		if err != nil {
			log.Fatalf("Error detaching group policy: %s", err)
			return "", fmt.Errorf("error detaching group policy: %s", err)
		}
	}

	// Delete the IAM group
	log.Printf("Deleting IAM group: %s", event.UserGroup)
	_, err = iamClient.DeleteGroup(ctx, &iam.DeleteGroupInput{
		GroupName: aws.String(event.UserGroup),
	})
	if err != nil {
		log.Fatalf("Error deleting IAM group: %s", err)
		return "", fmt.Errorf("error deleting IAM group: %s", err)
	}

	log.Printf("IAM user '%s' and group '%s' deleted successfully", event.UserName, event.UserGroup)
	return fmt.Sprintf("IAM user '%s' and group '%s' deleted successfully", event.UserName, event.UserGroup), nil
}

func main() {
	lambda.Start(HandleRequest)
}
