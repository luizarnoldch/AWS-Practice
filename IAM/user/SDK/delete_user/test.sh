#!/bin/bash

test_lambda() {
      cd "SDK/delete_user" || exit
      aws lambda invoke --function-name IAM_practice_DeleteUser --payload file://input.json --cli-binary-format raw-in-base64-out ./output.json
      echo -e "\n"
      cat ./output.json
      echo -e "\n"
}

test_lambda