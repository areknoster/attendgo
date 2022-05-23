package webui

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/areknoster/attendgo/domain"
)

type ServerConfig struct{
	Address string `env:"ADDRESS,default=:8080"`
}

type Handler interface{
	Register(r chi.Router)
}

func NewServer(config ServerConfig, handlers ...Handler) *Server {
	r := chi.NewRouter()
	for _, h := range handlers{
		h.Register(r)
	}

	return &Server{
		server: &http.Server{
			Addr:              config.Address,
			Handler:           r,
		},
	}
}

var _ domain.Runner = (*Server)(nil)

type Server struct {
	server *http.Server
}

// Run implements domain.Runner
func (s *Server) Run(ctx context.Context) error {
	go func(){
		<-ctx.Done()
		s.server.Close()
	}()
	return s.server.ListenAndServe()
}

