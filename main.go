package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	Replace()
	execRoot, err := os.Getwd()
	if err != nil {
		log.Println("err -> ", err)
	}
	log.Println(execRoot)
}

// Replace 覆蓋掉
func Replace() {
	// ls $K8sBaseKustomizationPath | while read -r filename; do sed -i "s/VERSION_ID/${IMAGEID}/g" $K8sBaseKustomizationPath/$filename; done
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(
		"/bin/bash", "job.sh",
		"AgOcean.InfraInfo/resources/infra-app/ag-ocean/base-app/cloud/accountsystem",
		"1.1.1.1",
		"accountsystem",
		"agocean",
		"tmp",
		".zip",
	)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("err -> %+v , %+v \n", err, string(stderr.Bytes()))
	}
	log.Println(string(stdout.Bytes()))
}

// GitPull 更新infra repo
func GitPull() {
	// Update CN site repo.
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("git", "-C", "AgOcean.InfraInfo", "pull")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("err -> %+v , %+v \n", err, string(stderr.Bytes()))
	}
	log.Println(string(stdout.Bytes()))
}

// Demo 範例
func Demo() {
	cmd := exec.Command("ls", "-lah")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
	return
}
