package devfile

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/redhat-developer/odo-fork/pkg/config"
	"gopkg.in/yaml.v2"
)

// IDP constants
const (
	DevfileYaml = "devfile.yaml"
)

// Load read Devfile from filename
func Load() (*Devfile, error) {

	// Retrieve the IDP.yaml file
	udoDir, err := config.GetUDOFolder("")
	if err != nil {
		return nil, fmt.Errorf("unabled to find .udo folder in current directory")
	}
	filepath := path.Join(udoDir, DevfileYaml)

	f, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var data Devfile
	err = yaml.Unmarshal(f, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
