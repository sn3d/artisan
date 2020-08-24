package api

// ImageRegistry hold and provide you images for factions.
type ImageRegistry interface {
	Build(f *Faction, srcDir string) (*Image, error)
}
