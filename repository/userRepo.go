package repository

import (
	"context"
	"log"

	"github.com/Yefhem/mongo/dictionary/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Insert(user models.User) error
	FindByKeyValue(key string, value interface{}) (models.User, error)
	CheckEmailPass(email, pass interface{}) (models.User, error)
	DeleteAllUser() error
}

type userConnect struct {
	ctx context.Context
	col *mongo.Collection
}

func NewUserRepository(ctx context.Context, col *mongo.Collection) UserRepository {
	return &userConnect{
		ctx: ctx,
		col: col,
	}
}

// --------> Methods
// --------> Insert to User
func (c *userConnect) Insert(user models.User) error {
	user.Password = PasswordHasher(user.Password)
	_, err := c.col.InsertOne(c.ctx, user)
	if err != nil {
		return err
	}

	return nil
}

// --------> Check Key and Value
func (c *userConnect) FindByKeyValue(key string, value interface{}) (models.User, error) {
	var user models.User

	filter := bson.D{primitive.E{Key: key, Value: value}}

	err := c.col.FindOne(c.ctx, filter).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

// --------> Check Email and Pass
func (c *userConnect) CheckEmailPass(email, pass interface{}) (models.User, error) {
	var user models.User

	filter := bson.D{primitive.E{Key: "email", Value: email}, {Key: "password", Value: pass}}

	err := c.col.FindOne(c.ctx, filter).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

// --------> Delete to User
func (c *userConnect) DeleteAllUser() error {
	filter := bson.D{{}}

	_, err := c.col.DeleteMany(c.ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

// --------> Hash to User Pass
func PasswordHasher(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}
