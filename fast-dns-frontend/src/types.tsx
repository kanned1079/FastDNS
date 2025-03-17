export type Cache = {
    cache_hint: number;
    cache_miss: number;
    cache_rate: number;
}

export type Config = {
    dns: DNS;
    management: Management;
    config: {
        cache_size: number;
        dnssec: boolean;
        edns: boolean;
        edns_custom_ip_addr: string;
        edns_custom_ip_enabled: boolean;
        ipv6_enabled: boolean;
        log_buffer: number;
        log_size_limit: number;
        speed_limit: number;
    }
}

export type DNS = {
    backup_secondary_list: string[];
    list: string[];
    mode: "balance" | "parallel"; // 根据实际值调整
    tls_enabled: boolean;
}

export type Management = {
    backend_listen_addr: string;
    dns_server_listen_addr: string;
    frontend_listen_addr: string;
}

export type Response = {
    cache: Cache;
    code: number;
    config: Config;
    dns: DNS

}