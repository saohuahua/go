package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTodoAPI(t *testing.T) {
	api := todoAPI(newTodoStore())

	t.Run("创建后可查询和删除", func(t *testing.T) {
		created := performRequest(api, http.MethodPost, "/todos", `{"text":"学 net/http"}`)
		if created.Code != http.StatusCreated {
			t.Fatalf("创建状态码 = %d，want %d", created.Code, http.StatusCreated)
		}
		if !strings.Contains(created.Body.String(), `"id":1`) {
			t.Fatalf("创建响应 = %s，缺少 id", created.Body.String())
		}

		listed := performRequest(api, http.MethodGet, "/todos", "")
		if listed.Code != http.StatusOK || !strings.Contains(listed.Body.String(), "学 net/http") {
			t.Fatalf("列表响应 = %d %s", listed.Code, listed.Body.String())
		}

		deleted := performRequest(api, http.MethodDelete, "/todos?id=1", "")
		if deleted.Code != http.StatusNoContent {
			t.Fatalf("删除状态码 = %d，want %d", deleted.Code, http.StatusNoContent)
		}

		missing := performRequest(api, http.MethodGet, "/todos?id=1", "")
		if missing.Code != http.StatusNotFound {
			t.Fatalf("删除后查询状态码 = %d，want %d", missing.Code, http.StatusNotFound)
		}
	})

	t.Run("无效输入返回 400", func(t *testing.T) {
		invalidJSON := performRequest(api, http.MethodPost, "/todos", `{`)
		if invalidJSON.Code != http.StatusBadRequest {
			t.Fatalf("无效 JSON 状态码 = %d，want %d", invalidJSON.Code, http.StatusBadRequest)
		}

		emptyText := performRequest(api, http.MethodPost, "/todos", `{"text":"   "}`)
		if emptyText.Code != http.StatusBadRequest {
			t.Fatalf("空 text 状态码 = %d，want %d", emptyText.Code, http.StatusBadRequest)
		}

		invalidID := performRequest(api, http.MethodGet, "/todos?id=zero", "")
		if invalidID.Code != http.StatusBadRequest {
			t.Fatalf("无效 id 状态码 = %d，want %d", invalidID.Code, http.StatusBadRequest)
		}
	})

	t.Run("路由和方法错误分别返回 404 与 405", func(t *testing.T) {
		missingRoute := performRequest(api, http.MethodGet, "/missing", "")
		if missingRoute.Code != http.StatusNotFound {
			t.Fatalf("不存在路由状态码 = %d，want %d", missingRoute.Code, http.StatusNotFound)
		}

		wrongMethod := performRequest(api, http.MethodPut, "/todos", "")
		if wrongMethod.Code != http.StatusMethodNotAllowed {
			t.Fatalf("错误方法状态码 = %d，want %d", wrongMethod.Code, http.StatusMethodNotAllowed)
		}
		if got := wrongMethod.Header().Get("Allow"); got != "GET, POST, DELETE" {
			t.Fatalf("Allow = %q，want %q", got, "GET, POST, DELETE")
		}
	})
}

func TestTodoAPIConcurrentCreate(t *testing.T) {
	api := todoAPI(newTodoStore())
	const requestCount = 50

	results := make(chan *httptest.ResponseRecorder, requestCount)
	for range requestCount {
		go func() {
			results <- performRequest(api, http.MethodPost, "/todos", `{"text":"并发创建"}`)
		}()
	}
	for range requestCount {
		response := <-results
		if response.Code != http.StatusCreated {
			t.Errorf("并发创建状态码 = %d，want %d", response.Code, http.StatusCreated)
		}
	}
}

func performRequest(handler http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, "http://example.com"+path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}
