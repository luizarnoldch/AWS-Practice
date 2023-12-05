CREATE_GROUP_CLI := CLI/create_group.sh

CREATE_GROUP_CF := CF/main.yml
CREATE_CF_STACK_NAME := create-CLI-IAM-test-group

CREATE_GROUP_SDK := SDK/create_group/main.yml
CREATE_SDK_STACK_NAME := create-SDK-IAM-test-group

# Create User
create_cli:
	./${CREATE_GROUP_CLI}

deploy_create_cf:
	sam deploy --template-file $(CREATE_GROUP_CF) --stack-name $(CREATE_CF_STACK_NAME) --capabilities CAPABILITY_NAMED_IAM --resolve-s3
destroy_create_cf:
	aws cloudformation delete-stack --stack-name $(CREATE_CF_STACK_NAME)

deploy_create_sdk:
	sam deploy --template-file $(CREATE_GROUP_SDK) --stack-name $(CREATE_SDK_STACK_NAME) --capabilities CAPABILITY_NAMED_IAM --resolve-s3
destroy_create_sdk:
	aws cloudformation delete-stack --stack-name $(CREATE_SDK_STACK_NAME)
build_create_sdk:
	./SDK/create_group/build.sh
test_create_sdk:
	./SDK/create_group/test.sh

create_sdk:
	make destroy_create_sdk
	sleep 10
	make build_create_sdk
	make deploy_create_sdk
	sleep 3
	make test_create_sdk
	make destroy_create_sdk

create_cf:
	make deploy_create_cf
	sleep 10
	make destroy_create_cf

# Delete Group

DELETE_GROUP_CLI := CLI/delete_group.sh

DELETE_GROUP_CF := CF/main.yml
DELETE_CF_STACK_NAME := delete-CLI-IAM-test-group

DELETE_GROUP_SDK := SDK/delete_group/main.yml
DELETE_SDK_STACK_NAME := delete-SDK-IAM-test-group

delete_cli:
	./${DELETE_GROUP_CLI}

deploy_delete_cf:
	sam deploy --template-file $(DELETE_GROUP_CF) --stack-name $(DELETE_CF_STACK_NAME) --capabilities CAPABILITY_NAMED_IAM --resolve-s3
destroy_delete_cf:
	aws cloudformation delete-stack --stack-name $(DELETE_CF_STACK_NAME)

deploy_delete_sdk:
	sam deploy --template-file $(DELETE_GROUP_SDK) --stack-name $(DELETE_SDK_STACK_NAME) --capabilities CAPABILITY_NAMED_IAM --resolve-s3
destroy_delete_sdk:
	aws cloudformation delete-stack --stack-name $(DELETE_SDK_STACK_NAME)
build_delete_sdk:
	./SDK/delete_group/build.sh
test_delete_sdk:
	./SDK/delete_group/test.sh

delete_sdk:
	make destroy_delete_sdk
	sleep 10
	make build_delete_sdk
	make deploy_delete_sdk
	sleep 3
	make test_delete_sdk
	make destroy_delete_sdk

delete_cf:
	make deploy_delete_cf
	sleep 10
	make destroy_delete_cf
	make destroy_delete_sdk


tt:
	make create_sdk
	make delete_sdk