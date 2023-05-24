package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Job represents a scheduled task
type Job struct {
	ID       int       `json:"id"`
	Priority int       `json:"priority"`
	Message  string    `json:"message"`
	Params   string    `json:"params"`
	NextRun  time.Time `json:"next_run"`
}

var jobs []Job

func main() {
	// Initialize Gin router
	router := gin.Default()

	// Define CRUD endpoints for job scheduling
	router.POST("/jobs", createJob)
	router.GET("/jobs", getJobs)
	router.GET("/jobs/:id", getJobByID)
	router.PUT("/jobs/:id", updateJob)
	router.DELETE("/jobs/:id", deleteJob)

	// Start the scheduler in a separate goroutine
	go scheduler()

	// Run the server
	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func createJob(c *gin.Context) {
	var job Job
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job.ID = len(jobs) + 1
	jobs = append(jobs, job)

	c.JSON(http.StatusCreated, job)
}

func getJobs(c *gin.Context) {
	c.JSON(http.StatusOK, jobs)
}

func getJobByID(c *gin.Context) {
	id := c.Param("id")

	for _, job := range jobs {
		if fmt.Sprintf("%d", job.ID) == id {
			c.JSON(http.StatusOK, job)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
}

func updateJob(c *gin.Context) {
	id := c.Param("id")

	var updatedJob Job
	if err := c.ShouldBindJSON(&updatedJob); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, job := range jobs {
		if fmt.Sprintf("%d", job.ID) == id {
			jobs[i].Priority = updatedJob.Priority
			jobs[i].Message = updatedJob.Message
			jobs[i].Params = updatedJob.Params
			jobs[i].NextRun = updatedJob.NextRun

			c.JSON(http.StatusOK, jobs[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
}

func deleteJob(c *gin.Context) {
	id := c.Param("id")

	for i, job := range jobs {
		if fmt.Sprintf("%d", job.ID) == id {
			jobs = append(jobs[:i], jobs[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Job deleted"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
}

func scheduler() {
	for {
		// Check if there are any jobs
		if len(jobs) == 0 {
			time.Sleep(time.Second)
			continue
		}

		// Get the job with the highest priority
		highestPriorityJob := jobs[0]
		for _, job := range jobs {
			if job.Priority > highestPriorityJob.Priority {
				highestPriorityJob = job
			}
		}

		// Randomize the execution time within a range
		min := 5  // minimum delay in seconds
		max := 15 // maximum delay in seconds
		delay := rand.Intn(max-min+1) + min
		nextRun := time.Now().Add(time.Duration(delay) * time.Second)

		// Execute the job (Python script) here
		fmt.Printf("Executing job with ID %d at %s\n", highestPriorityJob.ID, nextRun.String())
		// Call your Python script here using appropriate libraries or commands

		// Update the job's next run time
		for i, job := range jobs {
			if job.ID == highestPriorityJob.ID {
				jobs[i].NextRun = nextRun
				break
			}
		}

		time.Sleep(time.Second)
	}
}
