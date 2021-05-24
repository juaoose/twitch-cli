// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package mock_server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints"
)

const MOCK_NAMESPACE = "/mock"
const UNITS_NAMESPACE = "/units"
const AUTH_NAMESPACE = "/mock_auth"

func StartServer(port int) {
	m := http.NewServeMux()

	ctx := context.Background()

	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err.Error())
		return
	}

	ctx = context.WithValue(ctx, "db", db)

	RegisterHandlers(m)
	s := http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: m,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
		return
	}
}

func RegisterHandlers(m *http.ServeMux) {
	// all mock endpoints live in the /mock/ namespace
	for _, e := range endpoints.All() {
		m.Handle(MOCK_NAMESPACE+e.GetPath(), authentication.AuthenticationMiddleware(e))
	}
}
