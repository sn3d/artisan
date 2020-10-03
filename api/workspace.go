package api

// Workspace is root directory of repository and contains
// all modules and environments
type Workspace struct {
	// absolute path to root of the artisan. This is a dir where
	// artisan file is situated
	RootDir string

	// hold all available environments for artisan
	Environments Environments

	// mainModule is module with reference '//'
	MainModule *Module
}
