package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"go.uber.org/fx"
)

type Server struct {
	server *http.Server

	logger logger.Logger
	config config.Config
}

type Params struct {
	fx.In

	Logger logger.Logger
	Config config.Config
}

func NewServer(params Params) Server {
	server := &http.Server{
		Addr: "0.0.0.0:3000",
	}

	return Server{
		server: server,
		logger: params.Logger,
		config: params.Config,
	}
}

func (s *Server) Start() error {
	router := http.NewServeMux()
	s.server.Handler = router

	router.Handle("/metrics", promhttp.Handler())

	s.logger.Info("http server is listening", slog.String("addr", s.server.Addr))

	if err := s.server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return fmt.Errorf("listen and serve: %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
