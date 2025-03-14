package config

type RootConfig struct {
	Details Details `yaml:"details"`
}

type Details struct {
	Dns struct {
		TlsEnabled          bool     `yaml:"tls_enabled"`
		BackupSecondaryList []string `yaml:"backup_secondary_list"`
		List                []string `yaml:"list"`
	} `yaml:"dns"`
	Config struct {
		Ipv6Enabled bool  `yaml:"ipv6_enabled"`
		SpeedLimit  int64 `yaml:"speed_limit"`
		CacheSize   int64 `yaml:"cache_size"` //bytes
	} `yaml:"config"`
	Management struct {
		DnsServerListenAddr string `yaml:"dns_server_listen_addr"`
		BackendListenAddr   string `yaml:"backend_listen_addr"`
	}
}
