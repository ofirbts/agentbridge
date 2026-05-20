package observability

import "time"

type Recorder struct {
	runID  string
	tracer *Tracer
	start  time.Time
}

func (r *Recorder) initStart() {
	if r.start.IsZero() {
		r.start = time.Now()
	}
}

func (r *Recorder) AddStep(step, message string) {
	r.initStart()
	r.tracer.logs = append(r.tracer.logs, StepLog{
		RunID:      r.runID,
		Step:       step,
		Status:     "ok",
		Retries:    0,
		DurationMS: time.Since(r.start).Milliseconds(),
		Message:    message,
	})
}

func (r *Recorder) AddStepWithMeta(step, status, provider, query string, retries int, durationMS int64) {
	r.initStart()
	r.tracer.logs = append(r.tracer.logs, StepLog{
		RunID:      r.runID,
		Step:       step,
		Status:     status,
		Retries:    retries,
		DurationMS: durationMS,
		Provider:   provider,
		Query:      query,
	})
}
