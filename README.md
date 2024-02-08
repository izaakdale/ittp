# ittp
ittp is a simple yet effective wrapper around net/http.ServeMux, intended to be used completely interchangably with it.


## Usage/Examples


```bash
  go get github.com/izaakdale/ittp
```

main.go

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/izaakdale/ittp"
)

func main() {
	mux := ittp.NewServeMux()

	// can still use the old way of registering methods
	mux.Handle("/oldway", http.HandlerFunc(pingpong))
	mux.HandleFunc("/oldwayFunc", pingpong)

	mux.Get("/ping", pingpong)
	mux.Post("/ping", pingpong)
	// mux.Head/Options/Trace.....
	mux.MethodHandleFunc(http.MethodPut, "/ping", pingpong)
	mux.MethodHandle(http.MethodPatch, "/ping", http.HandlerFunc(pingpong))

	// middleware is executed in order they are added
	mux.AddMiddleware(pingpongMiddleware1)
	mux.AddMiddleware(pingpongMiddleware2)

	http.ListenAndServe("localhost:8080", mux)
}

func pingpong(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func pingpongMiddleware1(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("ping1")
		next.ServeHTTP(w, r)
		fmt.Println("pong1")
	})
}
func pingpongMiddleware2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("ping2")
		next.ServeHTTP(w, r)
		fmt.Println("pong2")
	})
}

```

