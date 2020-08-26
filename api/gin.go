//Download and install it:
//$ go get github.com/gin-gonic/gin

package main

import (
	"encoding/json"
	"fmt"
	"gihub.com/cronjob-scheduler/cron"
	"github.com/gin-gonic/gin"
)

var DB = make(map[string]string)

type Job struct {
	Id int
}

func main() {
	cr := cron.New(cron.WithParser(cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)))
	cr.Start()
	cr.AddFunc("* * * * * ?", func() { fmt.Println("job1") })

	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// Get job list
	r.GET("/jobs", func(c *gin.Context) {
		entries := cr.Entries()
		var jobs []Job
		for _, entry := range entries {
			fmt.Println(int(entry.ID))
			jobs = append(jobs, Job{
				Id: int(entry.ID),
			})
		}
		result, _ := json.Marshal(jobs)
		c.String(200, string(result))
	})

	// Remove job
	r.GET("/remove", func(c *gin.Context) {
		cr.Remove(1)
		c.String(200, "1")
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := DB[user]
		if ok {
			c.JSON(200, gin.H{"user": user, "value": value})
		} else {
			c.JSON(200, gin.H{"user": user, "status": "no value"})
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			DB[user] = json.Value
			c.JSON(200, gin.H{"status": "ok"})
		}
	})

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
