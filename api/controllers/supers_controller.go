package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/matancredi/superheroes-api/api/models"
	"github.com/matancredi/superheroes-api/api/responses"
	"github.com/matancredi/superheroes-api/api/utils/formaterror"
)

func (server *Server) CreateSuper(w http.ResponseWriter, r *http.Request) {

	// Gets the name of hero to be searched and registered
	bodyName, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	superToSearch := models.Super{}
	err = json.Unmarshal(bodyName, &superToSearch)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Sets up the URL used to search for the superhero
	response, err := http.Get(server.ApiUrl + superToSearch.Name)

	// Gets the response from URL
	body, err := ioutil.ReadAll(response.Body)

	results := models.Results{}
	err = json.Unmarshal(body, &results)

	// Returns error if no heroes is found
	if len(results.Supers) == 0 {
		responses.ERROR(w, http.StatusNotFound, err)
		return
	}

	// Saves the first result for the superhero searched
	superhero := results.Supers[0]

	superhero.Prepare()
	err = superhero.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	superCreated, err := superhero.SaveSuper(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	group := models.Group{}
	group.CreateGroup(server.DB, superCreated.Uuid, superhero.Connections.GroupAffiliation)

	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, superhero.Uuid))
	responses.JSON(w, http.StatusCreated, superCreated)

}

func (server *Server) GetSupers(w http.ResponseWriter, r *http.Request) {

	super := models.Super{}

	supers, err := super.FindAllSupers(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, supers)
}

func (server *Server) GetSuperById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Gets and converts uuid to int64
	uuid, err := strconv.ParseUint(vars["uuid"], 10, 64)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	super := models.Super{}

	superReceived, err := super.FindSuperByID(server.DB, uuid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, superReceived)

}

func (server *Server) GetSuperByName(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	name := vars["name"]

	super := models.Super{}

	superReceived, err := super.FindSuperByName(server.DB, name)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, superReceived)

}

func (server *Server) GetSuperByAlignment(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	alignment := vars["params"]

	super := models.Super{}

	superReceived, err := super.FindSuperByAlignment(server.DB, alignment)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, superReceived)

}

func (server *Server) DeleteSuper(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Gets and converts uuid to int64
	uuid, err := strconv.ParseUint(vars["uuid"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	super := models.Super{}
	err = server.DB.Debug().Model(models.Super{}).Where("uuid = ?", uuid).Take(&super).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("StatusNotFound"))
		return
	}

	_, err = super.DeleteASuper(server.DB, uuid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uuid))
	responses.JSON(w, http.StatusOK, uuid)
}
