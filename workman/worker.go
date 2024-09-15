package workman

import (
	"context"
	"github.com/go-per/simpkg/logger"

	"github.com/go-per/simpkg/cache"
	"github.com/go-per/simpkg/random"
	"github.com/go-per/simpkg/tasks"
	"github.com/imroc/req/v3"
)

// IWorker interface
type IWorker interface {
	Init([]byte, ...string) (string, error)
	Boot() error
	RegisterListeners()
	GetID() string
	GetDetails() any
	SetIndex(i int)
	Index() int
	SetCache(cache.ICache)
	Cache() cache.ICache
	SetTaskManager(tasks.IManager)
	TaskManager() tasks.IManager
	SetClient(*req.Client)
	Client() *req.Client
	SetLogger(logger.ILogger)
	Logger() logger.ILogger
	Start()
	Stop()
}

// Worker struct
type Worker struct {
	index       int
	cache       cache.ICache
	taskManager tasks.IManager
	client      *req.Client
	ctx         context.Context
	logger      logger.ILogger
	filePath    string
}

// Init initialize worker
func (w *Worker) Init([]byte, ...string) (string, error) { return "", nil }

// RegisterListeners trigger before boot
func (w *Worker) RegisterListeners() {}

// Boot method
func (w *Worker) Boot() error { return nil }

// SetIndex sets worker index
func (w *Worker) SetIndex(i int) {
	w.index = i
}

// Index returns worker index
func (w *Worker) Index() int {
	return w.index
}

// GetID returns worker id
func (w *Worker) GetID() string {
	return random.String(8)
}

// GetDetails returns worker details
func (w *Worker) GetDetails() any {
	return nil
}

// SetCache sets cache instance
func (w *Worker) SetCache(ic cache.ICache) {
	w.cache = ic
}

// Cache returns cache instance
func (w *Worker) Cache() cache.ICache {
	return w.cache
}

// SetLogger sets logger instance
func (w *Worker) SetLogger(l logger.ILogger) {
	w.logger = l
}

// Logger returns logger instance
func (w *Worker) Logger() logger.ILogger {
	return w.logger
}

// SetTaskManager sets task manager
func (w *Worker) SetTaskManager(tm tasks.IManager) {
	w.taskManager = tm
}

// TaskManager returns task manager
func (w *Worker) TaskManager() tasks.IManager {
	return w.taskManager
}

// SetClient set worker client
func (w *Worker) SetClient(client *req.Client) {
	w.client = client
}

// Client returns client instance
func (w *Worker) Client() *req.Client {
	return w.client
}

// Start worker
func (w *Worker) Start() {}

// Stop worker
func (w *Worker) Stop() {}
