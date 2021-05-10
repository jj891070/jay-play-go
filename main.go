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
	if imageName == "" || sourceImageTag == "" || addImageTag == "" {
		log.Fatalln("Please Input 'image name' and 'source image tag' and 'destination image tag'")
	}

	// _json_key
	username := os.Getenv("IMAGE_USERNAME")
	if username == "" {
		log.Fatalln("Please Input Env 'IMAGE_USERNAME'")
	}
	// gcp iam json file
	password := os.Getenv("IMAGE_PASSWORD")
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
	// resp, err := cli.ContainerCreate(ctx, &container.Config{
	// 	Image: "asia.gcr.io/ag-ocean-registry/jay/golang",
	// 	Cmd:   []string{"echo", "hello world"},
	// }, nil, nil, nil, "")
	// if err != nil {
	// 	panic(err)
	// }

	// if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
	// 	panic(err)
	// }

	// statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	// select {
	// case err := <-errCh:
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// case <-statusCh:
	// }

	// out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	// if err != nil {
	// 	panic(err)
	// }

	// stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}
