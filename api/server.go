package api

import (
	db "simple_bank/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	store  db.IStore
}

func NewServer(s db.IStore) *Server {
	server := &Server{store: s}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
