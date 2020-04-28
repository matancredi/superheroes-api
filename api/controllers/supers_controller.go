package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/matancredi/superheroes-api/api/models"
	"github.com/matancredi/superheroes-api/api/responses"
	"github.com/matancredi/superheroes-api/api/utils/formaterror"
)

func (server *Server) CreateSuper(w http.ResponseWriter, r *http.Request) {

	// Load values from .env file
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

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
	response, err := http.Get(os.Getenv("API_HEROES") + os.Getenv("ACCESS_TOKEN") + "/search/" + superToSearch.Name)

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
