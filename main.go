package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"encoding/json"
	"gorm.io/gorm"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

func main() {
	router := gin.Default()
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v \n", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	// router.Use(gin.Logger())
	// router.Use(gin.Recovery())
	
	router.Use(Cors())
	
	// db, _ := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/todo")
	db, _ := sql.Open("mysql", "root:123456@/todo")
	
	todo := router.Group("/todo")
	todo.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello golang")
	})
	
	todo.GET("/list", func(ctx *gin.Context) {
		completed := ctx.Query("completed")

		list := "select * from todo"
		
		var rows *sql.Rows
		if completed != "" {
			list += " where completed = ?"
			rows, _ = db.Query(list, completed)
		}
		if completed == "" {
			rows, _ = db.Query(list)
		}

		data := []Todo{}
		for rows.Next() {
			var todo Todo
			err := rows.Scan(&todo.Id, &todo.Text, &todo.Completed)

			if err != nil {
				fmt.Println("scan err: ", err)
			}
			
			data = append(data, todo)
		}

		ctx.JSON(http.StatusOK, gin.H{ "data": data})
	})
	
	todo.GET("/:id", func(ctx *gin.Context) {
		// id := ctx.Params.ByName("id")
		id := ctx.Param("id")
		detail := "select * from todo where id = ?"
		var todo Todo
		err := db.QueryRow(detail, id).Scan(&todo.Id, &todo.Text, &todo.Completed)

		if err != nil {
			fmt.Println("query err: ", err)
		}

		ctx.JSON(http.StatusOK, gin.H{ "data": todo })
	})

	todo.POST("/add", func (ctx *gin.Context) {
		var params Todo
		_ = ctx.Bind(&params)
		text := params.Text
		completed := params.Completed

		add := "insert into todo (text, completed) values (?, ?)"
		ret, err := db.Exec(add, text, completed)
		if err != nil {
			fmt.Println("add err: ", err)
		}

		id, err := ret.LastInsertId()
		if err != nil {
			fmt.Println("get id err: ", err)
		}

		ctx.JSON(http.StatusOK, gin.H{ "data": id })
	})

	todo.POST("/update", func (ctx *gin.Context) {
		var params Todo
		_ = ctx.Bind(&params)
		id := params.Id
		text := params.Text
		completed := params.Completed

		if text != "" {
			update := "update todo set text = ?, completed = ? where id = ?"
			ret, _ := db.Exec(update, text, completed, id)
			num, _ := ret.RowsAffected()
			ctx.JSON(http.StatusOK, gin.H{ "data": num })
		}
		if text == "" {
			update := "update todo set completed = ? where id = ?"
			ret, _ := db.Exec(update, completed, id)
			num, _ := ret.RowsAffected()
			ctx.JSON(http.StatusOK, gin.H{ "data": num })
		}
	})

	todo.POST("/delete", func (ctx *gin.Context) {
		var params map[string]int
		json.NewDecoder(ctx.Request.Body).Decode(&params)
		id := params["id"]

		del := "delete from todo where id = ?"
		ret, err := db.Exec(del, id)
		if err != nil {
			fmt.Println("delete err: ", err)
		}
		num, _ := ret.RowsAffected()

		ctx.JSON(http.StatusOK, gin.H{ "data": num })
	})

	todo.POST("/clear", func (ctx *gin.Context) {
		clear := "delete from todo"
		ret, _ := db.Exec(clear)
		num, _ := ret.RowsAffected()
		ctx.JSON(http.StatusOK, gin.H{ "data": num })
	})

	router.NoRoute(func (ctx *gin.Context) {
		
	})

	Test()

	router.Run(":8701")
}

type Todo struct {
	Id int `json:"id"`
	Text string	`json:"text"`
	Completed int `json:"completed"`
}

type TodoModel struct {
	gorm.Model
	Id int
	Text string
	Completed int
}

func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method

		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
	}
}

func Test() {
	data2 := [2]string{"a", "b"}
	fmt.Println(data2, data2[1])

	jsonStr := `
		{
			"id": 2,
			"text": "hello",
			"completed": 1
		}
	`
	data3 := []byte(jsonStr)
	var todo_ Todo
	json.Unmarshal(data3, &todo_)
	fmt.Println(todo_, todo_.Text)

	js11, _ := json.Marshal(todo_)
	js12 := string(js11)
	fmt.Println(js12, strings.Count(js12, "")-1, len([]rune(js12)))

	var mr map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &mr)
	fmt.Println(mr, mr["id"])

	js21, _ := json.Marshal(mr)
	js22 := string(js21)
	fmt.Println(js22)

	fmt.Println("aabbccdd\t:",strings.TrimLeft("aabbccdd","abcd"))  // 空字符串
	fmt.Println("aabbccdde\t:",strings.TrimLeft("aabbccdde","abcd")) // e
	fmt.Println("aabbeccdd\t:",strings.TrimLeft("aabbedcba","abcd")) // edcba
	fmt.Println("aabbccdd\t:",strings.TrimRight("aabbccdd","abcd"))  // 空字符串
	fmt.Println("aabbccdde\t:",strings.TrimRight("aabbccdde","abcd")) // aabbccdde
	fmt.Println("aabbeccdd\t:",strings.TrimRight("aabbedcba","abcd")) //aabbe

	strHaiCoder := "HaiCoder 嗨客网 HaiCoderHaiCoder"
	TrimRightStr := strings.TrimRight(strHaiCoder, "HaiCoder")
	fmt.Println("TrimRightStr =", TrimRightStr)

	fmt.Println(strings.TrimRight("cyeamblog.go", ".go"))
	fmt.Println(strings.TrimSuffix("cyeamblog.go", ".go"))
}