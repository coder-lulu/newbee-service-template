package svc

import (
	"github.com/coder-lulu/newbee-common/v2/orm/ent/hooks"
	"github.com/coder-lulu/newbee-{{SERVICE_NAME}}-rpc/ent"
	_ "github.com/coder-lulu/newbee-{{SERVICE_NAME}}-rpc/ent/runtime"
	"github.com/coder-lulu/newbee-{{SERVICE_NAME}}-rpc/internal/config"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config config.Config
	DB     *ent.Client
	Redis  redis.UniversalClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化Redis
	rds := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{c.RedisConf.Host},
		Password: c.RedisConf.Pass,
		DB:       c.RedisConf.Db,
	})

	// 初始化数据库客户端
	db := ent.NewClient(
		ent.Log(logx.Info), // logger
		ent.Driver(c.DatabaseConf.NewNoCacheDriver()),
		ent.Debug(), // debug mode
	)

	// ===========================================
	// 🎯 多租户集成 - 遵循 NewBee 编码准则
	// ===========================================

	// 使用统一Hook系统 - 一键设置租户和部门Hook
	// 根据服务需要，添加特定的系统表排除规则
	// hooks.AddExcludedTable("system_table1")
	// hooks.AddExcludedTable("system_table2")

	// 一键设置：初始化配置 + 注册所有hooks (租户Hook + 部门Hook)
	if err := hooks.QuickSetup(db); err != nil {
		logx.Errorw("Failed to setup unified hooks", logx.Field("error", err.Error()))
		panic("统一Hook初始化失败: " + err.Error())
	}

	// RPC服务层不注册数据权限拦截器
	// 数据权限控制应该在API层通过中间件处理，RPC层作为数据访问层不承担权限职责
	// 这样可以保持清晰的层次分离，避免跨服务的上下文传递问题

	logx.Infow("✅ {{SERVICE_NAME_UPPER}} RPC service: Unified hooks initialized successfully")

	return &ServiceContext{
		Config: c,
		DB:     db,
		Redis:  rds,
	}
}
