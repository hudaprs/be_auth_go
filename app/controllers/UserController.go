package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"be_auth_go/app/helpers"
	"be_auth_go/app/models"
	"be_auth_go/app/utils"
)

// CreateUser Create a new User
func (app *App) CreateUser(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{"Status": "success", "Message": "User created"}
	user := &models.User{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Check the user
	userData, _ := user.GetUser(app.DB)
	if userData != nil {
		response["Status"] = "error"
		response["Message"] = "User already registered"
		helpers.JSON(w, http.StatusUnauthorized, response)
		return
	}

	// Trim
	user.Prepare()

	// Validate the user
	err = user.Validate("REGISTER")
	if err != nil {
		helpers.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Save the user
	newUser, err := user.SaveUser(app.DB)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}
	response["data"] = map[string]interface{}{"ID": newUser.ID, "Name": newUser.Name, "Email": newUser.Email, "CreatedAt": newUser.CreatedAt}
	helpers.JSON(w, http.StatusCreated, response)
	return
}

// Login a user
func (app *App) Login(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{"Status": "success", "Message": "Successfully LoggedIn"}
	user := &models.User{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Trim
	user.Prepare()

	// Validate
	err = user.Validate("LOGIN")
	if err != nil {
		helpers.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Check user
	userData, err := user.GetUser(app.DB)

	if err != nil {
		helpers.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if userData == nil {
		response["Status"] = "error"
		response["Message"] = "User is not registered"
		helpers.JSON(w, http.StatusBadRequest, response)
		return
	}

	// Check the password
	err = models.CheckHashPassword(user.Password, userData.Password)
	if err != nil {
		response["Status"] = "error"
		response["Message"] = "Invalid credentials"
		helpers.JSON(w, http.StatusBadRequest, response)
		return
	}

	token, err := utils.EncodeAuthToken(userData.ID)
	if err != nil {
		helpers.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response["token"] = token
	helpers.JSON(w, http.StatusOK, response)
	return
}