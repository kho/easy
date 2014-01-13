package easy

import (
	"testing"
)

func TestParallelMap(t *testing.T) {
	cases := []struct {
		f     func(interface{}) interface{}
		input []int
		sum   int
	}{
		{intIdentity, []int{1, 2, 3, 4, 5}, 15},
		{intIdentity, nil, 0},
		{intSquare, []int{1, 2, 3}, 14},
	}
	for _, c := range cases {
		sum := 0
		source := make(chan interface{})
		go func() {
			for _, v := range c.input {
				source <- v
			}
			close(source)
		}()
		for v := range Parallel(2).Map(c.f, source) {
			sum += v.(int)
		}
		if sum != c.sum {
			t.Errorf("expected sum = %d; got %d from %v", c.sum, sum, c.input)
		}
	}
}

func TestParallelMemMap(t *testing.T) {
	cases := []struct {
		f     func(interface{}) interface{}
		input []int
		sum   int
	}{
		{intIdentity, []int{1, 2, 3, 4, 5}, 15},
		{intIdentity, nil, 0},
		{intSquare, []int{1, 2, 3}, 14},
	}
	for _, c := range cases {
		sum := 0
		source := make([]interface{}, len(c.input))
		for i, v := range c.input {
			source[i] = v
		}
		for v := range Parallel(2).MemMap(c.f, source) {
			sum += v.(int)
		}
		if sum != c.sum {
			t.Errorf("expected sum = %d; got %d from %v", c.sum, sum, c.input)
		}
	}
}

func intIdentity(x interface{}) interface{} {
	return x.(int)
}

func intSquare(x interface{}) interface{} {
	v := x.(int)
	return v * v
}
