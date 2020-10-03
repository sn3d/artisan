package yaml

import (
	"fmt"
	"github.com/unravela/artisan/api"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)
// LoadModule open the file on giiven path and load data
// into given api.Module structure.
//
// If file doesn't exist, or cannot be read, or there is syntax error,
// the function returns error
func LoadModule(path string, mod *api.Module) error {

	// load the YAML data
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read module file in %s", path)
	}

	yamlMod := make(map[interface{}]interface{})
	err = yaml.Unmarshal(data, &yamlMod)
	if err != nil {
		return err
	}

	// transform YAML data and fill into module
	yamlToModule(yamlMod, mod)
	return nil
}

func yamlToModule(yml map[interface{}]interface{}, mod *api.Module) {
	tasks := yml["tasks"].([]interface{})
	for _, values := range tasks {
		task := new(api.Task)
		yamlToTask(castToMap(values), task)
		mod.Tasks = append(mod.Tasks, task)
	}
}

func yamlToTask(yml map[interface{}]interface{}, t *api.Task) {
	t.Name = castToStr(yml["name"])
	t.EnvName = castToStr(yml["env"])
	t.Script = castToStr(yml["script"])
	t.Deps = castToStringArray(yml["deps"])
	t.Exclude = castToStringArray(yml["exclude"])
	t.Include=  castToStringArray(yml["include"])
}