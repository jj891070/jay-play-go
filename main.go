package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	/**
	./app -n ubuntu -s 1.7.8 -d jay
	*/
	var imageName, sourceImageTag, addImageTag string
	flag.StringVar(&imageName, "n", "", "input image name. ex: ubuntu")
	flag.StringVar(&sourceImageTag, "s", "", "input source image tag name. ex: 1.7.8")
	flag.StringVar(&addImageTag, "d", "", "input add image tag name. ex: jay")
	flag.Parse()
	if imageName == "" {
		log.Fatalln("Please Input '-n [image_name]'")
	}
	if sourceImageTag == "" {
		log.Fatalln("Please Input '-s [source_image_tag]'")
	}
	if addImageTag == "" {
		log.Fatalln("Please Input '-d [destination_image_tag]'")
	}

	// _json_key,AWS
	username := os.Getenv("IMAGE_USERNAME")
	if username == "" {
		log.Fatalln("Please Input Env 'IMAGE_USERNAME'")
	}

	// gcp iam json file
	var password string
	switch username {
	case "AWS":
		if os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
			log.Fatalln("Please Input Env 'AWS_SECRET_ACCESS_KEY'")
		}
		if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
			log.Fatalln("Please Input Env 'AWS_ACCESS_KEY_ID'")
		}
		if os.Getenv("AWS_ACCOUNT") == "" {
			log.Fatalln("Please Input Env 'AWS_ACCOUNT'")
		}
		region := os.Getenv("AWS_REGION")
		if region == "" {
			region = "ap-east-1"
		}

		password, err = AwsGetToken(region)
		if err != nil {
			log.Fatalln(fmt.Errorf("aws login failed. err: %v", err))
		}
	default:
		password = os.Getenv("IMAGE_PASSWORD")
	}

	if password == "" {
		log.Fatalln("Please Input Env 'IMAGE_PASSWORD'")
	}

	authConfig := types.AuthConfig{
		Username: username,
		Password: password,
	}

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		log.Fatalln(fmt.Errorf("error when encoding authConfig. err: %v", err))
	}

	pullTarget := imageName + ":" + sourceImageTag
	pushNewTarget := imageName + ":" + addImageTag
	reader, err := cli.ImagePull(
		ctx,
		pullTarget,
		types.ImagePullOptions{
			RegistryAuth: base64.URLEncoding.EncodeToString(encodedJSON),
		})
	if err != nil {
		log.Fatalln(fmt.Errorf("error when pull image. err: %v", err))
	}
	io.Copy(os.Stdout, reader)

	err = cli.ImageTag(
		ctx,
		pullTarget,
		pushNewTarget,
	)
	if err != nil {
		log.Fatalln(fmt.Errorf("error when tag image. err: %v", err))
	}

	reader, err = cli.ImagePush(
		ctx,
		pushNewTarget,
		types.ImagePushOptions{
			RegistryAuth: base64.URLEncoding.EncodeToString(encodedJSON),
		},
	)
	if err != nil {
		log.Fatalln(fmt.Errorf("error when push image. err: %v", err))
	}
	io.Copy(os.Stdout, reader)
	fmt.Println("---Finish---")
}

func AwsGetToken(region string) (token string, err error) {
	svc := ecr.New(
		session.Must(session.NewSession()),
		aws.NewConfig().WithRegion(region),
	)

	input := &ecr.GetAuthorizationTokenInput{}
	var result *ecr.GetAuthorizationTokenOutput
	result, err = svc.GetAuthorizationToken(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecr.ErrCodeServerException:
				err = fmt.Errorf("%v%v", ecr.ErrCodeServerException, aerr.Error())
				// fmt.Println(ecr.ErrCodeServerException, aerr.Error())
			case ecr.ErrCodeInvalidParameterException:
				err = fmt.Errorf("%v%v", ecr.ErrCodeInvalidParameterException, aerr.Error())
				// fmt.Println(ecr.ErrCodeInvalidParameterException, aerr.Error())
			default:
				err = fmt.Errorf("%v", aerr.Error())
			}
		}
		return
	}

	if len(result.AuthorizationData) != 1 {
		err = fmt.Errorf("token failed => %+v", result.AuthorizationData)
		return
	}

	var tmpUsernameANDToken []byte
	tmpUsernameANDToken, err = base64.URLEncoding.DecodeString(*result.AuthorizationData[0].AuthorizationToken)
	auth := strings.Split(string(tmpUsernameANDToken), ":")
	if len(auth) != 2 {
		err = fmt.Errorf("token failed => %+v", result.AuthorizationData)
		return
	}

	token = auth[1]
	return
}
