// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: types_v1.proto

package crud_api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Kind struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Major         string                 `protobuf:"bytes,1,opt,name=major,proto3" json:"major,omitempty"`
	Minor         string                 `protobuf:"bytes,2,opt,name=minor,proto3" json:"minor,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Kind) Reset() {
	*x = Kind{}
	mi := &file_types_v1_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Kind) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Kind) ProtoMessage() {}

func (x *Kind) ProtoReflect() protoreflect.Message {
	mi := &file_types_v1_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Kind.ProtoReflect.Descriptor instead.
func (*Kind) Descriptor() ([]byte, []int) {
	return file_types_v1_proto_rawDescGZIP(), []int{0}
}

func (x *Kind) GetMajor() string {
	if x != nil {
		return x.Major
	}
	return ""
}

func (x *Kind) GetMinor() string {
	if x != nil {
		return x.Minor
	}
	return ""
}

type TimeBasedValue struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	StartTime     string                 `protobuf:"bytes,1,opt,name=startTime,proto3" json:"startTime,omitempty"`
	EndTime       string                 `protobuf:"bytes,2,opt,name=endTime,proto3" json:"endTime,omitempty"`
	Value         *anypb.Any             `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"` // Storing any type of value
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TimeBasedValue) Reset() {
	*x = TimeBasedValue{}
	mi := &file_types_v1_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TimeBasedValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeBasedValue) ProtoMessage() {}

func (x *TimeBasedValue) ProtoReflect() protoreflect.Message {
	mi := &file_types_v1_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimeBasedValue.ProtoReflect.Descriptor instead.
func (*TimeBasedValue) Descriptor() ([]byte, []int) {
	return file_types_v1_proto_rawDescGZIP(), []int{1}
}

func (x *TimeBasedValue) GetStartTime() string {
	if x != nil {
		return x.StartTime
	}
	return ""
}

func (x *TimeBasedValue) GetEndTime() string {
	if x != nil {
		return x.EndTime
	}
	return ""
}

func (x *TimeBasedValue) GetValue() *anypb.Any {
	if x != nil {
		return x.Value
	}
	return nil
}

type Relationship struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	RelatedEntityId string                 `protobuf:"bytes,1,opt,name=relatedEntityId,proto3" json:"relatedEntityId,omitempty"`
	StartTime       string                 `protobuf:"bytes,2,opt,name=startTime,proto3" json:"startTime,omitempty"`
	EndTime         string                 `protobuf:"bytes,3,opt,name=endTime,proto3" json:"endTime,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *Relationship) Reset() {
	*x = Relationship{}
	mi := &file_types_v1_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Relationship) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Relationship) ProtoMessage() {}

func (x *Relationship) ProtoReflect() protoreflect.Message {
	mi := &file_types_v1_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Relationship.ProtoReflect.Descriptor instead.
func (*Relationship) Descriptor() ([]byte, []int) {
	return file_types_v1_proto_rawDescGZIP(), []int{2}
}

func (x *Relationship) GetRelatedEntityId() string {
	if x != nil {
		return x.RelatedEntityId
	}
	return ""
}

func (x *Relationship) GetStartTime() string {
	if x != nil {
		return x.StartTime
	}
	return ""
}

func (x *Relationship) GetEndTime() string {
	if x != nil {
		return x.EndTime
	}
	return ""
}

type Entity struct {
	state         protoimpl.MessageState         `protogen:"open.v1"`
	Id            string                         `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`                 // Read-only unique identifier
	Kind          *Kind                          `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`             // Read-only entity type
	Created       string                         `protobuf:"bytes,3,opt,name=created,proto3" json:"created,omitempty"`       // Read-only created timestamp
	Terminated    string                         `protobuf:"bytes,4,opt,name=terminated,proto3" json:"terminated,omitempty"` // Nullable terminated timestamp
	Name          *TimeBasedValue                `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	Metadata      map[string]*anypb.Any          `protobuf:"bytes,6,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`           // Metadata as a flexible key-value map
	Attributes    map[string]*TimeBasedValueList `protobuf:"bytes,7,rep,name=attributes,proto3" json:"attributes,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`       // Attributes as a time-based list
	Relationships map[string]*Relationship       `protobuf:"bytes,8,rep,name=relationships,proto3" json:"relationships,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"` // Relationships to other entities
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Entity) Reset() {
	*x = Entity{}
	mi := &file_types_v1_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Entity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Entity) ProtoMessage() {}

func (x *Entity) ProtoReflect() protoreflect.Message {
	mi := &file_types_v1_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Entity.ProtoReflect.Descriptor instead.
func (*Entity) Descriptor() ([]byte, []int) {
	return file_types_v1_proto_rawDescGZIP(), []int{3}
}

func (x *Entity) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Entity) GetKind() *Kind {
	if x != nil {
		return x.Kind
	}
	return nil
}

func (x *Entity) GetCreated() string {
	if x != nil {
		return x.Created
	}
	return ""
}

func (x *Entity) GetTerminated() string {
	if x != nil {
		return x.Terminated
	}
	return ""
}

func (x *Entity) GetName() *TimeBasedValue {
	if x != nil {
		return x.Name
	}
	return nil
}

func (x *Entity) GetMetadata() map[string]*anypb.Any {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *Entity) GetAttributes() map[string]*TimeBasedValueList {
	if x != nil {
		return x.Attributes
	}
	return nil
}

func (x *Entity) GetRelationships() map[string]*Relationship {
	if x != nil {
		return x.Relationships
	}
	return nil
}

// Wrapper for a repeated TimeBasedValue (since Protobuf does not support nested lists in maps)
type TimeBasedValueList struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Values        []*TimeBasedValue      `protobuf:"bytes,1,rep,name=values,proto3" json:"values,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TimeBasedValueList) Reset() {
	*x = TimeBasedValueList{}
	mi := &file_types_v1_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TimeBasedValueList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeBasedValueList) ProtoMessage() {}

func (x *TimeBasedValueList) ProtoReflect() protoreflect.Message {
	mi := &file_types_v1_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimeBasedValueList.ProtoReflect.Descriptor instead.
func (*TimeBasedValueList) Descriptor() ([]byte, []int) {
	return file_types_v1_proto_rawDescGZIP(), []int{4}
}

func (x *TimeBasedValueList) GetValues() []*TimeBasedValue {
	if x != nil {
		return x.Values
	}
	return nil
}

// Request message for deleting an entity by ID
type EntityId struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EntityId) Reset() {
	*x = EntityId{}
	mi := &file_types_v1_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EntityId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntityId) ProtoMessage() {}

func (x *EntityId) ProtoReflect() protoreflect.Message {
	mi := &file_types_v1_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntityId.ProtoReflect.Descriptor instead.
func (*EntityId) Descriptor() ([]byte, []int) {
	return file_types_v1_proto_rawDescGZIP(), []int{5}
}

func (x *EntityId) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// Request message for updating an entity
type UpdateEntityRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Entity        *Entity                `protobuf:"bytes,2,opt,name=entity,proto3" json:"entity,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateEntityRequest) Reset() {
	*x = UpdateEntityRequest{}
	mi := &file_types_v1_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateEntityRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateEntityRequest) ProtoMessage() {}

func (x *UpdateEntityRequest) ProtoReflect() protoreflect.Message {
	mi := &file_types_v1_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateEntityRequest.ProtoReflect.Descriptor instead.
func (*UpdateEntityRequest) Descriptor() ([]byte, []int) {
	return file_types_v1_proto_rawDescGZIP(), []int{6}
}

func (x *UpdateEntityRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateEntityRequest) GetEntity() *Entity {
	if x != nil {
		return x.Entity
	}
	return nil
}

// Empty message response
type Empty struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Empty) Reset() {
	*x = Empty{}
	mi := &file_types_v1_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_types_v1_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_types_v1_proto_rawDescGZIP(), []int{7}
}

var File_types_v1_proto protoreflect.FileDescriptor

var file_types_v1_proto_rawDesc = string([]byte{
	0x0a, 0x0e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x5f, 0x76, 0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x04, 0x63, 0x72, 0x75, 0x64, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x32, 0x0a, 0x04, 0x4b, 0x69, 0x6e, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x6d, 0x61, 0x6a,
	0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6d, 0x61, 0x6a, 0x6f, 0x72, 0x12,
	0x14, 0x0a, 0x05, 0x6d, 0x69, 0x6e, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x6d, 0x69, 0x6e, 0x6f, 0x72, 0x22, 0x74, 0x0a, 0x0e, 0x54, 0x69, 0x6d, 0x65, 0x42, 0x61, 0x73,
	0x65, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x54, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x12,
	0x2a, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x41, 0x6e, 0x79, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x70, 0x0a, 0x0c, 0x52,
	0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x68, 0x69, 0x70, 0x12, 0x28, 0x0a, 0x0f, 0x72,
	0x65, 0x6c, 0x61, 0x74, 0x65, 0x64, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x49, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x65, 0x64, 0x45, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54,
	0x69, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x22, 0xdb, 0x04,
	0x0a, 0x06, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1e, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x63, 0x72, 0x75, 0x64, 0x2e, 0x4b, 0x69,
	0x6e, 0x64, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x74, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x64,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x74,
	0x65, 0x64, 0x12, 0x28, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x63, 0x72, 0x75, 0x64, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x42, 0x61, 0x73, 0x65,
	0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x36, 0x0a, 0x08,
	0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x63, 0x72, 0x75, 0x64, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x12, 0x3c, 0x0a, 0x0a, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74,
	0x65, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x63, 0x72, 0x75, 0x64, 0x2e,
	0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0a, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74,
	0x65, 0x73, 0x12, 0x45, 0x0a, 0x0d, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x68,
	0x69, 0x70, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x63, 0x72, 0x75, 0x64,
	0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x68, 0x69, 0x70, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0d, 0x72, 0x65, 0x6c, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x68, 0x69, 0x70, 0x73, 0x1a, 0x51, 0x0a, 0x0d, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2a, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e,
	0x79, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x57, 0x0a, 0x0f,
	0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x2e, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x18, 0x2e, 0x63, 0x72, 0x75, 0x64, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x42, 0x61, 0x73, 0x65,
	0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x54, 0x0a, 0x12, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x68, 0x69, 0x70, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x28, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x63,
	0x72, 0x75, 0x64, 0x2e, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x68, 0x69, 0x70,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x42, 0x0a, 0x12, 0x54,
	0x69, 0x6d, 0x65, 0x42, 0x61, 0x73, 0x65, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x4c, 0x69, 0x73,
	0x74, 0x12, 0x2c, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x14, 0x2e, 0x63, 0x72, 0x75, 0x64, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x42, 0x61, 0x73,
	0x65, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x22,
	0x1a, 0x0a, 0x08, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x49, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x4b, 0x0a, 0x13, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x24, 0x0a, 0x06, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x63, 0x72, 0x75, 0x64, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x52, 0x06, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x32, 0xcb, 0x01, 0x0a, 0x0b, 0x43, 0x72, 0x75, 0x64, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x2a, 0x0a, 0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x45, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x12, 0x0c, 0x2e, 0x63, 0x72, 0x75, 0x64, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x1a,
	0x0c, 0x2e, 0x63, 0x72, 0x75, 0x64, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x2a, 0x0a,
	0x0a, 0x52, 0x65, 0x61, 0x64, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x0e, 0x2e, 0x63, 0x72,
	0x75, 0x64, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x49, 0x64, 0x1a, 0x0c, 0x2e, 0x63, 0x72,
	0x75, 0x64, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x37, 0x0a, 0x0c, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x19, 0x2e, 0x63, 0x72, 0x75, 0x64,
	0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x0c, 0x2e, 0x63, 0x72, 0x75, 0x64, 0x2e, 0x45, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x12, 0x2b, 0x0a, 0x0c, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x45, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x12, 0x0e, 0x2e, 0x63, 0x72, 0x75, 0x64, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x49, 0x64, 0x1a, 0x0b, 0x2e, 0x63, 0x72, 0x75, 0x64, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42,
	0x1c, 0x5a, 0x1a, 0x6c, 0x6b, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x66, 0x6f, 0x75, 0x6e, 0x64, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x63, 0x72, 0x75, 0x64, 0x2d, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_types_v1_proto_rawDescOnce sync.Once
	file_types_v1_proto_rawDescData []byte
)

func file_types_v1_proto_rawDescGZIP() []byte {
	file_types_v1_proto_rawDescOnce.Do(func() {
		file_types_v1_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_types_v1_proto_rawDesc), len(file_types_v1_proto_rawDesc)))
	})
	return file_types_v1_proto_rawDescData
}

var file_types_v1_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_types_v1_proto_goTypes = []any{
	(*Kind)(nil),                // 0: crud.Kind
	(*TimeBasedValue)(nil),      // 1: crud.TimeBasedValue
	(*Relationship)(nil),        // 2: crud.Relationship
	(*Entity)(nil),              // 3: crud.Entity
	(*TimeBasedValueList)(nil),  // 4: crud.TimeBasedValueList
	(*EntityId)(nil),            // 5: crud.EntityId
	(*UpdateEntityRequest)(nil), // 6: crud.UpdateEntityRequest
	(*Empty)(nil),               // 7: crud.Empty
	nil,                         // 8: crud.Entity.MetadataEntry
	nil,                         // 9: crud.Entity.AttributesEntry
	nil,                         // 10: crud.Entity.RelationshipsEntry
	(*anypb.Any)(nil),           // 11: google.protobuf.Any
}
var file_types_v1_proto_depIdxs = []int32{
	11, // 0: crud.TimeBasedValue.value:type_name -> google.protobuf.Any
	0,  // 1: crud.Entity.kind:type_name -> crud.Kind
	1,  // 2: crud.Entity.name:type_name -> crud.TimeBasedValue
	8,  // 3: crud.Entity.metadata:type_name -> crud.Entity.MetadataEntry
	9,  // 4: crud.Entity.attributes:type_name -> crud.Entity.AttributesEntry
	10, // 5: crud.Entity.relationships:type_name -> crud.Entity.RelationshipsEntry
	1,  // 6: crud.TimeBasedValueList.values:type_name -> crud.TimeBasedValue
	3,  // 7: crud.UpdateEntityRequest.entity:type_name -> crud.Entity
	11, // 8: crud.Entity.MetadataEntry.value:type_name -> google.protobuf.Any
	4,  // 9: crud.Entity.AttributesEntry.value:type_name -> crud.TimeBasedValueList
	2,  // 10: crud.Entity.RelationshipsEntry.value:type_name -> crud.Relationship
	3,  // 11: crud.CrudService.CreateEntity:input_type -> crud.Entity
	5,  // 12: crud.CrudService.ReadEntity:input_type -> crud.EntityId
	6,  // 13: crud.CrudService.UpdateEntity:input_type -> crud.UpdateEntityRequest
	5,  // 14: crud.CrudService.DeleteEntity:input_type -> crud.EntityId
	3,  // 15: crud.CrudService.CreateEntity:output_type -> crud.Entity
	3,  // 16: crud.CrudService.ReadEntity:output_type -> crud.Entity
	3,  // 17: crud.CrudService.UpdateEntity:output_type -> crud.Entity
	7,  // 18: crud.CrudService.DeleteEntity:output_type -> crud.Empty
	15, // [15:19] is the sub-list for method output_type
	11, // [11:15] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_types_v1_proto_init() }
func file_types_v1_proto_init() {
	if File_types_v1_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_types_v1_proto_rawDesc), len(file_types_v1_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_types_v1_proto_goTypes,
		DependencyIndexes: file_types_v1_proto_depIdxs,
		MessageInfos:      file_types_v1_proto_msgTypes,
	}.Build()
	File_types_v1_proto = out.File
	file_types_v1_proto_goTypes = nil
	file_types_v1_proto_depIdxs = nil
}
