package letitgo

// Config is the lowest level config.
type Config struct {
	Name        string                   `yaml:"name"`
	Description string                   `yaml:"description"`
	Actions     []map[string]interface{} `yaml:"actions"`
}
