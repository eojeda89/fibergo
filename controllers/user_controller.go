package controllers

import (
	"fibergo/models"
	"fibergo/repositories"
	"fibergo/responses"
	"fibergo/security"
	"fibergo/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/asaskevich/govalidator.v9"
	"time"
)

var validate = validator.New()

// Create an user godoc
//@Summary Create an user.
//@Description Create an user .
//@Param user body models.User true "The user object to be created"
//@Accept json
//@Produce json
//@Success 201 {object} responses.UserResponse{data=string}
//@Failure 400 {object} responses.UserResponse{data=utils.JError}
//@Failure 401 {object} responses.UserResponse{data=utils.JError}
//@Failure 500 {object} responses.UserResponse{data=utils.JError}
//@Router /users [post]
func CreateUser(c *fiber.Ctx) error {
	_, err := AuthRequestWithRole(c, []string{"admin", "manager"})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(responses.UserResponse{
				Status:  fiber.StatusUnauthorized,
				Message: "error",
				Data:    utils.NewError(utils.ErrUnauthorized)})
	}
	var user models.User

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(responses.UserResponse{
				Status:  fiber.StatusBadRequest,
				Message: "error",
				Data:    utils.NewError(err)})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(responses.UserResponse{
				Status:  fiber.StatusBadRequest,
				Message: "error",
				Data:    utils.NewError(validationErr)})
	}
	var userFound models.User
	if !govalidator.IsEmail(utils.NormalizeEmail(user.Email)) {
		return c.Status(fiber.StatusBadRequest).
			JSON(responses.UserResponse{
				Status:  fiber.StatusBadRequest,
				Message: "error",
				Data:    utils.NewError(utils.ErrInvalidEmail)})
	}
	userFound, err = repositories.FindOneUser(bson.M{"email": user.Email})
	if userFound.Email != "" {
		return c.Status(fiber.StatusBadRequest).
			JSON(responses.UserResponse{
				Status:  fiber.StatusBadRequest,
				Message: "error",
				Data:    utils.NewError(utils.ErrEmailAlreadyExists)})
	}

	user.Password, err = security.EncryptPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(responses.UserResponse{
				Status:  fiber.StatusBadRequest,
				Message: "error",
				Data:    utils.NewError(err)})
	}
	user.CreatedAt = primitive.Timestamp{
		T: uint32(time.Now().Unix()),
		I: 0,
	}
	user.UpdatedAt = user.CreatedAt

	newProfile := models.Profile{
		Age:      user.Profile.Age,
		Country:  user.Profile.Country,
		Image:    user.Profile.Image,
		FullName: user.Profile.FullName,
	}

	newUser := models.User{
		Id:        primitive.NewObjectID(),
		Email:     user.Email,
		Password:  user.Password,
		Status:    user.Status,
		Roles:     user.Roles,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Profile:   newProfile,
	}

	_, err = repositories.InsertOneUser(newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(responses.UserResponse{
				Status:  fiber.StatusInternalServerError,
				Message: "error",
				Data:    utils.NewError(err)})
	}

	return c.Status(fiber.StatusCreated).
		JSON(responses.UserResponse{
			Status:  fiber.StatusCreated,
			Message: "success",
			Data:    &fiber.Map{"id": newUser.Id}})
}

func CreateUser2(c *fiber.Ctx) error {
	_, err := AuthRequestWithRole(c, []string{"admin", "manager"})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(responses.UserResponse{
				Status:  fiber.StatusUnauthorized,
				Message: "error",
				Data:    utils.NewError(utils.ErrUnauthorized)})
	}
	var userPost models.UserPost

	//validate the request body
	if err := c.BodyParser(&userPost); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(responses.UserResponse{
				Status:  fiber.StatusBadRequest,
				Message: "error",
				Data:    utils.NewError(err)})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&userPost); validationErr != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(responses.UserResponse{
				Status:  fiber.StatusBadRequest,
				Message: "error",
				Data:    utils.NewError(validationErr)})
	}

	newProfile := models.Profile{
		Age:      userPost.Age,
		Country:  userPost.Country,
		Image:    userPost.Image,
		FullName: userPost.FullName,
	}

	newUser := models.User{
		Id:        primitive.NewObjectID(),
		Email:     userPost.Email,
		Password:  userPost.Password,
		Status:    userPost.Status,
		Roles:     userPost.Roles,
		CreatedAt: userPost.CreatedAt,
		UpdatedAt: userPost.UpdatedAt,
		Profile:   newProfile,
	}

	result, err := repositories.InsertOneUser(newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(responses.UserResponse{
				Status:  fiber.StatusInternalServerError,
				Message: "error",
				Data:    utils.NewError(err)})
	}

	return c.Status(fiber.StatusCreated).
		JSON(responses.UserResponse{
			Status:  fiber.StatusCreated,
			Message: "success",
			Data:    result})
}

// Get an user godoc
//@Summary Show an user.
//@Description Get an user by it´s ID.
//@Param id path string true "The user ID to be showed"
//@Produce json
//@Success 200 {object} responses.UserResponse{data=models.User}
//@Failure 401 {object} responses.UserResponse{data=utils.JError}
//@Failure 404 {object} responses.UserResponse{data=utils.JError}
//@Router /users/{id} [get]
func GetAUser(c *fiber.Ctx) error {
	payload, err := AuthRequestWithId(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(responses.UserResponse{
				Status:  fiber.StatusUnauthorized,
				Message: "error",
				Data:    utils.NewError(utils.ErrUnauthorized)})
	}
	userId := payload.Id
	objId, _ := primitive.ObjectIDFromHex(userId)

	user, err := repositories.FindOneUser(bson.M{"id": objId})
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(responses.UserResponse{
				Status:  fiber.StatusNotFound,
				Message: "error",
				Data:    utils.NewError(err)})
	}

	return c.Status(fiber.StatusOK).
		JSON(responses.UserResponse{
			Status:  fiber.StatusOK,
			Message: "success",
			Data:    user})
}

// Edit an user godoc
//@Summary Edit an user.
//@Description Edit an user by it´s ID.
//@Param id path string true "The user ID to be edited"
//@Param input body models.User true "The user object with the new values"
//@Produce json
//@Success 200 {object} responses.UserResponse{data=models.User}
//@Failure 400 {object} responses.UserResponse{data=utils.JError}
//@Failure 401 {object} responses.UserResponse{data=utils.JError}
//@Failure 500 {object} responses.UserResponse{data=utils.JError}
//@Router /users/{id} [put]
func EditAUser(c *fiber.Ctx) error {
	_, err := AuthRequestWithRole(c, []string{"admin", "manager"})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(responses.UserResponse{
				Status:  fiber.StatusUnauthorized,
				Message: "error",
				Data:    utils.NewError(utils.ErrUnauthorized)})
	}
	userId := c.Params("userId")
	var user models.User

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(responses.UserResponse{
				Status:  fiber.StatusBadRequest,
				Message: "error",
				Data:    utils.NewError(err)})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(responses.UserResponse{
				Status:  fiber.StatusBadRequest,
				Message: "error",
				Data:    utils.NewError(validationErr)})
	}

	update := bson.M{"email": user.Email, "password": user.Password, "status": user.Status, "roles": user.Roles, "profile": user.Profile, "updatedat": primitive.Timestamp{
		T: uint32(time.Now().Unix()),
		I: 0,
	}}

	result, err := repositories.EditUserById(userId, update)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(responses.UserResponse{
				Status:  fiber.StatusInternalServerError,
				Message: "error",
				Data:    utils.NewError(err)})
	}
	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusBadRequest).
			JSON(responses.UserResponse{
				Status:  fiber.StatusBadRequest,
				Message: "error",
				Data:    utils.NewError(utils.ErrUserNotFound)})
	}
	//get updated user details
	var updatedUser models.User
	objId, _ := primitive.ObjectIDFromHex(userId)
	if result.MatchedCount == 1 {
		updatedUser, err = repositories.FindOneUser(bson.M{"id": objId})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).
				JSON(responses.UserResponse{
					Status:  fiber.StatusInternalServerError,
					Message: "error",
					Data:    utils.NewError(err)})
		}
	}
	return c.Status(fiber.StatusOK).
		JSON(responses.UserResponse{
			Status:  fiber.StatusOK,
			Message: "success",
			Data:    updatedUser})
}

// Delete an user godoc
//@Summary Delete an user.
//@Description Delete an user by it´s ID.
//@Param id path string true "The user ID to be deleted"
//@Produce json
//@Success 200 {object} responses.UserResponse{data=string}
//@Failure 401 {object} responses.UserResponse{data=utils.JError}
//@Failure 404 {object} responses.UserResponse{data=utils.JError}
//@Failure 500 {object} responses.UserResponse{data=utils.JError}
//@Router /users/{id} [delete]
func DeleteAUser(c *fiber.Ctx) error {
	_, err := AuthRequestWithRole(c, []string{"admin", "manager"})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(responses.UserResponse{
				Status:  fiber.StatusUnauthorized,
				Message: "error",
				Data:    utils.NewError(utils.ErrUnauthorized)})
	}
	userId := c.Params("userId")
	result, err := repositories.DeleteUserById(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(responses.UserResponse{
				Status:  fiber.StatusInternalServerError,
				Message: "error",
				Data:    utils.NewError(err)})
	}
	if result.DeletedCount < 1 {
		return c.Status(fiber.StatusNotFound).
			JSON(responses.UserResponse{
				Status:  fiber.StatusNotFound,
				Message: "error",
				Data:    utils.NewError(utils.ErrUserNotFound)},
			)
	}
	return c.Status(fiber.StatusOK).
		JSON(responses.UserResponse{
			Status:  fiber.StatusOK,
			Message: "success",
			Data:    "User successfully deleted!"},
		)
}

// Get all users godoc
//@Summary Show a list of all users.
//@Description Get a list of all users.
//@Produce json
//@Success 200 {object} responses.UserResponse{data=[]models.User}
//@Failure 500 {object} responses.UserResponse{data=utils.JError}
//@Router /users [get]
func GetAllUsers(c *fiber.Ctx) error {
	users, err := repositories.FindAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(responses.UserResponse{
				Status:  fiber.StatusInternalServerError,
				Message: "error",
				Data:    utils.NewError(err)})
	}
	return c.Status(fiber.StatusOK).
		JSON(responses.UserResponse{
			Status:  fiber.StatusOK,
			Message: "success",
			Data:    users},
		)
}
