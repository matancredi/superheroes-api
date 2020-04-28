package controllers

import "github.com/matancredi/superheroes-api/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	//Posts routes
	s.Router.HandleFunc("/supers", middlewares.SetMiddlewareJSON(s.CreateSuper)).Methods("POST")
	s.Router.HandleFunc("/supers", middlewares.SetMiddlewareJSON(s.GetSupers)).Methods("GET")
	s.Router.HandleFunc("/supers/{id}", middlewares.SetMiddlewareJSON(s.GetSuper)).Methods("GET")
	s.Router.HandleFunc("/supers/{id}", middlewares.SetMiddlewareJSON(s.DeleteSuper)).Methods("DELETE")
}
