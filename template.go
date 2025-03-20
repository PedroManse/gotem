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

type MiddlePlugin = func(w HttpWriter, r HttpReq, info map[string]any, next func(info map[string]any))
type EndPlugin = func(w HttpWriter, r HttpReq, info map[string]any)

type DynamicPage struct {
	info map[string]any
	plugins []MiddlePlugin
	end EndPlugin
}

func (dp DynamicPage) ServeHTTP(w HttpWriter, r HttpReq) {
	ps := dp.plugins
	if (len(ps) != 0) {
		var now MiddlePlugin
		now, ps = ps[0], ps[1:]
		now(
			w, r, dp.info,
			func(info map[string]any) {
				DynamicPage{info, ps, dp.end}.ServeHTTP(w, r)
			},
		)
	} else {
		dp.end(w, r, dp.info)
	}
}

func one(w HttpWriter, r HttpReq, info map[string]any, next func(info map[string]any)) {
	fmt.Println("one-b")
	if (r) {
		next(info)
		fmt.Println("one-a")
	}
}

func two(w HttpWriter, r HttpReq, info map[string]any, next func(info map[string]any)) {
	fmt.Println("two-b")
	next(info)
	fmt.Println("two-a")
}

func last(w HttpWriter, r HttpReq, info map[string]any) {
	fmt.Println("last")
}

func main() {
	var dp = DynamicPage{
		info: make(map[string]any),
		plugins: []MiddlePlugin{one, two},
		end: last,
	}
	dp.ServeHTTP(1, false)
}

