package api

import (
	"fmt"

	"expe-plugin/api/server/restapi/operations"

	"github.com/bluele/gcache"
	"github.com/go-openapi/runtime/middleware"
)

func set(gc gcache.Cache, val []byte) {
	err := gc.Set("key", val)
	if err != nil {
		panic(err)
	}
}

func get(gc gcache.Cache) []byte {
	value, err := gc.Get("key")
	if err != nil {
		fmt.Println("not found")
		return nil
	}

	valueBytes := value.([]byte)
	vc := make([]byte, len(valueBytes))
	copy(vc, valueBytes)

	return vc
}

func NewExpe(gc gcache.Cache) operations.ExpeHandler {
	return &ExpeEndpoint{gc: gc}
}

type ExpeEndpoint struct {
	gc gcache.Cache
}

func (e *ExpeEndpoint) Handle(params operations.ExpeParams) middleware.Responder {
	value := []byte{1, 2, 3, 4, 5}

	val := get(e.gc)
	if val == nil {
		set(e.gc, value)
		return operations.NewExpeOK().WithPayload(
			&operations.ExpeOKBody{Message: "val added to cache"},
		)
	}

	return operations.NewExpeOK().WithPayload(
		&operations.ExpeOKBody{Message: fmt.Sprintf("val is %v", val)},
	)
}
