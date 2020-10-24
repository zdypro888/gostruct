package gostruct

import (
	"bytes"
	"reflect"
	"time"

	"github.com/stoewer/go-strcase"
)

type property struct {
	Name string
	Type string
	Key  string
}

type gostruct map[string]*property

type goplist map[string]gostruct

func (gp goplist) String() string {
	gwriter := &bytes.Buffer{}
	for name, object := range gp {
		writer := &bytes.Buffer{}
		writer.WriteString("type ")
		writer.WriteString(name)
		writer.WriteString(" struct {\n")
		for _, prop := range object {
			writer.WriteString("\t")
			writer.WriteString(prop.Name)
			writer.WriteString(" ")
			writer.WriteString(prop.Type)
			writer.WriteString(" `plist:\"")
			writer.WriteString(prop.Key)
			writer.WriteString("\"`\n")
		}
		writer.WriteString("}\n")
		gwriter.WriteString(writer.String())
	}
	return gwriter.String()
}

func (gp goplist) item(key string, value interface{}) *property {
	prop := &property{
		Name: strcase.UpperCamelCase(key),
		Key:  key,
	}
	switch object := value.(type) {
	case int:
		prop.Type = "int"
	case int64:
		prop.Type = "int64"
	case uint64:
		prop.Type = "uint64"
	case float32:
		prop.Type = "float32"
	case float64:
		prop.Type = "float64"
	case string:
		prop.Type = "string"
	case time.Time:
		prop.Type = "time.Time"
	case []interface{}:
		prop.Type = gp.array(key, object)
	case map[string]interface{}:
		prop.Type = gp.dict(strcase.UpperCamelCase(key), object)
	default:
		panic(reflect.TypeOf(object))
	}
	return prop
}

func (gp goplist) array(key string, values []interface{}) string {
	props := make([]*property, len(values))
	for i, value := range values {
		props[i] = gp.item(key, value)
	}
	return "[]" + props[0].Name
}

func (gp goplist) dict(name string, messages map[string]interface{}) string {
	object := make(gostruct)
	for key, value := range messages {
		prop := gp.item(key, value)
		object[prop.Name] = prop
	}
	if oldobj, ok := gp[name]; ok {
		for key, value := range object {
			if _, pok := oldobj[key]; !pok {
				oldobj[key] = value
			}
		}
	} else {
		gp[name] = object
	}
	return name
}

//PlistToGo plist 到 go 结构体
func PlistToGo(name string, plistmap map[string]interface{}) string {
	gp := make(goplist)
	gp.dict(name, plistmap)
	return gp.String()
}
