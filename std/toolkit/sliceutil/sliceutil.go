package sliceutil

import (
	"encoding/json"
	"fmt"
	"strings"
)

func RemoveEmptyStrings(xs []string) []string {
	ys := []string{}
	for _, x := range xs {
		if x != "" {
			ys = append(ys, x)
		}
	}
	return ys
}

func RemoveDuplicateStrings(xs []string) []string {
	ys := []string{}
	visited := map[string]bool{}
	for _, x := range xs {
		if visited[x] == false {
			visited[x] = true
			ys = append(ys, x)
		}
	}
	return ys
}

func DeltaStrings(xs []string, ys []string) []string {
	zs := []string{}
L0:
	for _, x := range xs {
		for _, y := range ys {
			if x == y {
				continue L0
			}
		}
		zs = append(zs, x)
	}
	return zs
}

func IntersectionStrings(xs []string, ys []string) []string {
	zs := []string{}

	for _, x := range xs {
		for _, y := range ys {
			if x == y {
				zs = append(zs, x)
			}
		}
	}
	return zs
}

func IntersectionInt64(xs []int64, ys []int64) []int64 {
	var zs []int64

	for _, x := range xs {
		for _, y := range ys {
			if x == y {
				zs = append(zs, x)
			}
		}
	}
	return zs
}

func UnIntersectionInt64(xs []int64, ys []int64) []int64 {
	var zs []int64
	for _, v := range ys {
		if !ContainsInt64(xs, v) {
			zs = append(zs, v)
		}
	}
	return zs
}

func CartesianProductStrings(xs ...[]string) [][]string {
	product := [][]string{}

	switch {
	case len(xs) == 0:
	case len(xs) == 1:
		m := xs[0]
		for _, v := range m {
			product = append(product, []string{v})
		}
	default:
		m := xs[0]
		n := CartesianProductStrings(xs[1:]...)
		for _, v := range m {
			if len(n) > 0 {
				for _, w := range n {
					product = append(product, append([]string{v}, w...))
				}
			}
		}
	}

	return product
}

func ContainsString(xs []string, x string) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}

func ContainsStrings(xs []string, strs ...string) bool {
	return ContainsStringArr(xs, strs)
}

func ContainsStringArr(xs []string, arr []string) bool {
	if len(xs) == 0 || len(arr) == 0 {
		return false
	}

	for _, v := range arr {
		if ContainsString(xs, v) {
			return true
		}
	}
	return false
}

func ContainsInt64(xs []int64, x int64) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}

func ContainsInt32(xs []int32, x int32) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}

func RemoveDuplicateInt64s(xs []int64) []int64 {
	ys := []int64{}
	visited := map[int64]bool{}
	for _, x := range xs {
		if visited[x] == false {
			visited[x] = true
			ys = append(ys, x)
		}
	}
	return ys
}

func RemoveOnString(arr []string, key string) []string {
	res := []string{}
	for _, val := range arr {
		if val != key {
			res = append(res, val)
		}
	}
	return res
}

func SliceInt64Arr(arr []int64, start, limit int64) []int64 {
	if start > int64(len(arr)-1) {
		return []int64{}
	}
	if (start + limit) >= int64(len(arr)-1) {
		return arr[start:]
	}
	return arr[start:(start + limit)]
}

func StringsToInterfaceArr(strs []string) []interface{} {

	res := []interface{}{}
	for _, c := range strs {
		res = append(res, c)
	}

	return res
}

func DiffInt64Arr(arr1 []int64, arr2 []int64) []int64 {
	var res []int64
	for _, a := range arr1 {
		if !ContainsInt64(arr2, a) {
			res = append(res, a)
		}
	}

	return res
}

func FromStringArray(arr []string) string {
	data, _ := json.Marshal(arr)
	return string(data)
}

func ToStringArray(str string) (arr []string) {
	json.Unmarshal([]byte(str), &arr)
	return
}

func FromInt64Array(arr []int64) string {
	data, _ := json.Marshal(arr)
	return string(data)
}

func JoinInt64Slice(slice []int64, seg string) string {
	var strs []string
	for _, i := range slice{
		strs = append(strs, fmt.Sprintf("%d", i))
	}

	return strings.Join(strs, seg)
}

func ToInt64Array(str string) (arr []int64) {
	json.Unmarshal([]byte(str), &arr)
	return
}

func InsertInt64Arr(d []int64, el int64, index int) []int64 {
	if len(d) < index {
		d = append(d, el)
		return d
	}

	rear := append([]int64{}, d[index:]...)
	d = append(d[0:index], el)
	d = append(d, rear...)
	return d
}

func InsertMapInterface(d []map[string]interface{}, el map[string]interface{}, index int) []map[string]interface{} {

	if len(d) < index {
		d = append(d, el)
	}
	rear := append([]map[string]interface{}{}, d[index:]...)
	d = append(d[0:index], el)
	d = append(d, rear...)
	return d
}
