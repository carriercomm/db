// *
// * Copyright 2013, Scott Cagno, All rights reserved.
// * BSD Licensed - sites.google.com/site/bsdc3license
// *
// * Db :: Database
// *

package db

import (
	"fmt"
	"reflect"
)

// clean interface value, return string slice
func MatchCleaner(v interface{}) []string {
	var ss []string
	switch v.(type) {
	case []interface{}:
		for _, val := range v.([]interface{}) {
			ss = append(ss, fmt.Sprintf("%v", val))
		}
	case []string:
		return v.([]string)
	default:
		ss = append(ss, fmt.Sprintf("%v", v))
	}
	return ss
}

func Match(m map[string]interface{}, ss []string, n int) interface{} {
	if len(ss) == 0 || ss[0] == "" {
		return m
	}
	if v, ok := m[ss[n]]; ok {
		if reflect.TypeOf(v).Kind() == reflect.Map && n < len(ss)-1 {
			n++
			return Match(v.(map[string]interface{}), ss, n)
		}
		return v
	}
	return nil
}
