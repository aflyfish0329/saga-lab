package user

import (
	"errors"
	"saga/saga"

	uuid "github.com/satori/go.uuid"
)

type CreateUserSagaInput struct {
	UserName string
}

type CreateUserSagaOutput struct {
	UserId    string
	ExecuteTx []saga.ExecuteResult
	ExecuteRx []saga.ExecuteResult
}

func NewCreateUserSaga() CreateUserSaga {
	return CreateUserSaga{}
}

type CreateUserSaga struct {
}

func (cus CreateUserSaga) Run(input CreateUserSagaInput) (CreateUserSagaOutput, error) {
	type CreateUserSagaData struct {
		UserName string
		UserId   string
	}

	data := CreateUserSagaData{
		UserName: input.UserName,
	}

	s := saga.NewSaga("createUserSaga")
	s.AddStep(
		saga.NewCommand("create user in crm", func() error {
			data.UserId = uuid.NewV4().String()
			return nil
		}),
		saga.NewCommand("delete user crm", func() error {
			return errors.New("delete user in crm failed")
		}),
	)
	s.AddStep(
		saga.NewCommand("create user in acs", func() error {
			return nil
		}),
		saga.NewCommand("delete user in acs", func() error {
			return nil
		}),
	)
	s.AddStep(
		saga.NewCommand("create user in billing", func() error {
			if data.UserName == "hello" {
				return nil
			}

			return errors.New("create user in billing failed")
		}),
		saga.NewCommand("delete user in billing", func() error {
			return nil
		}),
	)
	s.AddStep(
		saga.NewCommand("nothing", func() error {
			return nil
		}),
		saga.NewCommand("nothing", func() error {
			return nil
		}),
	)

	c := saga.NewCoordinator(s)
	result := c.Run()

	return CreateUserSagaOutput{
		UserId:    data.UserId,
		ExecuteTx: c.GetExecuteTx(),
		ExecuteRx: c.GetExecuteRx(),
	}, result.TxErr
}
