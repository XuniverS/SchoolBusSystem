package backend

import (
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestQueryUser(t *testing.T) {
	// 使用 sqlmock 创建一个 mock 数据库连接
	mockDB, mock, err := sqlmock.New()
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

	// 定义要查询的用户
	userToQuery := &User{
		UserName: "user123",
	}

	// 定义期望的查询 SQL 及其结果
	rows := sqlmock.NewRows([]string{"userId", "username", "password"}).
		AddRow("1", "user123", "hashed_password")

	// 模拟查询操作
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ? LIMIT ?")).
		WithArgs(userToQuery.UserName, 1).
		WillReturnRows(rows)

	// 调用 queryUser 函数进行查询
	queriedUser, err := queryUser(userToQuery)

	// 验证查询结果
	assert.NoError(t, err)
	assert.NotNil(t, queriedUser)
	assert.Equal(t, "user123", queriedUser.UserName)
	assert.Equal(t, "1", queriedUser.UserID)
	assert.Equal(t, "hashed_password", queriedUser.Password)

	// 确保所有预期的操作都已执行
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("There were unmet expectations: %s", err)
	}
}

// 单元测试
func TestUserLogin(t *testing.T) {
	// 初始化 Gin engine
	r := gin.Default()

	// 注册登录接口
	r.POST("/login", userLogin)

	// 正常登录测试
	t.Run("Valid Login", func(t *testing.T) {
		// 使用 sqlmock 创建一个 mock 数据库连接
		mockDB, mock, err := sqlmock.New()
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

		// 定义要查询的用户
		userToQuery := &User{
			UserName: "user123",
			UserID:   "1",
		}

		// 定义期望的查询 SQL 及其结果
		rows := sqlmock.NewRows([]string{"userId", "userType", "username", "password", "Is_first_login"}).
			AddRow("1", "admin", "user123", "b17e1e0450dac425ea318253f6f852972d69731d6c7499c001468b695b6da219", "1")

		// 模拟查询操作
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ? LIMIT ?")).
			WithArgs(userToQuery.UserName, 1).
			WillReturnRows(rows)
		// 模拟更新操作
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `is_first_login` = ? WHERE `userId` = ?")).
			WithArgs(0, userToQuery.UserID).
			WillReturnResult(sqlmock.NewResult(1, 1)) // 假设更新了 1 条记录
		// 构造请求数据
		input := map[string]string{
			"username": "user123",
			"password": "123456Aa",
		}

		// 转换成 JSON
		jsonValue, _ := json.Marshal(input)
		//fmt.Printf("POST: %+v\n", bytes.NewReader(jsonValue))
		// 创建 POST 请求
		req, _ := http.NewRequest("POST", "/login", bytes.NewReader(jsonValue))
		w := httptest.NewRecorder()

		// 执行请求
		r.ServeHTTP(w, req)

		// 断言响应状态码为 200
		assert.Equal(t, http.StatusOK, w.Code)

		// 解析响应体
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err) // 确保 JSON 解析没有错误

		// 验证响应内容
		assert.Equal(t, "success", response["status"])
		assert.Equal(t, "1", response["userid"])
		assert.Equal(t, "admin", response["usertype"])
		assert.Equal(t, float64(1), response["isfirstlogin"])
	})
}
