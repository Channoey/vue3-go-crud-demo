package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strconv"
	"time"
)

func main() {
	//数据库连接
	dsn := "root:123456@tcp(127.0.0.1:3306)/go-vue-demo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //后面创建表时不自动加s
		},
	})
	fmt.Println(db)
	fmt.Println(err)

	//数据库配置
	sqlDB, err := db.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	//数据结构体
	type List struct {
		gorm.Model
		Name    string `gorm:"type:varchar(20); not null" json:"name" binding:"required"`
		State   string `gorm:"type:varchar(20); not null" json:"state" binding:"required"`
		Phone   string `gorm:"type:varchar(20); not null" json:"phone" binding:"required"`
		Email   string `gorm:"type:varchar(40); not null" json:"email" binding:"required"`
		Address string `gorm:"type:varchar(200); not null" json:"address" binding:"required"`
	}

	//根据结构体自动创建表
	db.AutoMigrate(&List{})

	//接口初始化
	r := gin.Default()

	//POST 增加接口
	r.POST("/user/add", func(c *gin.Context) {
		var data List
		err := c.ShouldBindJSON(&data) //将传入数据与定义的data绑定

		if err != nil {
			c.JSON(200, gin.H{
				"msg":  "添加失败",
				"data": gin.H{},
				"code": 400,
			})
		} else {
			db.Create(&data) //自动插入List对应的表
			c.JSON(200, gin.H{
				"msg":  "添加成功",
				"data": data,
				"code": 200,
			})
		}
	})

	//DELETE 删除接口
	r.DELETE("/user/delete/:id", func(c *gin.Context) {
		var data []List
		id := c.Param("id")                // "/id=1"用query "/:id"用Param
		db.Where("id = ?", id).Find(&data) //判断id是否存在
		if len(data) == 0 {
			c.JSON(200, gin.H{
				"msg":  "删除失败，id未找到",
				"code": 400,
			})
		} else {
			db.Delete(&data) //删除对应Id的数据
			c.JSON(200, gin.H{
				"msg":  "删除成功",
				"code": 200,
			})
		}

	})

	//PUT 修改接口
	r.PUT("/user/update/:id", func(c *gin.Context) {
		var data List
		id := c.Param("id")
		db.Where("id = ?", id).Find(&data) //判断id是否存在
		if data.ID == 0 {
			c.JSON(200, gin.H{
				"msg":  "用户id未找到",
				"code": 400,
			})
		} else {
			err := c.ShouldBindJSON(&data)
			if err != nil {
				c.JSON(200, gin.H{
					"msg":  "修改失败",
					"code": 400,
				})
			} else {
				db.Where("id = ?", id).Updates(&data) //修改操作
				c.JSON(200, gin.H{
					"msg":  "修改成功",
					"code": 200,
				})
			}

		}
	})

	//GET 查询接口
	//1.条件查询（根据姓名模糊查询）
	r.GET("/user/list/:name", func(c *gin.Context) {
		name := c.Param("name")
		var dataList []List
		db.Where("name LIKE ?", "%%"+name+"%%").Find(&dataList)
		if len(dataList) == 0 {
			c.JSON(200, gin.H{
				"msg":  "未查询到数据",
				"code": 400,
				"data": gin.H{},
			})
		} else {
			c.JSON(200, gin.H{
				"msg":  "查询成功",
				"code": 200,
				"data": dataList,
			})
		}

	})

	//2.分页查询
	r.GET("/user/list", func(c *gin.Context) {
		var dataList []List
		pageSize, _ := strconv.Atoi(c.Query("pageSize"))
		pageNum, _ := strconv.Atoi(c.Query("pageNum"))

		//判断是否需要分页，pageSize和pageNum均为0的话返回全部
		if pageSize == 0 {
			pageSize = -1
		}
		if pageNum == 0 {
			pageNum = -1
		}

		//从数据库中分页查询
		offsetVal := (pageNum - 1) * pageSize
		if pageNum == -1 && pageSize == -1 {
			offsetVal = -1
		}
		var total int64
		db.Model(dataList).Count(&total).Limit(pageSize).Offset(offsetVal).Find(&dataList)

		if len(dataList) == 0 {
			c.JSON(200, gin.H{
				"msg":  "未查询到数据",
				"code": 400,
				"data": gin.H{},
			})
		} else {
			c.JSON(200, gin.H{
				"msg":  "查询成功",
				"code": 200,
				"data": gin.H{
					"list":     dataList,
					"total":    total,
					"pageNum":  pageNum,
					"pageSize": pageSize,
				},
			})
		}
	})

	PORT := "3000"
	r.Run(":" + PORT)
}
