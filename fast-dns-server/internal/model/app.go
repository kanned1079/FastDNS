package model

import (
	"fast-dns-server/internal/resolver"
	"fast-dns-server/internal/router"
)

type App struct {
	Id          int64
	ApiGateway  *router.ApiInstance
	DnsResolver *resolver.DnsServerInst
}
