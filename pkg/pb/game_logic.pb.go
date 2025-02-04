// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.3
// 	protoc        v5.29.3
// source: game_logic.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ActionRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"` // Исправлено: добавлена точка с запятой
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ActionRequest) Reset() {
	*x = ActionRequest{}
	mi := &file_game_logic_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ActionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActionRequest) ProtoMessage() {}

func (x *ActionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_game_logic_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActionRequest.ProtoReflect.Descriptor instead.
func (*ActionRequest) Descriptor() ([]byte, []int) {
	return file_game_logic_proto_rawDescGZIP(), []int{0}
}

func (x *ActionRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type ActionResponse struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	Id              int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId          int64                  `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	AreaId          int64                  `protobuf:"varint,3,opt,name=area_id,json=areaId,proto3" json:"area_id,omitempty"`
	ObjectSourceId  int64                  `protobuf:"varint,4,opt,name=object_source_id,json=objectSourceId,proto3" json:"object_source_id,omitempty"`
	ObjectDestId    int64                  `protobuf:"varint,5,opt,name=object_dest_id,json=objectDestId,proto3" json:"object_dest_id,omitempty"`
	ActionType      string                 `protobuf:"bytes,6,opt,name=action_type,json=actionType,proto3" json:"action_type,omitempty"`
	Characteristics string                 `protobuf:"bytes,7,opt,name=characteristics,proto3" json:"characteristics,omitempty"`
	StartTime       *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	Duration        *durationpb.Duration   `protobuf:"bytes,9,opt,name=duration,proto3" json:"duration,omitempty"`
	Status          string                 `protobuf:"bytes,10,opt,name=status,proto3" json:"status,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *ActionResponse) Reset() {
	*x = ActionResponse{}
	mi := &file_game_logic_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ActionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActionResponse) ProtoMessage() {}

func (x *ActionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_game_logic_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActionResponse.ProtoReflect.Descriptor instead.
func (*ActionResponse) Descriptor() ([]byte, []int) {
	return file_game_logic_proto_rawDescGZIP(), []int{1}
}

func (x *ActionResponse) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ActionResponse) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *ActionResponse) GetAreaId() int64 {
	if x != nil {
		return x.AreaId
	}
	return 0
}

func (x *ActionResponse) GetObjectSourceId() int64 {
	if x != nil {
		return x.ObjectSourceId
	}
	return 0
}

func (x *ActionResponse) GetObjectDestId() int64 {
	if x != nil {
		return x.ObjectDestId
	}
	return 0
}

func (x *ActionResponse) GetActionType() string {
	if x != nil {
		return x.ActionType
	}
	return ""
}

func (x *ActionResponse) GetCharacteristics() string {
	if x != nil {
		return x.Characteristics
	}
	return ""
}

func (x *ActionResponse) GetStartTime() *timestamppb.Timestamp {
	if x != nil {
		return x.StartTime
	}
	return nil
}

func (x *ActionResponse) GetDuration() *durationpb.Duration {
	if x != nil {
		return x.Duration
	}
	return nil
}

func (x *ActionResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type LogicRequest struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	UserId          int64                  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"` // Исправлено: поле должно начинаться с 1
	AreaId          int64                  `protobuf:"varint,2,opt,name=area_id,json=areaId,proto3" json:"area_id,omitempty"`
	ObjectSourceId  int64                  `protobuf:"varint,3,opt,name=object_source_id,json=objectSourceId,proto3" json:"object_source_id,omitempty"`
	ObjectDestId    int64                  `protobuf:"varint,4,opt,name=object_dest_id,json=objectDestId,proto3" json:"object_dest_id,omitempty"`
	ActionType      string                 `protobuf:"bytes,5,opt,name=action_type,json=actionType,proto3" json:"action_type,omitempty"`
	Characteristics string                 `protobuf:"bytes,6,opt,name=characteristics,proto3" json:"characteristics,omitempty"`
	StartTime       *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	Duration        *durationpb.Duration   `protobuf:"bytes,8,opt,name=duration,proto3" json:"duration,omitempty"`
	Status          string                 `protobuf:"bytes,9,opt,name=status,proto3" json:"status,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *LogicRequest) Reset() {
	*x = LogicRequest{}
	mi := &file_game_logic_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LogicRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogicRequest) ProtoMessage() {}

func (x *LogicRequest) ProtoReflect() protoreflect.Message {
	mi := &file_game_logic_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogicRequest.ProtoReflect.Descriptor instead.
func (*LogicRequest) Descriptor() ([]byte, []int) {
	return file_game_logic_proto_rawDescGZIP(), []int{2}
}

func (x *LogicRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *LogicRequest) GetAreaId() int64 {
	if x != nil {
		return x.AreaId
	}
	return 0
}

func (x *LogicRequest) GetObjectSourceId() int64 {
	if x != nil {
		return x.ObjectSourceId
	}
	return 0
}

func (x *LogicRequest) GetObjectDestId() int64 {
	if x != nil {
		return x.ObjectDestId
	}
	return 0
}

func (x *LogicRequest) GetActionType() string {
	if x != nil {
		return x.ActionType
	}
	return ""
}

func (x *LogicRequest) GetCharacteristics() string {
	if x != nil {
		return x.Characteristics
	}
	return ""
}

func (x *LogicRequest) GetStartTime() *timestamppb.Timestamp {
	if x != nil {
		return x.StartTime
	}
	return nil
}

func (x *LogicRequest) GetDuration() *durationpb.Duration {
	if x != nil {
		return x.Duration
	}
	return nil
}

func (x *LogicRequest) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

var File_game_logic_proto protoreflect.FileDescriptor

var file_game_logic_proto_rawDesc = []byte{
	0x0a, 0x10, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0a, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x1a, 0x1b,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x1f, 0x0a, 0x0d,
	0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x22, 0xf7, 0x02,
	0x0a, 0x0e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x61, 0x72, 0x65,
	0x61, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x61, 0x72, 0x65, 0x61,
	0x49, 0x64, 0x12, 0x28, 0x0a, 0x10, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0e, 0x6f, 0x62,
	0x6a, 0x65, 0x63, 0x74, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0e,
	0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x64, 0x65, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x44, 0x65, 0x73, 0x74,
	0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x28, 0x0a, 0x0f, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72,
	0x69, 0x73, 0x74, 0x69, 0x63, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x63, 0x68,
	0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x69, 0x73, 0x74, 0x69, 0x63, 0x73, 0x12, 0x39, 0x0a,
	0x0a, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x73,
	0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x35, 0x0a, 0x08, 0x64, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0xe5, 0x02, 0x0a, 0x0c, 0x4c, 0x6f, 0x67, 0x69,
	0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x12, 0x17, 0x0a, 0x07, 0x61, 0x72, 0x65, 0x61, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x06, 0x61, 0x72, 0x65, 0x61, 0x49, 0x64, 0x12, 0x28, 0x0a, 0x10, 0x6f, 0x62,
	0x6a, 0x65, 0x63, 0x74, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0e, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x53, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0e, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x64,
	0x65, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x6f, 0x62,
	0x6a, 0x65, 0x63, 0x74, 0x44, 0x65, 0x73, 0x74, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x28, 0x0a, 0x0f, 0x63,
	0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x69, 0x73, 0x74, 0x69, 0x63, 0x73, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x69,
	0x73, 0x74, 0x69, 0x63, 0x73, 0x12, 0x39, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x74,
	0x69, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65,
	0x12, 0x35, 0x0a, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x64,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x32,
	0x95, 0x01, 0x0a, 0x10, 0x47, 0x61, 0x6d, 0x65, 0x4c, 0x6f, 0x67, 0x69, 0x63, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x42, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x41, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x19, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x2e, 0x41,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x67,
	0x61, 0x6d, 0x65, 0x5f, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3d, 0x0a, 0x09, 0x41, 0x64, 0x64, 0x41,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x6c, 0x6f, 0x67,
	0x69, 0x63, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x0e, 0x5a, 0x0c, 0x63, 0x79, 0x62, 0x65, 0x72,
	0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_game_logic_proto_rawDescOnce sync.Once
	file_game_logic_proto_rawDescData = file_game_logic_proto_rawDesc
)

func file_game_logic_proto_rawDescGZIP() []byte {
	file_game_logic_proto_rawDescOnce.Do(func() {
		file_game_logic_proto_rawDescData = protoimpl.X.CompressGZIP(file_game_logic_proto_rawDescData)
	})
	return file_game_logic_proto_rawDescData
}

var file_game_logic_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_game_logic_proto_goTypes = []any{
	(*ActionRequest)(nil),         // 0: game_logic.ActionRequest
	(*ActionResponse)(nil),        // 1: game_logic.ActionResponse
	(*LogicRequest)(nil),          // 2: game_logic.LogicRequest
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
	(*durationpb.Duration)(nil),   // 4: google.protobuf.Duration
	(*emptypb.Empty)(nil),         // 5: google.protobuf.Empty
}
var file_game_logic_proto_depIdxs = []int32{
	3, // 0: game_logic.ActionResponse.start_time:type_name -> google.protobuf.Timestamp
	4, // 1: game_logic.ActionResponse.duration:type_name -> google.protobuf.Duration
	3, // 2: game_logic.LogicRequest.start_time:type_name -> google.protobuf.Timestamp
	4, // 3: game_logic.LogicRequest.duration:type_name -> google.protobuf.Duration
	0, // 4: game_logic.GameLogicService.GetAction:input_type -> game_logic.ActionRequest
	2, // 5: game_logic.GameLogicService.AddAction:input_type -> game_logic.LogicRequest
	1, // 6: game_logic.GameLogicService.GetAction:output_type -> game_logic.ActionResponse
	5, // 7: game_logic.GameLogicService.AddAction:output_type -> google.protobuf.Empty
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_game_logic_proto_init() }
func file_game_logic_proto_init() {
	if File_game_logic_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_game_logic_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_game_logic_proto_goTypes,
		DependencyIndexes: file_game_logic_proto_depIdxs,
		MessageInfos:      file_game_logic_proto_msgTypes,
	}.Build()
	File_game_logic_proto = out.File
	file_game_logic_proto_rawDesc = nil
	file_game_logic_proto_goTypes = nil
	file_game_logic_proto_depIdxs = nil
}
