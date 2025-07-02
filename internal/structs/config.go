package structs

type Config struct {
	PruneLimit   int  `yaml:"prune_limit"`
	AutoCompress bool `yaml:"auto_compress"`
}
