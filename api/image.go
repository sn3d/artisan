package api

// Image entity is output of every faction. For task execution we need
// to have image (docker image), not faction. The ImageRegistry provide you
// images fo factions
type Image struct {
	// ID is unique image ID that is associated to Faction
	ID string
}
