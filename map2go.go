package gostruct

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
)

func writeValue(writer io.StringWriter, value interface{}) {
	switch object := value.(type) {
	case bool:
		writer.WriteString(fmt.Sprintf("%v", object))
	case int:
		writer.WriteString(strconv.Itoa(object))
	case int64:
		writer.WriteString(strconv.Itoa(int(object)))
	case uint64:
		writer.WriteString(strconv.FormatUint(object, 10))
	case float32:
		writer.WriteString(strconv.FormatFloat(float64(object), 'g', 10, 32))
	case float64:
		writer.WriteString(strconv.FormatFloat(float64(object), 'g', 10, 64))
	case string:
		writer.WriteString("\"")
		writer.WriteString(object)
		writer.WriteString("\"")
	case []uint8:
		writeData(writer, object)
	case []interface{}:
		writeSlice(writer, object)
	case map[string]interface{}:
		writeMap(writer, object)
	default:
		panic(reflect.TypeOf(object))
	}
}

func writeData(writer io.StringWriter, data []uint8) {
	writer.WriteString("[]byte{\n")
	for _, value := range data {
		writer.WriteString(fmt.Sprintf("%#x", value))
	}
	writer.WriteString("}")
}
func writeSlice(writer io.StringWriter, slice []interface{}) {
	writer.WriteString("[]interface{}{\n")
	for _, value := range slice {
		writeValue(writer, value)
		writer.WriteString(",\n")
	}
	writer.WriteString("}")
}

func writeMap(writer io.StringWriter, mapval map[string]interface{}) {
	writer.WriteString("map[string]interface{}{\n")
	keys := make([]string, len(mapval))
	i := 0
	for key := range mapval {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	for _, key := range keys {
		writer.WriteString("\"")
		writer.WriteString(key)
		writer.WriteString("\" : ")
		writeValue(writer, mapval[key])
		writer.WriteString(",\n")
	}
	writer.WriteString("}")
}

//GOString 把Map信息写出到Go定义格式 map[string]interface { "a" : "b", "c" : 0x0d }
func GOString(mapval map[string]interface{}) string {
	writer := &bytes.Buffer{}
	writeMap(writer, mapval)
	return writer.String()
}
