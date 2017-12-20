package timer

import (
	"encoding/json"
	"fmt"
	"gopkg.in/robfig/cron.v2"
	"log"
	"strings"
)

const Bucket = "timers"

type Job struct {
	ID        string `json:"id"`
	Minute    string `json:"minute"`
	Day       string `json:"day"`
	Hour      string `json:"hour"`
	Second    string `json:"second"`
	Equipment string `json:"equipment"`
	On        bool   `json:"on"`
	Name      string `json:"name"`
}

func (j *Job) CronSpec() string {
	return strings.Join([]string{j.Second, j.Minute, j.Hour, j.Day, "*", "?"}, " ")
}

func (j *Job) Validate() error {
	_, err := cron.Parse(j.CronSpec())
	return err
}

func (c *Controller) Get(id string) (Job, error) {
	var job Job
	return job, c.store.Get(Bucket, id, &job)
}

func (c *Controller) List() ([]Job, error) {
	jobs := []Job{}
	fn := func(v []byte) error {
		var job Job
		if err := json.Unmarshal(v, &job); err != nil {
			return err
		}
		jobs = append(jobs, job)
		return nil
	}
	return jobs, c.store.List(Bucket, fn)
}

func (c *Controller) Create(job Job) error {
	if err := job.Validate(); err != nil {
		return fmt.Errorf("Ivalid cronspec: %s", job.CronSpec())
	}
	if job.Equipment == "" {
		return fmt.Errorf("Missing equipment")
	}
	_, err := c.equipment.Get(job.Equipment)
	if err != nil {
		return err
	}
	fn := func(id string) interface{} {
		job.ID = id
		return job
	}
	if err := c.store.Create(Bucket, fn); err != nil {
		return err
	}
	return c.addToCron(job)
}

func (c *Controller) Update(id string, payload Job) error {
	payload.ID = id
	return c.store.Update(Bucket, id, &payload)
}

func (c *Controller) Delete(id string) error {
	if err := c.store.Delete(Bucket, id); err != nil {
		return err
	}
	return c.deleteFromCron(id)
}

func (c *Controller) loadAllJobs() error {
	jobs, err := c.List()
	if err != nil {
		return err
	}
	if jobs == nil {
		log.Printf("No jobs present")
		return nil
	}
	for _, job := range jobs {
		if err := c.addToCron(job); err != nil {
			log.Println("ERROR: Failed to add job in cron runner. Error:", err)
		}
	}
	return nil
}

func (c *Controller) addToCron(job Job) error {
	runner, err := c.Runner(job.Equipment, job.On)
	if err != nil {
		return err
	}
	cronID, err := c.runner.AddJob(job.CronSpec(), runner)
	if err != nil {
		return err
	}
	log.Println("Successfully added cron entry. ID:", cronID)
	c.cronIDs[job.ID] = cronID
	return nil
}

func (c *Controller) deleteFromCron(jobID string) error {
	id, ok := c.cronIDs[jobID]
	if !ok {
		return fmt.Errorf("Cron ID not found for job ID:%s", jobID)
	}
	if c.runner != nil {
		c.runner.Remove(id)
	}
	return nil
}
