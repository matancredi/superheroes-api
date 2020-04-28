package controllers

import "github.com/matancredi/superheroes-api/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	//Posts routes
	s.Router.HandleFunc("/supers", middlewares.SetMiddlewareJSON(s.CreateSuper)).Methods("POST")
	s.Router.HandleFunc("/supers", middlewares.SetMiddlewareJSON(s.GetSupers)).Methods("GET")
	s.Router.HandleFunc("/supers/{uuid}", middlewares.SetMiddlewareJSON(s.GetSuperById)).Methods("GET")
	s.Router.HandleFunc("/supers/search/{name}", middlewares.SetMiddlewareJSON(s.GetSuperByName)).Methods("GET")
	s.Router.HandleFunc("/supers/alignment/{params}", middlewares.SetMiddlewareJSON(s.GetSuperByAlignment)).Methods("GET")
	s.Router.HandleFunc("/supers/{uuid}", middlewares.SetMiddlewareJSON(s.DeleteSuper)).Methods("DELETE")
}
