package asynchook

import (
	"io"
	"os"
	"sync"

	"github.com/LyricTian/queue"
	"github.com/sirupsen/logrus"
)

var defaultOptions = options{
	maxQueues:  512,
	maxWorkers: 2,
	levels: []logrus.Level{
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	},
	out: os.Stderr,
}

// ExecCloser write the logrus entry to the store and close the store
type ExecCloser interface {
	Exec(entry *logrus.Entry) error
	Close() error
}

// FilterHandle a filter handler
type FilterHandle func(*logrus.Entry) *logrus.Entry

type options struct {
	maxQueues  int
	maxWorkers int
	extra      map[string]interface{}
	exec       ExecCloser
	filter     FilterHandle
	levels     []logrus.Level
	out        io.Writer
}

// SetMaxQueues set the number of buffers
func SetMaxQueues(maxQueues int) Option {
	return func(o *options) {
		o.maxQueues = maxQueues
	}
}

// SetMaxWorkers set the number of worker threads
func SetMaxWorkers(maxWorkers int) Option {
	return func(o *options) {
		o.maxWorkers = maxWorkers
	}
}

// SetExtra set extended parameters
func SetExtra(extra map[string]interface{}) Option {
	return func(o *options) {
		o.extra = extra
	}
}

// SetExec set the Execer interface
func SetExec(exec ExecCloser) Option {
	return func(o *options) {
		o.exec = exec
	}
}

// SetFilter set the entry filter
func SetFilter(filter FilterHandle) Option {
	return func(o *options) {
		o.filter = filter
	}
}

// SetLevels set the available log level
func SetLevels(levels ...logrus.Level) Option {
	return func(o *options) {
		if len(levels) == 0 {
			return
		}
		o.levels = levels
	}
}

// SetOut set error output
func SetOut(out io.Writer) Option {
	return func(o *options) {
		o.out = out
	}
}

// Option a hook parameter options
type Option func(*options)

// New creates a hook to be added to an instance of logger
func New(opt ...Option) *Hook {
	opts := defaultOptions
	for _, o := range opt {
		o(&opts)
	}

	if opts.exec == nil {
		panic("Unknown Execer interface implementation")
	}

	q := queue.NewQueue(opts.maxQueues, opts.maxWorkers)
	q.Run()

	return &Hook{
		opts: opts,
		q:    q,
		pool: sync.Pool{
			New: func() interface{} {
				return &job{
					out:    opts.out,
					extra:  opts.extra,
					exec:   opts.exec,
					filter: opts.filter,
				}
			},
		},
	}
}

// Hook to send logs to a mongo database
type Hook struct {
	opts options
	q    *queue.Queue
	pool sync.Pool
}

// Levels returns the available logging levels
func (h *Hook) Levels() []logrus.Level {
	return h.opts.levels
}

// Fire is called when a log event is fired
func (h *Hook) Fire(entry *logrus.Entry) error {
	job := &*(h.pool.Get().(*job))
	job.Reset(entry)
	h.q.Push(job)
	h.pool.Put(job)
	return nil
}

// Flush waits for the log queue to be empty
func (h *Hook) Flush() {
	h.q.Terminate()
	h.opts.exec.Close()
}
