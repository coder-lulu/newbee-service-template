package svc

import (
	"context"

	"github.com/coder-lulu/newbee-common/v2/i18n"
	"github.com/coder-lulu/newbee-common/v2/middleware/audit"
	"github.com/coder-lulu/newbee-common/v2/middleware/integration"
	"github.com/coder-lulu/newbee-common/v2/middleware/keys"
	"github.com/coder-lulu/newbee-common/v2/orm/ent/hooks"
	"github.com/coder-lulu/newbee-core/rpc/coreclient"
	"github.com/coder-lulu/newbee-{{SERVICE_NAME}}-api/internal/config"
	i18n2 "github.com/coder-lulu/newbee-{{SERVICE_NAME}}-api/internal/i18n"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	ContextManager *keys.ContextManager
	CoreRpc        coreclient.Core
	Trans          *i18n.Translator
	// 🎯 统一中间件集成结果
	IntegrationResult *integration.Result
	// 统一中间件链
	ManagedMiddlewareChain []rest.Middleware
}

// RpcApiResourceProvider 通过RPC获取API资源名称的提供器
type RpcApiResourceProvider struct {
	coreRpc coreclient.Core
}

func NewRpcApiResourceProvider(coreRpc coreclient.Core) *RpcApiResourceProvider {
	return &RpcApiResourceProvider{coreRpc: coreRpc}
}

func (p *RpcApiResourceProvider) GetApiResourceName(ctx context.Context, method, path string) (string, error) {
	// 使用Core RPC服务获取API资源名称
	// 这里可以实现更复杂的资源名称解析逻辑
	return "{{SERVICE_NAME}}:" + method + ":" + path, nil
}

// GetCoreRpcClient 实现audit.AuditSvcProvider接口，为审计插件提供Core RPC客户端
func (svc *ServiceContext) GetCoreRpcClient() interface{} {
	return svc.CoreRpc
}

func NewServiceContext(c config.Config) *ServiceContext {
	// ===========================================
	// 🎉 统一中间件框架集成 - IO服务
	// ===========================================

	// 1. 初始化基础服务
	rds := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{c.RedisConf.Host},
		Password: c.RedisConf.Pass,
		DB:       c.RedisConf.Db,
	})

	trans := i18n.NewTranslator(c.I18nConf, i18n2.LocaleFS)

	// 2. 初始化RPC客户端 - 使用SystemContext拦截器支持系统级操作
	coreRpcClient, err := zrpc.NewClient(c.CoreRpc, zrpc.WithUnaryClientInterceptor(hooks.SystemContextClientInterceptor()))
	if err != nil {
		panic("Failed to create Core RPC client: " + err.Error())
	}
	coreRpc := coreclient.NewCore(coreRpcClient)

	// 3. 获取JWT密钥 - 优先使用Middleware配置，向后兼容Auth配置
	jwtSecret := c.Auth.AccessSecret
	if c.Middleware.Auth != nil && c.Middleware.Auth.AccessSecret != "" {
		jwtSecret = c.Middleware.Auth.AccessSecret
	}

	// 4. 创建服务上下文实例（需要先创建以便传递给审计插件）
	svcCtx := &ServiceContext{
		Config:  c,
		CoreRpc: coreRpc,
		Trans:   trans,
	}

	// 5. 🎯 使用统一中间件集成API，创建审计写入器
	ioAuditWriter := audit.NewBuiltinAuditWriter(svcCtx)

	result, err := integration.Setup(&integration.Config{
		Redis:               rds,
		JWTSecret:           jwtSecret,
		Mode:                integration.Production,
		ApiResourceProvider: NewRpcApiResourceProvider(coreRpc),
		// 使用common包的BuiltinAuditWriter通过Core RPC写入审计日志
		AuditWriter:        ioAuditWriter,
		TenantInfoProvider: NewRpcTenantInfoProvider(coreRpc),
		Middleware:         &c.Middleware,
	})
	if err != nil {
		panic("IO统一中间件集成失败: " + err.Error())
	}

	// 6. 完善服务上下文
	svcCtx.ContextManager = result.ContextManager
	svcCtx.IntegrationResult = result
	svcCtx.ManagedMiddlewareChain = result.Middlewares

	return svcCtx
}
