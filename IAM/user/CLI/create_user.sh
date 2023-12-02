# Crear usuario IAM
aws iam create-user --user-name nuevoUsuario_CLI

# Crear y descargar las credenciales de acceso (clave de acceso y clave secreta)
aws iam create-access-key --user-name nuevoUsuario_CLI > ./credenciales.json

# Asignar política de acceso a S3
aws iam attach-user-policy --user-name nuevoUsuario_CLI --policy-arn arn:aws:iam::aws:policy/AmazonS3FullAccess
