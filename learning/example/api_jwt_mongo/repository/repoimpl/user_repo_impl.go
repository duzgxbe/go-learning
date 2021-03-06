package repoimpl

/**
NOTE: implement repo should only contain wrap to mongo, no business logic here

i.e: Insert() -> NO logic check exists or hashed password here; move it to handler
*/

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"testProject/pkg/utils"

	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	models "testProject/learning/example/api_jwt_mongo/model"
	repo "testProject/learning/example/api_jwt_mongo/repository"
)

const colName = "users"

type UserRepoImpl struct {
	Db *mongo.Database
}

func NewUserRepo(db *mongo.Database) repo.UserRepo {
	return &UserRepoImpl{
		Db: db,
	}
}

func (ins *UserRepoImpl) FindAll(queryData map[string]interface{}) ([]models.User, error) {
	var users []models.User

	query := bson.M{}
	for k, v := range queryData {
		query[k] = v
	}

	cur, err := ins.Db.Collection(colName).Find(context.Background(), query)
	if err != nil {
		return users, err
	}

	for cur.Next(context.Background()) {
		var user models.User
		err := cur.Decode(&user)
		if err != nil {
			return users, err
		}
		user.Password = "<hidden>"
		users = append(users, user)
	}

	if err := cur.Err(); err != nil {
		return users, err
	}

	return users, nil
}

func (ins *UserRepoImpl) FindOne(queryData map[string]interface{}) (models.User, error) {
	query := bson.M{}
	for k, v := range queryData {
		query[k] = v
	}
	result := ins.Db.Collection(colName).
		FindOne(context.Background(), query)
	// FindOne(context.Background(), bson.M{"email": u.Email})

	var user models.User
	err := result.Decode(&user)
	if err != nil {
		return user, err
	}

	// user.Password = "<hidden>"
	return user, nil
}

// TODO: refactor code like Insert(): move logic code to handler!!!

func (ins *UserRepoImpl) Insert(u *models.User) (string, error) {
	// encode data
	bbytes, err := bson.Marshal(u)
	if err != nil {
		return "", err
	}

	// insert db
	result, err := ins.Db.Collection(colName).InsertOne(context.Background(), bbytes)
	if err != nil {
		return "", err
	}

	_id := result.InsertedID.(string)
	fmt.Println("Inserted user ", _id, u.Email, result)
	return _id, nil
}

// check email + password
func (ins *UserRepoImpl) CheckLogin(email, password string) (models.User, error) {
	queryExists := map[string]interface{}{
		"email": email,
	}

	var emptyUser models.User
	user, err := ins.FindOne(queryExists)
	if err != nil {
		return emptyUser, err
	}

	// compare pwd
	valid := utils.CheckPwd(user.Password, password)
	if !valid {
		return emptyUser, errors.New("Invalid pwd")
	}

	user.Password = "<hidden>"
	return user, nil
}
