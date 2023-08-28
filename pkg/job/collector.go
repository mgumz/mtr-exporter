package job

import (
	"sync"
)

type Collector struct {
	jobs []JobMeta

	mu sync.Mutex

	opts collectorOpts
}

type collectorOpts struct {
	doRenderDeprecatedMetrics bool
}

func NewCollector() *Collector {
	return new(Collector)
}

func (c *Collector) SetRenderDeprecatedMetrics(f bool) *Collector {
	c.opts.doRenderDeprecatedMetrics = f
	return c
}

func (c *Collector) RemoveJob(label string) bool {

	c.mu.Lock()
	defer c.mu.Unlock()

	jobs := []JobMeta{}

	for i := range c.jobs {
		if c.jobs[i].Label != label {
			jobs = append(jobs, c.jobs[i])
		}
	}

	if len(jobs) < len(c.jobs) {
		c.jobs = jobs
		return true
	}
	return false
}

func (c *Collector) AddJob(job JobMeta) bool {

	c.mu.Lock()
	defer c.mu.Unlock()

	for i := range c.jobs {
		if job.Label == c.jobs[i].Label {
			return false
		}
	}
	c.jobs = append(c.jobs, job)

	return true
}

func (c *Collector) UpdateJob(job JobMeta) bool {

	c.mu.Lock()
	defer c.mu.Unlock()

	for i := range c.jobs {
		if c.jobs[i].Label == job.Label {
			c.jobs[i] = job
			return true
		}
	}

	return false
}
