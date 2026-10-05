package base

import (
	"context"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/coder-lulu/newbee-common/v2/msg/errormsg"
	"github.com/coder-lulu/newbee-common/v2/msg/logmsg"
	"github.com/zeromicro/go-zero/core/errorx"

	"github.com/coder-lulu/newbee-{{SERVICE_NAME}}-rpc/internal/svc"
	"github.com/coder-lulu/newbee-{{SERVICE_NAME}}-rpc/types/{{SERVICE_NAME}}"

	"github.com/zeromicro/go-zero/core/logx"
)

type InitDatabaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInitDatabaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitDatabaseLogic {
	return &InitDatabaseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *InitDatabaseLogic) InitDatabase(in *{{SERVICE_NAME}}.Empty) (*{{SERVICE_NAME}}.BaseResp, error) {
    if err := l.svcCtx.DB.Schema.Create(l.ctx, schema.WithForeignKeys(false)); err != nil {
        logx.Errorw(logmsg.DatabaseError, logx.Field("detail", err.Error()))
        return nil, errorx.NewInternalError(err.Error())
    }

	return &{{SERVICE_NAME}}.BaseResp{Msg: errormsg.Success}, nil
}
