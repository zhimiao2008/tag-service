package server

import (
	"context"
	"github.com/weizhimiao/tag-service/pkg/errcode"
	"google.golang.org/grpc/metadata"
)

type Auth struct {
	AppKey    string
	AppSecret string
}

func (a *Auth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"app_key":    a.AppKey,
		"app_secret": a.AppSecret,
	}, nil
}

func (a *Auth) RequireTransportSecurity() bool {
	return false
}

func (a *Auth) GetAppKey() string {
	return "appkeysss"
}
func (a *Auth) GetAppSecret() string {
	return "secretsecret"
}

func (a *Auth) Check(ctx context.Context) error {
	md,_:= metadata.FromIncomingContext(ctx)

	var appKey, appSecret string
	if value, ok:= md["app_key"]; ok {
		appKey = value[0]
	}

	if value, ok:= md["app_secret"]; ok {
		appSecret = value[0]
	}

	if appKey != a.GetAppKey() || appSecret != a.GetAppSecret() {
		return errcode.TogRPCError(errcode.Unauthorized)
	}
	return nil
}

