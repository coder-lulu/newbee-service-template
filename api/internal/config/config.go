package config

import (
	"github.com/coder-lulu/newbee-common/v2/config"
	"github.com/coder-lulu/newbee-common/v2/i18n"
	"github.com/coder-lulu/newbee-common/v2/middleware/framework"
	"github.com/coder-lulu/newbee-common/v2/plugins/casbin"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth         rest.AuthConf             // Deprecated: Use Middleware.Auth instead, kept for routes.go compatibility
	Middleware   framework.UnifiedConfig `json:",optional"`
	RedisConf    config.RedisConf
	CoreRpc      zrpc.RpcClientConf
	DatabaseConf config.DatabaseConf
	CasbinConf   casbin.CasbinConf
	I18nConf     i18n.Conf
	CROSConf     config.CROSConf
}
