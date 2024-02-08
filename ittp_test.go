package ittp_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/izaakdale/ittp"
)

type testCriteria struct {
	StatusCode int
	Request    *http.Request
}

func TestMethodHandleFunc(t *testing.T) {
	mux := ittp.NewServeMux()
	mux.MethodHandleFunc(http.MethodGet, "/GET", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("success"))
	})

	for _, c := range getTestCreteria {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, c.Request)

		if rec.Result().StatusCode != c.StatusCode {
			t.Errorf("method %s gave a status code %d but was expecting %d", c.Request.Method, rec.Result().StatusCode, c.StatusCode)
		}
	}
}

func TestMethodHandle(t *testing.T) {
	innerMux := ittp.NewServeMux()
	innerMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("success"))
	})

	outterMux := ittp.NewServeMux()
	outterMux.MethodHandle(http.MethodGet, "/GET", innerMux)

	for _, c := range getTestCreteria {
		rec := httptest.NewRecorder()
		outterMux.ServeHTTP(rec, c.Request)

		if rec.Result().StatusCode != c.StatusCode {
			t.Errorf("method %s gave a status code %d but was expecting %d", c.Request.Method, rec.Result().StatusCode, c.StatusCode)
		}
	}
}

func TestHandle(t *testing.T) {

	innerMux := ittp.NewServeMux()
	innerMux.MethodHandleFunc(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("success"))
	})

	outterMux := ittp.NewServeMux()
	outterMux.Handle("/GET", innerMux)

	for _, c := range getTestCreteria {
		rec := httptest.NewRecorder()
		outterMux.ServeHTTP(rec, c.Request)

		if rec.Result().StatusCode != c.StatusCode {
			t.Errorf("method %s gave a status code %d but was expecting %d", c.Request.Method, rec.Result().StatusCode, c.StatusCode)
		}
	}
}

type (
	ctxKey string
	ctxVal string
)

var (
	testCtxKey ctxKey = "TESTCTXKEY"
	testCtxVal ctxVal = "TESTCTXKEY"
)

func TestMiddleware(t *testing.T) {
	mux := ittp.NewServeMux()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := r.Context().Value(testCtxKey)
		if val == nil {
			t.Error("ctxkey not present")
		}
		valStr, ok := val.(ctxVal)
		if !ok {
			t.Error("not a ctx value")
		}
		if valStr != testCtxVal {
			t.Error("wrong ctx value")
		}
	})

	mux.HandleFunc("/middleware", nextHandler)

	mux.AddMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), testCtxKey, testCtxVal))
			next.ServeHTTP(w, r)
		})
	})

	req := httptest.NewRequest("GET", "/middleware", nil)
	mux.ServeHTTP(httptest.NewRecorder(), req)
}

func TestHandler(t *testing.T) {
	mux := http.NewServeMux()
	imux := ittp.NewServeMux()

	req, _ := http.NewRequest("GET", "/handler", nil)

	_, pattern := mux.Handler(req)
	_, iPattern := imux.Handler(req)

	if pattern != iPattern {
		t.Error("pattern from request handler do not match")
	}
}

func getTestMux(t *testing.T) *ittp.ServeMux {
	mux := ittp.NewServeMux()

	handlers := map[string]func(pattern string, handler http.HandlerFunc){
		http.MethodGet:     mux.Get,
		http.MethodHead:    mux.Head,
		http.MethodPost:    mux.Post,
		http.MethodPut:     mux.Put,
		http.MethodPatch:   mux.Patch,
		http.MethodDelete:  mux.Delete,
		http.MethodConnect: mux.Connect,
		http.MethodOptions: mux.Options,
		http.MethodTrace:   mux.Trace,
	}

	for method, handler := range handlers {
		path := fmt.Sprintf("/%s", method)
		handler(path, func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("success"))
		})
	}

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Path: "/notfound"},
	})

	if rec.Result().StatusCode != http.StatusNotFound {
		t.Errorf("handler not reporting not found")
	}

	return mux
}
func TestMethodFuncs(t *testing.T) {
	mux := getTestMux(t)

	criteria := []testCriteria{}
	criteria = append(criteria, getTestCreteria...)
	criteria = append(criteria, headTestCreteria...)
	criteria = append(criteria, postTestCreteria...)
	criteria = append(criteria, putTestCreteria...)
	criteria = append(criteria, patchTestCreteria...)
	criteria = append(criteria, deleteTestCreteria...)
	criteria = append(criteria, connectTestCreteria...)
	criteria = append(criteria, optionsTestCreteria...)
	criteria = append(criteria, traceTestCreteria...)

	for _, c := range criteria {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, c.Request)

		if rec.Result().StatusCode != c.StatusCode {
			t.Errorf("method %s gave a status code %d but was expecting %d", c.Request.Method, rec.Result().StatusCode, c.StatusCode)
		}
	}
}

var getTestCreteria = []testCriteria{
	{
		StatusCode: http.StatusOK,
		Request: &http.Request{
			URL:    &url.URL{Path: "/GET"},
			Method: http.MethodGet,
		},
	},
	{
		StatusCode: http.StatusMethodNotAllowed,
		Request: &http.Request{
			URL:    &url.URL{Path: "/GET"},
			Method: http.MethodPost,
		},
	},
}
var headTestCreteria = []testCriteria{
	{
		StatusCode: http.StatusOK,
		Request: &http.Request{
			URL:    &url.URL{Path: "/HEAD"},
			Method: http.MethodHead,
		},
	},
	{
		StatusCode: http.StatusMethodNotAllowed,
		Request: &http.Request{
			URL:    &url.URL{Path: "/HEAD"},
			Method: http.MethodPost,
		},
	},
}
var postTestCreteria = []testCriteria{
	{
		StatusCode: http.StatusOK,
		Request: &http.Request{
			URL:    &url.URL{Path: "/POST"},
			Method: http.MethodPost,
		},
	},
	{
		StatusCode: http.StatusMethodNotAllowed,
		Request: &http.Request{
			URL:    &url.URL{Path: "/POST"},
			Method: http.MethodGet,
		},
	},
}
var putTestCreteria = []testCriteria{
	{
		StatusCode: http.StatusOK,
		Request: &http.Request{
			URL:    &url.URL{Path: "/PUT"},
			Method: http.MethodPut,
		},
	},
	{
		StatusCode: http.StatusMethodNotAllowed,
		Request: &http.Request{
			URL:    &url.URL{Path: "/PUT"},
			Method: http.MethodGet,
		},
	},
}
var patchTestCreteria = []testCriteria{
	{
		StatusCode: http.StatusOK,
		Request: &http.Request{
			URL:    &url.URL{Path: "/PATCH"},
			Method: http.MethodPatch,
		},
	},
	{
		StatusCode: http.StatusMethodNotAllowed,
		Request: &http.Request{
			URL:    &url.URL{Path: "/PATCH"},
			Method: http.MethodGet,
		},
	},
}
var deleteTestCreteria = []testCriteria{
	{
		StatusCode: http.StatusOK,
		Request: &http.Request{
			URL:    &url.URL{Path: "/DELETE"},
			Method: http.MethodDelete,
		},
	},
	{
		StatusCode: http.StatusMethodNotAllowed,
		Request: &http.Request{
			URL:    &url.URL{Path: "/DELETE"},
			Method: http.MethodGet,
		},
	},
}
var connectTestCreteria = []testCriteria{
	{
		StatusCode: http.StatusOK,
		Request: &http.Request{
			URL:    &url.URL{Path: "/CONNECT"},
			Method: http.MethodConnect,
		},
	},
	{
		StatusCode: http.StatusMethodNotAllowed,
		Request: &http.Request{
			URL:    &url.URL{Path: "/CONNECT"},
			Method: http.MethodGet,
		},
	},
}
var optionsTestCreteria = []testCriteria{
	{
		StatusCode: http.StatusOK,
		Request: &http.Request{
			URL:    &url.URL{Path: "/OPTIONS"},
			Method: http.MethodOptions,
		},
	},
	{
		StatusCode: http.StatusMethodNotAllowed,
		Request: &http.Request{
			URL:    &url.URL{Path: "/OPTIONS"},
			Method: http.MethodGet,
		},
	},
}
var traceTestCreteria = []testCriteria{
	{
		StatusCode: http.StatusOK,
		Request: &http.Request{
			URL:    &url.URL{Path: "/TRACE"},
			Method: http.MethodTrace,
		},
	},
	{
		StatusCode: http.StatusMethodNotAllowed,
		Request: &http.Request{
			URL:    &url.URL{Path: "/TRACE"},
			Method: http.MethodGet,
		},
	},
}
