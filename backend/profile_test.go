package backend

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestQueryUsersWithUserID(t *testing.T) {
	// 使用 sqlmock 创建一个 mock 数据库连接
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	// 模拟 SELECT VERSION() 查询（gorm 会执行这个查询来检查数据库连接）
	mock.ExpectQuery("SELECT VERSION()").
		WillReturnRows(sqlmock.NewRows([]string{"Version"}).AddRow("5.7.33"))

	// 使用 gorm.Open 配合 sqlmock 模拟数据库连接
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: mockDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open gorm DB: %v", err)
	}
	db = gormDB // 将 db 赋值为 mock gormDB

	// 创建 Gin 路由实例
	router := gin.Default()
	router.POST("/profile/queryUser", queryUsersWithUserID)

	// 测试用例 1: 用户存在
	t.Run("User Exists", func(t *testing.T) {
		// 模拟查询操作
		userID := "1"
		queriedUser := &User{
			UserID:   userID,
			UserType: "admin",
			UserName: "user123",
			Email:    "user123@example.com",
		}
		rows := sqlmock.NewRows([]string{"userId", "username", "userType", "Email"}).
			AddRow(queriedUser.UserID, queriedUser.UserName, queriedUser.UserType, queriedUser.Email)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE userId = ? LIMIT ?")).
			WithArgs(userID, 1).
			WillReturnRows(rows)

		// 准备请求数据
		reqBody := `{"userId": "1"}`
		req := httptest.NewRequest(http.MethodPost, "/profile/queryUser", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		// 创建响应记录器
		rec := httptest.NewRecorder()

		// 执行请求
		router.ServeHTTP(rec, req)

		// 验证响应
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"status":"success"`)
		assert.Contains(t, rec.Body.String(), `"usertype":"admin"`)
		assert.Contains(t, rec.Body.String(), `"username":"user123"`)
		assert.Contains(t, rec.Body.String(), `"email":"user123@example.com"`)

		// 确保所有预期的操作都已执行
		err := mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("There were unmet expectations: %s", err)
		}
	})
}
