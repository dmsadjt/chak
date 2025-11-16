package config

type Profile struct {
	SzName string `json:"name"`
	SzDescription string `json:"description"`
	SzDirectories []string `json:"directories"`
	SzMemoryFile string `json:"memory_file"`
	SzIndexFile string `json:"index_file"`
	Extensions []string `json:"extensions"`
	InMaxSizeFile int64 `json:"max_file_size"`
}

type ConfigInterface interface {
	GetActiveProfile() Profile 
	SwitchProfile(szName string) error
	ListProfile() []Profile
	GetProfile(szName string) (Profile, error)
}

type Config struct {
	SzActiveProfile string `json:"active_profile"`
	Profiles map[string]Profile `json:"profiles"`
}

