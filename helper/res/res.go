package res

func Ok() interface{} {
	return map[string]interface{}{
		"error": 0,
		"msg":   "ok",
	}
}

func Data(value interface{}) interface{} {
	return map[string]interface{}{
		"error": 0,
		"data":  value,
	}
}

func Error(msg interface{}) interface{} {
	if val, ok := msg.(error); ok {
		msg = val.Error()
	}
	return map[string]interface{}{
		"error": 1,
		"msg":   msg,
	}
}
