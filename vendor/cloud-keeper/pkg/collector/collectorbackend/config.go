package collectorbackend

const (
	//OperatorsVultr vultr
	OperatorsVultr = "Vultr"
	//OperatorsDigOC digitalocean
	OperatorsDigOC = "DigitalOcean"
)

// Config is configuration for creating a storage backend.
type Config struct {
	// Type defines the type of operators,like as vultr, oc and so on.
	Type string

	//accesse key for request to operators
	APIKey string
}
