package logic

import (
	"context"
	"shorturl/rpc/transform/transform"

	"shorturl/api/internal/svc"
	"shorturl/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExpandLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExpandLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExpandLogic {
	return &ExpandLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExpandLogic) Expand(req *types.ExpandReq) (resp *types.ExpandResp, err error) {
	// todo: add your logic here and delete this line
	// 手动代码开始
	var urlResp *transform.ExpandResp
	urlResp, err = l.svcCtx.Transformer.Expand(l.ctx, &transform.ExpandReq{
		Shorten: req.Shorten,
	})
	if err != nil {
		return &types.ExpandResp{}, err
	}

	return &types.ExpandResp{
		Url: urlResp.Url,
	}, nil
	// 手动代码结束
}