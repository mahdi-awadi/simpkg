package workman

import (
	"errors"
	"path"
	"path/filepath"
	"sort"
	"sync"

	"github.com/go-per/simpkg/cache"
	"github.com/go-per/simpkg/client"
	"github.com/go-per/simpkg/events"
	"github.com/go-per/simpkg/helpers"
	"github.com/go-per/simpkg/logger"
	"github.com/go-per/simpkg/str"
	"github.com/go-per/simpkg/tasks"
)

// WorkerEvent type
type WorkerEvent string

const (
	EventWorkerAddRemove WorkerEvent = "worker.add_remove"
	EventWorkerLoad      WorkerEvent = "workers.load"
)

// IManager interface
type IManager interface {
	SetDebug(isDebug bool)
	IsDebug() bool
	SetExtension(ext string)
	GetExtension() string
	SetRootPath(path string)
	SetWorkersDir(path string)
	GetWorkersPath() string
	GetWorkerFilePath(id string) string
	RootPath() string
	CachePath() string
	Add(index int, config []byte, fp ...string) (IWorker, error)
	Remove(id string) bool
	Get(id string) (IWorker, bool)
	Next() IWorker
	SelectedWorker() IWorker
	Workers() []IWorker
	WorkersCount() int
	Load() error
	Reload() error
	Initialize(any) error
	Eventbus() events.IEventbus
}

// Manager struct
type Manager struct {
	isDebug       bool
	rootPath      string
	workersDir    string
	workersExt    string
	workerBuilder WorkerBuilderFunc
	workers       []IWorker
	eventbus      events.IEventbus
	locker        sync.RWMutex

	selectWorkerLocker sync.RWMutex
	selectedWorker     IWorker
}

// WorkerBuilderFunc is worker builder function
type WorkerBuilderFunc func() IWorker

// WorkerOnAddRemove struct
type WorkerOnAddRemove struct {
	Added   bool
	Removed bool
	Worker  IWorker
}

// WorkerOnUpdate struct
type WorkerOnUpdate struct {
	Workers []IWorker
}

// NewManager returns new manager instance
func NewManager(builder WorkerBuilderFunc) *Manager {
	return &Manager{
		workerBuilder:      builder,
		workers:            make([]IWorker, 0),
		eventbus:           events.New(),
		workersDir:         "workers",
		workersExt:         ".json",
		locker:             sync.RWMutex{},
		selectedWorker:     nil,
		selectWorkerLocker: sync.RWMutex{},
	}
}

// Initialize initializes manager
func (m *Manager) Initialize(any) error { return nil }

// SetDebug sets debug mode
func (m *Manager) SetDebug(isDebug bool) {
	m.isDebug = isDebug
}

// IsDebug returns debug mode
func (m *Manager) IsDebug() bool {
	return m.isDebug
}

// SetWorkersDir sets workers directory
func (m *Manager) SetWorkersDir(path string) {
	m.workersDir = path
}

// SetExtension sets workers extension
func (m *Manager) SetExtension(ext string) {
	m.workersExt = "." + str.Strip(ext)
}

// GetExtension returns workers extension
func (m *Manager) GetExtension() string {
	return m.workersExt
}

// SetRootPath sets root path
func (m *Manager) SetRootPath(path string) {
	m.rootPath = path
}

// RootPath returns root path
func (m *Manager) RootPath() string {
	return m.rootPath
}

// CachePath returns cache path
func (m *Manager) CachePath() string {
	return filepath.Join(m.RootPath(), "cache")
}

// Add adds worker to manager
func (m *Manager) Add(index int, data []byte, fileName ...string) (IWorker, error) {
	if m.workerBuilder == nil {
		return nil, errors.New("worker type is not set")
	}

	// create new worker
	worker := m.workerBuilder()

	// parse worker config
	id, err := worker.Init(data, fileName...)
	if err != nil {
		return nil, err
	}

	// if already exists
	if _, ok := m.Get(id); ok {
		return nil, errors.New("worker already exists")
	}

	// set worker props
	worker.SetIndex(index)
	worker.SetTaskManager(tasks.New())
	worker.SetClient(client.New())

	// configure and set cache
	cachePath := path.Join(m.CachePath(), "_workers", worker.GetID())
	fileCache := cache.New()
	fileCache.SetRoot(cachePath)
	worker.SetCache(fileCache)

	// configure and set logger
	l := logger.New()
	l.SetRootPath(cachePath)
	worker.SetLogger(l)

	// register listeners
	worker.RegisterListeners()

	// initialize worker
	err = worker.Boot()
	if err != nil {
		return nil, err
	}

	// add worker to manager
	m.locker.Lock()
	m.workers = append(m.workers, worker)
	m.locker.Unlock()

	// dispatch event
	m.eventbus.Dispatch(string(EventWorkerAddRemove), WorkerOnAddRemove{
		Added:  true,
		Worker: worker,
	})

	return worker, nil
}

// Remove removes worker from manager
func (m *Manager) Remove(id string) bool {
	m.locker.Lock()
	removed := false

	var wk IWorker
	for i, worker := range m.workers {
		if worker.GetID() == id {
			worker.Stop()
			m.workers = append(m.workers[:i], m.workers[i+1:]...)
			removed = true
			wk = worker
			break
		}
	}
	m.locker.Unlock()

	if !removed {
		return true
	}

	// dispatch event
	m.eventbus.Dispatch(string(EventWorkerAddRemove), WorkerOnAddRemove{
		Removed: true,
		Worker:  wk,
	})

	return true
}

// Get returns worker by id
func (m *Manager) Get(id string) (IWorker, bool) {
	var worker IWorker
	m.locker.Lock()
	for _, wk := range m.workers {
		if wk.GetID() == id {
			worker = wk
			break
		}
	}
	m.locker.Unlock()

	return worker, worker != nil
}

// Next select workers as circular selection
func (m *Manager) Next() IWorker {
	//m.selectWorkerLocker.Lock()
	//if m.selectedWorker != nil && m.selectedWorker.Index() >= m.WorkersCount() {
	//	m.selectedWorker = nil
	//}
	//
	//for _, worker := range m.workers {
	//	if m.selectedWorker == nil || worker.Index() > m.selectedWorker.Index() {
	//		m.selectedWorker = worker
	//		break
	//	}
	//}
	//
	//if m.selectedWorker != nil {
	//	return m.selectedWorker
	//}
	//
	//if m.WorkersCount() > 0 {
	//	m.selectedWorker = m.workers[0]
	//}
	//m.selectWorkerLocker.Unlock()

	return m.selectedWorker
}

// SelectedWorker returns selected worker
func (m *Manager) SelectedWorker() IWorker {
	return m.selectedWorker
}

// Workers returns all workers
func (m *Manager) Workers() []IWorker {
	return m.workers
}

// WorkersCount returns workers count
func (m *Manager) WorkersCount() int {
	return len(m.workers)
}

// Load loads workers from config files in directory
func (m *Manager) Load() error {
	workersPath := filepath.Join(m.rootPath, m.workersDir, "*"+m.workersExt)
	files, err := filepath.Glob(workersPath)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return errors.New("no worker file found at " + m.RootPath())
	}

	for i, file := range files {
		content, err := helpers.ReadFile(file)
		if err != nil {
			return err
		}

		if err != nil {
			return err
		}

		_, err = m.Add(i, content, file)
		if err != nil {
			return err
		}
	}

	// check for workers count
	if m.workers == nil || len(m.workers) == 0 {
		return errors.New("no worker found")
	}

	// sort workers
	sort.Slice(m.workers, func(i, j int) bool { return m.workers[i].Index() < m.workers[j].Index() })

	// dispatch event
	m.eventbus.Dispatch(string(EventWorkerLoad), WorkerOnUpdate{
		Workers: m.workers,
	})

	return nil
}

// GetWorkersPath returns workers path
func (m *Manager) GetWorkersPath() string {
	return filepath.Join(m.rootPath, m.workersDir)
}

// GetWorkerFilePath returns worker file path
func (m *Manager) GetWorkerFilePath(id string) string {
	return filepath.Join(m.rootPath, m.workersDir, id+m.workersExt)
}

// Reload workers
func (m *Manager) Reload() error {
	return m.Load()
}

// Eventbus returns eventbus instance
func (m *Manager) Eventbus() events.IEventbus {
	return m.eventbus
}
