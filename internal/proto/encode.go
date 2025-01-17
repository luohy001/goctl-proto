package proto

import (
	"bytes"
	"strings"
	"text/template"
)

func (f *File) Marshal() ([]byte, error) {
	buf := bytes.NewBufferString("")
	if err := template.Must(template.New("").Funcs(funcMap).Parse(fileTemplate)).
		Execute(buf, f); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Refine refine messages and rpcs
func (f *File) Refine(includeRpcs, excludeRpcs []string) *File {
	if f == nil || len(f.Services) == 0 {
		return f
	}
	if len(includeRpcs) > 0 {
		includes := make(map[string]bool, len(includeRpcs))
		for _, rpc := range includeRpcs {
			includes[rpc] = true
		}
		for _, service := range f.Services {
			newRpcs := make([]*ServiceRpc, 0, len(includeRpcs))
			for _, rpc := range service.Rpcs {
				if !includes[rpc.Name] {
					continue
				}
				newRpcs = append(newRpcs, rpc)
			}
			service.Rpcs = newRpcs
		}
	} else {
		excludes := make(map[string]bool, len(excludeRpcs))
		for _, rpc := range excludeRpcs {
			excludes[rpc] = true
		}
		// find rpc which desc contains @goctl-proto
		noFlagService := make(map[string]bool)
		for _, service := range f.Services {
			preRpcs := make([]*ServiceRpc, 0, len(service.Rpcs))
			for _, rpc := range service.Rpcs {
				if strings.Contains(strings.Join(rpc.Descs, " "), "@goctl-proto") {
					for i := range rpc.Descs {
						rpc.Descs[i] = strings.TrimSpace(strings.ReplaceAll(rpc.Descs[i], "@goctl-proto", ""))
					}
					preRpcs = append(preRpcs, rpc)
				}
			}
			if len(preRpcs) > 0 {
				service.Rpcs = preRpcs
			} else {
				noFlagService[service.Name] = true
			}
		}
		// clean service rpcs which not contains @goctl-proto
		if len(noFlagService) < len(f.Services) {
			for i := 0; i < len(f.Services); i++ {
				if !noFlagService[f.Services[i].Name] {
					continue
				}
				f.Services = append(f.Services[:i:cap(f.Services)-1], f.Services[i+1:]...)
				i--
			}
		}

		for _, service := range f.Services {
			for i := 0; i < len(service.Rpcs); i++ {
				if !excludes[service.Rpcs[i].Name] {
					continue
				}
				service.Rpcs = append(service.Rpcs[:i:cap(service.Rpcs)-1], service.Rpcs[i+1:]...)
				i--
			}
		}
	}

	messageCustomTypeNames := make(map[string][]string, len(f.Messages))
	for _, message := range f.Messages {
		for _, field := range message.Fields {
			messageCustomTypeNames[message.Name] = append(messageCustomTypeNames[message.Name], field.CustomTypeNames...)
		}
	}

	usedMessages := make(map[string]bool)
	for _, service := range f.Services {
		for _, rpc := range service.Rpcs {
			usedMessages[rpc.Request.Name], usedMessages[rpc.Response.Name] = true, true
			var (
				customTypeNames    = append(messageCustomTypeNames[rpc.Request.Name], messageCustomTypeNames[rpc.Response.Name]...)
				tmpCustomTypeNames = make([]string, 0)
				hasNewMessage      bool
			)
		pickUsedCustomMessageLoop:
			for _, customTypeName := range customTypeNames {
				if usedMessages[customTypeName] {
					continue
				}
				hasNewMessage = true
				usedMessages[customTypeName] = true
				tmpCustomTypeNames = append(tmpCustomTypeNames, messageCustomTypeNames[customTypeName]...)
			}
			if hasNewMessage && len(tmpCustomTypeNames) > 0 {
				hasNewMessage = false
				customTypeNames = tmpCustomTypeNames
				tmpCustomTypeNames = make([]string, 0)
				goto pickUsedCustomMessageLoop
			}
		}
	}

	for i := 0; i < len(f.Messages); i++ {
		if usedMessages[f.Messages[i].Name] {
			continue
		}
		f.Messages = append(f.Messages[:i:cap(f.Messages)-1], f.Messages[i+1:]...)
		i--
	}
	return f
}
