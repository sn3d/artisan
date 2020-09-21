package api

// Workspace is root directory of repository and contains
// all modules and factions
type Workspace struct {
	// absolute path to root of the artisan. This is a dir where
	// artisan file is situated
	RootDir string

	// hold all available factions for artisan
	Factions Factions

	// mainModule is module with reference '//'
	MainModule *Module
}
