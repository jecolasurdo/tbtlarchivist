// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: protobuf/contracts.proto

package contracts

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type ClipInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InitialDateCurated *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=initial_date_curated,json=initialDateCurated,proto3" json:"initial_date_curated,omitempty"`
	LastDateCurated    *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=last_date_curated,json=lastDateCurated,proto3" json:"last_date_curated,omitempty"`
	CuratorInformation string                 `protobuf:"bytes,3,opt,name=curator_information,json=curatorInformation,proto3" json:"curator_information,omitempty"`
	Title              string                 `protobuf:"bytes,4,opt,name=title,proto3" json:"title,omitempty"`
	Description        string                 `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`
	MediaUri           string                 `protobuf:"bytes,6,opt,name=media_uri,json=mediaUri,proto3" json:"media_uri,omitempty"`
	MediaType          string                 `protobuf:"bytes,7,opt,name=media_type,json=mediaType,proto3" json:"media_type,omitempty"`
	Priority           int32                  `protobuf:"varint,8,opt,name=priority,proto3" json:"priority,omitempty"`
}

func (x *ClipInfo) Reset() {
	*x = ClipInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_contracts_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClipInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClipInfo) ProtoMessage() {}

func (x *ClipInfo) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_contracts_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClipInfo.ProtoReflect.Descriptor instead.
func (*ClipInfo) Descriptor() ([]byte, []int) {
	return file_protobuf_contracts_proto_rawDescGZIP(), []int{0}
}

func (x *ClipInfo) GetInitialDateCurated() *timestamppb.Timestamp {
	if x != nil {
		return x.InitialDateCurated
	}
	return nil
}

func (x *ClipInfo) GetLastDateCurated() *timestamppb.Timestamp {
	if x != nil {
		return x.LastDateCurated
	}
	return nil
}

func (x *ClipInfo) GetCuratorInformation() string {
	if x != nil {
		return x.CuratorInformation
	}
	return ""
}

func (x *ClipInfo) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *ClipInfo) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *ClipInfo) GetMediaUri() string {
	if x != nil {
		return x.MediaUri
	}
	return ""
}

func (x *ClipInfo) GetMediaType() string {
	if x != nil {
		return x.MediaType
	}
	return ""
}

func (x *ClipInfo) GetPriority() int32 {
	if x != nil {
		return x.Priority
	}
	return 0
}

type EpisodeInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InitialDateCurated *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=initial_date_curated,json=initialDateCurated,proto3" json:"initial_date_curated,omitempty"`
	LastDateCurated    *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=last_date_curated,json=lastDateCurated,proto3" json:"last_date_curated,omitempty"`
	CuratorInformation string                 `protobuf:"bytes,3,opt,name=curator_information,json=curatorInformation,proto3" json:"curator_information,omitempty"`
	DateAired          *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=date_aired,json=dateAired,proto3" json:"date_aired,omitempty"`
	Title              string                 `protobuf:"bytes,5,opt,name=title,proto3" json:"title,omitempty"`
	Description        string                 `protobuf:"bytes,6,opt,name=description,proto3" json:"description,omitempty"`
	MediaUri           string                 `protobuf:"bytes,7,opt,name=media_uri,json=mediaUri,proto3" json:"media_uri,omitempty"`
	MediaType          string                 `protobuf:"bytes,8,opt,name=media_type,json=mediaType,proto3" json:"media_type,omitempty"`
	Priority           int32                  `protobuf:"varint,9,opt,name=priority,proto3" json:"priority,omitempty"`
}

func (x *EpisodeInfo) Reset() {
	*x = EpisodeInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_contracts_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EpisodeInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EpisodeInfo) ProtoMessage() {}

func (x *EpisodeInfo) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_contracts_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EpisodeInfo.ProtoReflect.Descriptor instead.
func (*EpisodeInfo) Descriptor() ([]byte, []int) {
	return file_protobuf_contracts_proto_rawDescGZIP(), []int{1}
}

func (x *EpisodeInfo) GetInitialDateCurated() *timestamppb.Timestamp {
	if x != nil {
		return x.InitialDateCurated
	}
	return nil
}

func (x *EpisodeInfo) GetLastDateCurated() *timestamppb.Timestamp {
	if x != nil {
		return x.LastDateCurated
	}
	return nil
}

func (x *EpisodeInfo) GetCuratorInformation() string {
	if x != nil {
		return x.CuratorInformation
	}
	return ""
}

func (x *EpisodeInfo) GetDateAired() *timestamppb.Timestamp {
	if x != nil {
		return x.DateAired
	}
	return nil
}

func (x *EpisodeInfo) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *EpisodeInfo) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *EpisodeInfo) GetMediaUri() string {
	if x != nil {
		return x.MediaUri
	}
	return ""
}

func (x *EpisodeInfo) GetMediaType() string {
	if x != nil {
		return x.MediaType
	}
	return ""
}

func (x *EpisodeInfo) GetPriority() int32 {
	if x != nil {
		return x.Priority
	}
	return 0
}

type PendingResearchItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LeaseId string       `protobuf:"bytes,1,opt,name=lease_id,json=leaseId,proto3" json:"lease_id,omitempty"`
	Episode *EpisodeInfo `protobuf:"bytes,2,opt,name=episode,proto3" json:"episode,omitempty"`
	Clips   []*ClipInfo  `protobuf:"bytes,3,rep,name=clips,proto3" json:"clips,omitempty"`
}

func (x *PendingResearchItem) Reset() {
	*x = PendingResearchItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_contracts_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PendingResearchItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PendingResearchItem) ProtoMessage() {}

func (x *PendingResearchItem) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_contracts_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PendingResearchItem.ProtoReflect.Descriptor instead.
func (*PendingResearchItem) Descriptor() ([]byte, []int) {
	return file_protobuf_contracts_proto_rawDescGZIP(), []int{2}
}

func (x *PendingResearchItem) GetLeaseId() string {
	if x != nil {
		return x.LeaseId
	}
	return ""
}

func (x *PendingResearchItem) GetEpisode() *EpisodeInfo {
	if x != nil {
		return x.Episode
	}
	return nil
}

func (x *PendingResearchItem) GetClips() []*ClipInfo {
	if x != nil {
		return x.Clips
	}
	return nil
}

type CompletedResearchItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResearchDate    *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=research_date,json=researchDate,proto3" json:"research_date,omitempty"`
	EpisodeInfo     *EpisodeInfo           `protobuf:"bytes,2,opt,name=episode_info,json=episodeInfo,proto3" json:"episode_info,omitempty"`
	ClipInfo        *ClipInfo              `protobuf:"bytes,3,opt,name=clip_info,json=clipInfo,proto3" json:"clip_info,omitempty"`
	EpisodeDuration int64                  `protobuf:"varint,4,opt,name=episode_duration,json=episodeDuration,proto3" json:"episode_duration,omitempty"`
	EpisodeHash     string                 `protobuf:"bytes,5,opt,name=episode_hash,json=episodeHash,proto3" json:"episode_hash,omitempty"`
	ClipDuration    int64                  `protobuf:"varint,6,opt,name=clip_duration,json=clipDuration,proto3" json:"clip_duration,omitempty"`
	ClipHash        string                 `protobuf:"bytes,7,opt,name=clip_hash,json=clipHash,proto3" json:"clip_hash,omitempty"`
	ClipOffsets     []int64                `protobuf:"varint,8,rep,packed,name=clip_offsets,json=clipOffsets,proto3" json:"clip_offsets,omitempty"`
	LeaseId         string                 `protobuf:"bytes,9,opt,name=lease_id,json=leaseId,proto3" json:"lease_id,omitempty"`
	RevokeLease     bool                   `protobuf:"varint,10,opt,name=revoke_lease,json=revokeLease,proto3" json:"revoke_lease,omitempty"`
}

func (x *CompletedResearchItem) Reset() {
	*x = CompletedResearchItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_contracts_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CompletedResearchItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CompletedResearchItem) ProtoMessage() {}

func (x *CompletedResearchItem) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_contracts_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CompletedResearchItem.ProtoReflect.Descriptor instead.
func (*CompletedResearchItem) Descriptor() ([]byte, []int) {
	return file_protobuf_contracts_proto_rawDescGZIP(), []int{3}
}

func (x *CompletedResearchItem) GetResearchDate() *timestamppb.Timestamp {
	if x != nil {
		return x.ResearchDate
	}
	return nil
}

func (x *CompletedResearchItem) GetEpisodeInfo() *EpisodeInfo {
	if x != nil {
		return x.EpisodeInfo
	}
	return nil
}

func (x *CompletedResearchItem) GetClipInfo() *ClipInfo {
	if x != nil {
		return x.ClipInfo
	}
	return nil
}

func (x *CompletedResearchItem) GetEpisodeDuration() int64 {
	if x != nil {
		return x.EpisodeDuration
	}
	return 0
}

func (x *CompletedResearchItem) GetEpisodeHash() string {
	if x != nil {
		return x.EpisodeHash
	}
	return ""
}

func (x *CompletedResearchItem) GetClipDuration() int64 {
	if x != nil {
		return x.ClipDuration
	}
	return 0
}

func (x *CompletedResearchItem) GetClipHash() string {
	if x != nil {
		return x.ClipHash
	}
	return ""
}

func (x *CompletedResearchItem) GetClipOffsets() []int64 {
	if x != nil {
		return x.ClipOffsets
	}
	return nil
}

func (x *CompletedResearchItem) GetLeaseId() string {
	if x != nil {
		return x.LeaseId
	}
	return ""
}

func (x *CompletedResearchItem) GetRevokeLease() bool {
	if x != nil {
		return x.RevokeLease
	}
	return false
}

var File_protobuf_contracts_proto protoreflect.FileDescriptor

var file_protobuf_contracts_proto_rawDesc = []byte{
	0x0a, 0x18, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72,
	0x61, 0x63, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x61, 0x63, 0x74, 0x73, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe1, 0x02, 0x0a, 0x08, 0x43, 0x6c, 0x69, 0x70, 0x49,
	0x6e, 0x66, 0x6f, 0x12, 0x4c, 0x0a, 0x14, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x5f, 0x64,
	0x61, 0x74, 0x65, 0x5f, 0x63, 0x75, 0x72, 0x61, 0x74, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x12, 0x69,
	0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x44, 0x61, 0x74, 0x65, 0x43, 0x75, 0x72, 0x61, 0x74, 0x65,
	0x64, 0x12, 0x46, 0x0a, 0x11, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x63,
	0x75, 0x72, 0x61, 0x74, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0f, 0x6c, 0x61, 0x73, 0x74, 0x44, 0x61,
	0x74, 0x65, 0x43, 0x75, 0x72, 0x61, 0x74, 0x65, 0x64, 0x12, 0x2f, 0x0a, 0x13, 0x63, 0x75, 0x72,
	0x61, 0x74, 0x6f, 0x72, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x63, 0x75, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x49,
	0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69,
	0x74, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65,
	0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x5f, 0x75, 0x72, 0x69, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x55, 0x72, 0x69, 0x12,
	0x1d, 0x0a, 0x0a, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x22, 0x9f, 0x03, 0x0a, 0x0b, 0x45,
	0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x4c, 0x0a, 0x14, 0x69, 0x6e,
	0x69, 0x74, 0x69, 0x61, 0x6c, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x63, 0x75, 0x72, 0x61, 0x74,
	0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x12, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x44, 0x61, 0x74,
	0x65, 0x43, 0x75, 0x72, 0x61, 0x74, 0x65, 0x64, 0x12, 0x46, 0x0a, 0x11, 0x6c, 0x61, 0x73, 0x74,
	0x5f, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x63, 0x75, 0x72, 0x61, 0x74, 0x65, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x0f, 0x6c, 0x61, 0x73, 0x74, 0x44, 0x61, 0x74, 0x65, 0x43, 0x75, 0x72, 0x61, 0x74, 0x65, 0x64,
	0x12, 0x2f, 0x0a, 0x13, 0x63, 0x75, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x69, 0x6e, 0x66, 0x6f,
	0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x63,
	0x75, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x39, 0x0a, 0x0a, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x61, 0x69, 0x72, 0x65, 0x64, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x09, 0x64, 0x61, 0x74, 0x65, 0x41, 0x69, 0x72, 0x65, 0x64, 0x12, 0x14, 0x0a, 0x05,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74,
	0x6c, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x5f, 0x75, 0x72,
	0x69, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x55, 0x72,
	0x69, 0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x22, 0x8d, 0x01, 0x0a,
	0x13, 0x50, 0x65, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x49, 0x74, 0x65, 0x6d, 0x12, 0x19, 0x0a, 0x08, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x49, 0x64, 0x12,
	0x30, 0x0a, 0x07, 0x65, 0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x16, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x73, 0x2e, 0x45, 0x70, 0x69,
	0x73, 0x6f, 0x64, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x07, 0x65, 0x70, 0x69, 0x73, 0x6f, 0x64,
	0x65, 0x12, 0x29, 0x0a, 0x05, 0x63, 0x6c, 0x69, 0x70, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x13, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x73, 0x2e, 0x43, 0x6c, 0x69,
	0x70, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x05, 0x63, 0x6c, 0x69, 0x70, 0x73, 0x22, 0xb6, 0x03, 0x0a,
	0x15, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x52, 0x65, 0x73, 0x65, 0x61, 0x72,
	0x63, 0x68, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x3f, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x65, 0x61, 0x72,
	0x63, 0x68, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x65, 0x61,
	0x72, 0x63, 0x68, 0x44, 0x61, 0x74, 0x65, 0x12, 0x39, 0x0a, 0x0c, 0x65, 0x70, 0x69, 0x73, 0x6f,
	0x64, 0x65, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x73, 0x2e, 0x45, 0x70, 0x69, 0x73, 0x6f, 0x64,
	0x65, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0b, 0x65, 0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x49, 0x6e,
	0x66, 0x6f, 0x12, 0x30, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x70, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74,
	0x73, 0x2e, 0x43, 0x6c, 0x69, 0x70, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x70,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x29, 0x0a, 0x10, 0x65, 0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x5f,
	0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0f,
	0x65, 0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x21, 0x0a, 0x0c, 0x65, 0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x65, 0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x48, 0x61,
	0x73, 0x68, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x6c, 0x69, 0x70, 0x5f, 0x64, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x63, 0x6c, 0x69, 0x70, 0x44,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x70, 0x5f,
	0x68, 0x61, 0x73, 0x68, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x70,
	0x48, 0x61, 0x73, 0x68, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x6c, 0x69, 0x70, 0x5f, 0x6f, 0x66, 0x66,
	0x73, 0x65, 0x74, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28, 0x03, 0x52, 0x0b, 0x63, 0x6c, 0x69, 0x70,
	0x4f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x73, 0x12, 0x19, 0x0a, 0x08, 0x6c, 0x65, 0x61, 0x73, 0x65,
	0x5f, 0x69, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6c, 0x65, 0x61, 0x73, 0x65,
	0x49, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65, 0x76, 0x6f, 0x6b, 0x65, 0x5f, 0x6c, 0x65, 0x61,
	0x73, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x72, 0x65, 0x76, 0x6f, 0x6b, 0x65,
	0x4c, 0x65, 0x61, 0x73, 0x65, 0x42, 0x17, 0x5a, 0x15, 0x67, 0x6f, 0x2f, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x73, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protobuf_contracts_proto_rawDescOnce sync.Once
	file_protobuf_contracts_proto_rawDescData = file_protobuf_contracts_proto_rawDesc
)

func file_protobuf_contracts_proto_rawDescGZIP() []byte {
	file_protobuf_contracts_proto_rawDescOnce.Do(func() {
		file_protobuf_contracts_proto_rawDescData = protoimpl.X.CompressGZIP(file_protobuf_contracts_proto_rawDescData)
	})
	return file_protobuf_contracts_proto_rawDescData
}

var file_protobuf_contracts_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_protobuf_contracts_proto_goTypes = []interface{}{
	(*ClipInfo)(nil),              // 0: contracts.ClipInfo
	(*EpisodeInfo)(nil),           // 1: contracts.EpisodeInfo
	(*PendingResearchItem)(nil),   // 2: contracts.PendingResearchItem
	(*CompletedResearchItem)(nil), // 3: contracts.CompletedResearchItem
	(*timestamppb.Timestamp)(nil), // 4: google.protobuf.Timestamp
}
var file_protobuf_contracts_proto_depIdxs = []int32{
	4,  // 0: contracts.ClipInfo.initial_date_curated:type_name -> google.protobuf.Timestamp
	4,  // 1: contracts.ClipInfo.last_date_curated:type_name -> google.protobuf.Timestamp
	4,  // 2: contracts.EpisodeInfo.initial_date_curated:type_name -> google.protobuf.Timestamp
	4,  // 3: contracts.EpisodeInfo.last_date_curated:type_name -> google.protobuf.Timestamp
	4,  // 4: contracts.EpisodeInfo.date_aired:type_name -> google.protobuf.Timestamp
	1,  // 5: contracts.PendingResearchItem.episode:type_name -> contracts.EpisodeInfo
	0,  // 6: contracts.PendingResearchItem.clips:type_name -> contracts.ClipInfo
	4,  // 7: contracts.CompletedResearchItem.research_date:type_name -> google.protobuf.Timestamp
	1,  // 8: contracts.CompletedResearchItem.episode_info:type_name -> contracts.EpisodeInfo
	0,  // 9: contracts.CompletedResearchItem.clip_info:type_name -> contracts.ClipInfo
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_protobuf_contracts_proto_init() }
func file_protobuf_contracts_proto_init() {
	if File_protobuf_contracts_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protobuf_contracts_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClipInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protobuf_contracts_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EpisodeInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protobuf_contracts_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PendingResearchItem); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protobuf_contracts_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CompletedResearchItem); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_protobuf_contracts_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_protobuf_contracts_proto_goTypes,
		DependencyIndexes: file_protobuf_contracts_proto_depIdxs,
		MessageInfos:      file_protobuf_contracts_proto_msgTypes,
	}.Build()
	File_protobuf_contracts_proto = out.File
	file_protobuf_contracts_proto_rawDesc = nil
	file_protobuf_contracts_proto_goTypes = nil
	file_protobuf_contracts_proto_depIdxs = nil
}
