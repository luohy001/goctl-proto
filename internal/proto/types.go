package proto

const (
	Version2 = "proto2"
	Version3 = "proto3"
)

type MessageFieldType int32

const (
	MessageFieldTypeNormal MessageFieldType = iota + 1
	MessageFieldTypeSlice
	MessageFieldTypeMap
)

// File
type File struct {
	Syntax   string
	Package  string
	Options  []*Option
	Messages []*Message
	Services []*Service
}

type Option struct {
	Name  string
	Value string
}

type Message struct {
	Name   string
	Descs  []string
	Fields []*MessageField
}

type MessageField struct {
	Name            string
	Descs           []string
	TypeName        string
	Repeated        bool
	CustomTypeNames []string
	Tags            []string
}

type Service struct {
	Name  string
	Descs []string
	Rpcs  []*ServiceRpc
}

type ServiceRpc struct {
	Name     string
	Descs    []string
	Request  *Message
	Response *Message
}
