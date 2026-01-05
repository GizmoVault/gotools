package redisx

import (
	"fmt"
	"strings"

	"github.com/GizmoVault/gotools/base/errorx"
)

func FTGenTextQuery(name, value string) string {
	return fmt.Sprintf("@%s:%s", name, value)
}

func FTGenTagsQuery(name string, vs []string) string {
	if len(vs) == 0 {
		return ""
	}

	var query string

	query += "@" + name + ":{"
	query += strings.Join(vs, "|")
	query += "} "

	return query
}

func FTGenTagsQuery2[T any](name string, vs []T) string {
	if len(vs) == 0 {
		return ""
	}

	var query string

	vss := make([]string, len(vs))
	for i, v := range vs {
		vss[i] = fmt.Sprintf("%v", v)
	}

	query += "@" + name + ":{"
	query += strings.Join(vss, "|")
	query += "} "

	return query
}

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type NumericRangeFlag int

const (
	// BoundInclusive represents a closed boundary that includes the endpoint value.
	// Example: [5, 10] — both 5 and 10 are included.
	BoundInclusive NumericRangeFlag = iota

	// BoundExclusive represents an open boundary that excludes the endpoint value.
	// Example: (5, 10) — neither 5 nor 10 is included.
	BoundExclusive

	// BoundNegInf represents negative infinity (-∞).
	// Used to indicate no lower bound in a numeric range.
	BoundNegInf

	// BoundPosInf represents positive infinity (+∞).
	// Used to indicate no upper bound in a numeric range.
	BoundPosInf
)

func (f NumericRangeFlag) String() string {
	switch f {
	case BoundInclusive:
		return "Inclusive"
	case BoundExclusive:
		return "Exclusive"
	case BoundNegInf:
		return "-Inf"
	case BoundPosInf:
		return "+Inf"
	default:
		return fmt.Sprintf("NumericRangeFlag(%d)", f)
	}
}

func FTGenNumericRangeQueryIgnoreError[T Integer](name string, from, to T, fromFlag, toFlag NumericRangeFlag) (query string) {
	query, _ = FTGenNumericRangeQuery[T](name, from, to, fromFlag, toFlag)

	return
}
func FTGenNumericRangeQuery[T Integer](name string, from, to T, fromFlag, toFlag NumericRangeFlag) (query string, err error) {
	query = "@" + name + ":["
	switch fromFlag {
	case BoundNegInf:
		query += "-Inf"
	case BoundInclusive:
		query += fmt.Sprintf("%v", from)
	case BoundExclusive:
		query += fmt.Sprintf("(%v", from)
	default:
		err = errorx.ErrInvalidArgs

		return
	}

	query += " "
	switch toFlag {
	case BoundPosInf:
		query += "+Inf"
	case BoundInclusive:
		query += fmt.Sprintf("%v", to)
	case BoundExclusive:
		query += fmt.Sprintf("(%v", to)
	default:
		err = errorx.ErrInvalidArgs

		return
	}

	query += "]"

	return
}

func FTGenNumericTagsQuery[T Integer](name string, vs []T) string {
	if len(vs) == 0 {
		return ""
	}

	subQueries := make([]string, 0, len(vs))

	for _, v := range vs {
		subQuery := "@" + name + ":["
		subQuery += fmt.Sprintf("%d %d", v, v)
		subQuery += "]"

		subQueries = append(subQueries, subQuery)
	}

	query := "("
	query += strings.Join(subQueries, " | ")
	query += ")" + " "

	return query
}

func FTGenNumericTagsQueryEx[T Integer](name string, vs []T, useNot, useAnd bool) string {
	if len(vs) == 0 {
		return ""
	}

	subQueries := make([]string, 0, len(vs))

	var n string
	if useNot {
		n = "-"
	}

	for _, v := range vs {
		subQuery := n + "@" + name + ":["
		subQuery += fmt.Sprintf("%d %d", v, v)
		subQuery += "]"

		subQueries = append(subQueries, subQuery)
	}

	query := "("
	if useAnd {
		query += strings.Join(subQueries, " ")
	} else {
		query += strings.Join(subQueries, " | ")
	}

	query += ")" + " "

	return query
}

func FTFinalQuery(query string) string {
	query = strings.Trim(query, " ")
	if query == "" {
		query = "*"
	}

	return query
}
