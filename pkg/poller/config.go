package poller

// Config represents the poller configuration struct
type Config struct {
	Nodes []NodeConfig `yaml:"nodes"`
}

// NodeConfig represents a node entry in the poller
type NodeConfig struct {
	Name         string       `yaml:"name"`
	Label        string       `yaml:"label"`
	Address      string       `yaml:"remote_addr"`
	Links        []LinkConfig `yaml:"links"`
	PollInterval int          `yaml:"poll_interval"`
}

// LinkConfig represents a link entry in a poller's node
type LinkConfig struct {
	Label      string `yaml:"label"`
	HostURL    string `yaml:"host_url"`
	ServiceURL string `yaml:"service_url"`
}
