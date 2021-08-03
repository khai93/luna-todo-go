package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type List struct {
	Id    int              `json:"id"`
	Items [][2]interface{} `json:"items"`
}

func (list *List) AddItem(item [2]interface{}) [][2]interface{} {
	list.Items = append(list.Items, item)
	return list.Items
}

func (list *List) RemoveItem(itemId int) {
	itemId -= 1
	list.Items = append(list.Items[:itemId], list.Items[itemId+1:]...)
}

type CreateListItemRequestBody struct {
	Todo_Desc string `json:"todo_desc"`
}

var DB = make(map[int]*List)

func findListById(id int) *List {
	if val, ok := DB[id]; ok {
		return val
	} else {
		return nil
	}
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/api/lists", func(c *gin.Context) {
		c.JSON(http.StatusOK, &DB)
	})

	router.GET("/api/list/:Id", func(c *gin.Context) {
		IdParam, err := strconv.Atoi(c.Param("Id"))
		if err != nil {
			log.Fatal(err)
		}

		list := findListById(IdParam)
		if list == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "List was not found."})
			return
		}

		c.JSON(http.StatusOK, list)
	})

	router.POST("/api/list/create", func(c *gin.Context) {
		// append to List to DB
		var createId = len(DB) + 1
		var emptyItems [][2]interface{}

		DB[createId] = &List{createId, emptyItems}

		c.JSON(http.StatusCreated, gin.H{"id": createId})
	})

	router.POST("/api/list/:Id/createItem", func(c *gin.Context) {
		IdParam, err := strconv.Atoi(c.Param("Id"))
		if err != nil {
			log.Fatal(err)
		}

		list := findListById(IdParam)
		if list == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "List was not found."})
			return
		}

		var body CreateListItemRequestBody

		c.BindJSON(&body)

		if body.Todo_Desc == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "'todo_desc' was not found in the body"})
		} else {
			DB[IdParam].AddItem([2]interface{}{body.Todo_Desc, false})

			c.JSON(http.StatusCreated, gin.H{"message": "Created."})
		}
	})

	router.DELETE("/api/list/:Id/items/:itemId", func(c *gin.Context) {
		ListIdParam, err := strconv.Atoi(c.Param("Id"))
		if err != nil {
			log.Fatal(err)
		}

		ItemIdParam, err := strconv.Atoi(c.Param("Id"))
		if err != nil {
			log.Fatal(err)
		}

		list := findListById(ListIdParam)
		if list == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "List was not found."})
			return
		}

		if len(list.Items) < ItemIdParam {
			c.JSON(http.StatusNotFound, gin.H{"message": "Item could not be located."})
			return
		}

		DB[ListIdParam].RemoveItem(ItemIdParam)

		c.JSON(http.StatusOK, gin.H{"message": "Deleted."})
	})

	return router
}

func main() {
	DB[1] = &List{}

	r := setupRouter()

	StartService()

	r.Run(":4000")
}
