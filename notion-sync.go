package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	lambda "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"os"
)

type NotionSyncStackProps struct {
	awscdk.StackProps
}

func NewNotionSyncStack(scope constructs.Construct, id string, props *NotionSyncStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	functionName := jsii.String("rvx-lbd-as2-notion-sync")

	notionApiKey := os.Getenv("NOTION_API_KEY")
	notionDatabaseId := os.Getenv("NOTION_DATABASE_ID")
	githubUser := os.Getenv("GITHUB_USER")

	environmentVars := map[string]*string{
		"NOTION_API_KEY":     &notionApiKey,
		"NOTION_DATABASE_ID": &notionDatabaseId,
		"GITHUB_USER":        &githubUser,
	}

	lambda.NewDockerImageFunction(stack, functionName, &lambda.DockerImageFunctionProps{
		Architecture: lambda.Architecture_ARM_64(),
		FunctionName: functionName,
		Code:         lambda.DockerImageCode_FromImageAsset(jsii.String("./lambda"), nil),
		Environment:  &environmentVars,
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewNotionSyncStack(app, "NotionSyncStack", &NotionSyncStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
