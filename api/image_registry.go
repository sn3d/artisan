package api

type ImageRegistry interface {
	Build(f *Faction, srcDir string) (*Image, error)
}
