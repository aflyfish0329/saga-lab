package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateUser_Succeed(t *testing.T) {
	input := CreateUserSagaInput{
		UserName: "hello",
	}

	manager := NewCreateUserSaga()
	output, err := manager.Run(input)
	assert.NoError(t, err)
	assert.NotEmpty(t, output.UserId)
	assert.Equal(t, 4, len(output.ExecuteTx))
	assert.Equal(t, 0, len(output.ExecuteRx))
}

func Test_CreateUser_Failed(t *testing.T) {
	input := CreateUserSagaInput{
		UserName: "user1",
	}

	manager := NewCreateUserSaga()
	output, err := manager.Run(input)
	assert.ErrorContains(t, err, "create user in billing failed")
	assert.Equal(t, 3, len(output.ExecuteTx))
	assert.Equal(t, 2, len(output.ExecuteRx))
}
