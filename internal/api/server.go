package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zillalikestocode/summarize-api/internal/api/routes"
)

type Server struct {
	addr   string
	router *gin.Engine
}

func NewServer(addr string) *Server {
	s := &Server{
		addr:   addr,
		router: gin.Default(),
	}

	s.setupRoutes()
	return s
}

func (s *Server) Run() error {
	return s.router.Run(s.addr)
}

func (s *Server) setupRoutes() {
	s.router.Use(cors.New(cors.Config{
		AllowMethods: []string{"POST", "GET"},
		AllowOrigins: []string{"*"}, AllowHeaders: []string{"Content-type"}}))
	routes.SetupRoutes(s.router)
}
