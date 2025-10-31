package internal

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/kaptinlin/jsonschema"
	"github.com/paulhalleux/workflow-engine-go/proto"
	tjs "github.com/swaggest/jsonschema-go"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
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
	engineConnection       *grpc.ClientConn
	taskDefinitionRegistry *TaskDefinitionRegistry
	taskQueue              chan *TaskExecution
	sem                    chan struct{}
}

func NewTaskExecutor(
	taskDefinitionRegistry *TaskDefinitionRegistry,
	engineConnection *grpc.ClientConn,
) *TaskExecutor {
	return &TaskExecutor{
		engineConnection:       engineConnection,
		taskDefinitionRegistry: taskDefinitionRegistry,
		taskQueue:              make(chan *TaskExecution, 100),
		sem:                    make(chan struct{}, 1),
	}
}

func (te *TaskExecutor) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case exec := <-te.taskQueue:
			te.sem <- struct{}{}
			te.startTask(exec)

			var err error
			taskDef, found := te.taskDefinitionRegistry.Get(exec.TaskDefName)
			if !found {
				err = errors.New("task definition not found")
				println("Task definition not found for task:", exec.TaskID)
			}

			err = te.validateParameters(exec.Input, taskDef.InputParameters)
			if err != nil {
				println("Invalid input parameters for task:", exec.TaskID, "Error:", err.Error())
			}

			if err == nil {
				go func(exec *TaskExecution) {
					defer func() { <-te.sem }()
					te.handle(exec, taskDef)
				}(exec)
			} else {
				te.failTask(exec, err)
				<-te.sem
			}
		}
	}
}

func (te *TaskExecutor) handle(execCtx *TaskExecution, taskDef TaskDefinition) {
	req := &TaskExecutionRequest{
		Input: execCtx.Input,
	}

	result := taskDef.Handle(req)
	if result.Error != nil {
		println("Error executing task:", execCtx.TaskID, "Error:", (*result.Error).Error())
		te.failTask(execCtx, *result.Error)
	} else {
		log.Printf("Task %s executed successfully with output: %v", execCtx.TaskID, *result.Output)
		te.completeTask(execCtx, *result.Output)
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

func (te *TaskExecutor) startTask(exec *TaskExecution) {
	log.Printf("Task %s started", exec.TaskID)
	client := proto.NewTaskServiceClient(te.engineConnection)
	_, _ = client.NotifyTaskStatus(
		context.Background(),
		&proto.NotifyTaskStatusRequest{
			TaskId: exec.TaskID,
			Status: proto.TaskStatus_RUNNING,
		},
	)
}

func (te *TaskExecutor) failTask(exec *TaskExecution, err error) {
	log.Printf("Task %s failed with error: %s", exec.TaskID, err.Error())
	client := proto.NewTaskServiceClient(te.engineConnection)
	_, _ = client.NotifyTaskStatus(
		context.Background(),
		&proto.NotifyTaskStatusRequest{
			TaskId:  exec.TaskID,
			Status:  proto.TaskStatus_FAILED,
			Message: err.Error(),
		},
	)
}

func (te *TaskExecutor) completeTask(exec *TaskExecution, output map[string]interface{}) {
	log.Printf("Task %s completed successfully", exec.TaskID)

	outputStruct, err := structpb.NewStruct(output)
	if err != nil {
		log.Printf("Failed to convert output to structpb for task %s: %s", exec.TaskID, err.Error())
		te.failTask(exec, err)
		return
	}

	client := proto.NewTaskServiceClient(te.engineConnection)
	_, _ = client.NotifyTaskStatus(
		context.Background(),
		&proto.NotifyTaskStatusRequest{
			TaskId:           exec.TaskID,
			Status:           proto.TaskStatus_COMPLETED,
			OutputParameters: outputStruct,
		},
	)
}

func (te *TaskExecutor) EnqueueTask(exec *TaskExecution) {
	log.Printf("Enqueuing task %s of type %s", exec.TaskID, exec.TaskDefName)
	te.taskQueue <- exec
}
