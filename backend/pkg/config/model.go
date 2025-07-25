package config

type Config struct {
	Server   ConfigServer   `yaml:"server"`
	Logs     ConfigLogs     `yaml:"logs"`
	Database ConfigDatabase `yaml:"database"`
	App      ConfigApp      `yaml:"app"`
}

type ConfigServer struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

type ConfigLogs struct {
	Level  string `yaml:"level"`
	File   string `yaml:"file"`
	Pretty bool   `yaml:"pretty"`
}

type ConfigDatabase struct {
	Dialect          string `yaml:"dialect"`
	ConnectionString string `yaml:"connection_string"`
}

type ConfigApp struct {
	Kubernetes  ConfigAppKubernetes  `yaml:"kubernetes"`
	Helm        ConfigAppHelm        `yaml:"helm"`
	Initializer ConfigAppInitializer `yaml:"initializer"`
}

type ConfigAppKubernetes struct {
	AuthURL   string `yaml:"auth_url"`
	AuthToken string `yaml:"auth_token"`
	AuthCert  string `yaml:"auth_cert"`
	AuthCA    string `yaml:"auth_ca"`
	Namespace string `yaml:"namespace"`
}

type ConfigAppHelm struct {
	Repository        string `yaml:"repository"`
	BasicAuthUsername string `yaml:"basic_auth_username"`
	BasicAuthPassword string `yaml:"basic_auth_password"`
}

type ConfigAppInitializer struct {
	AdminPassword string `yaml:"admin_password"`
}
