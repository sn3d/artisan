package delvin

type ClassBuilder interface {
	Build(c *Class, ws *Workspace) (*ClassImage, error)
}
