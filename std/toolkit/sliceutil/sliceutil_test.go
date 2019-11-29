package sliceutil_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/PKUJohnson/solar/std/toolkit/sliceutil"
)

func TestRemoveEmptyStrings(t *testing.T) {
	tcs := []struct {
		in []string
		eo []string
	}{
		{in: []string{}, eo: []string{}},
		{in: []string{""}, eo: []string{}},
		{in: []string{"foo", ""}, eo: []string{"foo"}},
		{in: []string{"", "bar", "qux"}, eo: []string{"bar", "qux"}},
		{in: []string{"", "", ""}, eo: []string{}},
		{in: []string{"", "foo", ""}, eo: []string{"foo"}},
	}
	for i, tc := range tcs {
		ao := sliceutil.RemoveEmptyStrings(tc.in)
		if !reflect.DeepEqual(ao, tc.eo) {
			t.Errorf("case %d: expected: %#v, actual: %#v", i, tc.eo, ao)
		}
	}
}

func TestRemoveDuplicateStrings(t *testing.T) {
	tcs := []struct {
		in []string
		eo []string
	}{
		{in: []string{}, eo: []string{}},
		{in: []string{""}, eo: []string{""}},
		{in: []string{"foo", ""}, eo: []string{"foo", ""}},
		{in: []string{"", "bar", "qux", "qux"}, eo: []string{"", "bar", "qux"}},
		{in: []string{"", "", ""}, eo: []string{""}},
		{in: []string{"", "foo", ""}, eo: []string{"", "foo"}},
	}
	for i, tc := range tcs {
		ao := sliceutil.RemoveDuplicateStrings(tc.in)
		if !reflect.DeepEqual(ao, tc.eo) {
			t.Errorf("case %d: expected: %#v, actual: %#v", i, tc.eo, ao)
		}
	}
}

func TestDeltaStrings(t *testing.T) {
	tcs := []struct {
		xs []string
		ys []string
		eo []string
	}{
		{xs: []string{}, ys: []string{}, eo: []string{}},
		{xs: []string{""}, ys: []string{}, eo: []string{""}},
		{xs: []string{}, ys: []string{""}, eo: []string{}},
		{xs: []string{"foo"}, ys: []string{"foo"}, eo: []string{}},
		{xs: []string{"foo", "bar"}, ys: []string{"foo"}, eo: []string{"bar"}},
		{xs: []string{"foo", "bar", "qux"}, ys: []string{"foo"}, eo: []string{"bar", "qux"}},
		{xs: []string{"foo", "bar", "qux"}, ys: []string{"foo", "bar"}, eo: []string{"qux"}},
		{xs: []string{"foo", "bar", "bar"}, ys: []string{"foo", "bar"}, eo: []string{}},
		{xs: []string{"foo", "bar", "bar"}, ys: []string{"foo"}, eo: []string{"bar", "bar"}},
	}
	for i, tc := range tcs {
		ao := sliceutil.DeltaStrings(tc.xs, tc.ys)
		if !reflect.DeepEqual(ao, tc.eo) {
			t.Errorf("case %d: expected: %#v, actual: %#v", i, tc.eo, ao)
		}
	}
}

func TestIntersectionStrings(t *testing.T) {
	tcs := []struct {
		xs []string
		ys []string
		eo []string
	}{
		{xs: []string{}, ys: []string{}, eo: []string{}},
		{xs: []string{""}, ys: []string{}, eo: []string{}},
		{xs: []string{}, ys: []string{""}, eo: []string{}},
		{xs: []string{"foo"}, ys: []string{"foo"}, eo: []string{"foo"}},
		{xs: []string{"foo", "bar"}, ys: []string{"foo"}, eo: []string{"foo"}},
		{xs: []string{"foo", "bar", "qux"}, ys: []string{"foo"}, eo: []string{"foo"}},
		{xs: []string{"foo", "bar", "qux"}, ys: []string{"foo", "bar"}, eo: []string{"foo", "bar"}},
		{xs: []string{"foo", "bar", "foo"}, ys: []string{"foo"}, eo: []string{"foo", "foo"}},
	}
	for i, tc := range tcs {
		ao := sliceutil.IntersectionStrings(tc.xs, tc.ys)
		if !reflect.DeepEqual(ao, tc.eo) {
			t.Errorf("case %d: expected: %#v, actual: %#v", i, tc.eo, ao)
		}
	}
}

func TestCartesianProductStrings(t *testing.T) {
	tcs := []struct {
		in [][]string
		eo [][]string
	}{
		{in: [][]string{}, eo: [][]string{}},
		{in: [][]string{{""}}, eo: [][]string{{""}}},
		{in: [][]string{{"foo"}, {"foo"}}, eo: [][]string{{"foo", "foo"}}},
		{in: [][]string{{"foo"}, {"foo", "bar"}}, eo: [][]string{{"foo", "foo"}, {"foo", "bar"}}},
		{in: [][]string{{"foo", "bar"}, {"foo", "bar"}}, eo: [][]string{
			{"foo", "foo"}, {"foo", "bar"}, {"bar", "foo"}, {"bar", "bar"}}},
		{in: [][]string{{"foo", "bar"}, {}}, eo: [][]string{}},
		{in: [][]string{{"foo", "bar"}, {}, {"foo"}}, eo: [][]string{}},
		{in: [][]string{{}, {"foo"}, {"foo"}}, eo: [][]string{}},
		{in: [][]string{{"foo", "bar"}, {"foo"}, {"foo"}}, eo: [][]string{
			{"foo", "foo", "foo"}, {"bar", "foo", "foo"},
		}},
	}
	for i, tc := range tcs {
		ao := sliceutil.CartesianProductStrings(tc.in...)
		if !reflect.DeepEqual(ao, tc.eo) {
			t.Errorf("case %d: expected: %#v, actual: %#v", i, tc.eo, ao)
		}
	}
}

func TestRemoveDuplicateInt64s(t *testing.T) {
	a := sliceutil.RemoveDuplicateInt64s([]int64{2, 3, 43, 4, 5, 3, 3, 2, 2})
	fmt.Println(a)
}

//func TestInsertSlice(t *testing.T) {
//	a := sliceutil.InsertArray([]int{1, 2, 3, 4, 5}, 8, 2)
//}
