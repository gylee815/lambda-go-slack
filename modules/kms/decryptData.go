package decryptData

import (
	"context"
	b64 "encoding/base64"
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

// lambda functio name is required which encryoted data with encryptioncontext
type KMSLambdaPayload struct {
	Data               string `json:"data"`
	LambdaFunctionName string `json:"lambda_func_name"`
}

// KMSDecryptAPI defines the interface for the Decrypt function.
// We use this interface to test the function using a mocked service.
type KMSDecryptAPI interface {
	Decrypt(ctx context.Context,
		params *kms.DecryptInput,
		optFns ...func(*kms.Options)) (*kms.DecryptOutput, error)
}

// DecodeData decrypts some text that was encrypted with an AWS Key Management Service (AWS KMS) key (KMS key).
// Inputs:
//
//  c is the context of the method call, which includes the AWS Region.
//  api is the interface that defines the method call.
//  input defines the input arguments to the service call.
//
// Output:
//
//  If success, a DecryptOutput object containing the result of the service call and nil.
//  Otherwise, nil and an error from the call to Decrypt.
func DecodeData(c context.Context, api KMSDecryptAPI, input *kms.DecryptInput) (*kms.DecryptOutput, error) {
	return api.Decrypt(c, input)
}
func DecryptData(payload KMSLambdaPayload) (string, bool) {
	if payload.Data == "" {
		fmt.Println("You must supply the encrypted data")
		return "Empty input data", false
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	client := kms.NewFromConfig(cfg)
	blob, err := b64.StdEncoding.DecodeString(payload.Data)
	if err != nil {
		panic("error converting string to blob, " + err.Error())
	}
	var encryptionContext map[string]string
	encryptionContext = make(map[string]string)
	encryptionContext["LambdaFunctionName"] = payload.LambdaFunctionName
	input := &kms.DecryptInput{
		CiphertextBlob:    blob,
		EncryptionContext: encryptionContext,
	}
	result, err := DecodeData(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got error decrypting data: ", err)
		return "Fail to decode data", false
	}
	return string(result.Plaintext), true
}
func DecryptDataWithInput() string {
	data := flag.String("d", "", "The encrypted data, as a string")
	lambdaFuncName := flag.String("c", "", "encryptionContext as a string")
	flag.Parse()
	fmt.Printf("data type: %T, *data: %v, data: %v\n", data, *data, data)
	if *data == "" {
		fmt.Println("You must supply the encrypted data as a string")
		fmt.Println("-d DATA")
		return ""
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	client := kms.NewFromConfig(cfg)
	blob, err := b64.StdEncoding.DecodeString(*data)
	if err != nil {
		panic("error converting string to blob, " + err.Error())
	}
	var encryptionContext map[string]string
	encryptionContext = make(map[string]string)
	encryptionContext["LambdaFunctionName"] = *lambdaFuncName
	input := &kms.DecryptInput{
		CiphertextBlob:    blob,
		EncryptionContext: encryptionContext,
	}
	result, err := DecodeData(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got error decrypting data: ", err)
		return ""
	}
	return string(result.Plaintext)
}
