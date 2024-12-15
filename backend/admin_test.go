package backend

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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
	// 创建一个 gin 实例
	r := gin.Default()

	// 将 addBus 函数路由到 /addBus
	r.POST("/addBus", addBus)

	t.Run("Valid AddBus", func(t *testing.T) {
		// 构造请求数据
		input := map[string]interface{}{
			"origin":      "清水河",
			"destination": "沙河",
			"busType":     "师生车",
			"date":        "2024-12-12",
			"time":        "11:11",
			"plate":       "川A12345",
			"seats":       "40",
		}
		// 将请求数据转换成 JSON 格式
		jsonValue, _ := json.Marshal(input)

		// 创建 HTTP 请求
		req, err := http.NewRequest(http.MethodPost, "/addBus", bytes.NewReader(jsonValue))
		if err != nil {
			t.Fatal(err)
		}

		// 创建 ResponseRecorder 以便捕获响应
		w := httptest.NewRecorder()

		// 执行请求
		r.ServeHTTP(w, req)

		// 验证响应状态码和内容
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"status": "success", "message": "添加成功"}`, w.Body.String())
	})

	t.Run("Invalid BusType", func(t *testing.T) {
		// 构造错误的请求数据，传入无效的 busType
		input := map[string]interface{}{
			"origin":      "清水河",
			"destination": "沙河",
			"busType":     "无效车",
			"date":        "2024-12-12",
			"time":        "11:11",
			"plate":       "川A12345",
			"seats":       "40",
		}
		// 将请求数据转换成 JSON 格式
		jsonValue, _ := json.Marshal(input)

		// 创建 HTTP 请求
		req, err := http.NewRequest(http.MethodPost, "/addBus", bytes.NewReader(jsonValue))
		if err != nil {
			t.Fatal(err)
		}

		// 创建 ResponseRecorder 以便捕获响应
		w := httptest.NewRecorder()

		// 执行请求
		r.ServeHTTP(w, req)

		// 验证响应状态码和内容
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"status": "fail", "message": "班车类型错误！仅限为师生车或教职工车"}`, w.Body.String())
	})

	t.Run("Invalid Date", func(t *testing.T) {
		// 构造错误的请求数据，日期为过去的时间
		input := map[string]interface{}{
			"origin":      "清水河",
			"destination": "沙河",
			"busType":     "师生车",
			"date":        time.Now().Add(-24 * time.Hour).Format("2006-01-02"), // 过去的时间
			"time":        "11:11",
			"plate":       "川A12345",
			"seats":       "40",
		}
		// 将请求数据转换成 JSON 格式
		jsonValue, _ := json.Marshal(input)

		// 创建 HTTP 请求
		req, err := http.NewRequest(http.MethodPost, "/addBus", bytes.NewReader(jsonValue))
		if err != nil {
			t.Fatal(err)
		}

		// 创建 ResponseRecorder 以便捕获响应
		w := httptest.NewRecorder()

		// 执行请求
		r.ServeHTTP(w, req)

		// 验证响应状态码和内容
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"status": "fail", "message": "错误的时间！"}`, w.Body.String())
	})

	t.Run("Database Insert Fail", func(t *testing.T) {
		/*// 模拟 insertBus 返回错误
		insertBus = func(bus *Bus) error {
			return fmt.Errorf("数据库写入失败")
		}*/

		// 构造正确的请求数据
		input := map[string]interface{}{
			"origin":      "清水河",
			"destination": "沙河",
			"busType":     "师生车",
			"date":        "2024-12-12",
			"time":        "11:11",
			"plate":       "川A12345",
			"seats":       "40",
		}
		// 将请求数据转换成 JSON 格式
		jsonValue, _ := json.Marshal(input)

		// 创建 HTTP 请求
		req, err := http.NewRequest(http.MethodPost, "/addBus", bytes.NewReader(jsonValue))
		if err != nil {
			t.Fatal(err)
		}

		// 创建 ResponseRecorder 以便捕获响应
		w := httptest.NewRecorder()

		// 执行请求
		r.ServeHTTP(w, req)

		// 验证响应状态码和内容
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"status": "fail", "message": "数据库写入失败"}`, w.Body.String())
	})
}
