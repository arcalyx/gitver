package constants

const (
	VersionFileName  = ".version"
	ConfigName       = "config"
	ConfigType       = "yaml"
	ConfigFolderName = ".gitver"
	ProgramName      = "gitver"
	TagMessage       = "Tagged by gitver"
	CommitMessage    = "Bump Version [%s] -> [%s]"
	ReleaseTag       = "r%s"
	VersionTag       = "v%s"
)

// Package manager related constants
const (
	MavenPomFile      = "pom.xml"
	NPMPackageFile    = "package.json"
	GradleBuildFile   = "build.gradle"
	KotlinBuildFile   = "build.gradle.kts"
	GoModFile         = "go.mod"
	PythonSetupFile   = "setup.py"
	PyProjectTomlFile = "pyproject.toml"

	PackageUpdateSuccess = "Updated version in %s to %s"
	PackageUpdateFailed  = "Failed to update version in %s: %v"
)
