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
	"time"
)

// Mock insertBus 函数，模拟数据库操作
/*func insertBus(bus *Bus) error {
	// 模拟成功插入数据库
	return nil
}*/

// 测试 addBus 函数
func TestAddBus(t *testing.T) {
	// 初始化 Gin engine
	r := gin.Default()

	// 注册 addBus 接口
	r.POST("/addBus", addBus)

	// 正常添加班车测试
	t.Run("Valid Add Bus", func(t *testing.T) {
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

		busToInsert := &Bus{
			Origin:         "清水河",
			Destination:    "沙河",
			BusType:        "师生车",
			Date:           time.Now().Add(24 * time.Hour), // 设置日期为明天
			Time:           "08:00",
			Plate:          "京A12345",
			TotalSeats:     50,
			AvailableSeats: 30,
			CreatedAt:      time.Now(),
		}

		// 模拟事务开始
		mock.ExpectBegin()

		// 定义期望的插入 SQL 及其结果
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `buses` (`origin`, `destination`, `busType`, `date`, `time`, `plate`, `total_seats`, `available_seats`, `created_at`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")).
			WithArgs(busToInsert.Origin, busToInsert.Destination, busToInsert.BusType, busToInsert.Date, busToInsert.Time, busToInsert.Plate, busToInsert.TotalSeats, busToInsert.AvailableSeats, busToInsert.CreatedAt).
			WillReturnResult(sqlmock.NewResult(1, 1)) // 假设插入了 1 条记录

		mock.ExpectCommit()
		// 构造请求数据
		input := map[string]interface{}{
			"origin":         "清水河",
			"destination":    "沙河",
			"busType":        "师生车",
			"date":           busToInsert.Date.Format(time.RFC3339),
			"time":           busToInsert.Time,
			"plate":          busToInsert.Plate,
			"totalSeats":     busToInsert.TotalSeats,
			"availableSeats": busToInsert.AvailableSeats,
			"createdTime":    busToInsert.CreatedAt.Format(time.RFC3339),
		}

		// 转换成 JSON
		jsonValue, _ := json.Marshal(input)

		// 创建 POST 请求
		req, _ := http.NewRequest("POST", "/addBus", bytes.NewReader(jsonValue))
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
		assert.Equal(t, "添加成功", response["message"])
	})

	/*// 测试班车类型错误
	t.Run("Invalid Bus Type", func(t *testing.T) {
		// 构造请求数据，使用无效的班车类型
		input := map[string]interface{}{
			"bus_type": "无效车",
			"date":     time.Now().Add(time.Hour).Format(time.RFC3339),
		}

		// 转换成 JSON
		jsonValue, _ := json.Marshal(input)

		// 创建 POST 请求
		req, _ := http.NewRequest("POST", "/addBus", bytes.NewReader(jsonValue))
		w := httptest.NewRecorder()

		// 执行请求
		r.ServeHTTP(w, req)

		// 断言响应状态码为 500
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// 解析响应体
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err) // 确保 JSON 解析没有错误

		// 验证错误消息
		assert.Equal(t, "fail", response["status"])
		assert.Equal(t, "班车类型错误！仅限为师生车或教职工车", response["message"])
	})

	// 测试时间错误
	t.Run("Invalid Date", func(t *testing.T) {
		// 构造请求数据，使用过去的时间
		input := map[string]interface{}{
			"bus_type": "师生车",
			"date":     time.Now().Add(-time.Hour).Format(time.RFC3339), // 设置为过去的时间
		}

		// 转换成 JSON
		jsonValue, _ := json.Marshal(input)

		// 创建 POST 请求
		req, _ := http.NewRequest("POST", "/addBus", bytes.NewReader(jsonValue))
		w := httptest.NewRecorder()

		// 执行请求
		r.ServeHTTP(w, req)

		// 断言响应状态码为 500
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// 解析响应体
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err) // 确保 JSON 解析没有错误

		// 验证错误消息
		assert.Equal(t, "fail", response["status"])
		assert.Equal(t, "错误的时间！", response["message"])
	})

	// 测试数据库写入失败
	t.Run("Database Insert Fail", func(t *testing.T) {
		// 使用 sqlmock 创建一个 mock 数据库连接
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock database: %v", err)
		}
		defer mockDB.Close()

		// 使用 gorm.Open 配合 sqlmock 模拟数据库连接
		gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: mockDB}), &gorm.Config{})
		if err != nil {
			t.Fatalf("Failed to open gorm DB: %v", err)
		}
		db = gormDB // 将 db 赋值为 mock gormDB

		// 模拟插入失败
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `buses` (`bus_type`, `date`, `created_at`) VALUES (?, ?, ?)")).
			WithArgs("师生车", time.Now().Add(time.Hour), time.Now()). // 假设插入数据
			WillReturnError(errors.New("Database error"))

		// 构造请求数据
		input := map[string]interface{}{
			"bus_type": "师生车",
			"date":     time.Now().Add(time.Hour).Format(time.RFC3339),
		}

		// 转换成 JSON
		jsonValue, _ := json.Marshal(input)

		// 创建 POST 请求
		req, _ := http.NewRequest("POST", "/addBus", bytes.NewReader(jsonValue))
		w := httptest.NewRecorder()

		// 执行请求
		r.ServeHTTP(w, req)

		// 断言响应状态码为 500
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// 解析响应体
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err) // 确保 JSON 解析没有错误

		// 验证错误消息
		assert.Equal(t, "fail", response["status"])
		assert.Equal(t, "数据库写入失败", response["message"])
	})*/
}
