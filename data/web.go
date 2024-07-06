package data

type WebInfo struct {
	Version string `json:"version"`
}

type WebIndex struct {
	Version string
}

type WebReload struct {
	Version string
	Source  string
	Success bool
}

type WebShowConfig struct {
	Version  string
	Messages []string
	Content  string
	Diff     bool
	Success  bool
}
