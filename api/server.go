package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/pawpaw2022/simplebank/db/postgresql"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// Params: endpoint, *middleware* ,handler
	router.POST("/accounts", server.createAccount)
	// router.GET("/accounts/:id", server.getAccount)
	// router.DELETE("/accounts/:id", server.deleteAccount)
	// router.GET("/accounts/:id/entries", server.listAccountEntries)
	// router.POST("/transfers", server.createTransfer)

	server.router = router
	return server
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// errorResponse is a helper to map an error to a JSON response
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}