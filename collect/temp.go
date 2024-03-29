package collect

type Temp struct {
	data map[string]interface{}
}

// 返回临时缓存数据
func (t *Temp) Get(key string) interface{} {
	return t.data[key]
}

func (t *Temp) Set(key string, value interface{}) error {
	if t.data == nil {
		t.data = make(map[string]interface{})
	}

	t.data[key] = value
	return nil
}
