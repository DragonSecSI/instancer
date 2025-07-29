package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"go.yaml.in/yaml/v3"

	"github.com/DragonSecSI/instancer/backend/pkg/errors"
)

func LoadConfig() (*Config, error) {
	config := &Config{}

	config.LoadDefaults()

	configFile := flag.String("config", "", "Path to the configuration file")
	flag.Parse()
	if *configFile != "" {
		if err := config.LoadFile(*configFile); err != nil {
			return nil, err
		}
	}

	defaultConfigPaths := []string{
		"./config.yml",
		"./configs/config.yml",
		"/etc/instancer/config.yml",
	}
	for _, path := range defaultConfigPaths {
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				continue // File does not exist, skip
			} else {
				return nil, &errors.ConfigFileError{
					FilePath: path,
					Err:      err,
				}
			}
		} else {
			if err := config.LoadFile(path); err != nil {
				return nil, err
			}
			break
		}
	}

	if err := config.LoadEnv(); err != nil {
		return nil, err
	}

	if err := config.LoadArgs(); err != nil {
		return nil, err
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) LoadDefaults() {
	// Server configuration
	c.Server.Address = "127.0.0.1"
	c.Server.Port = 8080

	// Logs configuration
	c.Logs.Level = "info"
	c.Logs.File = ""
	c.Logs.Pretty = true

	// Database configuration
	c.Database.Dialect = "postgres"
	c.Database.ConnectionString = ""

	// App configuration
	c.App.Kubernetes.AuthURL = ""
	c.App.Kubernetes.AuthToken = ""
	c.App.Kubernetes.AuthCert = ""
	c.App.Kubernetes.AuthCA = ""
	c.App.Kubernetes.Namespace = "default"

	c.App.Helm.Repository = ""
	c.App.Helm.BasicAuthUsername = ""
	c.App.Helm.BasicAuthPassword = ""

	c.App.Initializer.AdminPassword = ""

	c.App.Meta.WebSuffix = ""
	c.App.Meta.SocketSuffix = ""
	c.App.Meta.SocketPort = 443
}

func (c *Config) LoadFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return &errors.ConfigFileError{
			FilePath: path,
			Err:      err,
		}
	}

	if err := yaml.Unmarshal(data, c); err != nil {
		return &errors.ConfigFileError{
			FilePath: path,
			Err:      err,
		}
	}

	return nil
}

func (c *Config) LoadEnv() error {
	// Server configuration
	if address := os.Getenv("SERVER_ADDRESS"); address != "" {
		c.Server.Address = address
	}

	if port := os.Getenv("SERVER_PORT"); port != "" {
		portInt, err := strconv.Atoi(port)
		if err != nil {
			return &errors.ConfigValueError{
				Key: "SERVER_PORT",
				Err: err,
			}
		}
		c.Server.Port = portInt
	}

	// Logs configuration
	if level := os.Getenv("LOGS_LEVEL"); level != "" {
		c.Logs.Level = level
	}

	if file := os.Getenv("LOGS_FILE"); file != "" {
		c.Logs.File = file
	}

	if pretty := os.Getenv("LOGS_PRETTY"); pretty != "" {
		prettyBool, err := strconv.ParseBool(pretty)
		if err != nil {
			return &errors.ConfigValueError{
				Key: "LOGS_PRETTY",
				Err: err,
			}
		}
		c.Logs.Pretty = prettyBool
	}

	// Database configuration
	if dialect := os.Getenv("DATABASE_DIALECT"); dialect != "" {
		c.Database.Dialect = dialect
	}

	if connStr := os.Getenv("DATABASE_CONNECTION_STRING"); connStr != "" {
		c.Database.ConnectionString = connStr
	}

	// App configuration
	if authURL := os.Getenv("APP_KUBERNETES_AUTH_URL"); authURL != "" {
		c.App.Kubernetes.AuthURL = authURL
	}

	if authToken := os.Getenv("APP_KUBERNETES_AUTH_TOKEN"); authToken != "" {
		c.App.Kubernetes.AuthToken = authToken
	}

	if authCert := os.Getenv("APP_KUBERNETES_AUTH_CERT"); authCert != "" {
		c.App.Kubernetes.AuthCert = authCert
	}

	if authCA := os.Getenv("APP_KUBERNETES_AUTH_CA"); authCA != "" {
		c.App.Kubernetes.AuthCA = authCA
	}

	if namespace := os.Getenv("APP_KUBERNETES_NAMESPACE"); namespace != "" {
		c.App.Kubernetes.Namespace = namespace
	}

	if repo := os.Getenv("APP_HELM_REPOSITORY"); repo != "" {
		c.App.Helm.Repository = repo
	}

	if basicAuthUsername := os.Getenv("APP_HELM_BASIC_AUTH_USERNAME"); basicAuthUsername != "" {
		c.App.Helm.BasicAuthUsername = basicAuthUsername
	}

	if basicAuthPassword := os.Getenv("APP_HELM_BASIC_AUTH_PASSWORD"); basicAuthPassword != "" {
		c.App.Helm.BasicAuthPassword = basicAuthPassword
	}

	if adminPassword := os.Getenv("APP_INITIALIZER_ADMIN_PASSWORD"); adminPassword != "" {
		c.App.Initializer.AdminPassword = adminPassword
	}

	if webSuffix := os.Getenv("APP_META_WEB_SUFFIX"); webSuffix != "" {
		c.App.Meta.WebSuffix = webSuffix
	}

	if socketSuffix := os.Getenv("APP_META_SOCKET_SUFFIX"); socketSuffix != "" {
		c.App.Meta.SocketSuffix = socketSuffix
	}

	if socketPort := os.Getenv("APP_META_SOCKET_PORT"); socketPort != "" {
		socketPortInt, err := strconv.Atoi(socketPort)
		if err != nil {
			return &errors.ConfigValueError{
				Key: "APP_META_SOCKET_PORT",
				Err: err,
			}
		}
		c.App.Meta.SocketPort = socketPortInt
	}

	return nil
}

func (c *Config) LoadArgs() error {
	// Server configuration
	flag.StringVar(&c.Server.Address, "server.address", c.Server.Address, "Server bind address")
	flag.IntVar(&c.Server.Port, "server.port", c.Server.Port, "Server port")

	// Logs configuration
	flag.StringVar(&c.Logs.Level, "logs.level", c.Logs.Level, "Logging level (debug, info, warn, error)")
	flag.StringVar(&c.Logs.File, "logs.file", c.Logs.File, "Path to the log file")
	flag.BoolVar(&c.Logs.Pretty, "logs.pretty", c.Logs.Pretty, "Enable pretty logging output")

	// Database configuration
	flag.StringVar(&c.Database.Dialect, "database.dialect", c.Database.Dialect, "Database dialect (e.g., postgres, mysql)")
	flag.StringVar(&c.Database.ConnectionString, "database.connection_string", c.Database.ConnectionString, "Database connection string")

	// App configuration
	flag.StringVar(&c.App.Kubernetes.AuthURL, "app.kubernetes.auth_url", c.App.Kubernetes.AuthURL, "Kubernetes authentication URL")
	flag.StringVar(&c.App.Kubernetes.AuthToken, "app.kubernetes.auth_token", c.App.Kubernetes.AuthToken, "Kubernetes authentication token")
	flag.StringVar(&c.App.Kubernetes.AuthCert, "app.kubernetes.auth_cert", c.App.Kubernetes.AuthCert, "Kubernetes authentication certificate")
	flag.StringVar(&c.App.Kubernetes.AuthCA, "app.kubernetes.auth_ca", c.App.Kubernetes.AuthCA, "Kubernetes authentication CA certificate")
	flag.StringVar(&c.App.Kubernetes.Namespace, "app.kubernetes.namespace", c.App.Kubernetes.Namespace, "Kubernetes namespace")

	flag.StringVar(&c.App.Helm.Repository, "app.helm.repository", c.App.Helm.Repository, "Helm chart repository URL")
	flag.StringVar(&c.App.Helm.BasicAuthUsername, "app.helm.basic_auth_username", c.App.Helm.BasicAuthUsername, "Helm chart repository basic auth username")
	flag.StringVar(&c.App.Helm.BasicAuthPassword, "app.helm.basic_auth_password", c.App.Helm.BasicAuthPassword, "Helm chart repository basic auth password")

	flag.StringVar(&c.App.Initializer.AdminPassword, "app.initializer.admin_password", c.App.Initializer.AdminPassword, "Admin password for the initializer")

	flag.StringVar(&c.App.Meta.WebSuffix, "app.meta.web_suffix", c.App.Meta.WebSuffix, "Web suffix for the application")
	flag.StringVar(&c.App.Meta.SocketSuffix, "app.meta.socket_suffix", c.App.Meta.SocketSuffix, "Socket suffix for the application")
	flag.IntVar(&c.App.Meta.SocketPort, "app.meta.socket_port", c.App.Meta.SocketPort, "Socket port for the application")

	flag.Parse()

	return nil
}

func (c *Config) Validate() error {
	// Server configuration
	if c.Server.Address == "" {
		return fmt.Errorf("Server address is required")
	}

	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("Server port must be between 1 and 65535")
	}

	// Logs configuration
	if c.Logs.Level == "" {
		return fmt.Errorf("Logs level is required")
	}

	// Database configuration
	if c.Database.Dialect != "postgres" && c.Database.Dialect != "mysql" && c.Database.Dialect != "sqlite" {
		return fmt.Errorf("Unsupported database dialect: %s", c.Database.Dialect)
	}

	if c.Database.ConnectionString == "" {
		return fmt.Errorf("Database connection string is required")
	}

	// App configuration
	if c.App.Kubernetes.AuthURL == "" {
		return fmt.Errorf("Kubernetes authentication URL is required")
	}

	if c.App.Kubernetes.AuthToken == "" {
		return fmt.Errorf("Kubernetes authentication token is required")
	}

	if c.App.Kubernetes.AuthCert == "" {
		return fmt.Errorf("Kubernetes authentication certificate is required")
	}

	if c.App.Kubernetes.AuthCA == "" {
		return fmt.Errorf("Kubernetes authentication CA certificate is required")
	}

	if c.App.Kubernetes.Namespace == "" {
		return fmt.Errorf("Kubernetes namespace is required")
	}

	if c.App.Helm.Repository == "" {
		return fmt.Errorf("Helm chart repository URL is required")
	}

	if c.App.Initializer.AdminPassword == "" {
		return fmt.Errorf("Admin password for the initializer is required")
	}

	if c.App.Meta.WebSuffix == "" {
		return fmt.Errorf("Web suffix for the application is required")
	}

	if c.App.Meta.SocketSuffix == "" {
		return fmt.Errorf("Socket suffix for the application is required")
	}

	if c.App.Meta.SocketPort <= 0 || c.App.Meta.SocketPort > 65535 {
		return fmt.Errorf("Socket port for the application must be between 1 and 65535")
	}

	return nil
}
