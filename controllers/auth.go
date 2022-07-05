package controllers

import (
	"context"
	"fibergo/models"
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

//var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "user")

func SingUp(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var newUser models.User
	var userFound models.User
	defer cancel()
	err := c.BodyParser(&newUser)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewError(err))
	}
	if !govalidator.IsEmail(utils.NormalizeEmail(newUser.Email)) {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewError(utils.ErrInvalidEmail))
	}
	err = userCollection.FindOne(ctx, bson.M{"email": newUser.Email}).Decode(&userFound)
	if err != nil {
		if strings.TrimSpace(newUser.Password) == "" {
			return c.Status(fiber.StatusBadRequest).JSON(utils.NewError(utils.ErrEmptyPassword))
		}
		newUser.Password, err = security.EncryptPassword(newUser.Password)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.NewError(err))
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
		result, err := userCollection.InsertOne(ctx, newUser)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(responses.UserResponse{Status: fiber.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
		return c.Status(fiber.StatusCreated).JSON(responses.UserResponse{Status: fiber.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
	}
	if userFound.Email == newUser.Email {
		err = utils.ErrEmailAlreadyExists
	}
	return c.Status(fiber.StatusBadRequest).JSON(utils.NewError(err))
}

func SingIn(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var input models.UserLogin
	var user models.User
	defer cancel()
	err := c.BodyParser(&input)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewError(err))
	}
	log.Println("Input email: ", input.Email, "-Input password: ", input.Password)
	if !govalidator.IsEmail(utils.NormalizeEmail(input.Email)) {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewError(utils.ErrInvalidEmail))
	}
	err = userCollection.FindOne(ctx, bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		log.Printf("%s signin failed: %v\n", input.Email, err.Error())
		return c.Status(fiber.StatusUnauthorized).JSON(utils.NewError(utils.ErrInvalidCredentials))
	}
	if strings.TrimSpace(input.Password) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewError(utils.ErrEmptyPassword))
	}
	err = security.VerifyPassword(user.Password, input.Password)
	if err != nil {
		log.Printf("%s signin failed: %v\n", input.Email, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewError(err))
	}
	token, err := security.NewToken(user.Id.Hex())
	if err != nil {
		log.Printf("%s signin failed: %v\n", input.Email, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewError(err))
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user":  user,
		"token": fmt.Sprintf("Bearer %s", token),
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	token := c.Locals("user").(*jwt.Token)
	var user models.User
	defer cancel()
	payload, err := security.ParseToken(token.Raw)
	log.Println("token id: ", payload.Id, " - token issuer: ", payload.Issuer)
	//log.Println("el id es: ", id)
	if err != nil {
		return nil, err
	}
	objId, _ := primitive.ObjectIDFromHex(payload.Id)
	err = userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
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
