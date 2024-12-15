package backend

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockDB 用于模拟数据库
type MockDB struct {
	mock.Mock
}

func (mdb *MockDB) Where(query string, args ...interface{}) *MockDB {
	mdb.Called(query, args)
	return mdb
}

func (mdb *MockDB) Take(dest interface{}) *MockDB {
	mdb.Called(dest)
	return mdb
}

func (mdb *MockDB) Update(column string, value interface{}) *MockDB {
	mdb.Called(column, value)
	return mdb
}

// 单元测试
func TestUserLogin(t *testing.T) {
	// 初始化 Gin engine
	r := gin.Default()

	// 注册登录接口
	r.POST("/login", userLogin)

	// 正常登录测试
	t.Run("Valid Login", func(t *testing.T) {
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
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err) // 确保 JSON 解析没有错误

		// 验证响应内容
		assert.Equal(t, "success", response["status"])
		assert.Equal(t, "1", response["userid"])
		assert.Equal(t, "admin", response["usertype"])
		assert.Equal(t, float64(1), response["isfirstlogin"])
	})

	/*// 错误的密码测试
	t.Run("Invalid Password", func(t *testing.T) {
		input := map[string]string{
			"username": "user123",
			"password": "wrongpassword",
		}

		jsonValue, _ := json.Marshal(input)
		req, _ := http.NewRequest("POST", "/login", bytes.NewReader(jsonValue))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// 断言响应状态码为 400
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "fail", response["status"])
	})

	// 用户不存在测试
	t.Run("User Not Found", func(t *testing.T) {
		input := map[string]string{
			"username": "nonexistentUser",
			"password": "123456Aa",
		}

		jsonValue, _ := json.Marshal(input)
		req, _ := http.NewRequest("POST", "/login", bytes.NewReader(jsonValue))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// 断言响应状态码为 500
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "fail", response["status"])
	})

	// 用户名为空测试
	t.Run("Empty Username", func(t *testing.T) {
		input := map[string]string{
			"username": "",
			"password": "123456Aa",
		}

		jsonValue, _ := json.Marshal(input)
		req, _ := http.NewRequest("POST", "/login", bytes.NewReader(jsonValue))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// 断言响应状态码为 500
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "fail", response["status"])
	})*/
}
