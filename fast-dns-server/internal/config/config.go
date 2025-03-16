package config

type RootConfig struct {
	Details Details `yaml:"details" json:"details"`
}

type Details struct {
	Dns struct {
		TlsEnabled          bool     `yaml:"tls_enabled" json:"tls_enabled"`
		BackupSecondaryList []string `yaml:"backup_secondary_list" json:"backup_secondary_list"`
		List                []string `yaml:"list" json:"list"`
		Mode                string   `yaml:"mode" json:"mode"`
	} `yaml:"dns" json:"dns"`
	Config struct {
		Ipv6Enabled         bool   `yaml:"ipv6_enabled" json:"ipv6_enabled"`
		SpeedLimit          int64  `yaml:"speed_limit" json:"speed_limit"`
		CacheSize           int64  `yaml:"cache_size" json:"cache_size"` //bytes
		Edns                bool   `yaml:"edns" json:"edns"`
		EdnsCustomIpEnabled bool   `yaml:"edns_custom_ip_enabled" json:"edns_custom_ip_enabled"`
		EdnsCustomIpAddr    string `yaml:"edns_custom_ip_addr" json:"edns_custom_ip_addr"`
		Dnssec              bool   `yaml:"dnssec" json:"dnssec"`
		LogBuffer           int    `yaml:"log_buffer" json:"log_buffer"`
		LogSizeLimit        int64  `yaml:"log_size_limit" json:"log_size_limit"`
	} `yaml:"config" json:"config"`
	Management struct {
		DnsServerListenAddr string `yaml:"dns_server_listen_addr" json:"dns_server_listen_addr"`
		BackendListenAddr   string `yaml:"backend_listen_addr" json:"backend_listen_addr"`
		FrontendListenAddr  string `yaml:"frontend_listen_addr" json:"frontend_listen_addr"`
	} `yaml:"management" json:"management"`
}
