package proto

import "testing"

func TestFile_Marshal(t *testing.T) {
	reqMessage := Message{
		Name:  "Request",
		Descs: []string{"This is Request"},
		Fields: []*MessageField{
			{
				Name:     "id",
				Descs:    []string{"This is id"},
				TypeName: "int64",
			},
			{
				Name:     "name",
				Descs:    []string{"This is name"},
				TypeName: "string",
			},
			{
				Name:     "tags",
				Descs:    []string{"This is tag slice"},
				TypeName: "string",
				Repeated: true,
			},
		},
	}
	respMessage := Message{
		Name:  "Response",
		Descs: []string{"Response", "This is Response"},
		Fields: []*MessageField{
			{
				Name:     "code",
				Descs:    []string{"code", "This is code"},
				TypeName: "int32",
			},
			{
				Name:     "message",
				Descs:    []string{"This is message"},
				TypeName: "string",
			},
		},
	}
	f := File{
		Syntax:   "proto3",
		Messages: []*Message{&reqMessage, &respMessage},
		Services: []*Service{
			{
				Name:  "fileEncoder",
				Descs: []string{"FileEncoder Description", "This is a mock service"},
				Rpcs: []*ServiceRpc{
					{
						Name:     "doRequest",
						Descs:    []string{"Do Api Description"},
						Request:  &reqMessage,
						Response: &respMessage,
					},
				},
			},
			{
				Name: "fileDecoder",
				//Descs: []string{"FileDecoder Description", "This is a mock service"},
				Rpcs: []*ServiceRpc{
					{
						Name:     "doRequest",
						Descs:    []string{"Do Api Description"},
						Request:  &reqMessage,
						Response: &respMessage,
					},
				},
			},
		},
	}
	data, err := f.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}
