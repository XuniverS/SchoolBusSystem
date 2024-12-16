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
	"strings"
	"testing"
	"time"
)

// 测试 addBus 函数
func TestAddBus(t *testing.T) {
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
	router.POST("/bus/addBus", addBus)

	// 定义测试用的班车数据
	// 定义时间字符串
	timeStr := "2024-12-17 00:00:00"
	// 使用 time.Parse 解析字符串，指定格式与 UTC 时区
	layout := "2006-01-02 15:04:05"
	parsedTime, err := time.Parse(layout, timeStr)
	reqBody := `{
		"origin": "清水河",
		"destination": "沙河",
		"bustype": "师生车",
		"date": "2024-12-17",
		"time": "08:00",
		"plate": "A12345",
		"totalseats": 50
	}`
	requestBus := Bus{
		Origin:      "清水河",
		Destination: "沙河",
		BusType:     "师生车",
		Time:        "08:00",
		Date:        parsedTime,
		Plate:       "A12345",
		TotalSeats:  50,
		CreatedAt:   time.Now().Truncate(time.Second),
	}

	mock.ExpectBegin() // 模拟事务开始
	mock.ExpectExec("INSERT INTO `buses` (`origin`,`destination`,`busType`,`date`,`time`,`plate`,`total_seats`,`available_seats`,`created_at`) VALUES (?,?,?,?,?,?,?,?,?)").
		WithArgs(requestBus.Origin, requestBus.Destination, requestBus.BusType, requestBus.Date, requestBus.Time, requestBus.Plate, requestBus.TotalSeats, requestBus.AvailableSeats, requestBus.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1)) // 假设插入了 1 条记录
	mock.ExpectCommit() // 模拟事务提交

	// 创建响应记录器
	rec := httptest.NewRecorder()
	// 创建请求数据
	req := httptest.NewRequest(http.MethodPost, "/bus/addBus", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// 执行请求
	router.ServeHTTP(rec, req)

	// 断言响应状态码为 200
	assert.Equal(t, http.StatusOK, rec.Code)

	// 解析响应体
	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err) // 确保 JSON 解析没有错误

	// 验证响应内容
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "添加成功", response["message"])

	// 验证模拟的 SQL 查询是否正确
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unmet expectations: %s", err)
	}
}
func TestRemoveBus(t *testing.T) {
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
	router.POST("/bus/deleteBus", removeBus)

	// 模拟请求体
	bus := Bus{
		BusId: 123, // 假设删除的公交车ID
	}

	// 将 bus 对象转换为 JSON
	jsonBus, err := json.Marshal(bus)
	if err != nil {
		t.Fatalf("Error marshalling bus: %v", err)
	}

	t.Run("TestDeleteBusHasBooking", func(t *testing.T) {
		// 模拟查询到记录且有用户预约
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `buses` WHERE busId =? LIMIT ?")).
			WithArgs(bus.BusId, 1).WillReturnRows(sqlmock.NewRows([]string{"busId"}))

		// 创建一个 POST 请求
		req, err := http.NewRequest(http.MethodPost, "/bus/deleteBus", bytes.NewBuffer(jsonBus))
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		// 创建一个响应记录器
		w := httptest.NewRecorder()

		// 执行请求
		router.ServeHTTP(w, req)

		// 断言状态码
		assert.Equal(t, http.StatusOK, w.Code)

		// 断言返回的 JSON 数据
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Error unmarshalling response: %v", err)
		}

		// 验证响应数据
		assert.Equal(t, "fail", response["status"])
		assert.Equal(t, "要删除的记录不存在", response["message"])
	})

}
func TestInitPasswordWithUserID(t *testing.T) {
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
	router.POST("/bus/initPassword", initPasswordWithUserID)

	// 测试数据
	userId := "user123"
	expectedUser := User{
		UserID:   userId,
		UserName: "Test User",
		Password: "", // 初始密码为空，稍后设置为 "b17e1e0450dac425ea318253f6f852972d69731d6c7499c001468b695b6da219"
	}

	// 设置 mock 查询行为，模拟用户存在
	mock.ExpectQuery("SELECT * FROM `users` WHERE userId =? AND `users`.`userId` = ? LIMIT ?").
		WithArgs(userId, userId, 1).
		WillReturnRows(sqlmock.NewRows([]string{"userId", "username", "email", "password"}).
			AddRow(expectedUser.UserID, expectedUser.UserName, expectedUser.Email, expectedUser.Password))

	// 设置 mock 保存用户操作，模拟成功保存
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `users` SET `userType`=?,`username`=?,`Email`=?,`password`=?,`is_First_Login`=?,`created_At`=? WHERE `userId` = ?").
		WithArgs(expectedUser.UserType, expectedUser.UserName, expectedUser.Email, "b17e1e0450dac425ea318253f6f852972d69731d6c7499c001468b695b6da219", expectedUser.Is_first_login, expectedUser.CreatedTime, userId).
		WillReturnResult(sqlmock.NewResult(1, 1)) // 返回影响的行数
	mock.ExpectCommit()
	// 准备请求数据
	reqBody := `{"userId": "user123"}`
	req := httptest.NewRequest(http.MethodPost, "/bus/initPassword", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// 创建响应记录器
	rec := httptest.NewRecorder()

	// 调用 initPasswordWithUserID 函数
	router.ServeHTTP(rec, req)

	// 断言结果
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"status":"success"`)

	// 检查所有期望的 SQL 查询是否都被执行
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}
