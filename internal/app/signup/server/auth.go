/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc/peer"
)

//AuthData deals with the certificate authentication process
type AuthData struct {
	ClientSecret string
}

//Authenticate validates the client certificate in every gRPC request
func (a AuthData) Authenticate(ctx context.Context) (context.Context, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return ctx, status.Error(codes.Unauthenticated, "no peer found")
	}

	tlsAuth, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return ctx, status.Error(codes.Unauthenticated, "unexpected peer credentials")
	}

	if len(tlsAuth.State.VerifiedChains) == 0 || len(tlsAuth.State.VerifiedChains[0]) == 0 {
		return ctx, status.Error(codes.Unauthenticated, "invalid certificate")
	}

	if tlsAuth.State.VerifiedChains[0][0].Subject.CommonName != a.ClientSecret {
		return ctx, status.Error(codes.Unauthenticated, "invalid client certificate secret")
	}

	return ctx, nil
}
