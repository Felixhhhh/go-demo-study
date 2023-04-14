package logic

import (
	"context"
	"shorturl/rpc/transform/transform"

	"shorturl/api/internal/svc"
	"shorturl/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShortenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShortenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShortenLogic {
	return &ShortenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShortenLogic) Shorten(req *types.ShortenReq) (resp *types.ShortenResp, err error) {
	// todo: add your logic here and delete this line
	// 手动代码开始
	var shortenResp *transform.ShortenResp
	shortenResp, err = l.svcCtx.Transformer.Shorten(l.ctx, &transform.ShortenReq{
		Url: req.Url,
	})
	if err != nil {
		return &types.ShortenResp{}, err
	}

	return &types.ShortenResp{
		Shorten: shortenResp.Shorten,
	}, nil
	// 手动代码结束
	return
}
