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

// MockDB 结构体模拟数据库
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

// queryUser 函数，模拟数据库查询
/*func queryUser(user *User) (*User, error) {
	var queriedUser User
	// 这里用 db 的 Where 和 Take 方法模拟查询
	result := db.Where("username = ?", user.UserName).Take(&queriedUser)
	if result.Error != nil {
		return &User{}, result.Error
	}
	return &queriedUser, nil
}*/

// 测试函数
/*func TestQueryUser(t *testing.T) {
	// 创建 MockDB 实例
	mdb := new(MockDB)

	// 创建一个模拟的 User
	mockUser := &User{
		UserID:   "123",
		UserName: "user123",
		Password: "password123",
	}

	// 设置期望：调用 db.Where 和 db.Take 时返回 mockUser
	mdb.On("Where", "username = ?", "user123").Return(mdb) // 模拟查询条件
	mdb.On("Take", mock.Anything).Run(func(args mock.Arguments) {
		// 模拟 db.Take 将 mockUser 填充到 dest 参数中
		*args.Get(0).(*User) = *mockUser
	}).Return(mdb)

	// 使用 mockDB 替代真实 db
	//db = mdb

	// 执行 queryUser，应该返回 mockUser
	user := &User{UserName: "user123"}
	queriedUser, err := queryUser(user)

	// 断言结果
	assert.NoError(t, err)
	assert.Equal(t, mockUser.UserID, queriedUser.UserID)
	assert.Equal(t, mockUser.UserName, queriedUser.UserName)
	assert.Equal(t, mockUser.Password, queriedUser.Password)

	// 验证期望的调用
	mdb.AssertExpectations(t)
}

/*
// 模拟的 shaEncode 函数
func shaEncode(password string) string {
	if password == "123456Aa" {
		return "hashed_password" // 模拟加密后的密码
	}
	return ""
}

// 模拟的 updateUserIsFirstLogin 函数
func updateUserIsFirstLogin(user *User) int {
	if user.Is_first_login {
		return 1
	}
	return 0
}

// 模拟的 queryUser 函数
func queryUser(user *User) (*User, error) {
	// 这里返回一个固定的用户，模拟数据库查询
	if user.UserName == "user123" {
		return &User{
			UserID:         "1",
			UserType:       "admin",
			Password:       "hashed_password",
			Is_first_login: true,
			CreatedTime:    time.Now(),
		}, nil
	}
	return nil, errors.New("user not found")
}
*/
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
		assert.Equal(t, 1, response["isfirstlogin"])
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
