// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.0
// 	protoc        v5.29.1
// source: dataq.proto

package rpc

import (
	pb "go.quinn.io/dataq/pb"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// ExtractRequest contains information about what data to extract
type ExtractRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Hash          string                 `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"` // Content address of the object responsible for creating the Extract
	Kind          string                 `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"` // Operation to be performed that will produce data
	Metadata      map[string]string      `protobuf:"bytes,3,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ExtractRequest) Reset() {
	*x = ExtractRequest{}
	mi := &file_dataq_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ExtractRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExtractRequest) ProtoMessage() {}

func (x *ExtractRequest) ProtoReflect() protoreflect.Message {
	mi := &file_dataq_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExtractRequest.ProtoReflect.Descriptor instead.
func (*ExtractRequest) Descriptor() ([]byte, []int) {
	return file_dataq_proto_rawDescGZIP(), []int{0}
}

func (x *ExtractRequest) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *ExtractRequest) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *ExtractRequest) GetMetadata() map[string]string {
	if x != nil {
		return x.Metadata
	}
	return nil
}

// ExtractResponse contains the result of an extraction
type ExtractResponse struct {
	state         protoimpl.MessageState       `protogen:"open.v1"`
	Hash          string                       `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`                                  // Content address of the produced data
	Kind          string                       `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`                                  // Kind of content being stored
	RequestHash   string                       `protobuf:"bytes,3,opt,name=request_hash,json=requestHash,proto3" json:"request_hash,omitempty"` // Address of the request
	Transforms    []*ExtractResponse_Transform `protobuf:"bytes,4,rep,name=transforms,proto3" json:"transforms,omitempty"`                      // List of transform requests to be created
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ExtractResponse) Reset() {
	*x = ExtractResponse{}
	mi := &file_dataq_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ExtractResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExtractResponse) ProtoMessage() {}

func (x *ExtractResponse) ProtoReflect() protoreflect.Message {
	mi := &file_dataq_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExtractResponse.ProtoReflect.Descriptor instead.
func (*ExtractResponse) Descriptor() ([]byte, []int) {
	return file_dataq_proto_rawDescGZIP(), []int{1}
}

func (x *ExtractResponse) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *ExtractResponse) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *ExtractResponse) GetRequestHash() string {
	if x != nil {
		return x.RequestHash
	}
	return ""
}

func (x *ExtractResponse) GetTransforms() []*ExtractResponse_Transform {
	if x != nil {
		return x.Transforms
	}
	return nil
}

// TransformRequest represents an action to be performed on data
type TransformRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Hash          string                 `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"` // Address of the data to transform
	Kind          string                 `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"` // Kind of transform to be applied
	Metadata      map[string]string      `protobuf:"bytes,3,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TransformRequest) Reset() {
	*x = TransformRequest{}
	mi := &file_dataq_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TransformRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransformRequest) ProtoMessage() {}

func (x *TransformRequest) ProtoReflect() protoreflect.Message {
	mi := &file_dataq_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransformRequest.ProtoReflect.Descriptor instead.
func (*TransformRequest) Descriptor() ([]byte, []int) {
	return file_dataq_proto_rawDescGZIP(), []int{2}
}

func (x *TransformRequest) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *TransformRequest) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *TransformRequest) GetMetadata() map[string]string {
	if x != nil {
		return x.Metadata
	}
	return nil
}

// TransformResponse contains the result of a transform operation
type TransformResponse struct {
	state         protoimpl.MessageState         `protogen:"open.v1"`
	Hash          string                         `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`                                  // Value copied from RequestHash
	Kind          string                         `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`                                  // Kind from request
	RequestHash   string                         `protobuf:"bytes,3,opt,name=request_hash,json=requestHash,proto3" json:"request_hash,omitempty"` // Address of request
	Extracts      []*TransformResponse_Extract   `protobuf:"bytes,4,rep,name=extracts,proto3" json:"extracts,omitempty"`                          // List of extracts to be performed
	Permanodes    []*TransformResponse_Permanode `protobuf:"bytes,5,rep,name=permanodes,proto3" json:"permanodes,omitempty"`                      // List of permanodes to be managed
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TransformResponse) Reset() {
	*x = TransformResponse{}
	mi := &file_dataq_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TransformResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransformResponse) ProtoMessage() {}

func (x *TransformResponse) ProtoReflect() protoreflect.Message {
	mi := &file_dataq_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransformResponse.ProtoReflect.Descriptor instead.
func (*TransformResponse) Descriptor() ([]byte, []int) {
	return file_dataq_proto_rawDescGZIP(), []int{3}
}

func (x *TransformResponse) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *TransformResponse) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *TransformResponse) GetRequestHash() string {
	if x != nil {
		return x.RequestHash
	}
	return ""
}

func (x *TransformResponse) GetExtracts() []*TransformResponse_Extract {
	if x != nil {
		return x.Extracts
	}
	return nil
}

func (x *TransformResponse) GetPermanodes() []*TransformResponse_Permanode {
	if x != nil {
		return x.Permanodes
	}
	return nil
}

// DataSource identifies an object within a specific plugin
type DataSource struct {
	state                 protoimpl.MessageState `protogen:"open.v1"`
	PermanodeHash         string                 `protobuf:"bytes,1,opt,name=permanode_hash,json=permanodeHash,proto3" json:"permanode_hash,omitempty"`
	TransformResponseHash string                 `protobuf:"bytes,2,opt,name=transform_response_hash,json=transformResponseHash,proto3" json:"transform_response_hash,omitempty"`
	Plugin                string                 `protobuf:"bytes,3,opt,name=plugin,proto3" json:"plugin,omitempty"` // Unique identifier for the plugin
	Key                   string                 `protobuf:"bytes,4,opt,name=key,proto3" json:"key,omitempty"`       // Value to uniquely identify plugin's internal representation
	unknownFields         protoimpl.UnknownFields
	sizeCache             protoimpl.SizeCache
}

func (x *DataSource) Reset() {
	*x = DataSource{}
	mi := &file_dataq_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DataSource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataSource) ProtoMessage() {}

func (x *DataSource) ProtoReflect() protoreflect.Message {
	mi := &file_dataq_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataSource.ProtoReflect.Descriptor instead.
func (*DataSource) Descriptor() ([]byte, []int) {
	return file_dataq_proto_rawDescGZIP(), []int{4}
}

func (x *DataSource) GetPermanodeHash() string {
	if x != nil {
		return x.PermanodeHash
	}
	return ""
}

func (x *DataSource) GetTransformResponseHash() string {
	if x != nil {
		return x.TransformResponseHash
	}
	return ""
}

func (x *DataSource) GetPlugin() string {
	if x != nil {
		return x.Plugin
	}
	return ""
}

func (x *DataSource) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

// PermanodeVersion represents a versioned permanode
type PermanodeVersion struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	PermanodeHash string                 `protobuf:"bytes,1,opt,name=permanode_hash,json=permanodeHash,proto3" json:"permanode_hash,omitempty"` // Hash of the permanode
	Timestamp     int64                  `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`                             // Version timestamp
	Deleted       bool                   `protobuf:"varint,4,opt,name=deleted,proto3" json:"deleted,omitempty"`                                 // Deletion status
	Source        *DataSource            `protobuf:"bytes,5,opt,name=source,proto3" json:"source,omitempty"`                                    // Source of the permanode version
	// Types that are valid to be assigned to Payload:
	//
	//	*PermanodeVersion_Email
	//	*PermanodeVersion_FinancialTransaction
	Payload       isPermanodeVersion_Payload `protobuf_oneof:"payload"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PermanodeVersion) Reset() {
	*x = PermanodeVersion{}
	mi := &file_dataq_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PermanodeVersion) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PermanodeVersion) ProtoMessage() {}

func (x *PermanodeVersion) ProtoReflect() protoreflect.Message {
	mi := &file_dataq_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PermanodeVersion.ProtoReflect.Descriptor instead.
func (*PermanodeVersion) Descriptor() ([]byte, []int) {
	return file_dataq_proto_rawDescGZIP(), []int{5}
}

func (x *PermanodeVersion) GetPermanodeHash() string {
	if x != nil {
		return x.PermanodeHash
	}
	return ""
}

func (x *PermanodeVersion) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *PermanodeVersion) GetDeleted() bool {
	if x != nil {
		return x.Deleted
	}
	return false
}

func (x *PermanodeVersion) GetSource() *DataSource {
	if x != nil {
		return x.Source
	}
	return nil
}

func (x *PermanodeVersion) GetPayload() isPermanodeVersion_Payload {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *PermanodeVersion) GetEmail() *pb.Email {
	if x != nil {
		if x, ok := x.Payload.(*PermanodeVersion_Email); ok {
			return x.Email
		}
	}
	return nil
}

func (x *PermanodeVersion) GetFinancialTransaction() *pb.FinancialTransaction {
	if x != nil {
		if x, ok := x.Payload.(*PermanodeVersion_FinancialTransaction); ok {
			return x.FinancialTransaction
		}
	}
	return nil
}

type isPermanodeVersion_Payload interface {
	isPermanodeVersion_Payload()
}

type PermanodeVersion_Email struct {
	Email *pb.Email `protobuf:"bytes,6,opt,name=email,proto3,oneof"`
}

type PermanodeVersion_FinancialTransaction struct {
	FinancialTransaction *pb.FinancialTransaction `protobuf:"bytes,7,opt,name=financial_transaction,json=financialTransaction,proto3,oneof"`
}

func (*PermanodeVersion_Email) isPermanodeVersion_Payload() {}

func (*PermanodeVersion_FinancialTransaction) isPermanodeVersion_Payload() {}

// Transform defines a transform operation to be performed
type ExtractResponse_Transform struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Kind          string                 `protobuf:"bytes,1,opt,name=kind,proto3" json:"kind,omitempty"`                                                                                   // Name of the transform
	Metadata      map[string]string      `protobuf:"bytes,2,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"` // Additional information for the transform
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ExtractResponse_Transform) Reset() {
	*x = ExtractResponse_Transform{}
	mi := &file_dataq_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ExtractResponse_Transform) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExtractResponse_Transform) ProtoMessage() {}

func (x *ExtractResponse_Transform) ProtoReflect() protoreflect.Message {
	mi := &file_dataq_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExtractResponse_Transform.ProtoReflect.Descriptor instead.
func (*ExtractResponse_Transform) Descriptor() ([]byte, []int) {
	return file_dataq_proto_rawDescGZIP(), []int{1, 0}
}

func (x *ExtractResponse_Transform) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *ExtractResponse_Transform) GetMetadata() map[string]string {
	if x != nil {
		return x.Metadata
	}
	return nil
}

// Extract defines an extract operation to be performed
type TransformResponse_Extract struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Kind          string                 `protobuf:"bytes,1,opt,name=kind,proto3" json:"kind,omitempty"`                                                                                   // Operation to perform
	Metadata      map[string]string      `protobuf:"bytes,2,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"` // Data necessary to perform the operation
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TransformResponse_Extract) Reset() {
	*x = TransformResponse_Extract{}
	mi := &file_dataq_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TransformResponse_Extract) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransformResponse_Extract) ProtoMessage() {}

func (x *TransformResponse_Extract) ProtoReflect() protoreflect.Message {
	mi := &file_dataq_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransformResponse_Extract.ProtoReflect.Descriptor instead.
func (*TransformResponse_Extract) Descriptor() ([]byte, []int) {
	return file_dataq_proto_rawDescGZIP(), []int{3, 0}
}

func (x *TransformResponse_Extract) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *TransformResponse_Extract) GetMetadata() map[string]string {
	if x != nil {
		return x.Metadata
	}
	return nil
}

type TransformResponse_Permanode struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Kind  string                 `protobuf:"bytes,1,opt,name=kind,proto3" json:"kind,omitempty"` // Kind of permanode to be managed
	Key   string                 `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`   // Value to uniquely identify the permanode within the plugin
	// Types that are valid to be assigned to Payload:
	//
	//	*TransformResponse_Permanode_Email
	//	*TransformResponse_Permanode_FinancialTransaction
	Payload       isTransformResponse_Permanode_Payload `protobuf_oneof:"payload"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TransformResponse_Permanode) Reset() {
	*x = TransformResponse_Permanode{}
	mi := &file_dataq_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TransformResponse_Permanode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransformResponse_Permanode) ProtoMessage() {}

func (x *TransformResponse_Permanode) ProtoReflect() protoreflect.Message {
	mi := &file_dataq_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransformResponse_Permanode.ProtoReflect.Descriptor instead.
func (*TransformResponse_Permanode) Descriptor() ([]byte, []int) {
	return file_dataq_proto_rawDescGZIP(), []int{3, 1}
}

func (x *TransformResponse_Permanode) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *TransformResponse_Permanode) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *TransformResponse_Permanode) GetPayload() isTransformResponse_Permanode_Payload {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *TransformResponse_Permanode) GetEmail() *pb.Email {
	if x != nil {
		if x, ok := x.Payload.(*TransformResponse_Permanode_Email); ok {
			return x.Email
		}
	}
	return nil
}

func (x *TransformResponse_Permanode) GetFinancialTransaction() *pb.FinancialTransaction {
	if x != nil {
		if x, ok := x.Payload.(*TransformResponse_Permanode_FinancialTransaction); ok {
			return x.FinancialTransaction
		}
	}
	return nil
}

type isTransformResponse_Permanode_Payload interface {
	isTransformResponse_Permanode_Payload()
}

type TransformResponse_Permanode_Email struct {
	Email *pb.Email `protobuf:"bytes,6,opt,name=email,proto3,oneof"`
}

type TransformResponse_Permanode_FinancialTransaction struct {
	FinancialTransaction *pb.FinancialTransaction `protobuf:"bytes,7,opt,name=financial_transaction,json=financialTransaction,proto3,oneof"`
}

func (*TransformResponse_Permanode_Email) isTransformResponse_Permanode_Payload() {}

func (*TransformResponse_Permanode_FinancialTransaction) isTransformResponse_Permanode_Payload() {}

var File_dataq_proto protoreflect.FileDescriptor

var file_dataq_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x64, 0x61, 0x74, 0x61, 0x71, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x64,
	0x61, 0x74, 0x61, 0x71, 0x1a, 0x0c, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xb6, 0x01, 0x0a, 0x0e, 0x45, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x69, 0x6e,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x3f, 0x0a,
	0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x23, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x71, 0x2e, 0x45, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x3b,
	0x0a, 0x0d, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xc9, 0x02, 0x0a, 0x0f,
	0x45, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68,
	0x61, 0x73, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x61, 0x73, 0x68, 0x12, 0x40, 0x0a, 0x0a, 0x74, 0x72,
	0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20,
	0x2e, 0x64, 0x61, 0x74, 0x61, 0x71, 0x2e, 0x45, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d,
	0x52, 0x0a, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x73, 0x1a, 0xa8, 0x01, 0x0a,
	0x09, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x69,
	0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x4a,
	0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x2e, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x71, 0x2e, 0x45, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f,
	0x72, 0x6d, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x3b, 0x0a, 0x0d, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xba, 0x01, 0x0a, 0x10, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x66, 0x6f, 0x72, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68,
	0x12, 0x12, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6b, 0x69, 0x6e, 0x64, 0x12, 0x41, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x71, 0x2e, 0x54,
	0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x6d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x3b, 0x0a, 0x0d, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x22, 0xc2, 0x04, 0x0a, 0x11, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f,
	0x72, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61,
	0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x12,
	0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x69,
	0x6e, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x68, 0x61,
	0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x48, 0x61, 0x73, 0x68, 0x12, 0x3c, 0x0a, 0x08, 0x65, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74,
	0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x71, 0x2e,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x2e, 0x45, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74, 0x52, 0x08, 0x65, 0x78, 0x74, 0x72, 0x61,
	0x63, 0x74, 0x73, 0x12, 0x42, 0x0a, 0x0a, 0x70, 0x65, 0x72, 0x6d, 0x61, 0x6e, 0x6f, 0x64, 0x65,
	0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x71, 0x2e,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x2e, 0x50, 0x65, 0x72, 0x6d, 0x61, 0x6e, 0x6f, 0x64, 0x65, 0x52, 0x0a, 0x70, 0x65, 0x72,
	0x6d, 0x61, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x1a, 0xa6, 0x01, 0x0a, 0x07, 0x45, 0x78, 0x74, 0x72,
	0x61, 0x63, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x4a, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x64, 0x61, 0x74, 0x61,
	0x71, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x2e, 0x45, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74, 0x2e, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x1a, 0x3b, 0x0a, 0x0d, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x1a, 0xb6, 0x01, 0x0a, 0x09, 0x50, 0x65, 0x72, 0x6d, 0x61, 0x6e, 0x6f, 0x64, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x69,
	0x6e, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x24, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x71, 0x2e, 0x45, 0x6d, 0x61, 0x69,
	0x6c, 0x48, 0x00, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x52, 0x0a, 0x15, 0x66, 0x69,
	0x6e, 0x61, 0x6e, 0x63, 0x69, 0x61, 0x6c, 0x5f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x64, 0x61, 0x74, 0x61,
	0x71, 0x2e, 0x46, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x69, 0x61, 0x6c, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52, 0x14, 0x66, 0x69, 0x6e, 0x61, 0x6e, 0x63,
	0x69, 0x61, 0x6c, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x09,
	0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x95, 0x01, 0x0a, 0x0a, 0x44, 0x61,
	0x74, 0x61, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x70, 0x65, 0x72, 0x6d,
	0x61, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0d, 0x70, 0x65, 0x72, 0x6d, 0x61, 0x6e, 0x6f, 0x64, 0x65, 0x48, 0x61, 0x73, 0x68, 0x12,
	0x36, 0x0a, 0x17, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x5f, 0x72, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x15, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x48, 0x61, 0x73, 0x68, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x22, 0xa1, 0x02, 0x0a, 0x10, 0x50, 0x65, 0x72, 0x6d, 0x61, 0x6e, 0x6f, 0x64, 0x65, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x25, 0x0a, 0x0e, 0x70, 0x65, 0x72, 0x6d, 0x61, 0x6e,
	0x6f, 0x64, 0x65, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x70, 0x65, 0x72, 0x6d, 0x61, 0x6e, 0x6f, 0x64, 0x65, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1c, 0x0a,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x64,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x64, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x64, 0x12, 0x29, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x71, 0x2e, 0x44, 0x61,
	0x74, 0x61, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x12, 0x24, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0c, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x71, 0x2e, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x48, 0x00, 0x52,
	0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x52, 0x0a, 0x15, 0x66, 0x69, 0x6e, 0x61, 0x6e, 0x63,
	0x69, 0x61, 0x6c, 0x5f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x71, 0x2e, 0x46, 0x69,
	0x6e, 0x61, 0x6e, 0x63, 0x69, 0x61, 0x6c, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x48, 0x00, 0x52, 0x14, 0x66, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x69, 0x61, 0x6c, 0x54,
	0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x09, 0x0a, 0x07, 0x70, 0x61,
	0x79, 0x6c, 0x6f, 0x61, 0x64, 0x32, 0x8b, 0x01, 0x0a, 0x0b, 0x44, 0x61, 0x74, 0x61, 0x51, 0x50,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x12, 0x3a, 0x0a, 0x07, 0x45, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74,
	0x12, 0x15, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x71, 0x2e, 0x45, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x71, 0x2e,
	0x45, 0x78, 0x74, 0x72, 0x61, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x40, 0x0a, 0x09, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x12, 0x17,
	0x2e, 0x64, 0x61, 0x74, 0x61, 0x71, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x71, 0x2e,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x42, 0x1c, 0x5a, 0x1a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x71, 0x75, 0x69, 0x6e, 0x6e, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x71, 0x2f, 0x72, 0x70,
	0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_dataq_proto_rawDescOnce sync.Once
	file_dataq_proto_rawDescData = file_dataq_proto_rawDesc
)

func file_dataq_proto_rawDescGZIP() []byte {
	file_dataq_proto_rawDescOnce.Do(func() {
		file_dataq_proto_rawDescData = protoimpl.X.CompressGZIP(file_dataq_proto_rawDescData)
	})
	return file_dataq_proto_rawDescData
}

var file_dataq_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_dataq_proto_goTypes = []any{
	(*ExtractRequest)(nil),              // 0: dataq.ExtractRequest
	(*ExtractResponse)(nil),             // 1: dataq.ExtractResponse
	(*TransformRequest)(nil),            // 2: dataq.TransformRequest
	(*TransformResponse)(nil),           // 3: dataq.TransformResponse
	(*DataSource)(nil),                  // 4: dataq.DataSource
	(*PermanodeVersion)(nil),            // 5: dataq.PermanodeVersion
	nil,                                 // 6: dataq.ExtractRequest.MetadataEntry
	(*ExtractResponse_Transform)(nil),   // 7: dataq.ExtractResponse.Transform
	nil,                                 // 8: dataq.ExtractResponse.Transform.MetadataEntry
	nil,                                 // 9: dataq.TransformRequest.MetadataEntry
	(*TransformResponse_Extract)(nil),   // 10: dataq.TransformResponse.Extract
	(*TransformResponse_Permanode)(nil), // 11: dataq.TransformResponse.Permanode
	nil,                                 // 12: dataq.TransformResponse.Extract.MetadataEntry
	(*pb.Email)(nil),                    // 13: dataq.Email
	(*pb.FinancialTransaction)(nil),     // 14: dataq.FinancialTransaction
}
var file_dataq_proto_depIdxs = []int32{
	6,  // 0: dataq.ExtractRequest.metadata:type_name -> dataq.ExtractRequest.MetadataEntry
	7,  // 1: dataq.ExtractResponse.transforms:type_name -> dataq.ExtractResponse.Transform
	9,  // 2: dataq.TransformRequest.metadata:type_name -> dataq.TransformRequest.MetadataEntry
	10, // 3: dataq.TransformResponse.extracts:type_name -> dataq.TransformResponse.Extract
	11, // 4: dataq.TransformResponse.permanodes:type_name -> dataq.TransformResponse.Permanode
	4,  // 5: dataq.PermanodeVersion.source:type_name -> dataq.DataSource
	13, // 6: dataq.PermanodeVersion.email:type_name -> dataq.Email
	14, // 7: dataq.PermanodeVersion.financial_transaction:type_name -> dataq.FinancialTransaction
	8,  // 8: dataq.ExtractResponse.Transform.metadata:type_name -> dataq.ExtractResponse.Transform.MetadataEntry
	12, // 9: dataq.TransformResponse.Extract.metadata:type_name -> dataq.TransformResponse.Extract.MetadataEntry
	13, // 10: dataq.TransformResponse.Permanode.email:type_name -> dataq.Email
	14, // 11: dataq.TransformResponse.Permanode.financial_transaction:type_name -> dataq.FinancialTransaction
	0,  // 12: dataq.DataQPlugin.Extract:input_type -> dataq.ExtractRequest
	2,  // 13: dataq.DataQPlugin.Transform:input_type -> dataq.TransformRequest
	1,  // 14: dataq.DataQPlugin.Extract:output_type -> dataq.ExtractResponse
	3,  // 15: dataq.DataQPlugin.Transform:output_type -> dataq.TransformResponse
	14, // [14:16] is the sub-list for method output_type
	12, // [12:14] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_dataq_proto_init() }
func file_dataq_proto_init() {
	if File_dataq_proto != nil {
		return
	}
	file_dataq_proto_msgTypes[5].OneofWrappers = []any{
		(*PermanodeVersion_Email)(nil),
		(*PermanodeVersion_FinancialTransaction)(nil),
	}
	file_dataq_proto_msgTypes[11].OneofWrappers = []any{
		(*TransformResponse_Permanode_Email)(nil),
		(*TransformResponse_Permanode_FinancialTransaction)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_dataq_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_dataq_proto_goTypes,
		DependencyIndexes: file_dataq_proto_depIdxs,
		MessageInfos:      file_dataq_proto_msgTypes,
	}.Build()
	File_dataq_proto = out.File
	file_dataq_proto_rawDesc = nil
	file_dataq_proto_goTypes = nil
	file_dataq_proto_depIdxs = nil
}