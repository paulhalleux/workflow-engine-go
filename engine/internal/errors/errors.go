package errors

type SimpleError string

func (e SimpleError) Error() string {
	return string(e)
}

const (
	ErrWorkflowDefinitionNoSteps SimpleError = "workflow definition has no steps"
)
