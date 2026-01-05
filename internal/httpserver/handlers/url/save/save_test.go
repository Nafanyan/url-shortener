package save_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"url-shortener/internal/httpserver/handlers/url/save"
	"url-shortener/internal/httpserver/handlers/url/save/mocks"
	sl "url-shortener/internal/lib/logger/handlers/slogdiscard"
	"url-shortener/internal/storage"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name         string // имя теста
		alias        string // отправляемый alias
		url          string // отправляемый url
		respError    string // какую ошибку должны получить
		respStatus   string // ожидаемый статус ответа
		expectedCode int    // ожидаемый HTTP код
		mockError    error  // ошибка, которую вернет мок
		checkAlias   bool   // нужно ли проверять alias в ответе
	}{
		{
			name:         "Success with alias",
			alias:        "test_alias",
			url:          "https://example.com",
			respError:    "",
			respStatus:   "Ok",
			expectedCode: http.StatusOK,
			mockError:    nil,
			checkAlias:   true,
		},
		{
			name:         "Success without alias (auto-generated)",
			alias:        "",
			url:          "https://google.com",
			respError:    "",
			respStatus:   "Ok",
			expectedCode: http.StatusOK,
			mockError:    nil,
			checkAlias:   true,
		},
		{
			name:         "Invalid URL format",
			url:          "not-a-url",
			alias:        "test",
			respError:    "field URL is not a valid URL",
			respStatus:   "Error",
			expectedCode: http.StatusOK,
			mockError:    nil,
			checkAlias:   false,
		},
		{
			name:         "URL already exists",
			alias:        "existing_alias",
			url:          "https://example.com",
			respError:    "url already exists",
			respStatus:   "Error",
			expectedCode: http.StatusOK,
			mockError:    storage.ErrURLExists,
			checkAlias:   false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlSaverMock := mocks.NewURLSaver(t)

			needsMock := tc.url != "" && (strings.HasPrefix(tc.url, "http://") || strings.HasPrefix(tc.url, "https://"))

			if needsMock {
				if tc.alias != "" {
					// Если alias указан, используем его
					urlSaverMock.On("SaveURL", tc.url, tc.alias).
						Return(int64(1), tc.mockError)
				} else {
					// Если alias не указан, он будет сгенерирован автоматически
					urlSaverMock.On("SaveURL", tc.url, mock.AnythingOfType("string")).
						Return(int64(1), tc.mockError)
				}
			}

			handler := save.New(sl.NewDiscardLogger(), urlSaverMock)

			var input string
			if tc.alias != "" {
				input = fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)
			} else {
				input = fmt.Sprintf(`{"url": "%s"}`, tc.url)
			}

			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code)

			body := rr.Body.String()
			var resp save.Response
			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
			require.Equal(t, tc.respStatus, resp.Status)

			if tc.checkAlias {
				if tc.alias != "" {
					require.Equal(t, tc.alias, resp.Alias)
				} else {
					// Если alias не был указан, проверяем что он был сгенерирован (длина 6)
					require.NotEmpty(t, resp.Alias)
					require.Len(t, resp.Alias, 6)
				}
			}
		})
	}
}
