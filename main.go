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
	var imageName, sourceImageTag, addImageTag, awsImageName string
	flag.StringVar(&imageName, "n", "", "input image name. ex: ubuntu")
	flag.StringVar(&awsImageName, "b", "", "input source image name(backup, now is aws). ex: 1.7.8")
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

	password := os.Getenv("IMAGE_PASSWORD")
	if password == "" {
		log.Fatalln("Please Input Env 'IMAGE_PASSWORD'")
	}

	awsUsername := os.Getenv("AWS_IMAGE_USERNAME")
	awsPassword := ""
	if awsUsername == "AWS" {
		if os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
			log.Fatalln("Please Input Env 'AWS_SECRET_ACCESS_KEY'")
		}
		if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
			log.Fatalln("Please Input Env 'AWS_ACCESS_KEY_ID'")
		}
		if os.Getenv("AWS_ACCOUNT") == "" {
			log.Fatalln("Please Input Env 'AWS_ACCOUNT'")
		}
		if awsImageName == "" {
			log.Fatalln("Please Input '-n [aws_image_name]'")
		}
		region := os.Getenv("AWS_REGION")
		if region == "" {
			region = "ap-east-1"
		}

		awsPassword, err = AwsGetToken(region)
		if err != nil {
			log.Fatalln(fmt.Errorf("aws login failed. err: %v", err))
		}
		fmt.Println("---aws backup start---")

	}

	var gcpAuth string
	gcpAuth, err = AuthEncode(username, password)
	if err != nil {
		log.Fatalln(fmt.Errorf("error when encoding authConfig. err: %v", err))
	}

	pullTarget := imageName + ":" + sourceImageTag
	pushNewTarget := imageName + ":" + addImageTag
	reader, err := cli.ImagePull(
		ctx,
		pullTarget,
		types.ImagePullOptions{
			RegistryAuth: gcpAuth,
		})
	if err != nil {
		log.Fatalln(fmt.Errorf("error when pull image. err: %v", err))
	}
	SeeUnauthorizedError(reader)

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
			RegistryAuth: gcpAuth,
		},
	)
	if err != nil {
		log.Fatalln(fmt.Errorf("error when push image. err: %v", err))
	}
	SeeUnauthorizedError(reader)

	if awsUsername == "AWS" {
		var awsAuth string
		awsAuth, err = AuthEncode(awsUsername, awsPassword)
		if err != nil {
			log.Fatalln(fmt.Errorf("error when encoding authConfig. err: %v", err))
		}

		pushAwsOldTarget := awsImageName + ":" + sourceImageTag
		pushAwsNewTarget := awsImageName + ":" + addImageTag
		err = cli.ImageTag(
			ctx,
			pullTarget,
			pushAwsOldTarget,
		)
		if err != nil {
			log.Fatalln(fmt.Errorf("error when tag image. err: %v", err))
		}
		err = cli.ImageTag(
			ctx,
			pullTarget,
			pushAwsNewTarget,
		)
		if err != nil {
			log.Fatalln(fmt.Errorf("error when tag image. err: %v", err))
		}

		reader, err = cli.ImagePush(
			ctx,
			pushAwsNewTarget,
			types.ImagePushOptions{
				RegistryAuth: awsAuth,
			},
		)
		if err != nil {
			log.Fatalln(fmt.Errorf("error when push image. err: %v", err))
		}

		SeeUnauthorizedError(reader)

		reader, err = cli.ImagePush(
			ctx,
			pushAwsOldTarget,
			types.ImagePushOptions{
				RegistryAuth: awsAuth,
			},
		)
		if err != nil {
			log.Fatalln(fmt.Errorf("error when push image. err: %v", err))
		}

		SeeUnauthorizedError(reader)
	}
	fmt.Println("---Finish---")
}

func AuthEncode(username, password string) (reulst string, err error) {
	authConfig := types.AuthConfig{
		Username: username,
		Password: password,
	}
	authConfigJSON, err := json.Marshal(authConfig)
	if err != nil {
		return
	}
	reulst = base64.URLEncoding.EncodeToString(authConfigJSON)
	return
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

func SeeUnauthorizedError(reader io.ReadCloser) {
	buf := new(strings.Builder)
	io.Copy(buf, reader)
	fmt.Println(buf.String())
	if strings.Contains(buf.String(), "unauthorized") {
		log.Fatalln("unauthorized")
	}
}
