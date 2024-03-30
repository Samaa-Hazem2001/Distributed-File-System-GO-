// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.0
// source: Ass/Allservices/dfs.proto

package AllServices

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	ClientService_Upload_FullMethodName   = "/AllServices.ClientService/Upload"
	ClientService_Download_FullMethodName = "/AllServices.ClientService/Download"
)

// ClientServiceClient is the client API for ClientService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ClientServiceClient interface {
	Upload(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error)
	Download(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (*DownloadResponse, error)
}

type clientServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewClientServiceClient(cc grpc.ClientConnInterface) ClientServiceClient {
	return &clientServiceClient{cc}
}

func (c *clientServiceClient) Upload(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error) {
	out := new(UpdateResponse)
	err := c.cc.Invoke(ctx, ClientService_Upload_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientServiceClient) Download(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (*DownloadResponse, error) {
	out := new(DownloadResponse)
	err := c.cc.Invoke(ctx, ClientService_Download_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClientServiceServer is the server API for ClientService service.
// All implementations must embed UnimplementedClientServiceServer
// for forward compatibility
type ClientServiceServer interface {
	Upload(context.Context, *UpdateRequest) (*UpdateResponse, error)
	Download(context.Context, *DownloadRequest) (*DownloadResponse, error)
	mustEmbedUnimplementedClientServiceServer()
}

// UnimplementedClientServiceServer must be embedded to have forward compatible implementations.
type UnimplementedClientServiceServer struct {
}

func (UnimplementedClientServiceServer) Upload(context.Context, *UpdateRequest) (*UpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Upload not implemented")
}
func (UnimplementedClientServiceServer) Download(context.Context, *DownloadRequest) (*DownloadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Download not implemented")
}
func (UnimplementedClientServiceServer) mustEmbedUnimplementedClientServiceServer() {}

// UnsafeClientServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClientServiceServer will
// result in compilation errors.
type UnsafeClientServiceServer interface {
	mustEmbedUnimplementedClientServiceServer()
}

func RegisterClientServiceServer(s grpc.ServiceRegistrar, srv ClientServiceServer) {
	s.RegisterService(&ClientService_ServiceDesc, srv)
}

func _ClientService_Upload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientServiceServer).Upload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClientService_Upload_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientServiceServer).Upload(ctx, req.(*UpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClientService_Download_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownloadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientServiceServer).Download(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClientService_Download_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientServiceServer).Download(ctx, req.(*DownloadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ClientService_ServiceDesc is the grpc.ServiceDesc for ClientService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ClientService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "AllServices.ClientService",
	HandlerType: (*ClientServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Upload",
			Handler:    _ClientService_Upload_Handler,
		},
		{
			MethodName: "Download",
			Handler:    _ClientService_Download_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "Ass/Allservices/dfs.proto",
}

const (
	UploadDownloadFileService_UploadFile_FullMethodName   = "/AllServices.UploadDownloadFileService/UploadFile"
	UploadDownloadFileService_DownloadFile_FullMethodName = "/AllServices.UploadDownloadFileService/DownloadFile"
)

// UploadDownloadFileServiceClient is the client API for UploadDownloadFileService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UploadDownloadFileServiceClient interface {
	UploadFile(ctx context.Context, in *UploadFileRequest, opts ...grpc.CallOption) (*UploadFileResponse, error)
	DownloadFile(ctx context.Context, in *DownloadFileRequest, opts ...grpc.CallOption) (*DownloadFileResponse, error)
}

type uploadDownloadFileServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUploadDownloadFileServiceClient(cc grpc.ClientConnInterface) UploadDownloadFileServiceClient {
	return &uploadDownloadFileServiceClient{cc}
}

func (c *uploadDownloadFileServiceClient) UploadFile(ctx context.Context, in *UploadFileRequest, opts ...grpc.CallOption) (*UploadFileResponse, error) {
	out := new(UploadFileResponse)
	err := c.cc.Invoke(ctx, UploadDownloadFileService_UploadFile_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uploadDownloadFileServiceClient) DownloadFile(ctx context.Context, in *DownloadFileRequest, opts ...grpc.CallOption) (*DownloadFileResponse, error) {
	out := new(DownloadFileResponse)
	err := c.cc.Invoke(ctx, UploadDownloadFileService_DownloadFile_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UploadDownloadFileServiceServer is the server API for UploadDownloadFileService service.
// All implementations must embed UnimplementedUploadDownloadFileServiceServer
// for forward compatibility
type UploadDownloadFileServiceServer interface {
	UploadFile(context.Context, *UploadFileRequest) (*UploadFileResponse, error)
	DownloadFile(context.Context, *DownloadFileRequest) (*DownloadFileResponse, error)
	mustEmbedUnimplementedUploadDownloadFileServiceServer()
}

// UnimplementedUploadDownloadFileServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUploadDownloadFileServiceServer struct {
}

func (UnimplementedUploadDownloadFileServiceServer) UploadFile(context.Context, *UploadFileRequest) (*UploadFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadFile not implemented")
}
func (UnimplementedUploadDownloadFileServiceServer) DownloadFile(context.Context, *DownloadFileRequest) (*DownloadFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DownloadFile not implemented")
}
func (UnimplementedUploadDownloadFileServiceServer) mustEmbedUnimplementedUploadDownloadFileServiceServer() {
}

// UnsafeUploadDownloadFileServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UploadDownloadFileServiceServer will
// result in compilation errors.
type UnsafeUploadDownloadFileServiceServer interface {
	mustEmbedUnimplementedUploadDownloadFileServiceServer()
}

func RegisterUploadDownloadFileServiceServer(s grpc.ServiceRegistrar, srv UploadDownloadFileServiceServer) {
	s.RegisterService(&UploadDownloadFileService_ServiceDesc, srv)
}

func _UploadDownloadFileService_UploadFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UploadDownloadFileServiceServer).UploadFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UploadDownloadFileService_UploadFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UploadDownloadFileServiceServer).UploadFile(ctx, req.(*UploadFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UploadDownloadFileService_DownloadFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownloadFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UploadDownloadFileServiceServer).DownloadFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UploadDownloadFileService_DownloadFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UploadDownloadFileServiceServer).DownloadFile(ctx, req.(*DownloadFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UploadDownloadFileService_ServiceDesc is the grpc.ServiceDesc for UploadDownloadFileService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UploadDownloadFileService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "AllServices.UploadDownloadFileService",
	HandlerType: (*UploadDownloadFileServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UploadFile",
			Handler:    _UploadDownloadFileService_UploadFile_Handler,
		},
		{
			MethodName: "DownloadFile",
			Handler:    _UploadDownloadFileService_DownloadFile_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "Ass/Allservices/dfs.proto",
}

const (
	NotifyMachineDataTransferService_NotifyMachineDataTransfer_FullMethodName = "/AllServices.NotifyMachineDataTransferService/NotifyMachineDataTransfer"
)

// NotifyMachineDataTransferServiceClient is the client API for NotifyMachineDataTransferService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NotifyMachineDataTransferServiceClient interface {
	NotifyMachineDataTransfer(ctx context.Context, in *NotifyMachineDataTransferRequest, opts ...grpc.CallOption) (*NotifyMachineDataTransferResponse, error)
}

type notifyMachineDataTransferServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNotifyMachineDataTransferServiceClient(cc grpc.ClientConnInterface) NotifyMachineDataTransferServiceClient {
	return &notifyMachineDataTransferServiceClient{cc}
}

func (c *notifyMachineDataTransferServiceClient) NotifyMachineDataTransfer(ctx context.Context, in *NotifyMachineDataTransferRequest, opts ...grpc.CallOption) (*NotifyMachineDataTransferResponse, error) {
	out := new(NotifyMachineDataTransferResponse)
	err := c.cc.Invoke(ctx, NotifyMachineDataTransferService_NotifyMachineDataTransfer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NotifyMachineDataTransferServiceServer is the server API for NotifyMachineDataTransferService service.
// All implementations must embed UnimplementedNotifyMachineDataTransferServiceServer
// for forward compatibility
type NotifyMachineDataTransferServiceServer interface {
	NotifyMachineDataTransfer(context.Context, *NotifyMachineDataTransferRequest) (*NotifyMachineDataTransferResponse, error)
	mustEmbedUnimplementedNotifyMachineDataTransferServiceServer()
}

// UnimplementedNotifyMachineDataTransferServiceServer must be embedded to have forward compatible implementations.
type UnimplementedNotifyMachineDataTransferServiceServer struct {
}

func (UnimplementedNotifyMachineDataTransferServiceServer) NotifyMachineDataTransfer(context.Context, *NotifyMachineDataTransferRequest) (*NotifyMachineDataTransferResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotifyMachineDataTransfer not implemented")
}
func (UnimplementedNotifyMachineDataTransferServiceServer) mustEmbedUnimplementedNotifyMachineDataTransferServiceServer() {
}

// UnsafeNotifyMachineDataTransferServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NotifyMachineDataTransferServiceServer will
// result in compilation errors.
type UnsafeNotifyMachineDataTransferServiceServer interface {
	mustEmbedUnimplementedNotifyMachineDataTransferServiceServer()
}

func RegisterNotifyMachineDataTransferServiceServer(s grpc.ServiceRegistrar, srv NotifyMachineDataTransferServiceServer) {
	s.RegisterService(&NotifyMachineDataTransferService_ServiceDesc, srv)
}

func _NotifyMachineDataTransferService_NotifyMachineDataTransfer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NotifyMachineDataTransferRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotifyMachineDataTransferServiceServer).NotifyMachineDataTransfer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotifyMachineDataTransferService_NotifyMachineDataTransfer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotifyMachineDataTransferServiceServer).NotifyMachineDataTransfer(ctx, req.(*NotifyMachineDataTransferRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NotifyMachineDataTransferService_ServiceDesc is the grpc.ServiceDesc for NotifyMachineDataTransferService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NotifyMachineDataTransferService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "AllServices.NotifyMachineDataTransferService",
	HandlerType: (*NotifyMachineDataTransferServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NotifyMachineDataTransfer",
			Handler:    _NotifyMachineDataTransferService_NotifyMachineDataTransfer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "Ass/Allservices/dfs.proto",
}

const (
	TransferFileService_TransferFile_FullMethodName = "/AllServices.TransferFileService/TransferFile"
)

// TransferFileServiceClient is the client API for TransferFileService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TransferFileServiceClient interface {
	TransferFile(ctx context.Context, in *TransferFileUploadRequest, opts ...grpc.CallOption) (*TransferFileUploadResponse, error)
}

type transferFileServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTransferFileServiceClient(cc grpc.ClientConnInterface) TransferFileServiceClient {
	return &transferFileServiceClient{cc}
}

func (c *transferFileServiceClient) TransferFile(ctx context.Context, in *TransferFileUploadRequest, opts ...grpc.CallOption) (*TransferFileUploadResponse, error) {
	out := new(TransferFileUploadResponse)
	err := c.cc.Invoke(ctx, TransferFileService_TransferFile_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TransferFileServiceServer is the server API for TransferFileService service.
// All implementations must embed UnimplementedTransferFileServiceServer
// for forward compatibility
type TransferFileServiceServer interface {
	TransferFile(context.Context, *TransferFileUploadRequest) (*TransferFileUploadResponse, error)
	mustEmbedUnimplementedTransferFileServiceServer()
}

// UnimplementedTransferFileServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTransferFileServiceServer struct {
}

func (UnimplementedTransferFileServiceServer) TransferFile(context.Context, *TransferFileUploadRequest) (*TransferFileUploadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransferFile not implemented")
}
func (UnimplementedTransferFileServiceServer) mustEmbedUnimplementedTransferFileServiceServer() {}

// UnsafeTransferFileServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TransferFileServiceServer will
// result in compilation errors.
type UnsafeTransferFileServiceServer interface {
	mustEmbedUnimplementedTransferFileServiceServer()
}

func RegisterTransferFileServiceServer(s grpc.ServiceRegistrar, srv TransferFileServiceServer) {
	s.RegisterService(&TransferFileService_ServiceDesc, srv)
}

func _TransferFileService_TransferFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TransferFileUploadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransferFileServiceServer).TransferFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TransferFileService_TransferFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransferFileServiceServer).TransferFile(ctx, req.(*TransferFileUploadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TransferFileService_ServiceDesc is the grpc.ServiceDesc for TransferFileService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TransferFileService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "AllServices.TransferFileService",
	HandlerType: (*TransferFileServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TransferFile",
			Handler:    _TransferFileService_TransferFile_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "Ass/Allservices/dfs.proto",
}

const (
	RepExstService_RepExst_FullMethodName = "/AllServices.RepExstService/RepExst"
)

// RepExstServiceClient is the client API for RepExstService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RepExstServiceClient interface {
	RepExst(ctx context.Context, in *RepExstRequest, opts ...grpc.CallOption) (*RepExstResponse, error)
}

type repExstServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRepExstServiceClient(cc grpc.ClientConnInterface) RepExstServiceClient {
	return &repExstServiceClient{cc}
}

func (c *repExstServiceClient) RepExst(ctx context.Context, in *RepExstRequest, opts ...grpc.CallOption) (*RepExstResponse, error) {
	out := new(RepExstResponse)
	err := c.cc.Invoke(ctx, RepExstService_RepExst_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RepExstServiceServer is the server API for RepExstService service.
// All implementations must embed UnimplementedRepExstServiceServer
// for forward compatibility
type RepExstServiceServer interface {
	RepExst(context.Context, *RepExstRequest) (*RepExstResponse, error)
	mustEmbedUnimplementedRepExstServiceServer()
}

// UnimplementedRepExstServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRepExstServiceServer struct {
}

func (UnimplementedRepExstServiceServer) RepExst(context.Context, *RepExstRequest) (*RepExstResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RepExst not implemented")
}
func (UnimplementedRepExstServiceServer) mustEmbedUnimplementedRepExstServiceServer() {}

// UnsafeRepExstServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RepExstServiceServer will
// result in compilation errors.
type UnsafeRepExstServiceServer interface {
	mustEmbedUnimplementedRepExstServiceServer()
}

func RegisterRepExstServiceServer(s grpc.ServiceRegistrar, srv RepExstServiceServer) {
	s.RegisterService(&RepExstService_ServiceDesc, srv)
}

func _RepExstService_RepExst_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RepExstRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RepExstServiceServer).RepExst(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RepExstService_RepExst_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RepExstServiceServer).RepExst(ctx, req.(*RepExstRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RepExstService_ServiceDesc is the grpc.ServiceDesc for RepExstService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RepExstService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "AllServices.RepExstService",
	HandlerType: (*RepExstServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RepExst",
			Handler:    _RepExstService_RepExst_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "Ass/Allservices/dfs.proto",
}

const (
	KeepersService_KeeperDone_FullMethodName = "/AllServices.KeepersService/KeeperDone"
	KeepersService_Alive_FullMethodName      = "/AllServices.KeepersService/Alive"
)

// KeepersServiceClient is the client API for KeepersService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KeepersServiceClient interface {
	KeeperDone(ctx context.Context, in *KeeperDoneRequest, opts ...grpc.CallOption) (*KeeperDoneResponse, error)
	Alive(ctx context.Context, in *AliveRequest, opts ...grpc.CallOption) (*AliveResponse, error)
}

type keepersServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewKeepersServiceClient(cc grpc.ClientConnInterface) KeepersServiceClient {
	return &keepersServiceClient{cc}
}

func (c *keepersServiceClient) KeeperDone(ctx context.Context, in *KeeperDoneRequest, opts ...grpc.CallOption) (*KeeperDoneResponse, error) {
	out := new(KeeperDoneResponse)
	err := c.cc.Invoke(ctx, KeepersService_KeeperDone_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keepersServiceClient) Alive(ctx context.Context, in *AliveRequest, opts ...grpc.CallOption) (*AliveResponse, error) {
	out := new(AliveResponse)
	err := c.cc.Invoke(ctx, KeepersService_Alive_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KeepersServiceServer is the server API for KeepersService service.
// All implementations must embed UnimplementedKeepersServiceServer
// for forward compatibility
type KeepersServiceServer interface {
	KeeperDone(context.Context, *KeeperDoneRequest) (*KeeperDoneResponse, error)
	Alive(context.Context, *AliveRequest) (*AliveResponse, error)
	mustEmbedUnimplementedKeepersServiceServer()
}

// UnimplementedKeepersServiceServer must be embedded to have forward compatible implementations.
type UnimplementedKeepersServiceServer struct {
}

func (UnimplementedKeepersServiceServer) KeeperDone(context.Context, *KeeperDoneRequest) (*KeeperDoneResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method KeeperDone not implemented")
}
func (UnimplementedKeepersServiceServer) Alive(context.Context, *AliveRequest) (*AliveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Alive not implemented")
}
func (UnimplementedKeepersServiceServer) mustEmbedUnimplementedKeepersServiceServer() {}

// UnsafeKeepersServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KeepersServiceServer will
// result in compilation errors.
type UnsafeKeepersServiceServer interface {
	mustEmbedUnimplementedKeepersServiceServer()
}

func RegisterKeepersServiceServer(s grpc.ServiceRegistrar, srv KeepersServiceServer) {
	s.RegisterService(&KeepersService_ServiceDesc, srv)
}

func _KeepersService_KeeperDone_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeeperDoneRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeepersServiceServer).KeeperDone(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KeepersService_KeeperDone_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeepersServiceServer).KeeperDone(ctx, req.(*KeeperDoneRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeepersService_Alive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AliveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeepersServiceServer).Alive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KeepersService_Alive_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeepersServiceServer).Alive(ctx, req.(*AliveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// KeepersService_ServiceDesc is the grpc.ServiceDesc for KeepersService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KeepersService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "AllServices.KeepersService",
	HandlerType: (*KeepersServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "KeeperDone",
			Handler:    _KeepersService_KeeperDone_Handler,
		},
		{
			MethodName: "Alive",
			Handler:    _KeepersService_Alive_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "Ass/Allservices/dfs.proto",
}

const (
	DoneUpService_DoneUp_FullMethodName = "/AllServices.DoneUpService/DoneUp"
)

// DoneUpServiceClient is the client API for DoneUpService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DoneUpServiceClient interface {
	DoneUp(ctx context.Context, in *DoneUpRequest, opts ...grpc.CallOption) (*DoneUpResponse, error)
}

type doneUpServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDoneUpServiceClient(cc grpc.ClientConnInterface) DoneUpServiceClient {
	return &doneUpServiceClient{cc}
}

func (c *doneUpServiceClient) DoneUp(ctx context.Context, in *DoneUpRequest, opts ...grpc.CallOption) (*DoneUpResponse, error) {
	out := new(DoneUpResponse)
	err := c.cc.Invoke(ctx, DoneUpService_DoneUp_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DoneUpServiceServer is the server API for DoneUpService service.
// All implementations must embed UnimplementedDoneUpServiceServer
// for forward compatibility
type DoneUpServiceServer interface {
	DoneUp(context.Context, *DoneUpRequest) (*DoneUpResponse, error)
	mustEmbedUnimplementedDoneUpServiceServer()
}

// UnimplementedDoneUpServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDoneUpServiceServer struct {
}

func (UnimplementedDoneUpServiceServer) DoneUp(context.Context, *DoneUpRequest) (*DoneUpResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DoneUp not implemented")
}
func (UnimplementedDoneUpServiceServer) mustEmbedUnimplementedDoneUpServiceServer() {}

// UnsafeDoneUpServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DoneUpServiceServer will
// result in compilation errors.
type UnsafeDoneUpServiceServer interface {
	mustEmbedUnimplementedDoneUpServiceServer()
}

func RegisterDoneUpServiceServer(s grpc.ServiceRegistrar, srv DoneUpServiceServer) {
	s.RegisterService(&DoneUpService_ServiceDesc, srv)
}

func _DoneUpService_DoneUp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DoneUpRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DoneUpServiceServer).DoneUp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DoneUpService_DoneUp_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DoneUpServiceServer).DoneUp(ctx, req.(*DoneUpRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DoneUpService_ServiceDesc is the grpc.ServiceDesc for DoneUpService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DoneUpService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "AllServices.DoneUpService",
	HandlerType: (*DoneUpServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DoneUp",
			Handler:    _DoneUpService_DoneUp_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "Ass/Allservices/dfs.proto",
}
