package util

// TryThen tries to 'try' parameter func get value, if error not null, will be invoked errorHandle
// decide again invoke then parameter
func TryThen[T interface{}](try func() (*T, error), then func() (*T, error), errorHandle func(error) bool) (*T, error) {
	v, err := try()
	if err != nil {
		again := errorHandle(err)
		if again {
			return then()
		}
		return nil, err
	}

	return v, nil
}
