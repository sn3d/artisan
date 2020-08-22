package workspace

import "github.com/unravela/delvin/api"

type Classes map[string]*api.Class

func arrayToClasses(arr []*api.Class) Classes {
	classMap := make(map[string]*api.Class)
	for _, c := range arr {
		classMap[c.Name] = c
	}
	return classMap
}
