package main

import (
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	eventBridge "github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	targets "github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	lambda "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"os"
)

type GHStars2NotionProps struct {
	awscdk.StackProps
}

func NewGHStars2NotionStack(scope constructs.Construct, id string, props *GHStars2NotionProps) awscdk.Stack {
	var sProps awscdk.StackProps

	if props != nil {
		sProps = props.StackProps
	}

	stack := awscdk.NewStack(scope, &id, &sProps)
	node := stack.Node()

	var notionApiKey, notionDatabaseId, githubUser string

	if val, err := getContextAsString(&node, "notionApiKey"); err != nil {
		panic(err)
	} else {
		notionApiKey = val
	}

	if val, err := getContextAsString(&node, "notionDatabaseId"); err != nil {
		panic(err)
	} else {
		notionDatabaseId = val
	}

	if val, err := getContextAsString(&node, "githubUser"); err != nil {
		panic(err)
	} else {
		githubUser = val
	}

	environmentVars := map[string]*string{
		"NOTION_API_KEY":     &notionApiKey,
		"NOTION_DATABASE_ID": &notionDatabaseId,
		"GITHUB_USER":        &githubUser,
	}

	functionName := jsii.String("rvx-lbd-as2-ghstars2notion-gh-notion-sync")

	lambdaFunc := lambda.NewDockerImageFunction(stack, functionName, &lambda.DockerImageFunctionProps{
		Architecture: lambda.Architecture_ARM_64(),
		Code:         lambda.DockerImageCode_FromImageAsset(jsii.String("./lambda"), nil),
		Description:  jsii.String("Lambda function that syncs GitHub stars to a Notion database"),
		Environment:  &environmentVars,
		FunctionName: functionName,
		Timeout:      awscdk.Duration_Seconds(jsii.Number(10)),
	})

	ebRuleName := jsii.String("rvx-evb-as2-ghstars2notion")

	eventBridge.NewRule(stack, ebRuleName, &eventBridge.RuleProps{
		Schedule: eventBridge.Schedule_Cron(&eventBridge.CronOptions{
			Hour:   jsii.String("0"),
			Minute: jsii.String("0"),
		}),
		Description: jsii.String("An EventBridge rule that fires the Lambda function at 00:00 every day"),
		Targets: &[]eventBridge.IRuleTarget{
			targets.NewLambdaFunction(lambdaFunc, nil),
		},
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewGHStars2NotionStack(app, "GHStars2Notion", &GHStars2NotionProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	account, region := os.Getenv("CDK_DEFAULT_ACCOUNT"), os.Getenv("CDK_DEFAULT_REGION")

	if account == "" || region == "" {
		return nil
	}

	return &awscdk.Environment{
		Account: jsii.String(account),
		Region:  jsii.String(region),
	}
}

func getContextAsString(node *constructs.Node, varName string) (string, error) {
	if val := (*node).TryGetContext(jsii.String(varName)); val != nil {
		return val.(string), nil
	}

	return "", fmt.Errorf("%s was not found in context", varName)
}
