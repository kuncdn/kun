package cors

import (
	"errors"
	"net/http"

	"github.com/justinas/alice"
	"tracfox.io/tracfox/internal/util"
)

const name = "cors"

type corsPlugin struct {
	allowOrigin  string
	allowMethods string
	allowHeaders string
}

type cors struct {
	plug *corsPlugin
	next http.Handler
}

func (c *cors) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", c.plug.allowOrigin)
	rw.Header().Set("Access-Control-Allow-Methods", c.plug.allowMethods)
	rw.Header().Set("Access-Control-Allow-Headers", c.plug.allowHeaders)
	if req.Method == "OPTIONS" {
		return
	}
	c.next.ServeHTTP(rw, req)
}

// Constructor .
func Constructor(args map[string]interface{}) (alice.Constructor, error) {
	var errs []error
	_, ok := args["allowOrigin"]
	if !ok {
		errs = append(errs, errors.New("cors plugin allowOrigin field is required"))
	}

	allowOrigin, ok := args["allowOrigin"].(string)
	if !ok {
		errs = append(errs, errors.New("cors plugin allowOrigin must be string"))
	}

	_, ok = args["allowMethods"]
	if !ok {
		errs = append(errs, errors.New("cors plugin allowMethods field is required"))
	}
	allowMethods, ok := args["allowMethods"].(string)
	if !ok {
		errs = append(errs, errors.New("cors plugin allowMethods must be string"))
	}

	_, ok = args["allowHeaders"]
	if !ok {
		errs = append(errs, errors.New("cors plugin allowHeaders field is required"))
	}
	allowHeaders, ok := args["allowHeaders"].(string)
	if !ok {
		errs = append(errs, errors.New("cors plugin allowHeaders must be string"))
	}

	if len(errs) != 0 {
		return nil, util.NewAggregateError(errs)
	}
	return func(next http.Handler) http.Handler {
		return &cors{
			next: next,
			plug: &corsPlugin{
				allowOrigin:  allowOrigin,
				allowMethods: allowMethods,
				allowHeaders: allowHeaders,
			},
		}
	}, nil
}
