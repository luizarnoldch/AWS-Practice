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
    log.Printf("Iniciando la creación del usuario IAM: %s", event.UserName)

    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        log.Fatalf("Error cargando la configuración de AWS: %s", err)
        return "", fmt.Errorf("error cargando la configuración de AWS: %s", err)
    }

    iamClient := iam.NewFromConfig(cfg)

    // Crear usuario IAM
    log.Printf("Creando el usuario IAM: %s", event.UserName)
    createUserOutput, err := iamClient.CreateUser(ctx, &iam.CreateUserInput{
        UserName: aws.String(event.UserName),
    })
    if err != nil {
        log.Fatalf("Error creando el usuario IAM: %s", err)
        return "", fmt.Errorf("error creando el usuario IAM: %s", err)
    }
    log.Printf("Usuario IAM creado: %s", *createUserOutput.User.UserName)

    // Asignar política de acceso a S3
    log.Printf("Asignando política de acceso a S3 al usuario: %s", event.UserName)
    _, err = iamClient.AttachUserPolicy(ctx, &iam.AttachUserPolicyInput{
        PolicyArn: aws.String("arn:aws:iam::aws:policy/AmazonS3FullAccess"),
        UserName:  aws.String(event.UserName),
    })
    if err != nil {
        log.Fatalf("Error asignando política de acceso a S3: %s", err)
        return "", fmt.Errorf("error asignando política de acceso a S3: %s", err)
    }
    log.Println("Política de acceso a S3 asignada correctamente.")

    return fmt.Sprintf("Usuario '%s' creado con acceso a S3", event.UserName), nil
}

func main() {
    lambda.Start(HandleRequest)
}
