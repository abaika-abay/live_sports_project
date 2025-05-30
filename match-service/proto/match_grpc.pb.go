// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v6.30.2
// source: match-service/proto/match.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	MatchService_GetMatchUpdates_FullMethodName   = "/match.MatchService/GetMatchUpdates"
	MatchService_CreateMatch_FullMethodName       = "/match.MatchService/CreateMatch"
	MatchService_UpdateMatchEvent_FullMethodName  = "/match.MatchService/UpdateMatchEvent"
	MatchService_GetAdminMatchList_FullMethodName = "/match.MatchService/GetAdminMatchList"
)

// MatchServiceClient is the client API for MatchService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MatchServiceClient interface {
	GetMatchUpdates(ctx context.Context, in *MatchRequest, opts ...grpc.CallOption) (*MatchResponse, error)
	CreateMatch(ctx context.Context, in *CreateMatchRequest, opts ...grpc.CallOption) (*MatchResponse, error)
	// New RPC for admin to update match events
	UpdateMatchEvent(ctx context.Context, in *UpdateMatchEventRequest, opts ...grpc.CallOption) (*MatchResponse, error)
	// Optional: RPC for getting a list of matches for admin panel
	GetAdminMatchList(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*MatchListResponse, error)
}

type matchServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMatchServiceClient(cc grpc.ClientConnInterface) MatchServiceClient {
	return &matchServiceClient{cc}
}

func (c *matchServiceClient) GetMatchUpdates(ctx context.Context, in *MatchRequest, opts ...grpc.CallOption) (*MatchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MatchResponse)
	err := c.cc.Invoke(ctx, MatchService_GetMatchUpdates_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchServiceClient) CreateMatch(ctx context.Context, in *CreateMatchRequest, opts ...grpc.CallOption) (*MatchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MatchResponse)
	err := c.cc.Invoke(ctx, MatchService_CreateMatch_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchServiceClient) UpdateMatchEvent(ctx context.Context, in *UpdateMatchEventRequest, opts ...grpc.CallOption) (*MatchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MatchResponse)
	err := c.cc.Invoke(ctx, MatchService_UpdateMatchEvent_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *matchServiceClient) GetAdminMatchList(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*MatchListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MatchListResponse)
	err := c.cc.Invoke(ctx, MatchService_GetAdminMatchList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MatchServiceServer is the server API for MatchService service.
// All implementations must embed UnimplementedMatchServiceServer
// for forward compatibility.
type MatchServiceServer interface {
	GetMatchUpdates(context.Context, *MatchRequest) (*MatchResponse, error)
	CreateMatch(context.Context, *CreateMatchRequest) (*MatchResponse, error)
	// New RPC for admin to update match events
	UpdateMatchEvent(context.Context, *UpdateMatchEventRequest) (*MatchResponse, error)
	// Optional: RPC for getting a list of matches for admin panel
	GetAdminMatchList(context.Context, *emptypb.Empty) (*MatchListResponse, error)
	mustEmbedUnimplementedMatchServiceServer()
}

// UnimplementedMatchServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMatchServiceServer struct{}

func (UnimplementedMatchServiceServer) GetMatchUpdates(context.Context, *MatchRequest) (*MatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMatchUpdates not implemented")
}
func (UnimplementedMatchServiceServer) CreateMatch(context.Context, *CreateMatchRequest) (*MatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMatch not implemented")
}
func (UnimplementedMatchServiceServer) UpdateMatchEvent(context.Context, *UpdateMatchEventRequest) (*MatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMatchEvent not implemented")
}
func (UnimplementedMatchServiceServer) GetAdminMatchList(context.Context, *emptypb.Empty) (*MatchListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAdminMatchList not implemented")
}
func (UnimplementedMatchServiceServer) mustEmbedUnimplementedMatchServiceServer() {}
func (UnimplementedMatchServiceServer) testEmbeddedByValue()                      {}

// UnsafeMatchServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MatchServiceServer will
// result in compilation errors.
type UnsafeMatchServiceServer interface {
	mustEmbedUnimplementedMatchServiceServer()
}

func RegisterMatchServiceServer(s grpc.ServiceRegistrar, srv MatchServiceServer) {
	// If the following call pancis, it indicates UnimplementedMatchServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&MatchService_ServiceDesc, srv)
}

func _MatchService_GetMatchUpdates_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchServiceServer).GetMatchUpdates(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MatchService_GetMatchUpdates_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchServiceServer).GetMatchUpdates(ctx, req.(*MatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchService_CreateMatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateMatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchServiceServer).CreateMatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MatchService_CreateMatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchServiceServer).CreateMatch(ctx, req.(*CreateMatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchService_UpdateMatchEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMatchEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchServiceServer).UpdateMatchEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MatchService_UpdateMatchEvent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchServiceServer).UpdateMatchEvent(ctx, req.(*UpdateMatchEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MatchService_GetAdminMatchList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchServiceServer).GetAdminMatchList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MatchService_GetAdminMatchList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchServiceServer).GetAdminMatchList(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// MatchService_ServiceDesc is the grpc.ServiceDesc for MatchService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MatchService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "match.MatchService",
	HandlerType: (*MatchServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetMatchUpdates",
			Handler:    _MatchService_GetMatchUpdates_Handler,
		},
		{
			MethodName: "CreateMatch",
			Handler:    _MatchService_CreateMatch_Handler,
		},
		{
			MethodName: "UpdateMatchEvent",
			Handler:    _MatchService_UpdateMatchEvent_Handler,
		},
		{
			MethodName: "GetAdminMatchList",
			Handler:    _MatchService_GetAdminMatchList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "match-service/proto/match.proto",
}
