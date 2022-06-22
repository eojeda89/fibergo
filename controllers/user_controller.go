package controllers

import (
	"context"
	"fibergo/configs"
	"fibergo/models"
	"fibergo/responses"
	"go.mongodb.org/mongo-driver/bson"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "user")
var validate = validator.New()

func CreateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{Status: fiber.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{Status: fiber.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

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

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.UserResponse{Status: fiber.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(fiber.StatusCreated).JSON(responses.UserResponse{Status: fiber.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}
func CreateUser2(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var userPost models.UserPost
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&userPost); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{Status: fiber.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&userPost); validationErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{Status: fiber.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
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

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.UserResponse{Status: fiber.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(fiber.StatusCreated).JSON(responses.UserResponse{Status: fiber.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(responses.UserResponse{Status: fiber.StatusNotFound, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(fiber.StatusOK).JSON(responses.UserResponse{Status: fiber.StatusOK, Message: "success", Data: &fiber.Map{"data": user}})
}

func EditAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{Status: fiber.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{Status: fiber.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{"email": user.Email, "password": user.Password, "status": user.Status, "roles": user.Roles, "profile": user.Profile}

	result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.UserResponse{Status: fiber.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	//get updated user details
	var updatedUser models.User
	if result.MatchedCount == 1 {
		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(responses.UserResponse{Status: fiber.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(fiber.StatusOK).JSON(responses.UserResponse{Status: fiber.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedUser}})
}

func DeleteAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.UserResponse{Status: fiber.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.Status(fiber.StatusNotFound).JSON(
			responses.UserResponse{Status: fiber.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "User with specified ID not found!"}},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		responses.UserResponse{Status: fiber.StatusOK, Message: "success", Data: &fiber.Map{"data": "User successfully deleted!"}},
	)
}

func GetAllUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	results, err := userCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.UserResponse{Status: fiber.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.User
		if err = results.Decode(&singleUser); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(responses.UserResponse{Status: fiber.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		users = append(users, singleUser)
	}

	return c.Status(fiber.StatusOK).JSON(
		responses.UserResponse{Status: fiber.StatusOK, Message: "success", Data: &fiber.Map{"data": users}},
	)
}

func HelloWorld(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("Hello, World")
}
