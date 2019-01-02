package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/buntdb"
)

// Task is the basic struct for handling tasks
type Task struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	IsCompleted bool   `json:"completed"`
}

// Tasks handles lists of tasks
type Tasks []Task

func main() {
	// Small DB to test the app
	db, err := buntdb.Open(":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initial data to play with
	firstTask := Task{ID: "1", Description: "Buy Stuff", IsCompleted: false}
	secondTask := Task{ID: "2", Description: "Make Stuff", IsCompleted: true}
	thirdTask := Task{ID: "3", Description: "Sell Stuff", IsCompleted: false}

	err = db.Update(func(tx *buntdb.Tx) error {
		task, err := json.Marshal(firstTask)
		if err != nil {
			fmt.Println("Marshalling error", err)
		}

		_, _, _ = tx.Set("1", string(task), nil)

		task, err = json.Marshal(secondTask)
		if err != nil {
			fmt.Println("Marshalling error", err)
		}

		_, _, _ = tx.Set("2", string(task), nil)

		task, err = json.Marshal(thirdTask)
		if err != nil {
			fmt.Println("Marshalling error", err)
		}

		_, _, _ = tx.Set("3", string(task), nil)
		return err
	})

	r := gin.Default()

	// Ping
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
		return
	})

	// CRUD
	r.GET("/api/tasks", func(c *gin.Context) {
		var response []*Task

		err := db.View(func(tx *buntdb.Tx) error {
			err := tx.Ascend("", func(key, value string) bool {
				task := &Task{}
				err = json.Unmarshal([]byte(value), task)

				response = append(response, task)
				return true
			})
			return err
		})

		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, response)
		return

	})

	r.PUT("/api/tasks", func(c *gin.Context) {
		var payload Task
		err := c.BindJSON(&payload)

		task, _ := json.Marshal(payload)

		err = db.Update(func(tx *buntdb.Tx) error {
			_, _, err = tx.Set(payload.ID, string(task), nil)
			return err
		})

		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, payload)
		return
	})

	r.GET("/api/tasks/:id", func(c *gin.Context) {
		ID := c.Param("id")

		response := &Task{}

		err = db.View(func(tx *buntdb.Tx) error {
			val, _ := tx.Get(ID)
			err = json.Unmarshal([]byte(val), response)
			return err
		})

		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, response)
		return
	})

	r.POST("/api/tasks/:id", func(c *gin.Context) {
		ID := c.Param("id")
		var payload Task
		err := c.BindJSON(&payload)

		task, _ := json.Marshal(payload)

		if ID != payload.ID {
			c.JSON(http.StatusBadRequest, "ID and Payload must match")
			return
		}

		err = db.Update(func(tx *buntdb.Tx) error {
			_, _, err = tx.Set(ID, string(task), nil)
			return err
		})

		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, payload)
		return
	})

	r.DELETE("/api/tasks/:id", func(c *gin.Context) {
		ID := c.Param("id")

		err = db.Update(func(tx *buntdb.Tx) error {
			if _, err = tx.Delete(ID); err != nil {
				return err
			}
			return err
		})

		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, "")
		return
	})

	// Run
	r.Run(":8080")
}
