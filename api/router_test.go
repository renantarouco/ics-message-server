package api

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

var routesAllowedMethods = map[string][]string{
	"join": []string{"POST"},
}

func TestRoutesAllowedMethods(t *testing.T) {
	err := APIRouter.Walk(func(route *mux.Route, router *mux.Router, acestors []*mux.Route) error {
		routeMethods, err := route.GetMethods()
		if err != nil {
			return err
		}
		allowedMethods := routesAllowedMethods[route.GetName()]
		if len(routeMethods) != len(allowedMethods) {
			return errors.New("different allowed methods count")
		}
		methodsMap := map[string]int{}
		for _, method := range routeMethods {
			methodsMap[method]++
		}
		for _, method := range allowedMethods {
			v, ok := methodsMap[method]
			if !ok {
				return fmt.Errorf("route %s should allow %s", route.GetName(), method)
			}
			v--
			if v == 0 {
				delete(methodsMap, method)
			}
			if len(methodsMap) > 0 {
				extraMethods := []string{}
				for k := range methodsMap {
					extraMethods = append(extraMethods, k)
				}
				return fmt.Errorf("route %s should not allow %s", route.GetName(), strings.Join(extraMethods, ","))
			}
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}
