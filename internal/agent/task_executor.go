package agent

type TaskExecutor struct {
	taskQueue *TaskQueue
	sem       chan struct{}
}

func NewTaskExecutor(taskQueue *TaskQueue, maxConcurrentTasks int) *TaskExecutor {
	return &TaskExecutor{
		taskQueue: taskQueue,
		sem:       make(chan struct{}, maxConcurrentTasks),
	}
}

func (e *TaskExecutor) Start() {
	go func() {
		for taskContext := range e.taskQueue.Dequeue() {
			e.sem <- struct{}{}
			go func(tc *TaskExecutionContext) {
				defer func() { <-e.sem }()
				e.handleTask(tc)
			}(taskContext)
		}
	}()
}

func (e *TaskExecutor) handleTask(taskContext *TaskExecutionContext) {

}
