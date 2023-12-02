# Listar las Claves de Acceso del Usuario
aws iam list-access-keys --user-name nuevoUsuario

# Eliminar las Claves de Acceso
aws iam delete-access-key --access-key-id [ID_CLAVE_ACCESO] --user-name nuevoUsuario

# Listar las Políticas Adjuntas al Usuario
aws iam delete-access-key --access-key-id [ID_CLAVE_ACCESO] --user-name nuevoUsuario

# Desasociar las Políticas del Usuario
aws iam detach-user-policy --user-name nuevoUsuario --policy-arn [ARN_POLITICA]

# Eliminar el Usuario IAM
aws iam delete-user --user-name nuevoUsuario
