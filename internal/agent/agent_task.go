package agent

type TaskExecutionContext struct {
	Input interface{}
}

type Task interface {
	Execute(context TaskExecutionContext) (interface{}, error)
}
