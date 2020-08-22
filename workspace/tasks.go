package workspace

import "github.com/unravela/delvin/api"

type Tasks []*api.Task

func (t Tasks) GetClasses(ws *Workspace) Classes {
	classMap := make(map[string]*api.Class)
	for _, task := range t {
		classDef := ws.Class(task.Class);
		classMap[task.Class] = classDef
	}
	return classMap
}
