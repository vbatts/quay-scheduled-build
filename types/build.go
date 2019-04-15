package types

// Build is the configuration for a scheduled build on quay.io
type Build struct {
	QuayRepo string `json:"quay_repo"` // org/repo on quay
	Schedule string `json:"schedule"`  // cron style schedule
	Token    string `json:"schedule"`  // quay oauth secret
	BuildRef        // inherit the quay structure needed pass as the body
}

// BuildRef is the particulars to fetch the Dockerfile
//
// This is all fairly predictable, and almost guessable for most github
// projects, but needs to be spelled out.
type BuildRef struct {
	PullRobot      string   `json:"pull_robot"`      // robot account's username
	Tags           []string `json:"docker_tags"`     // image name tags to apply to this build
	ArchiveUrl     string   `json:"archive_url"`     // this can be a tar[.gz] or zip archive. Must start with 'https' or 'http'
	DockerfilePath string   `json:"dockerfile_path"` // whole path within the archive to the Dockerfile
	Subdirectory   string   `json:"subdirectory"`    // path within the archive, to the build context
	Context        string   `json:"context"`         // path within the archive, to the build context
}
