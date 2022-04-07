package utils

func PanicNotNil(i interface{}) {
	if i != nil {
		panic(i)
	}
}
