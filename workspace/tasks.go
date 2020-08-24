package workspace

import "github.com/unravela/artisan/api"

type Tasks []*api.Task

func (t Tasks) GetFactions(ws *Workspace) api.Factions {
	factions := make(api.Factions)
	for _, task := range t {
		factionDef := ws.Faction(task.FactionName)
		factions[task.FactionName] = factionDef
	}
	return factions
}
