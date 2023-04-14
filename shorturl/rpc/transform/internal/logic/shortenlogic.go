package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/hash"
	"shorturl/rpc/transform/model"

	"shorturl/rpc/transform/internal/svc"
	"shorturl/rpc/transform/transform"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShortenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewShortenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShortenLogic {
	return &ShortenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ShortenLogic) Shorten(in *transform.ShortenReq) (*transform.ShortenResp, error) {
	// todo: add your logic here and delete this line
	// 手动代码开始，生成短链接
	key := hash.Md5Hex([]byte(in.Url))[:6]
	_, err := l.svcCtx.Model.Insert(l.ctx,&model.Shorturl{
		Shorten: key,
		Url:     in.Url,
	})
	if err != nil {
		return nil, err
	}

	return &transform.ShortenResp{
		Shorten: key,
	}, nil
	// 手动代码结束
	return &transform.ShortenResp{}, nil
}
