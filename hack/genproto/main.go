package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"text/template"

	"github.com/isaiahwong/gateway-go/hack/genproto/internal"
)

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

func format() error {
	format := exec.Command("go", "fmt", "github.com/isaiahwong/gateway-go/internal/server")
	if err := format.Run(); err != nil {
		return err
	}
	return nil
}

func main() {
	// Get working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	// Get options from flags
	descriptor := flag.String("d", dir+"/descriptor.json", "descriptor ")
	out := flag.String("o", dir+"/../../protogen/protos.pb.gw.go", "output of gwprotos")
	flag.Parse()
	// Parse JSON file
	var svcs []internal.ServiceDesc
	err = read(&svcs, *descriptor)
	if err != nil {
		panic(err)
	}
	f, err := os.Create(*out)
	if err != nil {
		panic(err)
	}
	t := template.Must(template.New("server").Parse(internal.ProtoTemplate))
	err = t.Execute(f, svcs)
	if err != nil {
		panic(err)
	}
	err = format()
	if err != nil {
		panic(err)
	}
}
