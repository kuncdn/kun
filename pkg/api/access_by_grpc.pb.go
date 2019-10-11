// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pkg/api/access_by_grpc.proto

package api

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type StatusCode int32

const (
	StatusCode_ACCEPT StatusCode = 0
	StatusCode_REJECT StatusCode = 1
	StatusCode_DROP   StatusCode = 4
)

var StatusCode_name = map[int32]string{
	0: "ACCEPT",
	1: "REJECT",
	4: "DROP",
}

var StatusCode_value = map[string]int32{
	"ACCEPT": 0,
	"REJECT": 1,
	"DROP":   4,
}

func (x StatusCode) String() string {
	return proto.EnumName(StatusCode_name, int32(x))
}

func (StatusCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_1de7a493db6b7348, []int{0}
}

type Userinfo struct {
	Username             string   `protobuf:"bytes,1,opt,name=Username,proto3" json:"Username,omitempty"`
	Password             string   `protobuf:"bytes,2,opt,name=Password,proto3" json:"Password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Userinfo) Reset()         { *m = Userinfo{} }
func (m *Userinfo) String() string { return proto.CompactTextString(m) }
func (*Userinfo) ProtoMessage()    {}
func (*Userinfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_1de7a493db6b7348, []int{0}
}

func (m *Userinfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Userinfo.Unmarshal(m, b)
}
func (m *Userinfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Userinfo.Marshal(b, m, deterministic)
}
func (m *Userinfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Userinfo.Merge(m, src)
}
func (m *Userinfo) XXX_Size() int {
	return xxx_messageInfo_Userinfo.Size(m)
}
func (m *Userinfo) XXX_DiscardUnknown() {
	xxx_messageInfo_Userinfo.DiscardUnknown(m)
}

var xxx_messageInfo_Userinfo proto.InternalMessageInfo

func (m *Userinfo) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *Userinfo) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type URL struct {
	Scheme               string    `protobuf:"bytes,1,opt,name=Scheme,proto3" json:"Scheme,omitempty"`
	Opaque               string    `protobuf:"bytes,2,opt,name=Opaque,proto3" json:"Opaque,omitempty"`
	User                 *Userinfo `protobuf:"bytes,3,opt,name=User,proto3" json:"User,omitempty"`
	Host                 string    `protobuf:"bytes,4,opt,name=Host,proto3" json:"Host,omitempty"`
	Path                 string    `protobuf:"bytes,5,opt,name=Path,proto3" json:"Path,omitempty"`
	RawPath              string    `protobuf:"bytes,6,opt,name=RawPath,proto3" json:"RawPath,omitempty"`
	ForceQuery           string    `protobuf:"bytes,7,opt,name=ForceQuery,proto3" json:"ForceQuery,omitempty"`
	RawQuery             string    `protobuf:"bytes,8,opt,name=RawQuery,proto3" json:"RawQuery,omitempty"`
	Fragment             string    `protobuf:"bytes,9,opt,name=Fragment,proto3" json:"Fragment,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *URL) Reset()         { *m = URL{} }
func (m *URL) String() string { return proto.CompactTextString(m) }
func (*URL) ProtoMessage()    {}
func (*URL) Descriptor() ([]byte, []int) {
	return fileDescriptor_1de7a493db6b7348, []int{1}
}

func (m *URL) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_URL.Unmarshal(m, b)
}
func (m *URL) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_URL.Marshal(b, m, deterministic)
}
func (m *URL) XXX_Merge(src proto.Message) {
	xxx_messageInfo_URL.Merge(m, src)
}
func (m *URL) XXX_Size() int {
	return xxx_messageInfo_URL.Size(m)
}
func (m *URL) XXX_DiscardUnknown() {
	xxx_messageInfo_URL.DiscardUnknown(m)
}

var xxx_messageInfo_URL proto.InternalMessageInfo

func (m *URL) GetScheme() string {
	if m != nil {
		return m.Scheme
	}
	return ""
}

func (m *URL) GetOpaque() string {
	if m != nil {
		return m.Opaque
	}
	return ""
}

func (m *URL) GetUser() *Userinfo {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *URL) GetHost() string {
	if m != nil {
		return m.Host
	}
	return ""
}

func (m *URL) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *URL) GetRawPath() string {
	if m != nil {
		return m.RawPath
	}
	return ""
}

func (m *URL) GetForceQuery() string {
	if m != nil {
		return m.ForceQuery
	}
	return ""
}

func (m *URL) GetRawQuery() string {
	if m != nil {
		return m.RawQuery
	}
	return ""
}

func (m *URL) GetFragment() string {
	if m != nil {
		return m.Fragment
	}
	return ""
}

type Request struct {
	Header               map[string]string `protobuf:"bytes,1,rep,name=Header,proto3" json:"Header,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Body                 []byte            `protobuf:"bytes,2,opt,name=Body,proto3" json:"Body,omitempty"`
	RemoteAddr           string            `protobuf:"bytes,3,opt,name=RemoteAddr,proto3" json:"RemoteAddr,omitempty"`
	Method               string            `protobuf:"bytes,5,opt,name=Method,proto3" json:"Method,omitempty"`
	Proto                string            `protobuf:"bytes,6,opt,name=Proto,proto3" json:"Proto,omitempty"`
	URL                  *URL              `protobuf:"bytes,7,opt,name=URL,proto3" json:"URL,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}
func (*Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_1de7a493db6b7348, []int{2}
}

func (m *Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Request.Unmarshal(m, b)
}
func (m *Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Request.Marshal(b, m, deterministic)
}
func (m *Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Request.Merge(m, src)
}
func (m *Request) XXX_Size() int {
	return xxx_messageInfo_Request.Size(m)
}
func (m *Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Request proto.InternalMessageInfo

func (m *Request) GetHeader() map[string]string {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Request) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

func (m *Request) GetRemoteAddr() string {
	if m != nil {
		return m.RemoteAddr
	}
	return ""
}

func (m *Request) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

func (m *Request) GetProto() string {
	if m != nil {
		return m.Proto
	}
	return ""
}

func (m *Request) GetURL() *URL {
	if m != nil {
		return m.URL
	}
	return nil
}

type Response struct {
	Header               map[string]string `protobuf:"bytes,1,rep,name=Header,proto3" json:"Header,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Body                 []byte            `protobuf:"bytes,2,opt,name=Body,proto3" json:"Body,omitempty"`
	Status               StatusCode        `protobuf:"varint,3,opt,name=Status,proto3,enum=proto.StatusCode" json:"Status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_1de7a493db6b7348, []int{3}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetHeader() map[string]string {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Response) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

func (m *Response) GetStatus() StatusCode {
	if m != nil {
		return m.Status
	}
	return StatusCode_ACCEPT
}

func init() {
	proto.RegisterEnum("proto.StatusCode", StatusCode_name, StatusCode_value)
	proto.RegisterType((*Userinfo)(nil), "proto.Userinfo")
	proto.RegisterType((*URL)(nil), "proto.URL")
	proto.RegisterType((*Request)(nil), "proto.Request")
	proto.RegisterMapType((map[string]string)(nil), "proto.Request.HeaderEntry")
	proto.RegisterType((*Response)(nil), "proto.Response")
	proto.RegisterMapType((map[string]string)(nil), "proto.Response.HeaderEntry")
}

func init() { proto.RegisterFile("pkg/api/access_by_grpc.proto", fileDescriptor_1de7a493db6b7348) }

var fileDescriptor_1de7a493db6b7348 = []byte{
	// 474 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x53, 0x4f, 0x6f, 0x12, 0x41,
	0x14, 0x77, 0x61, 0x59, 0xe0, 0x61, 0x5a, 0x9c, 0x18, 0x33, 0xc1, 0xc6, 0x10, 0x8c, 0x09, 0x7a,
	0x80, 0x64, 0x7b, 0x51, 0x6f, 0x40, 0x69, 0x1a, 0x53, 0xd3, 0x75, 0xda, 0x5e, 0xbc, 0x34, 0xd3,
	0xdd, 0x27, 0x90, 0xca, 0xce, 0x76, 0x66, 0x90, 0xec, 0x47, 0xf3, 0x7b, 0x79, 0xf0, 0x68, 0xe6,
	0xcf, 0x22, 0xf6, 0x64, 0xd2, 0x13, 0xbf, 0x3f, 0xf3, 0x1e, 0xef, 0xf7, 0x5e, 0x16, 0x8e, 0x8a,
	0xbb, 0xc5, 0x98, 0x17, 0xab, 0x31, 0x4f, 0x53, 0x54, 0xea, 0xe6, 0xb6, 0xbc, 0x59, 0xc8, 0x22,
	0x1d, 0x15, 0x52, 0x68, 0x41, 0x1a, 0xf6, 0x67, 0x30, 0x85, 0xd6, 0xb5, 0x42, 0xb9, 0xca, 0xbf,
	0x09, 0xd2, 0x73, 0x38, 0xe7, 0x6b, 0xa4, 0x41, 0x3f, 0x18, 0xb6, 0xd9, 0x8e, 0x1b, 0x2f, 0xe1,
	0x4a, 0x6d, 0x85, 0xcc, 0x68, 0xcd, 0x79, 0x15, 0x1f, 0xfc, 0x0a, 0xa0, 0x7e, 0xcd, 0xce, 0xc9,
	0x0b, 0x88, 0x2e, 0xd3, 0x25, 0xee, 0xaa, 0x3d, 0x33, 0xfa, 0x45, 0xc1, 0xef, 0x37, 0xe8, 0x2b,
	0x3d, 0x23, 0xaf, 0x21, 0x34, 0xfd, 0x69, 0xbd, 0x1f, 0x0c, 0x3b, 0xf1, 0xa1, 0x1b, 0x6c, 0x54,
	0x8d, 0xc3, 0xac, 0x49, 0x08, 0x84, 0x67, 0x42, 0x69, 0x1a, 0xda, 0x52, 0x8b, 0x8d, 0x96, 0x70,
	0xbd, 0xa4, 0x0d, 0xa7, 0x19, 0x4c, 0x28, 0x34, 0x19, 0xdf, 0x5a, 0x39, 0xb2, 0x72, 0x45, 0xc9,
	0x2b, 0x80, 0x53, 0x21, 0x53, 0xfc, 0xb2, 0x41, 0x59, 0xd2, 0xa6, 0x35, 0xf7, 0x14, 0x13, 0x8d,
	0xf1, 0xad, 0x73, 0x5b, 0x2e, 0x5a, 0xc5, 0x8d, 0x77, 0x2a, 0xf9, 0x62, 0x8d, 0xb9, 0xa6, 0x6d,
	0xe7, 0x55, 0x7c, 0xf0, 0x3b, 0x80, 0x26, 0xc3, 0xfb, 0x0d, 0x2a, 0x4d, 0x62, 0x88, 0xce, 0x90,
	0x67, 0x28, 0x69, 0xd0, 0xaf, 0x0f, 0x3b, 0x71, 0xcf, 0x87, 0xf1, 0xfe, 0xc8, 0x99, 0xf3, 0x5c,
	0xcb, 0x92, 0xf9, 0x97, 0x26, 0xc5, 0x54, 0x64, 0xa5, 0x5d, 0xca, 0x53, 0x66, 0xb1, 0x99, 0x95,
	0xe1, 0x5a, 0x68, 0x9c, 0x64, 0x99, 0x5b, 0x4c, 0x9b, 0xed, 0x29, 0x66, 0x95, 0x9f, 0x51, 0x2f,
	0x45, 0xe6, 0xb3, 0x7b, 0x46, 0x9e, 0x43, 0x23, 0x31, 0x7f, 0xe8, 0xb3, 0x3b, 0x42, 0x8e, 0xec,
	0x5d, 0x6c, 0xe4, 0x4e, 0x0c, 0xd5, 0x7e, 0xd9, 0x39, 0x33, 0x72, 0xef, 0x03, 0x74, 0xf6, 0xc6,
	0x22, 0x5d, 0xa8, 0xdf, 0x61, 0xe9, 0x4f, 0x67, 0xa0, 0x69, 0xfa, 0x83, 0x7f, 0xdf, 0x9d, 0xcd,
	0x91, 0x8f, 0xb5, 0xf7, 0xc1, 0xe0, 0x67, 0x00, 0x2d, 0x86, 0xaa, 0x10, 0xb9, 0x42, 0x72, 0xfc,
	0x20, 0xfb, 0xcb, 0x5d, 0x76, 0xf7, 0xe0, 0xbf, 0xc3, 0xbf, 0x85, 0xe8, 0x52, 0x73, 0xbd, 0x51,
	0x36, 0xf8, 0x41, 0xfc, 0xcc, 0x37, 0x72, 0xe2, 0x4c, 0x64, 0xc8, 0xfc, 0x83, 0x47, 0xcc, 0xfe,
	0x6e, 0x04, 0xf0, 0xb7, 0x21, 0x01, 0x88, 0x26, 0xb3, 0xd9, 0x3c, 0xb9, 0xea, 0x3e, 0x31, 0x98,
	0xcd, 0x3f, 0xcd, 0x67, 0x57, 0xdd, 0x80, 0xb4, 0x20, 0x3c, 0x61, 0x17, 0x49, 0x37, 0x8c, 0xc7,
	0x10, 0x4d, 0xec, 0x07, 0x44, 0xde, 0x40, 0xed, 0x44, 0x90, 0x83, 0x7f, 0x4f, 0xdb, 0x3b, 0x7c,
	0x10, 0x77, 0xda, 0xf8, 0x5a, 0xe7, 0xc5, 0xea, 0x36, 0xb2, 0xf2, 0xf1, 0x9f, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x4a, 0x3f, 0x9d, 0x91, 0x87, 0x03, 0x00, 0x00,
}
