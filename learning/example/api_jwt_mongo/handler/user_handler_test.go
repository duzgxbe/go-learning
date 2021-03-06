package handler

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strings"
	driverMongo "testProject/learning/example/api_jwt_mongo/driver/mongo"
	mocks2 "testProject/learning/example/api_jwt_mongo/driver/mongo/mocks"
	models "testProject/learning/example/api_jwt_mongo/model"
	"testProject/learning/example/api_jwt_mongo/repository/mocks"
	"testing"
)

// https://gowalker.org/github.com/stretchr/testify/mock#Call_On

func TestAddUser(t *testing.T) {
	userRepoMock := new(mocks.UserRepo)

	dataUsers := []models.User{
		{
			ID:       "0",
			Name:     "admin",
			Email:    "admin@gmail.com",
			Role:     "admin",
			Password: "1",
		},
		{
			ID:       "1",
			Name:     "user1",
			Email:    "user1@gmail.com",
			Role:     "user",
			Password: "1",
		},
	}

	// mock func return value each time called
	// i.e: err "user 'user1' existed"
	userRepoMock.
		On("FindOne", mock.Anything).Return(dataUsers[0], errors.New("not found")).Once(). // not found admin
		On("Insert", mock.Anything).Return("", nil).                                       // insert admin ok
		On("FindOne", mock.Anything).Return(dataUsers[1], nil)                             // found user1 -> err

	err := AddUser(userRepoMock)
	assert.Error(t, err)
	t.Log(err.Error())
	if !strings.Contains(err.Error(), "user1") {
		t.Fail()
	}
}

func TestDemo(t *testing.T) {
	userRepoMock := new(mocks.UserRepo)

	// orgUserRepo := extensions.GetInstance().UserRepo
	// extensions.GetInstance().UserRepo = userRepoMock
	// defer func() {
	// 	extensions.GetInstance().UserRepo = orgUserRepo
	// }()
	mongoMock := new(mocks2.Storage)
	mongoMock.
		On("GetInstance").
		On("GetUserRI").Return(userRepoMock)

	orgMongo := driverMongo.GetInstance()
	driverMongo.SetInstance(mongoMock)
	defer func() {
		driverMongo.SetInstance(orgMongo)
	}()

	dataUsers := []models.User{
		{
			ID:       "0",
			Name:     "admin",
			Email:    "admin@gmail.com",
			Role:     "admin",
			Password: "1",
		},
		{
			ID:       "1",
			Name:     "user1",
			Email:    "user1@gmail.com",
			Role:     "user",
			Password: "1",
		},
	}

	// mock func return value each time called
	// i.e: err "user 'user1' existed"
	userRepoMock.
		On("FindOne", mock.Anything).Return(dataUsers[0], errors.New("not found")).Once(). // not found admin
		On("Insert", mock.Anything).Return("", nil).                                       // insert admin ok
		On("FindOne", mock.Anything).Return(dataUsers[1], nil)                             // found user1 -> err

	err := Demo()
	assert.Error(t, err)
	t.Log(err.Error())
	if !strings.Contains(err.Error(), "user1") {
		t.Fail()
	}
}
