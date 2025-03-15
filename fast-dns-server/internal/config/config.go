package config

type RootConfig struct {
	Details Details `yaml:"details"`
}

type Details struct {
	Dns struct {
		TlsEnabled          bool     `yaml:"tls_enabled"`
		BackupSecondaryList []string `yaml:"backup_secondary_list"`
		List                []string `yaml:"list"`
		Mode                string   `yaml:"mode"`
	} `yaml:"dns"`
	Config struct {
		Ipv6Enabled         bool   `yaml:"ipv6_enabled"`
		SpeedLimit          int64  `yaml:"speed_limit"`
		CacheSize           int64  `yaml:"cache_size"` //bytes
		Ends                bool   `yaml:"ends"`
		EdnsCustomIpEnabled bool   `yaml:"edns_custom_ip_enabled"`
		EdnsCustomIpAddr    string `yaml:"edns_custom_ip_addr"`
		Dnssec              bool   `yaml:"dnssec"`
	} `yaml:"config"`
	Management struct {
		DnsServerListenAddr string `yaml:"dns_server_listen_addr"`
		BackendListenAddr   string `yaml:"backend_listen_addr"`
		FrontendListenAddr  string `yaml:"frontend_listen_addr"`
	}
}
