// Copyright (c) 2022 Orlov Boris onlycodergod@gmail.com.
// This file service.go is subject to the terms and
// conditions defined in file 'LICENSE', which is part of this project source code.
package server

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	srv *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.srv = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		WriteTimeout:   10 * time.Second,
		ReadTimeout:    10 * time.Second,
	}

	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
