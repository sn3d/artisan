package yaml

func castToMap(val interface{}) map[interface{}]interface{} {
	return val.(map[interface{}]interface{})
}

func castToStr(val interface{}) string {
	if val == nil {
		return ""
	}
	return val.(string)
}

func castToStringArray(val interface{}) []string {
	if val == nil {
		return []string{}
	}

	arr := val.([]interface{})
	res := make([]string, len(arr))
	for i, v := range arr {
		res[i] = v.(string)
	}
	return res
}

