package cmd

type ApolloClient struct {
	AppId      string   `yaml:"appId" json:"appId"`
	Namespaces []string `yaml:"namespaces" json:"namespaces"`
	Cluster    string   `yaml:"cluster" json:"cluster"`
	Server     string   `yaml:"server" json:"server"`
	Watch      bool     `yaml:"watch" json:"watch"`
}
