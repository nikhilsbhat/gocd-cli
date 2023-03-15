package render

import (
	"github.com/thedevsaddam/gojsonq/v2"
)

func (obj Object) GetBaseQuery() *gojsonq.JSONQ {
	return gojsonq.New().FromString(obj.String())
}

func (obj Object) GetQuery(query []string) interface{} {
	queryResponse := gojsonq.New().
		FromString(obj.String()).
		From(query[0])

	return queryResponse.Get()
}

func (obj Object) EqualQuery(query []string) interface{} {
	queryResponse := gojsonq.New().
		FromString(obj.String()).
		From(query[0]).
		WhereEqual(query[1], query[2])

	return queryResponse.Get()
}

func (obj Object) ContainsQuery(query []string) interface{} {
	queryResponse := gojsonq.New().
		FromString(obj.String()).
		From(query[0]).
		WhereContains(query[1], query[2])

	return queryResponse.Get()
}
