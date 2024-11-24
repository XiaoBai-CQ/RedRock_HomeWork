package main

import (
	"RR5/lv2andlv3/models"
	"context"
	"database/sql"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)
import "RR5/lv2andlv3/dao"

// 懒了点就没写utils hhhh
func main() {
	db := dao.InitDB()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	engine := server.Default()

	//添加学生
	engine.POST("/add", func(c context.Context, ctx *app.RequestContext) {
		var student models.Student
		err1 := ctx.Bind(&student)
		if err1 != nil {
			return
		}
		query := "INSERT INTO users (name,sex,birth,born) VALUES (?,?,?,?)"
		result, err2 := db.Exec(query, student.Name, student.Sex, student.Birth, student.Born)
		if err2 != nil {
			return
		}
		id, _ := result.LastInsertId()
		student.Id = int(id)
		ctx.JSON(http.StatusOK, student)
	})

	// 获取所有学生
	engine.GET("/messages", func(c context.Context, ctx *app.RequestContext) {
		rows, err1 := db.Query("SELECT id, name,sex,birth,born FROM users")
		if err1 != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
			return
		}
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {

			}
		}(rows)

		var messages []models.Student
		for rows.Next() {
			var student models.Student
			if err := rows.Scan(&student.Id, &student.Name, &student.Sex, &student.Born, &student.Birth); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan message"})
				return
			}
			messages = append(messages, student)
		}

		ctx.JSON(http.StatusOK, messages)
	})

	//删除学生
	engine.DELETE("/messages/:id", func(c context.Context, ctx *app.RequestContext) {
		id := ctx.Param("id")
		query := "DELETE FROM users WHERE id = ?"
		_, err := db.Exec(query, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete message"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
	})

	// 更新学生
	engine.PUT("/messages/:id", func(c context.Context, ctx *app.RequestContext) {
		id := ctx.Param("id")
		var student models.Student
		if err := ctx.Bind(&student); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := "UPDATE users SET name = ?, sex = ?,birth= ?,born = ? WHERE id = ?"
		_, err := db.Exec(query, student.Name, student.Sex, student.Birth, student.Born, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update message"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Message updated successfully"})
	})

	//精准查找
	engine.GET("/search/:id", func(c context.Context, ctx *app.RequestContext) {
		id := ctx.Param("id")
		var student models.Student
		query := "SELECT id, name,sex,birth,born FROM users WHERE id = ?"
		result := db.QueryRow(query, id)
		errSC := result.Scan(&student.Id, &student.Name, &student.Sex, &student.Birth, &student.Born)
		if errSC != nil {
			return
		}
		ctx.JSON(http.StatusOK, student)
	})
	errR := engine.Run()
	if errR != nil {
		return
	}
}
