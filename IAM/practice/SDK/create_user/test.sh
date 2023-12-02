#!/bin/bash

test_lambda() {
      cd "SDK/create_user" || exit
      aws lambda invoke --function-name IAM_practice_CreateUser --payload file://input.json ./output.json
      echo -e "\n"
      cat ./output.json
      echo -e "\n"
}

test_lambda