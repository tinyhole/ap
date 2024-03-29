// Code generated by protoc-gen-go. DO NOT EDIT.
// source: hello.proto

package hello

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

type SayHelloReq struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SayHelloReq) Reset()         { *m = SayHelloReq{} }
func (m *SayHelloReq) String() string { return proto.CompactTextString(m) }
func (*SayHelloReq) ProtoMessage()    {}
func (*SayHelloReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_61ef911816e0a8ce, []int{0}
}

func (m *SayHelloReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SayHelloReq.Unmarshal(m, b)
}
func (m *SayHelloReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SayHelloReq.Marshal(b, m, deterministic)
}
func (m *SayHelloReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SayHelloReq.Merge(m, src)
}
func (m *SayHelloReq) XXX_Size() int {
	return xxx_messageInfo_SayHelloReq.Size(m)
}
func (m *SayHelloReq) XXX_DiscardUnknown() {
	xxx_messageInfo_SayHelloReq.DiscardUnknown(m)
}

var xxx_messageInfo_SayHelloReq proto.InternalMessageInfo

func (m *SayHelloReq) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type SayHelloRsp struct {
	Reply                string   `protobuf:"bytes,1,opt,name=reply,proto3" json:"reply,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SayHelloRsp) Reset()         { *m = SayHelloRsp{} }
func (m *SayHelloRsp) String() string { return proto.CompactTextString(m) }
func (*SayHelloRsp) ProtoMessage()    {}
func (*SayHelloRsp) Descriptor() ([]byte, []int) {
	return fileDescriptor_61ef911816e0a8ce, []int{1}
}

func (m *SayHelloRsp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SayHelloRsp.Unmarshal(m, b)
}
func (m *SayHelloRsp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SayHelloRsp.Marshal(b, m, deterministic)
}
func (m *SayHelloRsp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SayHelloRsp.Merge(m, src)
}
func (m *SayHelloRsp) XXX_Size() int {
	return xxx_messageInfo_SayHelloRsp.Size(m)
}
func (m *SayHelloRsp) XXX_DiscardUnknown() {
	xxx_messageInfo_SayHelloRsp.DiscardUnknown(m)
}

var xxx_messageInfo_SayHelloRsp proto.InternalMessageInfo

func (m *SayHelloRsp) GetReply() string {
	if m != nil {
		return m.Reply
	}
	return ""
}

func init() {
	proto.RegisterType((*SayHelloReq)(nil), "hello.SayHelloReq")
	proto.RegisterType((*SayHelloRsp)(nil), "hello.SayHelloRsp")
}

func init() { proto.RegisterFile("hello.proto", fileDescriptor_61ef911816e0a8ce) }

var fileDescriptor_61ef911816e0a8ce = []byte{
	// 124 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xce, 0x48, 0xcd, 0xc9,
	0xc9, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x73, 0x94, 0x14, 0xb9, 0xb8, 0x83,
	0x13, 0x2b, 0x3d, 0x40, 0xec, 0xa0, 0xd4, 0x42, 0x21, 0x21, 0x2e, 0x96, 0xbc, 0xc4, 0xdc, 0x54,
	0x09, 0x46, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x30, 0x5b, 0x49, 0x19, 0x49, 0x49, 0x71, 0x81, 0x90,
	0x08, 0x17, 0x6b, 0x51, 0x6a, 0x41, 0x4e, 0x25, 0x54, 0x0d, 0x84, 0x63, 0x64, 0xcd, 0xc5, 0x0a,
	0x56, 0x21, 0x64, 0xc4, 0xc5, 0x01, 0x53, 0x2d, 0x24, 0xa4, 0x07, 0xb1, 0x11, 0xc9, 0x06, 0x29,
	0x0c, 0xb1, 0xe2, 0x02, 0x27, 0xf6, 0x28, 0x88, 0x6b, 0x92, 0xd8, 0xc0, 0x6e, 0x33, 0x06, 0x04,
	0x00, 0x00, 0xff, 0xff, 0x61, 0x5e, 0xb1, 0x25, 0xaa, 0x00, 0x00, 0x00,
}
