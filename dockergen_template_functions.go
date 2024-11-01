/*

The following code is borrowed from Jason Wilder's docker-gen (https://github.com/jwilder/docker-gen).

The MIT License (MIT)

Copyright (c) 2014 Jason Wilder

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

package kubegen

import (
	"fmt"
	"log"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// arrayClosest find the longest matching substring in values
// that matches input
func arrayClosest(values []string, input string) string {
	best := ""
	for _, v := range values {
		if strings.Contains(input, v) && len(v) > len(best) {
			best = v
		}
	}
	return best
}

// dirList returns a list of files in the specified path
func dirList(path string) ([]string, error) {
	names := []string{}
	files, err := os.ReadDir(path)
	if err != nil {
		log.Printf("Template error: %v", err)
		return names, nil
	}
	for _, f := range files {
		names = append(names, f.Name())
	}
	return names, nil
}

// Generalized groupBy function
func generalizedGroupBy(funcName string, entries any, getValue func(any) (any, error), addEntry func(map[string][]any, any, any)) (map[string][]any, error) {
	entriesVal, err := getArrayValues(funcName, entries)
	if err != nil {
		return nil, err
	}

	groups := make(map[string][]any)
	for i := 0; i < entriesVal.Len(); i++ {
		v := reflect.Indirect(entriesVal.Index(i)).Interface()
		value, err := getValue(v)
		if err != nil {
			return nil, err
		}
		if value != nil {
			addEntry(groups, value, v)
		}
	}
	return groups, nil
}

func generalizedGroupByKey(funcName string, entries any, key string, addEntry func(map[string][]any, any, any)) (map[string][]any, error) {
	getKey := func(v any) (any, error) {
		return deepGet(v, key), nil
	}
	return generalizedGroupBy(funcName, entries, getKey, addEntry)
}

func groupByMulti(entries any, key, sep string) (map[string][]any, error) {
	return generalizedGroupByKey("groupByMulti", entries, key, func(groups map[string][]any, value any, v any) {
		items := strings.Split(value.(string), sep) //nolint:forcetypeassert
		for _, item := range items {
			groups[item] = append(groups[item], v)
		}
	})
}

// groupBy groups a generic array or slice by the path property key
func groupBy(entries any, key string) (map[string][]any, error) {
	return generalizedGroupByKey("groupBy", entries, key, func(groups map[string][]any, value any, v any) {
		groups[value.(string)] = append(groups[value.(string)], v) //nolint:forcetypeassert
	})
}

// groupByKeys is the same as groupBy but only returns a list of keys
func groupByKeys(entries any, key string) ([]string, error) {
	keys, err := generalizedGroupByKey("groupByKeys", entries, key, func(groups map[string][]any, value any, v any) {
		groups[value.(string)] = append(groups[value.(string)], v) //nolint:forcetypeassert
	})
	if err != nil {
		return nil, err
	}

	ret := []string{}
	for k := range keys {
		ret = append(ret, k)
	}
	return ret, nil
}

// when returns the trueValue when the condition is true and the falseValue otherwise
func when(condition bool, trueValue, falseValue any) any {
	if condition {
		return trueValue
	} else {
		return falseValue
	}
}

// Generalized where function
func generalizedWhere(funcName string, entries any, key string, test func(any) bool) (any, error) {
	entriesVal, err := getArrayValues(funcName, entries)
	if err != nil {
		return nil, err
	}

	selection := make([]any, 0)
	for i := 0; i < entriesVal.Len(); i++ {
		v := reflect.Indirect(entriesVal.Index(i)).Interface()

		value := deepGet(v, key)
		if test(value) {
			selection = append(selection, v)
		}
	}

	return selection, nil
}

// selects entries based on key
func where(entries any, key string, cmp any) (any, error) {
	return generalizedWhere("where", entries, key, func(value any) bool {
		return reflect.DeepEqual(value, cmp)
	})
}

// selects entries where a key exists
func whereExist(entries any, key string) (any, error) {
	return generalizedWhere("whereExist", entries, key, func(value any) bool {
		return value != nil
	})
}

// selects entries where a key does not exist
func whereNotExist(entries any, key string) (any, error) {
	return generalizedWhere("whereNotExist", entries, key, func(value any) bool {
		return value == nil
	})
}

// selects entries based on key.  Assumes key is delimited and breaks it apart before comparing
func whereAny(entries any, key, sep string, cmp []string) (any, error) {
	return generalizedWhere("whereAny", entries, key, func(value any) bool {
		if value == nil {
			return false
		} else {
			items := strings.Split(value.(string), sep) //nolint:forcetypeassert
			return len(intersect(cmp, items)) > 0
		}
	})
}

// selects entries based on key.  Assumes key is delimited and breaks it apart before comparing
func whereAll(entries any, key, sep string, cmp []string) (any, error) {
	reqCount := len(cmp)
	return generalizedWhere("whereAll", entries, key, func(value any) bool {
		if value == nil {
			return false
		} else {
			items := strings.Split(value.(string), sep) //nolint:forcetypeassert
			return len(intersect(cmp, items)) == reqCount
		}
	})
}

func getArrayValues(funcName string, entries any) (*reflect.Value, error) {
	entriesVal := reflect.ValueOf(entries)

	kind := entriesVal.Kind()

	if kind == reflect.Ptr {
		entriesVal = reflect.Indirect(entriesVal)
		kind = entriesVal.Kind()
	}

	switch kind {
	case reflect.Array, reflect.Slice:
		break
	default:
		return nil, fmt.Errorf("must pass an array or slice to '%v'; received %v; kind %v", funcName, entries, kind)
	}
	return &entriesVal, nil
}

func intersect(l1, l2 []string) []string {
	m := make(map[string]bool)
	m2 := make(map[string]bool)
	for _, v := range l2 {
		m2[v] = true
	}
	for _, v := range l1 {
		if m2[v] {
			m[v] = true
		}
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func deepGetImpl(v reflect.Value, path []string) any {
	if !v.IsValid() {
		return nil
	}
	if len(path) == 0 {
		return v.Interface()
	}
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if v.Kind() == reflect.Pointer {
		log.Printf("unable to descend into pointer of a pointer\n")
		return nil
	}
	switch v.Kind() {
	case reflect.Struct:
		return deepGetImpl(v.FieldByName(path[0]), path[1:])
	case reflect.Map:
		return deepGetImpl(v.MapIndex(reflect.ValueOf(path[0])), path[1:])
	case reflect.Slice, reflect.Array:
		iu64, err := strconv.ParseUint(path[0], 10, 64)
		if err != nil {
			log.Printf("non-negative decimal number required for array/slice index, got %#v\n", path[0])
			return nil
		}
		if iu64 > math.MaxInt {
			iu64 = math.MaxInt
		}
		i := int(iu64)
		if i >= v.Len() {
			log.Printf("index %v out of bounds", i)
			return nil
		}
		return deepGetImpl(v.Index(i), path[1:])
	default:
		log.Printf("unable to index by %s (value %v, kind %s)\n", path[0], v, v.Kind())
		return nil
	}
}

func deepGet(item any, path string) any {
	var parts []string
	if path != "" {
		parts = strings.Split(strings.TrimPrefix(path, "."), ".")
	}
	return deepGetImpl(reflect.ValueOf(item), parts)
}
