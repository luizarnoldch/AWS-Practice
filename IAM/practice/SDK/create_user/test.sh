#!/bin/bash

LAMBDAS=("get_company_by_id" "get_companys")

test_lambda() {
  for folder in "${LAMBDAS[@]}"; do
    (
      aws lambda invoke --function-name IAM_practice_CreateUser --payload file://input.json ./output.json
      echo -e "\n"
      cat ./output.json
      echo -e "\n"
    )
  done
}

test_lambda