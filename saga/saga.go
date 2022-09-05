package saga

func NewCommand(name string, f func() error) Command {
	return Command{
		Name: name,
		Func: f,
	}
}

type Command struct {
	Name string
	Func func() error
}

type Step struct {
	Tx Command
	Rx Command
}

func NewSaga(name string) Saga {
	return Saga{
		Name: name,
	}
}

type Saga struct {
	Name  string
	steps []Step
}

func (s *Saga) AddStep(tx, rx Command) {
	s.steps = append(s.steps, Step{
		Tx: tx,
		Rx: rx,
	})
}
