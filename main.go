package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
)

var mapFile *string

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
	gen := buildCmd(
		"go", "run", "hack/genproto/main.go",
		"-d", "hack/genproto/descriptor.json",
		"-m", *mapFile,
		"-o", "api/api.pb.gw.go",
	)
	if err := gen.Run(); err != nil {
		return err
	}
	return nil
}

func server() error {
	server := buildCmd("go", "run", "cmd/main.go")
	if err := server.Run(); err != nil {
		return err
	}
	return nil
}

func main() {
	// Get options from flags
	mapFile = flag.String("m",  "api/map.json", "Directory of proto map file")
	build := flag.Bool("b", false, "Runs builder methods only")
	flag.Parse()
	pipeline := []func() error{
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
