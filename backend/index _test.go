package backend

import (
	"database/sql/driver"
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

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	// Match 方法中：判断字段值只要是time.Time 类型，就能验证通过
	_, ok := v.(time.Time)
	return ok
}

func TestQueryAll(t *testing.T) {
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
	router.POST("/index/queryAll", queryAll)

	// 输入的日期字符串
	dateStr := "2024-12-16"

	// 定义时间格式
	layout := "2006-01-02" // Go的时间格式必须使用固定的"2006-01-02"模式

	// 使用 time.Parse 解析日期字符串为 time.Time 类型
	parsedTime, err := time.Parse(layout, dateStr)

	// 模拟数据：假设查询结果为两条符合条件的 Bus 记录
	busRows := sqlmock.NewRows([]string{"busId", "available_seats", "busType", "date"}).
		AddRow(1, 10, "师生车", parsedTime).
		AddRow(2, 5, "教职工车", parsedTime)

	// 设置模拟查询：查询 buses 表
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `buses` WHERE date =? AND busType LIKE?")).
		WithArgs(parsedTime, "%师生车%").
		WillReturnRows(busRows)

	// 准备请求数据
	reqBody := `{"date": "2024-12-16", "usertype": "学生"}`
	req := httptest.NewRequest(http.MethodPost, "/index/queryAll", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// 创建响应记录器
	rec := httptest.NewRecorder()

	// 调用 queryAll 函数
	router.ServeHTTP(rec, req)

	// 断言结果
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"bustype":"师生车"`)
	assert.Contains(t, rec.Body.String(), `"bustype":"教职工车"`)
	//assert.Contains(t, rec.Body.String(), `"date":"2024-12-16"`)
}

func TestBook(t *testing.T) {
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

	// 创建 Gin 路由实例
	router := gin.Default()
	router.POST("/index/book", book)

	bookrows := sqlmock.NewRows([]string{"id", "userId", "busId", "status"}).
		AddRow(1, "user1", 1, "已预约") // 模拟一条已预约的记录

	// 设置模拟查询：查询 bookings 表
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `bookings` WHERE userId =? AND busId =? AND status =? LIMIT ?")).
		WithArgs("user1", 1, "已预约", 1).
		WillReturnRows(bookrows)

	// 模拟查询 buses 表（只查询 busId 和 availableSeats）
	rows := sqlmock.NewRows([]string{"busId", "availableSeats"}).
		AddRow(1, 10)

	// 设置模拟查询：查询 buses 表
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `buses` WHERE busId =? LIMIT ?")).
		WithArgs(1, 1).
		WillReturnRows(rows)

	// 准备请求数据
	reqBody := `{"userid": "user1", "busid": 1}`
	req := httptest.NewRequest(http.MethodPost, "/index/book", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// 创建响应记录器
	rec := httptest.NewRecorder()

	// 调用 book 函数
	router.ServeHTTP(rec, req)

	// 断言结果
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"status":"booked"`)
}

func TestPayed(t *testing.T) {
	// 使用 sqlmock 创建一个 mock 数据库连接
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
	router.POST("/index/payed", payed)

	busss := Bus{
		Origin:         "",
		Destination:    "",
		BusType:        "",
		AvailableSeats: 9,
	}
	/*bookinggg := Booking{
		UserId: "123",
		BusId:  1,
		Status: "已预约",
	}*/
	// 设置 mock 查询行为，模拟班车信息
	rows := sqlmock.NewRows([]string{"busId", "available_seats"}).
		AddRow(1, 10)

	mock.ExpectQuery("SELECT * FROM `buses` WHERE busId =? LIMIT ?").
		WithArgs(1, 1).
		WillReturnRows(rows)
	mock.ExpectBegin()
	// 模拟更新操作
	mock.ExpectExec("UPDATE `buses` SET `origin`=?,`destination`=?,`busType`=?,`date`=?,`time`=?,`plate`=?,`total_seats`=?,`available_seats`=?,`created_at`=? WHERE `busId` = ?").
		WithArgs(busss.Origin, busss.Destination, busss.BusType, busss.Date, busss.Time, busss.Plate, busss.TotalSeats, 9, busss.CreatedAt, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectBegin()                                                                                   // 开始事务
	mock.ExpectExec("INSERT INTO `bookings` (`userId`,`busId`,`status`,`created_at`) VALUES (?,?,?,?)"). // 模拟插入数据
														WithArgs("123", 1, "已预约", time.Now().Truncate(time.Second)).
														WillReturnResult(sqlmock.NewResult(1, 1)) // 返回影响的行数
	mock.ExpectCommit()                                                                                            // 提交事务
	mock.ExpectBegin()                                                                                             // 开始事务
	mock.ExpectExec("UPDATE `bookings` SET `userId`=?,`busId`=?,`status`=?,`created_at`=? WHERE `bookingId` = ?"). // 模拟插入数据
															WithArgs("123", 1, "已预约", time.Now().Truncate(time.Second), 1).
															WillReturnResult(sqlmock.NewResult(1, 1)) // 返回影响的行数
	mock.ExpectCommit() // 提交事务
	// 准备请求数据
	reqBody := `{"userid": "123", "busid": 1}`
	req := httptest.NewRequest(http.MethodPost, "/index/payed", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// 创建响应记录器
	rec := httptest.NewRecorder()

	// 调用 payed 函数
	router.ServeHTTP(rec, req)

	// 断言结果
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"status":"success"`)
}

func TestUnbook(t *testing.T) {
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
	router.POST("/index/unbook", unbook)

	// 设置 mock 查询行为
	rows := sqlmock.NewRows([]string{"userId", "busId", "status"}).
		AddRow("user1", 1, "已预约")

	mock.ExpectQuery("SELECT * FROM `bookings` WHERE userId =? AND busId =? LIMIT ?").
		WithArgs("user1", 1, 1).
		WillReturnRows(rows)
	rows = sqlmock.NewRows([]string{"userId", "busId", "status"}).
		AddRow("user1", 1, "已预约")
	mock.ExpectQuery("SELECT * FROM `bookings` WHERE userId =? AND busId =? AND status =? LIMIT ?").
		WithArgs("user1", 1, "已预约", 1).
		WillReturnRows(rows)

	// 模拟更新操作
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `bookings` (`userId`,`busId`,`status`,`created_at`) VALUES (?,?,?,?)").
		WithArgs("user1", 1, "已取消", time.Now().Truncate(time.Second)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery("SELECT * FROM `buses` WHERE busId =? LIMIT ?").
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"busId", "available_seats"}).
			AddRow(1, 9))
	mock.ExpectBegin()
	// 模拟班车座位数更新
	busss := Bus{
		Origin:         "",
		Destination:    "",
		BusType:        "",
		AvailableSeats: 9,
	}
	mock.ExpectExec("UPDATE `buses` SET `origin`=?,`destination`=?,`busType`=?,`date`=?,`time`=?,`plate`=?,`total_seats`=?,`available_seats`=?,`created_at`=? WHERE `busId` = ?").
		WithArgs(busss.Origin, busss.Destination, busss.BusType, busss.Date, busss.Time, busss.Plate, busss.TotalSeats, 10, busss.CreatedAt, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// 准备请求数据
	reqBody := `{"userid": "user1", "busid": 1}`
	req := httptest.NewRequest(http.MethodPost, "/index/unbook", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// 创建响应记录器
	rec := httptest.NewRecorder()

	// 调用 unbook 函数
	router.ServeHTTP(rec, req)

	// 断言结果
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"status":"success"`)
}

func TestQueryBooked(t *testing.T) {
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

	// 创建 Gin 路由实例
	router := gin.Default()
	router.POST("/index/queryBooked", queryBooked)

	// 设置 mock 查询行为，模拟查询已预约或已支付的预约记录
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `bookings` WHERE userId =? AND (status = '已预约' OR status = '已支付')")).
		WithArgs("user1").
		WillReturnRows(sqlmock.NewRows([]string{"userId", "busId", "status"}).
			AddRow("user1", 1, "已预约").
			AddRow("user1", 2, "已支付"))

	// 设置 mock 查询班车信息
	busRows := sqlmock.NewRows([]string{"busId", "origin", "destination", "busType", "availableSeats", "date"}).
		AddRow(1, "Origin A", "Destination B", "师生车", 10, time.Now()).
		AddRow(2, "Origin C", "Destination D", "教职工车", 5, time.Now())

	// 这里对于 busId=1 和 busId=2 都进行模拟返回
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `buses` WHERE busId =? LIMIT ?")).
		WithArgs(1, 1).
		WillReturnRows(busRows)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `buses` WHERE busId =? LIMIT ?")).
		WithArgs(2, 1).
		WillReturnRows(busRows)

	// 准备请求数据
	reqBody := `{"userid": "user1"}`
	req := httptest.NewRequest(http.MethodPost, "/index/queryBooked", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// 创建响应记录器
	rec := httptest.NewRecorder()

	// 调用 queryBooked 函数
	router.ServeHTTP(rec, req)

	// 断言结果
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"busid":1`)
	assert.Contains(t, rec.Body.String(), `"origin":"Origin A"`)
	assert.Contains(t, rec.Body.String(), `"destination":"Destination B"`)
	assert.Contains(t, rec.Body.String(), `"busid":2`)
	assert.Contains(t, rec.Body.String(), `"origin":"Origin C"`)
}
func TestQueryFinished(t *testing.T) {
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

	// 创建 Gin 路由实例
	router := gin.Default()
	RegisterIndexModule(router)

	// 设置 mock 查询行为，模拟查询已完成或已取消的预约记录
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `bookings` WHERE userId =? AND (status = '已完成' OR status = '已取消')")).
		WithArgs("user1").
		WillReturnRows(sqlmock.NewRows([]string{"userId", "busId", "status"}).
			AddRow("user1", 1, "已完成").
			AddRow("user1", 2, "已取消"))

	// 设置 mock 查询班车信息
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `buses` WHERE busId =? LIMIT ?")).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"busId", "origin", "destination", "busType", "availableSeats", "date"}).
			AddRow(1, "Origin A", "Destination B", "师生车", 10, time.Now()))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `buses` WHERE busId =? LIMIT ?")).
		WithArgs(2, 1).
		WillReturnRows(sqlmock.NewRows([]string{"busId", "origin", "destination", "busType", "availableSeats", "date"}).
			AddRow(2, "Origin C", "Destination D", "教职工车", 5, time.Now()))

	// 准备请求数据
	reqBody := `{"userid": "user1"}`
	req := httptest.NewRequest(http.MethodPost, "/index/queryFinished", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// 创建响应记录器
	rec := httptest.NewRecorder()

	// 调用 queryFinished 函数
	router.ServeHTTP(rec, req)

	// 断言结果
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"busid":1`)
	assert.Contains(t, rec.Body.String(), `"busid":2`)

	// 确保所有预期的操作都已执行
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("There were unmet expectations: %s", err)
	}
}
