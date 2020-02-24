package main

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"google.golang.org/appengine"
)

// Todo is todo model property
type Todo struct {
	gorm.Model
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

var db *gorm.DB
var err error

func main() {
	db, err = gorm.Open("mysql", os.Getenv("DATABASE_TOKEN"))
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	db.AutoMigrate(&Todo{})

	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/todos", GetTodos)
	router.GET("/todos/:id", GetTodo)
	router.POST("/todos", CreateTodo)
	router.PUT("/todos/:id", UpdateTodo)
	router.DELETE("/todos/:id", DeleteTodo)
	router.Run()
	appengine.Main()
}

// GetTodos is function to read all of the todos
func GetTodos(c *gin.Context) {
	var todo []Todo
	if err := db.Find(&todo).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, todo)
	}
}

//GetTodo is function to read a todo
func GetTodo(c *gin.Context) {
	id := c.Params.ByName("id")
	var todo Todo
	if err := db.Where("id = ?", id).First(&todo).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, todo)
	}
}

// CreateTodo is function to create a todo
func CreateTodo(c *gin.Context) {
	var todo Todo
	c.BindJSON(&todo)
	db.Create(&todo)
	c.JSON(200, todo)
}

// UpdateTodo is function to update a todo
func UpdateTodo(c *gin.Context) {
	var todo Todo
	id := c.Params.ByName("id")
	if err := db.Where("id = ?", id).First(&todo).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	c.BindJSON(&todo)
	db.Save(&todo)
	c.JSON(200, todo)
}

// DeleteTodo is function to delete a todo
func DeleteTodo(c *gin.Context) {
	var todo Todo
	id := c.Params.ByName("id")
	d := db.Where("id = ?", id).Delete(&todo)
	fmt.Println(d)
	c.JSON(200, gin.H{"id #" + id: "deleted!"})
}
