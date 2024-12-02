package main

import (
	"RR6/lv1Andlv2Andlv3/dao"
	"RR6/lv1Andlv2Andlv3/flag"
	"RR6/lv1Andlv2Andlv3/global"
	"RR6/lv1Andlv2Andlv3/models"
	"RR6/lv1Andlv2Andlv3/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func main() {
	global.DB = dao.InitMysql()

	engine := gin.Default()

	//参数创建表
	option := flag.Parse()
	if flag.IsWebStop(option) {
		flag.SwitchOption(option)
		return
	}

	engine.POST("/find", func(ctx *gin.Context) {
		sq := ctx.PostForm("SecurityQuestion")
		sa := ctx.PostForm("SecurityAnswer")
		username := ctx.PostForm("Username")
		newpassword := ctx.PostForm("Password")

		var user models.User
		err := global.DB.Where("security_question  = ? AND security_answer  = ? AND Username = ?", sq, sa, username).First(&user).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户名、密保问题或答案错误"})
			return
		}

		// 更新密码
		if err := global.DB.Model(&user).Update("Password", newpassword).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新密码失败"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "密码更新成功"})
	})

	//注册
	engine.POST("/register", func(ctx *gin.Context) {
		username := ctx.PostForm("username")
		password := ctx.PostForm("password")
		user := models.User{Username: username, Password: password}
		global.DB.Create(&user)
	})

	//登录
	engine.POST("/login", func(ctx *gin.Context) {
		username := ctx.PostForm("username")
		password := ctx.PostForm("password")

		user := models.User{Username: username, Password: password}

		err2 := global.DB.Where("username = ? AND password = ?", user.Username, user.Password).First(&user).Error
		if err2 != nil {
			if errors.Is(err2, gorm.ErrRecordNotFound) {
				// 用户名或密码错误
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户名或密码错误"})
			} else {
				// 数据库错误
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "数据库查询错误"})
			}
			return
		}

		// 生成 JWT
		token, err := utils.GenerateToken(user.Username)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  500,
				"message": "token generate error",
			})
			return
		}
		
		ctx.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": "login successful",
			"token":   token,
		})
	})

	protected := engine.Group("/")
	protected.Use(utils.JWTAuthMiddleware())

	// 添加学生
	protected.POST("/add", func(ctx *gin.Context) {
		var student models.Student
		err := ctx.Bind(&student)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if result := global.DB.Create(&student); result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add student"})
			return
		}

		ctx.JSON(http.StatusOK, student)
	})

	// 获取所有学生
	protected.GET("/messages", func(ctx *gin.Context) {
		var students []models.Student

		if result := global.DB.Find(&students); result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch students"})
			return
		}

		ctx.JSON(http.StatusOK, students)
	})

	// 删除学生
	protected.DELETE("Delete/deleteStudent/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		if result := global.DB.Delete(&models.Student{}, id); result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete student"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Student deleted successfully"})
	})

	// 更新学生
	protected.PUT("/messages/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var student models.Student

		if err := ctx.Bind(&student); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if result := global.DB.Model(&models.Student{}).Where("id = ?", id).Updates(student); result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Student updated successfully"})
	})

	// 精准查找学生
	protected.GET("/search/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var student models.Student

		if result := global.DB.First(&student, id); result.Error != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
			return
		}

		ctx.JSON(http.StatusOK, student)
	})

	err := engine.Run()
	if err != nil {
		return
	}
}
