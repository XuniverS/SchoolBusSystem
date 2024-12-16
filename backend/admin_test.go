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

// Mock insertBus 函数，模拟数据库操作
/*func insertBus(bus *Bus) error {
	// 模拟成功插入数据库
	return nil
}*/

// 测试 addBus 函数
func TestAddBus(t *testing.T) {
	// 初始化 Gin engine
	r := gin.Default()

	// 注册添加班车接口
	r.POST("/addBus", addBus)

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

	// 定义测试用的班车数据
	requestBus := struct {
		Origin      string `json:"origin"`
		Destination string `json:"destination"`
		BusType     string `json:"bustype"`
		Date        string `json:"date"`
		Time        string `json:"time"`
		Plate       string `json:"plate"`
		TotalSeats  int    `json:"totalseats"`
	}{
		Origin:      "清水河",
		Destination: "沙河",
		BusType:     "师生车",
		Date:        "2024-12-17",
		Time:        "08:00",
		Plate:       "京A12345",
		TotalSeats:  50,
	}

	// 转换成 JSON 格式
	jsonValue, _ := json.Marshal(requestBus)
	//busDate, err := time.Parse("2006-01-02", requestBus.Date)
	// 创建 POST 请求
	req, _ := http.NewRequest("POST", "/addBus", bytes.NewReader(jsonValue))
	w := httptest.NewRecorder()
	mock.ExpectBegin() // 模拟事务开始
	// 模拟数据库插入操作
	/*mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `buses` (`origin`, `destination`, `busType`, `date`, `time`, `plate`, `total_seats`, `available_seats`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")).
	WithArgs(requestBus.Origin, requestBus.Destination, requestBus.BusType, busDate, requestBus.Time, requestBus.Plate, requestBus.TotalSeats, requestBus.TotalSeats).
	WillReturnResult(sqlmock.NewResult(1, 1)) // 假设插入成功，返回 1 条记录*/

	//mock.ExpectQuery(regexp.QuoteMeta("SELECT").
	//	WithArgs("清水河", "沙河", "师生车", "2024-12-17 00:00:00", "08:00", "京A12345"). // 模拟查询条件
	//	WillReturnRows(sqlmock.NewRows([]string{"id", "origin", "destination", "busType", "date", "time", "plate", "total_seats", "available_seats", "created_at"}).
	//		AddRow(1, "清水河", "沙河", "师生车", "2024-12-17 00:00:00", "08:00", "京A12345", 50, 50, "2024-12-17 00:00:00")) // 模拟查询到一条记录
	//// 定义期望的 SQL 语句及参数
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `buses` (`origin`, `destination`, `busType`, `date`, `time`, `plate`, `total_seats`, `available_seats`, `created_at`) VALUES (?,?,?,?,?,?,?,?,?)")).
		WithArgs(
			"清水河",                 // origin
			"沙河",                  // destination
			"师生车",                 // busType
			"2024-12-17 00:00:00", // date
			"08:00",               // time
			"京A12345",             // plate
			50,                    // total_seats
			50,                    // available_seats
			"2024-12-17 00:00:00", // created_at (current time)
		).
		WillReturnResult(sqlmock.NewResult(1, 1)) // 假设插入了 1 条记录
	mock.ExpectCommit() // 模拟事务提交
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

	// 验证模拟的 SQL 查询是否正确
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unmet expectations: %s", err)
	}
}
