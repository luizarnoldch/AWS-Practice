CREATE_USER_CLI := CLI/create_user.sh

CREATE_USER_CF := CF/main.yml
CREATE_CF_STACK_NAME := create-CLI-IAM-test

CREATE_USER_SDK := SDK/main.yml
CREATE_SDK_STACK_NAME := create-SDK-IAM-test

create_cli:
	./${CREATE_USER_CLI}

deploy_create_cf:
	sam deploy --template-file $(CREATE_USER_CF) --stack-name $(CREATE_CF_STACK_NAME) --capabilities CAPABILITY_NAMED_IAM --resolve-s3
destroy_create_cf:
	aws cloudformation delete-stack --stack-name $(CREATE_CF_STACK_NAME)

deploy_create_sdk:
	sam deploy --template-file $(CREATE_USER_SDK) --stack-name $(CREATE_SDK_STACK_NAME) --capabilities CAPABILITY_NAMED_IAM --resolve-s3
destroy_create_sdk:
	aws cloudformation delete-stack --stack-name $(CREATE_SDK_STACK_NAME)
build_create_sdk:
	./SDK/build.sh
test_create_sdk:
	./SDK/test.sh

create_sdk:
	make destroy_create_sdk
	make build_create_sdk
	make deploy_create_sdk
	make test_create_sdk

create_cf:
	make deploy_create_cf
	make destroy_create_cf
