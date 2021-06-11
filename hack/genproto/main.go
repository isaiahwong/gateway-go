package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"text/template"

	"github.com/isaiahwong/gateway-go/hack/genproto/internal"
)

var descriptor, out, mapFile *string

func init() {
	// Get options from flags
	descriptor = flag.String("d", "", "location of descriptor file")
	mapFile = flag.String("m", "", "location of proto map file")
	out = flag.String("o", "", "Output of gwprotos. I.E protos.pb.gw.go")
	flag.Parse()

	if *descriptor == "" {
		panic("no directory specified for descriptor file.")
	}
	if *mapFile == "" {
		panic("no directory specified for proto map file.")
	}
	if *out == "" {
		panic("no directory specified for output protos")
	}
}

func buildCmd(name string, cmd ...string) *exec.Cmd {
	c := exec.Command(name, cmd...)
	c.Env = os.Environ()
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c
}

func npmInstall() error {
	cmd := buildCmd("npm", "i", "--prefix", "hack/genproto/js")
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func genDescriptor(mapFile []byte) error {
	cmd := buildCmd(
		"node",
		"hack/genproto/js/index.js",
		"--map",
		string(mapFile),
	)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func readMapFile() ([]byte, error) {
	f, err := os.Open(*mapFile)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func protoc(proto string, includes []string) error {
	cmdstr := "protoc "
	for _, in := range includes {
		cmdstr += fmt.Sprintf("-I%v ", in)
	}
	out := "./api/go/gen"
	cmdstr += fmt.Sprintf("--go_out ./%[1]v --go-grpc_out %[1]v "+
		"--grpc-gateway_out %[1]v "+
		"--grpc-gateway_opt logtostderr=true "+
		"--grpc-gateway_opt allow_repeated_fields_in_body=true "+
		"%[2]v/*.proto ", out, proto)

	cmd := buildCmd("/bin/sh", "-c", cmdstr)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func generateProtos(mapFile []byte) error {
	var protoMaps []internal.ProtoMap
	if err := json.Unmarshal(mapFile, &protoMaps); err != nil {
		log.Fatal(err)
	}
	// Create directory
	path := "api/go/gen"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0700)
	}
	for _, m := range protoMaps {
		for _, p := range m.Protos {
			err := protoc(p, m.Includes)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// readDescriptor reads generated descriptor json file and unmarshals it
func readDescriptor(svcs *[]internal.ServiceDesc, file string) error {
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
	mapFile, err := readMapFile()
	if err != nil {
		panic(err)
	}
	pipeline := []func() error{
		func() error {
			return generateProtos(mapFile)
		},
		npmInstall,
		func() error {
			return genDescriptor(mapFile)
		},
	}
	// Run Commands
	for _, fn := range pipeline {
		err := fn()
		if err != nil {
			log.Printf("Failed to start cmd: %v", err)
			break
		}
	}

	// Parse JSON file
	var svcs []internal.ServiceDesc
	err = readDescriptor(&svcs, *descriptor)
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
