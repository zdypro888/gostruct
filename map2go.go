package gostruct

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strconv"

	"github.com/zdypro888/go-plist"
)

func writeValue(writer io.StringWriter, value any) {
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
	case []any:
		writeSlice(writer, object)
	case map[string]any:
		writeMap(writer, object)
	default:
		panic(reflect.TypeOf(object))
	}
}

func writeData(writer io.StringWriter, data []uint8) {
	writer.WriteString("[]byte{\n")
	for _, value := range data {
		writer.WriteString(fmt.Sprintf("0x%02x, ", value))
	}
	writer.WriteString("}")
}
func writeSlice(writer io.StringWriter, slice []any) {
	writer.WriteString("[]any{\n")
	for _, value := range slice {
		writeValue(writer, value)
		writer.WriteString(",\n")
	}
	writer.WriteString("}")
}

func writeMap(writer io.StringWriter, mapval map[string]any) {
	writer.WriteString("map[string]any{\n")
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

// GOString 把Map信息写出到Go定义格式 map[string]interface { "a" : "b", "c" : 0x0d }
func GOString(mapval map[string]any) string {
	writer := &bytes.Buffer{}
	writeMap(writer, mapval)
	return writer.String()
}

func PlistFile2GOString(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	structmap := make(map[string]any)
	if _, err = plist.Unmarshal(data, structmap); err != nil {
		return err
	}
	return os.WriteFile("string.go", []byte(GOString(structmap)), 0644)
}
