package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
)

func buildCmd(name string, cmd ...string) *exec.Cmd {
	c := exec.Command(name, cmd...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c
}

func genDescriptor() error {
	cmd := buildCmd("node", "hack/genproto/js/index.js")
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// Generates service handlers dynamically
func genCode() error {
	// Get working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	gen := buildCmd(
		"go", "run", "hack/genproto/main.go",
		"-d", dir+"/hack/genproto/descriptor.json",
		"-o", dir+"/internal/server/proto.go",
	)
	if err := gen.Run(); err != nil {
		return err
	}
	return nil
}

func server() error {
	server := buildCmd("go", "run", "cmd/gateway/main.go")
	if err := server.Run(); err != nil {
		return err
	}
	return nil
}

func main() {
	// Get options from flags
	build := flag.Bool("b", false, "Runs builder methods only")
	flag.Parse()
	pipeline := []func() error{
		genDescriptor,
		genCode,
	}
	// Runs server if no build is specified
	if !*build {
		pipeline = append(pipeline, server)
	}
	flag.Parse()
	for _, fn := range pipeline {
		err := fn()
		if err != nil {
			log.Printf("Failed to start cmd: %v", err)
			break
		}
	}
}
