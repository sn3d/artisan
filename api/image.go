package api

// Image entity is output of every environment. For task execution we need
// to have image (docker image), not environment. The ImageRegistry provide you
// images fo environments
type Image struct {
	// ID is unique image ID that is associated to Environment
	ID string
}

// Images is used as image collection
type Images map[string]*Image
