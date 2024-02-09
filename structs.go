package web

import (
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type WebCheck struct {
	HomePage []HomePage `yaml:"homepage"`
}

type HomePage struct {
	Host string   `yaml:"host"`
	Path []string `yaml:"path"`
}

func GetWebCheck(fileName string) (*WebCheck, error) {

	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	buf, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	w := &WebCheck{}
	err = yaml.Unmarshal(buf, w)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return w, nil
}
