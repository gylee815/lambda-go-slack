// BUILD : GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o HandleLambdaEvent main.go
// ZIP : zip HandleLambdaEvent.zip HandleLambdaEvent iplist.json
package main

import (
	"context"
	"fmt"
	"os"

	// "os"
	"encoding/base64"
	"encoding/json"
	iplist "lambda/modules/iplist"
	slack "lambda/modules/slack"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Body struct {
	Hostname string `json:"host_name"`
	Nodename string `json:"node_name"`
}

// type MyEvent struct {
// 	Text		string `json:"text"`
// 	Username	string `json:"user_name"`
// 	HostName	string `json:"host_name"`
// 	NodeName	string `json:"node_name"`
// }

// type MyResponse struct {
// 	Message			string `json:"msg"`
// 	ResponseStatus	int `json:"resp_status"`
// }

// func HandleLambdaEvent(event map[string]interface{}) (MyResponse, error) {
// 	// Marshaling (Go object -> string)
// 	jsonDoc, _ := json.Marshal(event)
// 	fmt.Println(string(jsonDoc))

// 	// Unmarhaling (string -> Go object(struct))
// 	my_event := MyEvent{}
// 	json.Unmarshal(jsonDoc, &my_event)
// 	fmt.Printf("text: %s, user_name: %s\n", my_event.Text, my_event.Username)

// 	return MyResponse{ResponseStatus: 200}, nil
// }

func HandleLambdaEvent(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	fmt.Printf("Processing request data for request %s.\n", request.RequestContext.RequestID)
	fmt.Printf("Body size = %d.\n", len(request.Body))

	jsonHeader, _ := json.Marshal(request.Headers)
	fmt.Printf("Headers: %s\n", string(jsonHeader))

	// fmt.Println("Headers:")
	// for key, value := range request.Headers {
	// 	fmt.Printf("    %s: %s\n", key, value)
	// }

	rawDecodedBody, err := base64.StdEncoding.DecodeString(request.Body)
	if err != nil {
		panic(err)
	}

	data := rawDecodedBody
	bodyData := Body{}

	json.Unmarshal([]byte(data), &bodyData)
	fmt.Printf("Body: %s\n", rawDecodedBody)

	site_name := iplist.GetSite(request.RequestContext.HTTP.SourceIP)
	fmt.Printf("SourceIP: %s, SiteName: %s\n", request.RequestContext.HTTP.SourceIP, site_name)

	payload := slack.SlackPayload{Text: fmt.Sprintf("POD [ %s ] on NODE [ %s ] is terminated!\n SourceIP: %s\n", bodyData.Hostname, bodyData.Nodename, request.RequestContext.HTTP.SourceIP), Username: string(site_name)}
	url := os.Getenv("SLACK_ONPREM_K8S")

	// /* This will be active after slack url with encrypted with kms -> Terraform not ready yet */
	// url := os.Getenv("kmsEncryptedHookUrl")
	// lambdaFuncName := os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	// kmsPayload := kmsDecrypt.KMSLambdaPayload{
	// 	Data:               url,
	// 	LambdaFunctionName: lambdaFuncName,
	// }

	// if rst, success := kmsDecrypt.DecryptData(kmsPayload); success {
	// 	url = rst
	// 	fmt.Printf("Slack url: %s\n", rst)
	// } else {
	// 	fmt.Printf("Error msg: %s\n", rst)
	// }

	if err := slack.PostMessage(url, payload); err != nil {
		fmt.Printf("Error msg: %s\n", err)
	} else {
		fmt.Println("Send message to slcak webhook succeed\n")
	}

	return events.LambdaFunctionURLResponse{Body: request.Body, StatusCode: 200}, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
