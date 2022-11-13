set dotenv-load

# Deploy the stack
deploy:
    cdk deploy -c notionApiKey="$NOTION_API_KEY" -c notionDatabaseId="$NOTION_DATABASE_ID" -c githubUser="$GITHUB_USER"

# Format code
format:
    go fmt

# Run tests in the lambda directory
test-lambda:
    #!/usr/bin/env bash
    cd lambda
    go test ./...

# Run tests in the lambda directory and show coverage
test-lambda-cover:
    #!/usr/bin/env bash
    cd lambda
    go test -coverprofile=coverage.html ./...
    go tool cover -html=coverage.html

# Remove coverage files
clean-cover:
    rm lambda/coverage.html

# Remove CDK files
clean:
    rm -r ./cdk.out
