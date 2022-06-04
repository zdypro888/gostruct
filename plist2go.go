package gostruct

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"strings"
	"time"

	"github.com/stoewer/go-strcase"
	"github.com/zdypro888/go-plist"
)

type AnalyseProperty struct {
	Parent string
	Name   string
	Type   string
	Key    string
}

type AnalyseStruct struct {
	Name       string
	Properties []*AnalyseProperty
	Depth      int
}
type Analyzing struct {
	Structs []*AnalyseStruct
}

func (analyzing *Analyzing) String() string {
	gwriter := &bytes.Buffer{}
	for _, struct_ := range analyzing.Structs {
		writer := &bytes.Buffer{}
		writer.WriteString("type ")
		writer.WriteString(struct_.Name)
		writer.WriteString(" struct {\n")
		for _, field := range struct_.Properties {
			writer.WriteString("\t")
			writer.WriteString(field.Name)
			writer.WriteString(" ")
			writer.WriteString(field.Type)
			writer.WriteString(" `plist:\"")
			writer.WriteString(field.Key)
			writer.WriteString(",omitempty\" bson:\"")
			writer.WriteString(field.Key)
			writer.WriteString(",omitempty\"`\n")
		}
		writer.WriteString("}\n")
		gwriter.WriteString(writer.String())
	}
	return gwriter.String()
}

func (analyzing *Analyzing) fieldName(name string) string {
	name = strings.ReplaceAll(name, ".", "_")
	return strcase.UpperCamelCase(name)
}

func (analyzing *Analyzing) Analyse(depth int, property *AnalyseProperty, value interface{}) {
	switch object := value.(type) {
	case bool:
		property.Type = "bool"
	case int:
		property.Type = "int"
	case int64:
		property.Type = "int64"
	case uint64:
		property.Type = "uint64"
	case float32:
		property.Type = "float32"
	case float64:
		property.Type = "float64"
	case string:
		property.Type = "string"
	case time.Time:
		property.Type = "time.Time"
	case []uint8:
		property.Type = "[]byte"
	case []interface{}:
		if len(object) == 0 {
			property.Type = "[]any"
		} else {
			analyzing.AnalyseSlice(depth, property, object)
		}
	case map[string]interface{}:
		analyzing.AnalyseMap(depth, property, object)
	default:
		panic(reflect.TypeOf(object))
	}
}

func (analyzing *Analyzing) AnalyseSlice(depth int, property *AnalyseProperty, values []interface{}) {
	properties := make([]*AnalyseProperty, len(values))
	for i, value := range values {
		properties[i] = &AnalyseProperty{
			Parent: property.Name,
			Name:   property.Name,
			Key:    property.Key,
		}
		analyzing.Analyse(depth+1, properties[i], value)
	}
	property.Type = "[]" + properties[0].Type
}
func (analyzing *Analyzing) AnalyseMap(depth int, property *AnalyseProperty, value map[string]interface{}) {
	structName := property.Name
	var structProperties []*AnalyseProperty
	for name, field := range value {
		fieldProperty := &AnalyseProperty{
			Parent: structName,
			Name:   analyzing.fieldName(name),
			Key:    name,
		}
		analyzing.Analyse(depth+1, fieldProperty, field)
		structProperties = append(structProperties, fieldProperty)
	}
	for _, struct_ := range analyzing.Structs {
		if struct_.Name == structName {
			structName = property.Parent + structName
			break
		}
	}
	analyzing.Structs = append(analyzing.Structs, &AnalyseStruct{
		Name:       structName,
		Properties: structProperties,
		Depth:      depth,
	})
	property.Type = "*" + structName
}

//Plist2Go plist 到 go 结构体
func Plist2Go(name string, plistmap map[string]interface{}) string {
	analyzse := &Analyzing{}
	analyzse.Analyse(0, &AnalyseProperty{Name: name}, plistmap)
	return analyzse.String()
}

func PlistFile2Go(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	structmap := make(map[string]interface{})
	if _, err = plist.Unmarshal(data, structmap); err != nil {
		return err
	}
	return ioutil.WriteFile("request.go", []byte(Plist2Go("Request", structmap)), 0644)
}
