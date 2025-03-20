package main

import (
	"fmt"
)
	//"net/http"
	//"html/template"
	// unsafeTemplate "text/template"
	// "io"

//type HttpWriter = http.ResponseWriter
//type HttpReq = *http.Request
type HttpWriter = int
type HttpReq = bool

type MiddlePlugin[I any] = func(w HttpWriter, r HttpReq, info I, next func(info I))
type EndPlugin[I any] = func(w HttpWriter, r HttpReq, info I)

type DynamicPage[I any] struct {
	info I
	plugins []MiddlePlugin[I]
	end EndPlugin[I]
}

func (dp DynamicPage[I]) ServeHTTP(w HttpWriter, r HttpReq) {
	ps := dp.plugins
	if (len(ps) != 0) {
		var now MiddlePlugin[I]
		now, ps = ps[0], ps[1:]
		now(
			w, r, dp.info,
			func(info I) {
				DynamicPage[I]{info, ps, dp.end}.ServeHTTP(w, r)
			},
		)
	} else {
		dp.end(w, r, dp.info)
	}
}

func one[I any](w HttpWriter, r HttpReq, info I, next func(info I)) {
	fmt.Println("one-b")
	if (r) {
		next(info)
		fmt.Println("one-a")
	}
}

func two[I any](w HttpWriter, r HttpReq, info I, next func(info I)) {
	fmt.Println("two-b")
	next(info)
	fmt.Println("two-a")
}

func last[I any](w HttpWriter, r HttpReq, info I) {
	fmt.Println("last")
}

func main() {
	var dp = DynamicPage[string]{
		info: "hello",
		plugins: []MiddlePlugin[string]{one[string], two[string]},
		end: last[string],
	}
	dp.ServeHTTP(1, true)
}

