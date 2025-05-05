package installer

// GoVersions is a list of Go versions.
type GoVersions []GoVersion

// GoVersion represents a specific version of Go and its associated files.
type GoVersion struct {
	Version string    `json:"version"`
	Stable  bool      `json:"stable"`
	Files   []FileRef `json:"files"`
}

// FileRef represents a Go distribution file reference with its metadata.
type FileRef struct {
	Filename string `json:"filename"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Version  string `json:"version"`
	SHA256   string `json:"sha256"`
	Size     int64  `json:"size"`
	Kind     string `json:"kind"` // Can be "source", "archive", or "installer"
}
