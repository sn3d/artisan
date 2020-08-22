package api

type ImageRegistry interface {
	Build(c *Class, srcDir string) (*Image, error)
}
