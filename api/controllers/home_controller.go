package controllers

import (
	"net/http"

	"github.com/matancredi/superheroes-api/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Desafio Levpay - Mariana Tancredi")

}
