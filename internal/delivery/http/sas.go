package http

// import (
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	// "github.com/abstractbreazy/rate-limited-api/pkg/server"
// 	//"github.com/labstack/echo"
// 	"github.com/labstack/echo/v4"
// 	"github.com/stretchr/testify/require"
// )

// var (
// 	testIP         = "123.45.67.57"
// 	testBind       = "0.0.0.0"
// 	testPort       = "8000"
// 	testLimit      = 5
// 	testSubnetMask = 24
// 	testRedis      = "redis://localhost:6379"
// )

// func prepareServer(t *testing.T) *server.Server {

// }

// func TestPing(t *testing.T) {
// 	t.Skip()
// 	tests := []struct {
// 		ip   string
// 		code int
// 	}{
// 		{testIP, http.StatusOK},
// 		{testIP, http.StatusOK},
// 		{testIP, http.StatusOK},
// 		{testIP, http.StatusOK},
// 		{testIP, http.StatusOK},
// 		{testIP, http.StatusOK},
// 		{testIP, http.StatusOK},
// 		{testIP, http.StatusOK},
// 		{testIP, http.StatusOK},
// 		{testIP, http.StatusOK},
// 		{testIP, http.StatusOK},
// 		{testIP, http.StatusTooManyRequests},
// 		{testIP, http.StatusTooManyRequests},
// 	}
// 	for _, test := range tests {
// 		req := httptest.NewRequest(http.MethodGet, "/", nil)
// 		req.Header.Add(echo.HeaderXRealIP, test.ip)
// 		rec := httptest.NewRecorder()

// 		fmt.Println("code", rec.Code)

// 		require.Equal(t, test.code, rec.Code)
// 	}

// }
