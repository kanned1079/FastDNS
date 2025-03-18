// 定义 cache 对象的类型
type Cache = {
    cache_hint: number;
    cache_miss: number;
    cache_rate: number;
};

// 定义 dns 配置的类型
type DnsConfig = {
    tls_enabled: boolean;
    backup_secondary_list: string[];
    list: string[];
    mode: string;
};

// 定义 config 配置的类型
type GeneralConfig = {
    ipv6_enabled: boolean;
    speed_limit: number;
    cache_size: number;
    edns: boolean;
    edns_custom_ip_enabled: boolean;
    edns_custom_ip_addr: string;
    dnssec: boolean;
    log_buffer: number;
    log_size_limit: number;
};

// 定义 management 配置的类型
type ManagementConfig = {
    dns_server_listen_addr: string;
    backend_listen_addr: string;
    frontend_listen_addr: string;
};

// 定义完整的 config 对象的类型
type Config = {
    dns: DnsConfig;
    config: GeneralConfig;
    management: ManagementConfig;
};

// 定义 server 对象的类型
type Server = {
    arch: string;
    os: string;
    routine: number;
};

// 定义整个数据对象的类型
export type Response = {
    cache: Cache;
    code: number;
    config: Config;
    server: Server;
};