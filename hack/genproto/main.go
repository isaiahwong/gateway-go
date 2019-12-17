package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/isaiahwong/gateway-go/hack/genproto/internal"
)

type data struct {
	Name string
}

// TODO
func flags() {}

// read reads generated descriptor json file and unmarshals it
func read(svcs *[]internal.ServiceDesc, file string) error {
	if file == "" {
		file = "descriptor.json"
	}
	jsonFile, err := os.Open(file)
	if err != nil {
		return err
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(byteValue, svcs); err != nil {
		log.Fatal(err)
	}

	return nil
}

// TODO FILTER OFF NON HTTP ROUTES
func main() {
	var svcs []internal.ServiceDesc
	err := read(&svcs, "descriptor.json")
	if err != nil {
		panic(err)
	}

	// Get working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(dir + "/../../internal/server/proto.go")
	if err != nil {
		log.Fatalf("failed with %s\n", err)
	}
	t := template.Must(template.New("server").Parse(internal.ProtoTemplate))
	t.Execute(f, svcs)
}
