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
    UserName string `json:"userName"`
}

func HandleRequest(ctx context.Context, event Event) (string, error) {
    log.Printf("Iniciando la eliminación del usuario IAM: %s", event.UserName)

    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        log.Fatalf("Error cargando la configuración de AWS: %s", err)
        return "", fmt.Errorf("error cargando la configuración de AWS: %s", err)
    }

    iamClient := iam.NewFromConfig(cfg)

    // Listar y eliminar las claves de acceso del usuario
    log.Printf("Eliminando claves de acceso para el usuario: %s", event.UserName)
    accessKeys, err := iamClient.ListAccessKeys(ctx, &iam.ListAccessKeysInput{
        UserName: aws.String(event.UserName),
    })
    if err != nil {
        log.Fatalf("Error listando claves de acceso: %s", err)
        return "", fmt.Errorf("error listando claves de acceso: %s", err)
    }

    for _, accessKey := range accessKeys.AccessKeyMetadata {
        _, err = iamClient.DeleteAccessKey(ctx, &iam.DeleteAccessKeyInput{
            AccessKeyId: accessKey.AccessKeyId,
            UserName:    aws.String(event.UserName),
        })
        if err != nil {
            log.Fatalf("Error eliminando clave de acceso: %s", err)
            return "", fmt.Errorf("error eliminando clave de acceso: %s", err)
        }
    }

    // Desasociar políticas del usuario
    log.Printf("Desasociando políticas del usuario: %s", event.UserName)
    attachedPolicies, err := iamClient.ListAttachedUserPolicies(ctx, &iam.ListAttachedUserPoliciesInput{
        UserName: aws.String(event.UserName),
    })
    if err != nil {
        log.Fatalf("Error listando políticas adjuntas: %s", err)
        return "", fmt.Errorf("error listando políticas adjuntas: %s", err)
    }

    for _, policy := range attachedPolicies.AttachedPolicies {
        _, err = iamClient.DetachUserPolicy(ctx, &iam.DetachUserPolicyInput{
            PolicyArn: policy.PolicyArn,
            UserName:  aws.String(event.UserName),
        })
        if err != nil {
            log.Fatalf("Error desasociando política: %s", err)
            return "", fmt.Errorf("error desasociando política: %s", err)
        }
    }

    // Eliminar el usuario IAM
    log.Printf("Eliminando el usuario IAM: %s", event.UserName)
    _, err = iamClient.DeleteUser(ctx, &iam.DeleteUserInput{
        UserName: aws.String(event.UserName),
    })
    if err != nil {
        log.Fatalf("Error eliminando el usuario IAM: %s", err)
        return "", fmt.Errorf("error eliminando el usuario IAM: %s", err)
    }

    log.Printf("Usuario IAM eliminado: %s", event.UserName)
    return fmt.Sprintf("Usuario '%s' eliminado correctamente", event.UserName), nil
}

func main() {
    lambda.Start(HandleRequest)
}
