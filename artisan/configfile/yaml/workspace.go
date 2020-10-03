package yaml

import (
	"fmt"
	"github.com/unravela/artisan/api"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// TODO
func LoadWorkspace(path string, ws *api.Workspace) error {
	// load the YAML
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read module file in %s", path)
	}

	yamlWs := make(map[interface{}]interface{})
	err = yaml.Unmarshal(data, &yamlWs)
	if err != nil {
		return err
	}

	// initialize workspace values
	ws.Environments = make(api.Environments)
	ws.MainModule = &api.Module{}

	// translate YAML data and fill in Workspace
	envs := yamlWs["environments"].([]interface{})
	for _, values := range envs {
		valuesMap := castToMap(values)
		env := &api.Environment{
			Name: castToStr(valuesMap["name"]),
			Src: castToStr(valuesMap["src"]),
			Image: castToStr(valuesMap["image"]),
		}

		ws.Environments[env.Name] = env
	}

	return nil
}
