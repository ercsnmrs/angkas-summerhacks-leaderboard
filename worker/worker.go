package worker

import (
	"context"
	"fmt"
	"log/slog"
)

type Mode uint

const (
	// ModeAll run all worker types.
	ModeAll Mode = iota
	// ModeConsumerOnly flags worker to only run and execute consumer handlers.
	ModeConsumerOnly
	// ModeSchedulerOnly flags worker to only run and execute scheduled jobs.
	ModeSchedulerOnly
)

const defaultJobQueueSize = 10

// Job represents a task details for a worker.
type Job struct {
	Topic   string
	Payload []byte
	Done    func() error
}

// JobHandler represents worker handler functions
type JobHandler func(context.Context, Job) error

// MiddlewareFunc represents worker job middleware
type MiddlewareFunc func(JobHandler) JobHandler

// Worker represents a worker that waits for a job and process base on handler.
type Worker struct {
	Mode Mode

	queue       chan Job
	quit        chan struct{}
	router      map[string]JobHandler
	middlewares []MiddlewareFunc
	listener    jobListener
	schedules   []Schedule
	logger      *slog.Logger
}

// jobListener provides access to job producers.
type jobListener interface {
	// Listen starts subscription to topics and listen for the jobs.
	Listen(topics []string, q chan<- Job) (stop func(), err error)

	// Close stops listening or unsubscribing to jobs.
	Close() error
}

// New create new instance of worker.
func New(listener jobListener, queueSize int, logger *slog.Logger) *Worker {
	if queueSize == 0 {
		queueSize = defaultJobQueueSize
	}

	logger = logger.With("pkg", "worker")
	logger.Info("init", "queue-size", queueSize)
	return &Worker{
		Mode:     ModeAll,
		queue:    make(chan Job, queueSize),
		quit:     make(chan struct{}, 1),
		router:   map[string]JobHandler{},
		listener: listener,
		logger:   logger,
	}
}

func (w *Worker) Run() error {
	switch w.Mode {
	case ModeConsumerOnly:
		return w.runConsumers()
	case ModeSchedulerOnly:
		return w.runSchedulers()
	default:
		if err := w.runConsumers(); err != nil {
			return err
		}
		return w.runSchedulers()
	}
}

// Stop gracefully stops listener and closes job queue.
func (w *Worker) Stop() error {
	w.logger.Info("stopping worker...")
	for _, s := range w.schedules {
		s.done <- struct{}{}
	}

	if err := w.listener.Close(); err != nil {
		return fmt.Errorf("listener close: %s", err)
	}

	close(w.queue)

	return nil
}

// HandleFunc registers handler by topic that routes job to a handler.
func (w *Worker) HandleFunc(topic string, f JobHandler) {
	w.router[topic] = f
}

func (w *Worker) SetSchedule(s Schedule) {
	s.logger = w.logger
	s.done = make(chan struct{}, 1)
	w.schedules = append(w.schedules, s)
}

// Use registers middlewares for job handlers.
func (w *Worker) Use(mm ...MiddlewareFunc) {
	for _, m := range mm {
		w.middlewares = append(w.middlewares, m)
	}
}

func (w *Worker) runConsumers() error {
	var tt []string
	for t := range w.router {
		w.logger.Info("register topic", "topic", t)
		tt = append(tt, t)
	}
	stop, err := w.listener.Listen(tt, w.queue)
	if err != nil {
		return err
	}

	// Process job received from the job listener.
	for x := 1; x <= cap(w.queue); x++ {
		workerID := x // prevent un-predictable values of x when used on go-routine function body.
		go func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			for {
				job, ok := <-w.queue
				if !ok {
					cancel()
					stop()
					w.logger.Info("worker stopped and job queue closed", "worker_id", workerID)
					return
				}
				w.logger.Info("worker received job", "worker_id", workerID)

				handle, ok := w.router[job.Topic]
				if !ok {
					w.logger.Debug("topic not handled", "topic", job.Topic, "worker_id", workerID)
					continue
				}

				// Execute middlewares on router job handler.
				for _, m := range w.middlewares {
					handle = m(handle)
				}

				if err = handle(ctx, job); err != nil {
					continue
				}
				if err = job.Done(); err != nil {
					w.logger.Error("job done", "err", err, "topic", job.Topic, "worker_id", workerID)
					continue
				}
			}
		}()
	}

	w.logger.Info("worker running", "queue_size", cap(w.queue))
	return nil
}

func (w *Worker) runSchedulers() error {
	// Start running scheduled tasks.
	for _, s := range w.schedules {
		go func(s Schedule) {
			if err := s.run(context.Background()); err != nil {
				w.logger.Error("schedule run", "err", err)
			}
		}(s)
	}

	return nil
}
