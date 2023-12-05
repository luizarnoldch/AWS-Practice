#!/bin/bash

test_lambda() {
      cd "SDK/create_group" || exit
      aws lambda invoke --function-name IAM_practice_CreateGroup --payload file://input.json --cli-binary-format raw-in-base64-out ./output.json
      echo -e "\n"
      cat ./output.json
      echo -e "\n"
}

test_lambda