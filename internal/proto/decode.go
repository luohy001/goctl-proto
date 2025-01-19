package proto

import (
	"fmt"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"regexp"
	"strings"
)

type fieldTypeName struct {
	singular string
	plural   string
}

var (
	emptyMessage     = Message{"Empty", nil, nil}
	arrayReg         = regexp.MustCompile(`\[\d+]`)
	fieldTypeNameMap = func() map[string]fieldTypeName {
		basic := map[string]string{
			"float32":     "float",
			"float64":     "double",
			"complex64":   "float",
			"complex128":  "double",
			"int":         "int32",
			"int8":        "int32",
			"int16":       "int32",
			"int32":       "int32",
			"int64":       "int64",
			"uint":        "uint32",
			"uint8":       "uint32",
			"uint16":      "uint32",
			"uint32":      "uint32",
			"uint64":      "uint64",
			"uintptr":     "uint64",
			"bool":        "bool",
			"string":      "string",
			"any":         "bytes", // todo: use google.protobuf.Any
			"interface{}": "bytes",
		}
		result := make(map[string]fieldTypeName)
		for goType, protoType := range basic {
			result[goType] = fieldTypeName{
				singular: protoType,
				plural:   protoType,
			}
		}
		result["byte"] = fieldTypeName{
			singular: "uint32",
			plural:   "bytes",
		}
		result["rune"] = fieldTypeName{
			singular: "uint32",
			plural:   "bytes",
		}
		return result
	}()
	mapKeyTypeNameMap = map[string]string{
		"int":    "int32",
		"int8":   "int32",
		"int16":  "int32",
		"int32":  "int32",
		"int64":  "int64",
		"uint":   "uint32",
		"uint8":  "uint32",
		"uint16": "uint32",
		"uint32": "uint32",
		"uint64": "uint64",
		"string": "string",
		"byte":   "uint32",
		"rune":   "uint32",
	}
)

func Unmarshal(data any, multiple bool, apiFileName string) (f *File, err error) {
	switch val := data.(type) {
	case *spec.ApiSpec:
		f = &File{
			Syntax:  Version3,
			Package: strings.ToLower(apiFileName),
			Options: []*Option{
				{
					Name:  "go_package",
					Value: "/protoc-gen-go",
				},
			},
		}
		messageMap := make(map[string]*Message, len(val.Types))
		for _, typ := range val.Types {
			defineStruct, _ := typ.(spec.DefineStruct)
			var message Message
			message.Name = defineStruct.Name()
			message.Descs = defineStruct.Documents()
			for _, member := range defineStruct.Members {
				var field MessageField
				if err = field.Unmarshal(&member); err != nil {
					return nil, err
				}
				message.Fields = append(message.Fields, &field)
			}
			f.Messages = append(f.Messages, &message)
			messageMap[message.Name] = &message
		}
		serviceNameMap := make(map[string]int)
		for _, group := range val.Service.JoinPrefix().Groups {
			var serviceName string
			if groupName := group.GetAnnotation("group"); groupName != "" && multiple {
				serviceName = groupName
			} else {
				serviceName = apiFileName
			}
			var srv *Service
			if srvIndex, exist := serviceNameMap[serviceName]; exist {
				srv = f.Services[srvIndex]
			} else {
				srv = &Service{Name: serviceName}
				f.Services = append(f.Services, srv)
				serviceNameMap[serviceName] = len(f.Services) - 1
			}

			for _, route := range group.Routes {
				var rpc ServiceRpc
				rpc.Name = route.Handler
				rpc.Descs = []string{strings.Trim(route.JoinedDoc(), `"`)}
				if defineStruct, ok := route.RequestType.(spec.DefineStruct); ok {
					rpc.Request = messageMap[defineStruct.Name()]
				} else {
					if _, exist := messageMap[emptyMessage.Name]; !exist {
						messageMap[emptyMessage.Name] = &emptyMessage
						f.Messages = append([]*Message{&emptyMessage}, f.Messages...)
					}
					rpc.Request = &emptyMessage
				}

				if defineStruct, ok := route.ResponseType.(spec.DefineStruct); ok {
					rpc.Response = messageMap[defineStruct.Name()]
				} else {
					if _, exist := messageMap[emptyMessage.Name]; !exist {
						messageMap[emptyMessage.Name] = &emptyMessage
						f.Messages = append([]*Message{&emptyMessage}, f.Messages...)
					}
					rpc.Response = &emptyMessage
				}
				srv.Rpcs = append(srv.Rpcs, &rpc)
			}
		}
	default:
		return nil, fmt.Errorf("unsupported type %T, only supported *spec.ApiSpec", data)
	}
	return
}

func (v *MessageField) Unmarshal(data any) error {
	switch val := data.(type) {
	case *spec.Member:
		v.Name = val.Name
		v.Descs = val.Docs
		if comment := strings.TrimSpace(strings.TrimPrefix(val.GetComment(), "//")); comment != "" {
			v.Descs = append(v.Descs, comment)
		}
		// note: bug in parse spec.Member, val.Type.Name() may include comment information
		var valTypeName string
		if ft := strings.Fields(strings.TrimSpace(val.Type.Name())); len(ft) == 0 {
			return fmt.Errorf("invalid type [%s]", val.Name)
		} else {
			valTypeName = ft[0]
		}
		typeName, typ, customTypes := parseFieldType(valTypeName)
		v.CustomTypeNames = customTypes
		if typ == MessageFieldTypeSlice {
			v.Repeated = true
			v.TypeName = typeName.plural
		} else {
			v.TypeName = typeName.singular
		}
		// parse tag
		// todo: parse other tags
		for _, tag := range val.Tags() {
			switch tag.Key {
			case "json":
				v.Tags = append(v.Tags, fmt.Sprintf(`json_name = "%s"`, tag.Name))
			}
		}
	default:
		return fmt.Errorf("unsupported type %T, only supported *spec.Member", data)
	}
	return nil
}

func parseFieldType(typeStr string) (typeName fieldTypeName, typ MessageFieldType, customTypeNames []string) {
	typeStr = arrayReg.ReplaceAllString(typeStr, "[]")
	typeStr = strings.Trim(typeStr, "*")
	// slice
	if strings.HasPrefix(typeStr, "[]") {
		typeName, typ, customTypeNames = parseFieldType(strings.TrimPrefix(typeStr, "[]"))
		if typ != MessageFieldTypeNormal {
			typeName = fieldTypeName{
				singular: "bytes",
				plural:   "bytes",
			}
		}
		typ = MessageFieldTypeSlice
		return
	}

	// normal
	if !strings.HasPrefix(typeStr, "map[") {
		typ = MessageFieldTypeNormal
		var exist bool
		typeName, exist = fieldTypeNameMap[typeStr]
		if !exist {
			typeName = fieldTypeName{
				singular: typeStr,
				plural:   typeStr,
			}
			customTypeNames = append(customTypeNames, typeStr)
		}
		return
	}

	// map
	typ = MessageFieldTypeMap
	mapKey, mapVal, err := parseMapKeyAndValue(typeStr)
	if err != nil {
		panic(err)
	}

	// key
	keyFieldTypeName, keyTyp, keyCustomTypes := parseFieldType(mapKey)
	customTypeNames = append(customTypeNames, keyCustomTypes...)
	if keyTyp == MessageFieldTypeNormal && mapKeyTypeNameMap[keyFieldTypeName.singular] != "" {
		mapKey = mapKeyTypeNameMap[keyFieldTypeName.singular]
	} else {
		mapKey = "string"
	}

	// value
	valFieldTypeName, valTyp, valCustomTypes := parseFieldType(mapVal)
	customTypeNames = append(customTypeNames, valCustomTypes...)
	if valTyp == MessageFieldTypeNormal {
		mapVal = valFieldTypeName.singular
	} else {
		mapVal = "bytes"
	}
	typeName = fieldTypeName{
		singular: fmt.Sprintf("map<%s,%s>", mapKey, mapVal),
		plural:   fmt.Sprintf("map<%s,%s>", mapKey, mapVal),
	}
	return
}

func parseMapKeyAndValue(mapStr string) (string, string, error) {
	if !strings.HasPrefix(mapStr, "map[") {
		return "", "", fmt.Errorf("type %s is not a map type", mapStr)
	}
	if !strings.Contains(mapStr, "]") {
		return "", "", fmt.Errorf("unsupported field type %s", mapStr)
	}
	var bracketCount int
	for i, char := range mapStr {
		if char == '[' {
			bracketCount++
			continue
		}
		if char != ']' {
			continue
		}
		if bracketCount == 1 {
			return mapStr[len("map["):i], mapStr[i+1:], nil
		}
		bracketCount--
	}
	return "", "", fmt.Errorf("unreachable code")
}
