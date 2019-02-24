/*
Server implementation of the board game Go.
*/
package server

import (
	"github.com/camirmas/go_stop/models"
	"github.com/camirmas/go_stop/providers"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"log"
	"net/http"
)

type ServerConfig struct {
	Host       string
	Port       int
}

type Server struct {
	config *ServerConfig
	db     models.DB
	schema *graphql.Schema
	auth   *providers.Auth
}

func NewServer(config *ServerConfig, db models.DB, auth *providers.Auth, schema *graphql.Schema) *Server {
	return &Server{
		config,
		db,
		schema,
		auth,
	}
}

func (s *Server) Start() {
	h := handler.New(&handler.Config{
		Schema:   s.schema,
		Pretty:   true,
		GraphiQL: true,
	})

	log.Println("Listening on http://localhost:8000")
	http.Handle("/graphql", enableCorsMiddleware(s.VerifyTokenMiddleware(gqlMiddleware(h))))
	http.Handle("/login", s.LoginHandler())
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func gqlMiddleware(next *handler.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ContextHandler(r.Context(), w, r)
	})
}

func enableCorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		next.ServeHTTP(w, r)
	})
}
