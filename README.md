# GitHub Stars to Notion

An AWS Lambda project that loads starred repos from your GitHub account and pushes them to a Notion database.

This project is built with AWS CDK.

## Motivation

I wanted to be able to keep easy track of useful GitHub projects from an easy-to-use overview.
I found that this would make a good project to improve my skills around AWS, especially the serverless side and Cloudformation.

I also wanted to learn Go.

## AWS services used

- Lambda
- EventBridge
- ECR


## Features

- Repo diffing by checking which repos are starred and which are in the Notion database
- Adding and deleting rows

## Known issues

- Lambda will exit with a Runtime.ExitError, stating the runtime exited without providing a reason
  - The Lambda function will work just fine


## TODO

- Docs on building and deploying the function
