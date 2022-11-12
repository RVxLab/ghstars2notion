package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
)

func HandleRequest() (string, error) {
	return run()
}

func main() {
	if len(os.Args) > 1 && os.Args[1] != "run_local" {
		lambda.Start(HandleRequest)
	}

	out, err := run()

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(out)
	}

	return
}
