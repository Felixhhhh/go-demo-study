// Code generated by goctl. DO NOT EDIT!
// Source: transform.proto

package transformclient

import (
	"context"

	"shorturl/rpc/transform/transform"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	ExpandReq   = transform.ExpandReq
	ExpandResp  = transform.ExpandResp
	ShortenReq  = transform.ShortenReq
	ShortenResp = transform.ShortenResp

	Transform interface {
		Expand(ctx context.Context, in *ExpandReq, opts ...grpc.CallOption) (*ExpandResp, error)
		Shorten(ctx context.Context, in *ShortenReq, opts ...grpc.CallOption) (*ShortenResp, error)
	}

	defaultTransform struct {
		cli zrpc.Client
	}
)

func NewTransform(cli zrpc.Client) Transform {
	return &defaultTransform{
		cli: cli,
	}
}

func (m *defaultTransform) Expand(ctx context.Context, in *ExpandReq, opts ...grpc.CallOption) (*ExpandResp, error) {
	client := transform.NewTransformClient(m.cli.Conn())
	return client.Expand(ctx, in, opts...)
}

func (m *defaultTransform) Shorten(ctx context.Context, in *ShortenReq, opts ...grpc.CallOption) (*ShortenResp, error) {
	client := transform.NewTransformClient(m.cli.Conn())
	return client.Shorten(ctx, in, opts...)
}
