package serviceorderhandler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthorizeDocument(t *testing.T) {
	cases := []struct {
		name             string
		tokenDocument    string
		tokenRole        string
		requestedDoc     string
		expectedAuthized bool
		expectedStatus   int
	}{
		{"admin sem documento no token libera", "", "admin", "52998224725", true, http.StatusOK},
		{"cliente sem documento no token bloqueia", "", "client", "52998224725", false, http.StatusForbidden},
		{"token sem role nem documento bloqueia", "", "", "52998224725", false, http.StatusForbidden},
		{"documento igual libera", "52998224725", "client", "52998224725", true, http.StatusOK},
		{"documento diferente bloqueia", "52998224725", "client", "39053344705", false, http.StatusForbidden},
		{"documento pedido vazio bloqueia", "52998224725", "client", "", false, http.StatusForbidden},
	}

	gin.SetMode(gin.TestMode)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			if tc.tokenDocument != "" {
				c.Set("user_document", tc.tokenDocument)
			}
			if tc.tokenRole != "" {
				c.Set("user_role", tc.tokenRole)
			}

			if got := authorizeDocument(c, tc.requestedDoc); got != tc.expectedAuthized {
				t.Fatalf("esperava %v, veio %v", tc.expectedAuthized, got)
			}

			if !tc.expectedAuthized && rec.Code != tc.expectedStatus {
				t.Errorf("esperava status %d, veio %d", tc.expectedStatus, rec.Code)
			}
		})
	}
}
