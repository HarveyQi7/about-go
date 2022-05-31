package test

import (
	"about-go/tools"
	"reflect"
	"testing"
)

func TestStructToMap(t *testing.T) {
	var obj = struct {
		A1 string
		A2 int
		A3 struct {
			B1 string
			B2 int
		}
		A4 []string
		A5 map[string]string
	}{
		A1: "a1",
		A2: 12,
		A3: struct {
			B1 string
			B2 int
		}{
			B1: "b1",
			B2: 22,
		},
		A4: []string{"1", "2", "3"},
		A5: map[string]string{
			"k1": "v1",
			"k2": "v2",
		},
	}

	gotNotNested := tools.StructToMap(obj, false)
	wantNotNested := map[string]interface{}{
		"A1": "a1",
		"A2": 12,
		"B1": "b1",
		"B2": 22,
		"A4": []string{"1", "2", "3"},
		"A5": map[string]string{
			"k1": "v1",
			"k2": "v2",
		},
	}
	if !reflect.DeepEqual(wantNotNested, gotNotNested) {
		t.Errorf("excepted:%v, got:%v", wantNotNested, gotNotNested)
	}

	gotNested := tools.StructToMap(obj, true)
	wantNested := map[string]interface{}{
		"A1": "a1",
		"A2": 12,
		"A3": map[string]interface{}{
			"B1": "b1",
			"B2": 22,
		},
		"A4": []string{"1", "2", "3"},
		"A5": map[string]string{
			"k1": "v1",
			"k2": "v2",
		},
	}
	if !reflect.DeepEqual(wantNested, gotNested) {
		t.Errorf("excepted:%v, got:%v", wantNested, gotNested)
	}
}

func BenchmarkStructMap(b *testing.B) {
	b.StopTimer()
	var obj = struct {
		A1 string
		A2 int
		A3 struct {
			B1 string
			B2 int
		}
		A4 []string
		A5 map[string]string
	}{
		A1: "a1",
		A2: 12,
		A3: struct {
			B1 string
			B2 int
		}{
			B1: "b1",
			B2: 22,
		},
		A4: []string{"1", "2", "3"},
		A5: map[string]string{
			"k1": "v1",
			"k2": "v2",
		},
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tools.StructToMap(obj, false)
	}
}
