package controllers

import (
	"fibergo/models"
	"fibergo/repositories"
	"fibergo/responses"
	"fibergo/security"
	"fibergo/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/asaskevich/govalidator.v9"
	"log"
	"strings"
	"time"
)

// Register an user godoc
//@Summary Register an user.
//@Description Register an user .
//@Param user body models.User true "The user object to be registered"
//@Accept json
//@Produce json
//@Success 201 {object} responses.UserResponse{data=string}
//@Failure 400 {object} responses.UserResponse{data=utils.JError}
//@Failure 500 {object} responses.UserResponse{data=utils.JError}
//@Router /signup [post]
func SingUp(c *fiber.Ctx) error {
	var newUser models.User
	var userFound models.User
	err := c.BodyParser(&newUser)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewError(err))
	}
	if !govalidator.IsEmail(utils.NormalizeEmail(newUser.Email)) {
		return c.Status(fiber.StatusBadRequest).
			JSON(responses.UserResponse{
				Status:  fiber.StatusBadRequest,
				Message: "error",
				Data:    utils.NewError(utils.ErrInvalidEmail)})
	}
	userFound, err = repositories.FindOneUser(bson.M{"email": newUser.Email})
	if err != nil {
		if strings.TrimSpace(newUser.Password) == "" {
			return c.Status(fiber.StatusBadRequest).
				JSON(responses.UserResponse{
					Status:  fiber.StatusBadRequest,
					Message: "error",
					Data:    utils.NewError(utils.ErrEmptyPassword)})
		}
		newUser.Password, err = security.EncryptPassword(newUser.Password)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).
				JSON(responses.UserResponse{
					Status:  fiber.StatusBadRequest,
					Message: "error",
					Data:    utils.NewError(err)})
		}
		newUser.CreatedAt = primitive.Timestamp{
			T: uint32(time.Now().Unix()),
			I: 0,
		}
		newUser.UpdatedAt = newUser.CreatedAt
		//inserting the user
		newProfile := models.Profile{
			Age:      newUser.Profile.Age,
			Country:  newUser.Profile.Country,
			Image:    newUser.Profile.Image,
			FullName: newUser.Profile.FullName,
		}
		newUser := models.User{
			Id:        primitive.NewObjectID(),
			Email:     newUser.Email,
			Password:  newUser.Password,
			Status:    newUser.Status,
			Roles:     newUser.Roles,
			CreatedAt: newUser.CreatedAt,
			UpdatedAt: newUser.UpdatedAt,
			Profile:   newProfile,
		}
		_, err := repositories.InsertOneUser(newUser)
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
	if userFound.Email == newUser.Email {
		err = utils.ErrEmailAlreadyExists
	}
	return c.Status(fiber.StatusBadRequest).
		JSON(responses.UserResponse{
			Status:  fiber.StatusBadRequest,
			Message: "error",
			Data:    utils.NewError(err)})
}

// Singnin godoc
//@Summary Signin.
//@Description Signin an user .
//@Param input body models.UserLogin true "The signin object"
//@Accept json
//@Produce json
//@Success 200 {object} responses.UserResponse{data=responses.SigninResponse}
//@Failure 400 {object} responses.UserResponse{data=utils.JError}
//@Failure 401 {object} responses.UserResponse{data=utils.JError}
//@Failure 422 {object} responses.UserResponse{data=utils.JError}
//@Failure 500 {object} responses.UserResponse{data=utils.JError}
//@Router /signin [post]
func SingIn(c *fiber.Ctx) error {
	var input models.UserLogin
	var user models.User
	err := c.BodyParser(&input)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).
			JSON(responses.UserResponse{
				Status:  fiber.StatusUnprocessableEntity,
				Message: "error",
				Data:    utils.NewError(err)})
	}
	if !govalidator.IsEmail(utils.NormalizeEmail(input.Email)) {
		return c.Status(fiber.StatusBadRequest).
			JSON(responses.UserResponse{
				Status:  fiber.StatusBadRequest,
				Message: "error",
				Data:    utils.NewError(utils.ErrInvalidEmail)})
	}
	user, err = repositories.FindOneUser(bson.M{"email": input.Email})
	if err != nil {
		log.Printf("%s signin failed: %v\n", input.Email, err.Error())
		return c.Status(fiber.StatusUnauthorized).
			JSON(responses.UserResponse{
				Status:  fiber.StatusUnauthorized,
				Message: "error",
				Data:    utils.NewError(utils.ErrInvalidCredentials)})
	}
	if strings.TrimSpace(input.Password) == "" {
		return c.Status(fiber.StatusBadRequest).
			JSON(responses.UserResponse{
				Status:  fiber.StatusBadRequest,
				Message: "error",
				Data:    utils.NewError(utils.ErrEmptyPassword)})
	}
	err = security.VerifyPassword(user.Password, input.Password)
	if err != nil {
		log.Printf("%s signin failed: %v\n", input.Email, err.Error())
		return c.Status(fiber.StatusBadRequest).
			JSON(responses.UserResponse{
				Status:  fiber.StatusBadRequest,
				Message: "error",
				Data:    utils.NewError(err)})
	}
	token, err := security.NewToken(user.Id.Hex())
	if err != nil {
		log.Printf("%s signin failed: %v\n", input.Email, err.Error())
		return c.Status(fiber.StatusBadRequest).
			JSON(responses.UserResponse{
				Status:  fiber.StatusBadRequest,
				Message: "error",
				Data:    utils.NewError(err)})
	}
	return c.Status(fiber.StatusOK).
		JSON(responses.UserResponse{
			Status:  fiber.StatusOK,
			Message: "success",
			Data: responses.SigninResponse{
				User:  user,
				Token: fmt.Sprintf("Bearer %s", token),
			},
		})
}

func AuthRequestWithId(c *fiber.Ctx) (*jwt.StandardClaims, error) {
	id := c.Params("userId")
	token := c.Locals("user").(*jwt.Token)
	payload, err := security.ParseToken(token.Raw)
	if err != nil {
		return nil, err
	}
	if payload.Id != id || payload.Issuer != id {
		return nil, utils.ErrUnauthorized
	}
	return payload, nil
}

func AuthRequestWithRole(c *fiber.Ctx, roles []string) (*jwt.StandardClaims, error) {
	token := c.Locals("user").(*jwt.Token)
	var user models.User
	payload, err := security.ParseToken(token.Raw)
	if err != nil {
		return nil, err
	}
	objId, _ := primitive.ObjectIDFromHex(payload.Id)
	user, err = repositories.FindOneUser(bson.M{"id": objId})
	if err != nil {
		return nil, err
	}
	if !utils.ArrayContainsAny(user.Roles, roles) {
		return nil, utils.ErrUnauthorized
	}
	return payload, nil
}
