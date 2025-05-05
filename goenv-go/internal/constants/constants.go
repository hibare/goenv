package constants

const (
	ProjectName = "goenv"
)
const (
	DefaultVersion = "system"
)

const (
	LocalGoVersionFile             = ".go-version"
	GlobalGoVersionFile            = "version" // Default `${HOME}/.goenv/version`
	GlobalVersionFileLegacyDefault = "default" // Default `${HOME}/.goenv/default`
	GlobalVersionFileLegacyGlobal  = "global"  // Default `${HOME}/.goenv/global`
	GoModFile                      = "go.mod"
)

const (
	GoenvRootDir   = ".goenv"   // Default `${HOME}/.goenv`
	ShimsDir       = "shims"    // Default `${HOME}/.goenv/shims`
	VersionsDir    = "versions" // Default `${HOME}/.goenv/versions`
	VersionsBinDir = "bin"      // Default `${HOME}/.goenv/versions/bin`
)

const (
	GoSystemVersion = "system"
)

const (
	GoDevDl       = "https://go.dev/dl/?mode=json&include=all"
	ArchiveFormat = "tar.gz"
)

// Env variables
const (
	EnvGoenvRootDir = "GOENV_ROOT"
)
