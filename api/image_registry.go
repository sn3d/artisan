package api

// ImageRegistry hold and provide you images for Environment.
type ImageRegistry interface {
	Build(env *Environment, srcDir string) (*Image, error)
}
