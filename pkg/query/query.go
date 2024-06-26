package query

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/nikhilsbhat/common/content"
	"github.com/thedevsaddam/gojsonq/v2"
)

const (
	defaultLengthGetOrPluck = 1
	defaultLengthSort       = 2
	defaultLengthWhere      = 3
)

// Query holds information of queries to be executed.
type Query struct {
	object    string
	operator  string
	key       string
	value     string
	queryType string
	JSONQ     *gojsonq.JSONQ
}

// SetQuery processes the data into a format FromString of gojsonq understands and return finally processed query.
func SetQuery(data interface{}, query string) (*Query, error) {
	objectString, err := content.Marshal(data)
	if err != nil {
		return nil, err
	}

	queryObj := &Query{
		JSONQ: getBaseQuery(objectString),
	}

	queryObj.ConstructQuery(query)

	return queryObj, nil
}

// ConstructQuery process the query in string to a format that is understood by gojsonq.
func (q *Query) ConstructQuery(query string) {
	queryLHS := strings.Split(query, "|")
	queryLHSObject := strings.TrimRightFunc(queryLHS[0], unicode.IsSpace)

	switch len(queryLHS) {
	case defaultLengthGetOrPluck:
		q.object = queryLHSObject
		q.queryType = "get"
	case defaultLengthSort:
		q.object = queryLHSObject
		queryRHS := strings.TrimLeftFunc(queryLHS[1], unicode.IsSpace)
		queryRHSObject := strings.Split(queryRHS, " ")

		switch len(queryRHSObject) {
		case defaultLengthGetOrPluck:
			q.key = queryRHSObject[0]
			q.queryType = "pluck"
		case defaultLengthSort:
			q.key = queryRHSObject[0]
			q.value = queryRHSObject[1]
			q.queryType = "sort"
		case defaultLengthWhere:
			q.key = queryRHSObject[0]
			q.operator = Operator(queryRHSObject[1])
			q.value = queryRHSObject[2]
			q.queryType = "where"
		}
	}
}

func getBaseQuery(obj content.Object) *gojsonq.JSONQ {
	return gojsonq.New().FromString(obj.String())
}

// QueryGet returns response post the GET query.
func (q *Query) QueryGet() interface{} {
	queryResponse := q.JSONQ.From(q.object)

	return queryResponse.Get()
}

// QueryPluck returns response post the PLUCK query.
func (q *Query) QueryPluck() interface{} {
	queryResponse := q.JSONQ.From(q.object).Pluck(q.key)

	return queryResponse
}

// QueryWhere returns response post the WHERE query.
func (q *Query) QueryWhere() interface{} {
	queryResponse := q.JSONQ.From(q.object).Where(q.key, q.operator, q.value)

	return queryResponse.Get()
}

// QuerySort returns response post the SORT query.
func (q *Query) QuerySort() interface{} {
	queryResponse := q.JSONQ.From(q.object).SortBy(q.key, q.value)

	return queryResponse.Get()
}

// GetQueryType returns the type of query set.
func (q *Query) GetQueryType() string {
	return q.queryType
}

// GetQueryKey returns the key used for query if set.
func (q *Query) GetQueryKey() string {
	return q.key
}

// GetQueryValue returns the value used for query if set.
func (q *Query) GetQueryValue() string {
	return q.value
}

// GetQueryOperator returns the operator used for query if exists.
func (q *Query) GetQueryOperator() string {
	return q.operator
}

// GetQueryObject returns the object used for query if set.
func (q *Query) GetQueryObject() string {
	return q.object
}

// RunQuery identifies the query type based on the type set by ConstructQuery and return the result post applying the query.
func (q *Query) RunQuery() interface{} {
	switch q.queryType {
	case "get":
		return q.QueryGet()
	case "sort":
		return q.QuerySort()
	case "where":
		return q.QueryWhere()
	case "pluck":
		return q.QueryPluck()
	default:
		return nil
	}
}

// Print prints the Query object to string format, mostly used for debug message.
func (q *Query) Print() string {
	return fmt.Sprintf("processed query: 'object:%s type:%s ,key: %s ,value:%s ,operator:%s'",
		q.GetQueryObject(),
		q.GetQueryType(),
		q.GetQueryKey(),
		q.GetQueryValue(),
		q.GetQueryOperator())
}

// Operator identifies the operator passed in the query and matches to one from https://github.com/thedevsaddam/gojsonq/wiki/Queries#wherekey-op-val.
func Operator(query string) string {
	switch strings.ToLower(query) {
	case "equal", "eq", "=":
		return "eq"
	case "notequal", "nq", "neq", "!=":
		return "neq"
	case "gt", "greater-than", ">":
		return "gt"
	case "gte", "greater-than-or-equal", ">=":
		return "gte"
	case "lt", "lesser-than", "<":
		return "lt"
	case "lte", "lesser-than-or-equal", "<=":
		return "lte"
	case "starts-with", "sw":
		return "startsWith"
	case "ends-with", "ew":
		return "endsWith"
	case "contains":
		return "contains"
	case "strict-contains", "sc":
		return "strictContains"
	case "notin", "ni":
		return "notIn"
	case "exists", "ex":
		return "in"
	case "not-exists", "nex":
		return "in"
	default:
		return ""
	}
}
