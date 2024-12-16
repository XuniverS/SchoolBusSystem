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
func TestSubmitUserInfo(t *testing.T) {
	// 使用 sql mock 创建一个 mock 数据库连接
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	// 模拟 SELECT VERSION() 查询（gourm 会执行这个查询来检查数据库连接）
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
	router.POST("/profile/submitUser", submitUserInfo)

	// 测试数据
	userId := "user123"
	user := User{
		UserID:   userId,
		UserName: "New User",
		UserType: "admin",
		Email:    "newuser@example.com",
	}

	// 模拟 queryUserWithUserID 函数的行为
	rows := sqlmock.NewRows([]string{"userId", "username", "userType", "Email"}).
		AddRow(user.UserID, user.UserName, user.UserType, user.Email)

	mock.ExpectQuery("SELECT * FROM `users` WHERE userId = ? LIMIT ?").
		WithArgs(user.UserID, 1).
		WillReturnRows(rows)

	// 模拟更新用户信息
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `users` SET `userType`=?,`username`=?,`Email`=?,`password`=?,`is_First_Login`=?,`created_At`=? WHERE `userId` = ?").
		WithArgs(user.UserType, user.UserName, user.Email, user.Password, user.Is_first_login, user.CreatedTime, user.UserID).
		WillReturnResult(sqlmock.NewResult(1, 1)) // 返回影响的行数
	mock.ExpectCommit()

	// 准备请求数据
	reqBody := `{"userId": "user123", "userName": "New User", "email": "newuser@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/profile/submitUser", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// 创建响应记录器
	rec := httptest.NewRecorder()

	// 调用 submitUserInfo 函数
	router.ServeHTTP(rec, req)

	// 断言结果
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"status":"success"`)
	assert.Contains(t, rec.Body.String(), `"message":"User information updated"`)

	// 检查所有期望的 SQL 查询是否都被执行
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}
func TestChangePassword(t *testing.T) {
	// 使用 sql mock 创建一个 mock 数据库连接
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	// 模拟 SELECT VERSION() 查询（gourm 会执行这个查询来检查数据库连接）
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
	router.POST("/profile/changePassword", changePassword)

	// 测试数据
	changeRequest := UpdatePasswordRequest{
		UserID:         "user123",
		OriginPassword: "oldpassword",
		NewPassword:    "newpassword",
	}

	// 模拟查询用户存在
	existingUser := User{
		UserID:   "user123",
		UserName: "New User",
		UserType: "admin",
		Password: "f0a5cdf5a9b255d3a71acdee7bd29c6b320f27e71f105b86220696f21b67c6e9", // 这个密码应该与 OriginPassword 匹配
	}
	rows := sqlmock.NewRows([]string{"userId", "username", "userType", "password"}).
		AddRow(existingUser.UserID, existingUser.UserName, existingUser.UserType, existingUser.Password)
	mock.ExpectQuery("SELECT * FROM `users` WHERE userId = ? LIMIT ?").
		WithArgs(existingUser.UserID, 1).
		WillReturnRows(rows)

	// 模拟密码变更操作
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `users` SET `userType`=?,`username`=?,`Email`=?,`password`=?,`is_First_Login`=?,`created_At`=? WHERE `userId` = ?").
		WithArgs(existingUser.UserType, existingUser.UserName, existingUser.Email, "newpassword", 0, existingUser.CreatedTime, changeRequest.UserID).
		WillReturnResult(sqlmock.NewResult(1, 1)) // 返回影响的行数
	mock.ExpectCommit()

	// 准备请求数据
	reqBody := `{"userId": "user123", "originPassword": "oldpassword", "newPassword": "newpassword"}`
	req := httptest.NewRequest(http.MethodPost, "/profile/changePassword", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// 创建响应记录器
	rec := httptest.NewRecorder()

	// 调用 changePassword 函数
	router.ServeHTTP(rec, req)

	// 断言结果
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"status":"success"`)

	// 检查所有期望的 SQL 查询是否都被执行
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}
