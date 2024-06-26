package api

import (
	db "simple_bank/db/sqlc"
	"simple_bank/token"
	"simple_bank/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	maker  token.Maker
	router *gin.Engine
	store  db.IStore
	config util.Config
}

func NewServer(s db.IStore, c util.Config) (*Server, error) {
	server := &Server{store: s, config: c}

	maker, err := token.NewJWTMaker(c.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	server.maker = maker
	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}


	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)


	authGroup := router.Use(authMiddleware(server.maker))
	authGroup.POST("/accounts", server.createAccount)
	authGroup.GET("/accounts/:id", server.getAccount)
	authGroup.GET("/accounts", server.listAccounts)

	authGroup.POST("/transfers", server.createTransfer)


	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
