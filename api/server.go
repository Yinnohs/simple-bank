package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/yinnohs/simple-bank/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	// mainRouterGroup := router.Group("/api/v1/", gin.WrapH(http.DefaultServeMux))
	// accountRouterGroup := mainRouterGroup.Group("accounts", gin.WrapH(http.DefaultServeMux))

	router.POST("/api/v1/account/", server.createAccout)
	router.GET("/api/v1/account", server.listAccounts)
	router.GET("/api/v1/account/:id", server.getAccount)
	router.PUT("/api/v1/account/add-balance/", server.addAccountBalance)
	router.POST("/api/v1/transfer/", server.createNewTransfer)

	server.router = router

	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
