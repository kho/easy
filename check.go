package easy

func Check(v bool, message func() string) {
	if !v {
		panic("check failed: " + message())
	}
}
