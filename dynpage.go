package temṕlate

import (
	"net/http"
)
	//"html/template"
	// unsafeTemplate "text/template"
	// "io"

type HttpWriter = http.ResponseWriter
type HttpReq = *http.Request

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

func buildDynPage[I any](
	info I,
	last EndPlugin[I],
	plugins ...MiddlePlugin[I],
) DynamicPage[I] {
	return DynamicPage[I]{
		info,
		plugins,
		last,
	}
}

