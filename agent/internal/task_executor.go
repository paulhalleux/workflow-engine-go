package internal

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/kaptinlin/jsonschema"
	tjs "github.com/swaggest/jsonschema-go"
)

type TaskExecution struct {
	TaskID      string
	TaskDefName string
	Input       map[string]interface{}
}

type TaskExecutionResult struct {
	Output *map[string]interface{}
	Error  *error
}

type TaskExecutionRequest struct {
	Input map[string]interface{}
}

type TaskExecutor struct {
	taskDefinitionRegistry *TaskDefinitionRegistry
	taskQueue              chan *TaskExecution
	sem                    chan struct{}
}

func NewTaskExecutor(
	taskDefinitionRegistry *TaskDefinitionRegistry,
) *TaskExecutor {
	return &TaskExecutor{
		taskDefinitionRegistry: taskDefinitionRegistry,
		taskQueue:              make(chan *TaskExecution, 100),
		sem:                    make(chan struct{}, 1),
	}
}

func (te *TaskExecutor) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case exec := <-te.taskQueue:
				te.sem <- struct{}{}
				taskDef, found := te.taskDefinitionRegistry.Get(exec.TaskDefName)
				if !found {
					println("Task definition not found for task:", exec.TaskID)
				}

				err := te.validateParameters(exec.Input, taskDef.InputParameters)
				if err != nil {
					println("Invalid input parameters for task:", exec.TaskID, "Error:", err.Error())
				}

				if err == nil {
					go func(exec *TaskExecution) {
						defer func() { <-te.sem }()
						te.handle(exec, taskDef)
					}(exec)
				} else {
					<-te.sem
				}
			}
		}
	}()
}

func (te *TaskExecutor) handle(execCtx *TaskExecution, taskDef TaskDefinition) {
	req := &TaskExecutionRequest{
		Input: execCtx.Input,
	}

	result := taskDef.Handle(req)
	if result.Error != nil {
		println("Error executing task:", execCtx.TaskID, "Error:", (*result.Error).Error())
	} else {
		log.Printf("Task %s executed successfully with output: %v", execCtx.TaskID, *result.Output)
	}
}

func (te *TaskExecutor) validateParameters(params map[string]interface{}, schema *tjs.Schema) error {
	if schema == nil {
		return nil
	}

	compiler := jsonschema.NewCompiler()
	schemaBytes, err := json.Marshal(schema)
	if err != nil {
		return err
	}

	compiledSchema, err := compiler.Compile(schemaBytes)
	result := compiledSchema.Validate(params)
	if !result.Valid {
		return errors.New(result.Error())
	}

	return nil
}

func (te *TaskExecutor) EnqueueTask(exec *TaskExecution) {
	log.Printf("Enqueuing task %s of type %s", exec.TaskID, exec.TaskDefName)
	te.taskQueue <- exec
}
