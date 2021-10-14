package mocks

import "grpc-example/src/data/models"

type UserRepositoryMock struct {
	CreateResponse bool
	GetResponse    *models.User
	UpdateResponse bool
}

func (mock *UserRepositoryMock) Create(*models.User) bool {
	return mock.CreateResponse
}

func (mock *UserRepositoryMock) Update(*models.User) bool {
	return mock.UpdateResponse
}

func (mock *UserRepositoryMock) Get(*models.User) *models.User {
	return mock.GetResponse
}
