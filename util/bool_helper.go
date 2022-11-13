package util

// If 模拟的三元运算符
//
//	condition: 条件表达式
//	trueVal: 表达式为true时返回的值
//	falseVal: 表达式为false时返回的值
//
// return: 根据表达式的true/false，返回trueVal/falseVal
//
//	注意，由于返回的类型是interface{}，需要转换成trueVal/falseVal对应的类型
func If(condition bool, trueVal, falseVal any) any {
	if condition {
		return trueVal
	}
	return falseVal
}
