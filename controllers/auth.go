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
	if userFound.Email == newUser.Email {
		err = utils.ErrEmailAlreadyExists
	}
	return c.Status(fiber.StatusBadRequest).
		JSON(responses.UserResponse{
			Status:  fiber.StatusBadRequest,
			Message: "error",
			Data:    utils.NewError(err)})
}

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
	log.Println("Input email: ", input.Email, "-Input password: ", input.Password)
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
	log.Println("token id: ", payload.Id, " - token issuer: ", payload.Issuer)
	log.Println("el id es: ", id)
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
	log.Println("token id: ", payload.Id, " - token issuer: ", payload.Issuer)
	//log.Println("el id es: ", id)
	if err != nil {
		return nil, err
	}
	objId, _ := primitive.ObjectIDFromHex(payload.Id)
	user, err = repositories.FindOneUser(bson.M{"id": objId})
	log.Println("id del usuario: ", user.Id)
	log.Println("Profile del usuario: ", user.Profile)
	log.Println("Roles del usuario: ", user.Roles)
	if err != nil {
		log.Println("Hubo un error al encontrar el usuario")
		return nil, err
	}
	if !utils.ArrayContainsAny(user.Roles, roles) {
		return nil, utils.ErrUnauthorized
	}
	return payload, nil
}
