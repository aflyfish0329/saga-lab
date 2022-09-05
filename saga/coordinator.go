package saga

import (
	"os"

	"github.com/rs/zerolog"
)

func NewCoordinator(saga Saga) Coordinator {
	coordinator := Coordinator{
		saga:   saga,
		logger: zerolog.New(os.Stderr).With().Timestamp().Logger(),
	}

	return coordinator
}

type Coordinator struct {
	saga      Saga
	txErr     error
	logger    zerolog.Logger
	executeTx []ExecuteResult
	executeRx []ExecuteResult
}

type ExecuteResult struct {
	Name string
	Err  error
}

func (c *Coordinator) LogLevel(level int) {
	c.logger = c.logger.Level(zerolog.Level(level))
}

func (c *Coordinator) Run() Result {
	for index := range c.saga.steps {
		if c.txErr != nil {
			break
		}

		command := c.saga.steps[index].Tx
		err := command.Func()
		if err != nil {
			c.txErr = err
			c.logger.Error().Msgf("tx failed - saga: %s, step: %d, command: %s, error: %v", c.saga.Name, index, command.Name, err)
			c.rollback(index)
		} else {
			c.logger.Info().Msgf("tx succeed - saga: %s, step: %d, command: %s", c.saga.Name, index, command.Name)
		}

		c.executeTx = append(
			c.executeTx,
			ExecuteResult{
				Name: command.Name,
				Err:  err,
			},
		)
	}

	result := Result{
		TxErr: c.txErr,
	}

	return result
}

func (c *Coordinator) rollback(point int) {
	for index := point - 1; index >= 0; index-- {
		command := c.saga.steps[index].Rx
		err := command.Func()
		if err != nil {
			c.logger.Error().Msgf("rx failed - saga: %s, step: %d, command: %s, error: %v", c.saga.Name, index, command.Name, err)
		} else {
			c.logger.Info().Msgf("rx succeed - saga: %s, step: %d, command: %s", c.saga.Name, index, command.Name)
		}

		c.executeRx = append(
			c.executeRx,
			ExecuteResult{
				Name: command.Name,
				Err:  err,
			},
		)
	}
}

func (c Coordinator) GetExecuteTx() []ExecuteResult {
	return c.executeTx
}

func (c Coordinator) GetExecuteRx() []ExecuteResult {
	return c.executeRx
}
