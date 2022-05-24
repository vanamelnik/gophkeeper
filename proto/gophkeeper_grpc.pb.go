// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.0
// source: proto/gophkeeper.proto

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
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// GophkeeperClient is the client API for Gophkeeper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GophkeeperClient interface {
	// SignUp registers a new user and creates a new user session.
	SignUp(ctx context.Context, in *SignInData, opts ...grpc.CallOption) (*UserAuth, error)
	// LogIn creates a new session for the user provided.
	LogIn(ctx context.Context, in *SignInData, opts ...grpc.CallOption) (*UserAuth, error)
	// GetNewTokens generates a new AccessToken + RefreshToken pair.
	// If refresh token is expired, the session ends.
	GetNewTokens(ctx context.Context, in *RefreshToken, opts ...grpc.CallOption) (*UserAuth, error)
	// LogOut ends current user session.
	LogOut(ctx context.Context, in *RefreshToken, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// PublishLocalChanges applies the changes to the storage on the server.
	// This method is allowed only if the version of user's data on the client side is equal
	// to the version number on the server. Otherwise the error is returned and the client
	// must first update data from the server.
	PublishLocalChanges(ctx context.Context, in *PublishLocalChangesRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// WhatsNew compares provided Data Version with such one stored on the server. If they are the same,
	// OK status is returned. Otherwise, the error "update the data" is returned.
	WhatsNew(ctx context.Context, in *WhatsNewRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// DownloadUserData analyses existing versions of the local items and downloads latest updates of the user's data from the server.
	DownloadUserData(ctx context.Context, in *DownloadUserDataRequest, opts ...grpc.CallOption) (*UserData, error)
}

type gophkeeperClient struct {
	cc grpc.ClientConnInterface
}

func NewGophkeeperClient(cc grpc.ClientConnInterface) GophkeeperClient {
	return &gophkeeperClient{cc}
}

func (c *gophkeeperClient) SignUp(ctx context.Context, in *SignInData, opts ...grpc.CallOption) (*UserAuth, error) {
	out := new(UserAuth)
	err := c.cc.Invoke(ctx, "/proto.gophkeeper/SignUp", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) LogIn(ctx context.Context, in *SignInData, opts ...grpc.CallOption) (*UserAuth, error) {
	out := new(UserAuth)
	err := c.cc.Invoke(ctx, "/proto.gophkeeper/LogIn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) GetNewTokens(ctx context.Context, in *RefreshToken, opts ...grpc.CallOption) (*UserAuth, error) {
	out := new(UserAuth)
	err := c.cc.Invoke(ctx, "/proto.gophkeeper/GetNewTokens", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) LogOut(ctx context.Context, in *RefreshToken, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proto.gophkeeper/LogOut", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) PublishLocalChanges(ctx context.Context, in *PublishLocalChangesRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proto.gophkeeper/PublishLocalChanges", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) WhatsNew(ctx context.Context, in *WhatsNewRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proto.gophkeeper/WhatsNew", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) DownloadUserData(ctx context.Context, in *DownloadUserDataRequest, opts ...grpc.CallOption) (*UserData, error) {
	out := new(UserData)
	err := c.cc.Invoke(ctx, "/proto.gophkeeper/DownloadUserData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GophkeeperServer is the server API for Gophkeeper service.
// All implementations must embed UnimplementedGophkeeperServer
// for forward compatibility
type GophkeeperServer interface {
	// SignUp registers a new user and creates a new user session.
	SignUp(context.Context, *SignInData) (*UserAuth, error)
	// LogIn creates a new session for the user provided.
	LogIn(context.Context, *SignInData) (*UserAuth, error)
	// GetNewTokens generates a new AccessToken + RefreshToken pair.
	// If refresh token is expired, the session ends.
	GetNewTokens(context.Context, *RefreshToken) (*UserAuth, error)
	// LogOut ends current user session.
	LogOut(context.Context, *RefreshToken) (*emptypb.Empty, error)
	// PublishLocalChanges applies the changes to the storage on the server.
	// This method is allowed only if the version of user's data on the client side is equal
	// to the version number on the server. Otherwise the error is returned and the client
	// must first update data from the server.
	PublishLocalChanges(context.Context, *PublishLocalChangesRequest) (*emptypb.Empty, error)
	// WhatsNew compares provided Data Version with such one stored on the server. If they are the same,
	// OK status is returned. Otherwise, the error "update the data" is returned.
	WhatsNew(context.Context, *WhatsNewRequest) (*emptypb.Empty, error)
	// DownloadUserData analyses existing versions of the local items and downloads latest updates of the user's data from the server.
	DownloadUserData(context.Context, *DownloadUserDataRequest) (*UserData, error)
	mustEmbedUnimplementedGophkeeperServer()
}

// UnimplementedGophkeeperServer must be embedded to have forward compatible implementations.
type UnimplementedGophkeeperServer struct {
}

func (UnimplementedGophkeeperServer) SignUp(context.Context, *SignInData) (*UserAuth, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignUp not implemented")
}
func (UnimplementedGophkeeperServer) LogIn(context.Context, *SignInData) (*UserAuth, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LogIn not implemented")
}
func (UnimplementedGophkeeperServer) GetNewTokens(context.Context, *RefreshToken) (*UserAuth, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNewTokens not implemented")
}
func (UnimplementedGophkeeperServer) LogOut(context.Context, *RefreshToken) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LogOut not implemented")
}
func (UnimplementedGophkeeperServer) PublishLocalChanges(context.Context, *PublishLocalChangesRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublishLocalChanges not implemented")
}
func (UnimplementedGophkeeperServer) WhatsNew(context.Context, *WhatsNewRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WhatsNew not implemented")
}
func (UnimplementedGophkeeperServer) DownloadUserData(context.Context, *DownloadUserDataRequest) (*UserData, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DownloadUserData not implemented")
}
func (UnimplementedGophkeeperServer) mustEmbedUnimplementedGophkeeperServer() {}

// UnsafeGophkeeperServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GophkeeperServer will
// result in compilation errors.
type UnsafeGophkeeperServer interface {
	mustEmbedUnimplementedGophkeeperServer()
}

func RegisterGophkeeperServer(s grpc.ServiceRegistrar, srv GophkeeperServer) {
	s.RegisterService(&Gophkeeper_ServiceDesc, srv)
}

func _Gophkeeper_SignUp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignInData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).SignUp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.gophkeeper/SignUp",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).SignUp(ctx, req.(*SignInData))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_LogIn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignInData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).LogIn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.gophkeeper/LogIn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).LogIn(ctx, req.(*SignInData))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_GetNewTokens_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshToken)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).GetNewTokens(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.gophkeeper/GetNewTokens",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).GetNewTokens(ctx, req.(*RefreshToken))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_LogOut_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshToken)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).LogOut(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.gophkeeper/LogOut",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).LogOut(ctx, req.(*RefreshToken))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_PublishLocalChanges_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishLocalChangesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).PublishLocalChanges(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.gophkeeper/PublishLocalChanges",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).PublishLocalChanges(ctx, req.(*PublishLocalChangesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_WhatsNew_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WhatsNewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).WhatsNew(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.gophkeeper/WhatsNew",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).WhatsNew(ctx, req.(*WhatsNewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_DownloadUserData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownloadUserDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).DownloadUserData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.gophkeeper/DownloadUserData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).DownloadUserData(ctx, req.(*DownloadUserDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Gophkeeper_ServiceDesc is the grpc.ServiceDesc for Gophkeeper service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Gophkeeper_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.gophkeeper",
	HandlerType: (*GophkeeperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SignUp",
			Handler:    _Gophkeeper_SignUp_Handler,
		},
		{
			MethodName: "LogIn",
			Handler:    _Gophkeeper_LogIn_Handler,
		},
		{
			MethodName: "GetNewTokens",
			Handler:    _Gophkeeper_GetNewTokens_Handler,
		},
		{
			MethodName: "LogOut",
			Handler:    _Gophkeeper_LogOut_Handler,
		},
		{
			MethodName: "PublishLocalChanges",
			Handler:    _Gophkeeper_PublishLocalChanges_Handler,
		},
		{
			MethodName: "WhatsNew",
			Handler:    _Gophkeeper_WhatsNew_Handler,
		},
		{
			MethodName: "DownloadUserData",
			Handler:    _Gophkeeper_DownloadUserData_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/gophkeeper.proto",
}
