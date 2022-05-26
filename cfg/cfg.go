package cfg

type Config struct {
	Globals struct {
		Functions map[string]string `yaml:"functions" json:"functions"`
	} `yaml:"globals" json:"globals"`
}
