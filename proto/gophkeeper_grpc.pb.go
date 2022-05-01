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
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// GophkeeperClient is the client API for Gophkeeper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GophkeeperClient interface {
	SignUp(ctx context.Context, in *SignInData, opts ...grpc.CallOption) (*UserAuth, error)
	LogIn(ctx context.Context, in *SignInData, opts ...grpc.CallOption) (*UserAuth, error)
	GetAccessToken(ctx context.Context, in *RefreshToken, opts ...grpc.CallOption) (*AccessToken, error)
	StorePassword(ctx context.Context, in *StorePasswordRequest, opts ...grpc.CallOption) (*ItemID, error)
	StoreBlob(ctx context.Context, in *StoreBlobRequest, opts ...grpc.CallOption) (*ItemID, error)
	StoreText(ctx context.Context, in *StoreTextRequest, opts ...grpc.CallOption) (*ItemID, error)
	StoreCard(ctx context.Context, in *StoreCardRequest, opts ...grpc.CallOption) (*ItemID, error)
	GetData(ctx context.Context, in *AccessToken, opts ...grpc.CallOption) (*Data, error)
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

func (c *gophkeeperClient) GetAccessToken(ctx context.Context, in *RefreshToken, opts ...grpc.CallOption) (*AccessToken, error) {
	out := new(AccessToken)
	err := c.cc.Invoke(ctx, "/proto.gophkeeper/GetAccessToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) StorePassword(ctx context.Context, in *StorePasswordRequest, opts ...grpc.CallOption) (*ItemID, error) {
	out := new(ItemID)
	err := c.cc.Invoke(ctx, "/proto.gophkeeper/StorePassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) StoreBlob(ctx context.Context, in *StoreBlobRequest, opts ...grpc.CallOption) (*ItemID, error) {
	out := new(ItemID)
	err := c.cc.Invoke(ctx, "/proto.gophkeeper/StoreBlob", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) StoreText(ctx context.Context, in *StoreTextRequest, opts ...grpc.CallOption) (*ItemID, error) {
	out := new(ItemID)
	err := c.cc.Invoke(ctx, "/proto.gophkeeper/StoreText", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) StoreCard(ctx context.Context, in *StoreCardRequest, opts ...grpc.CallOption) (*ItemID, error) {
	out := new(ItemID)
	err := c.cc.Invoke(ctx, "/proto.gophkeeper/StoreCard", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) GetData(ctx context.Context, in *AccessToken, opts ...grpc.CallOption) (*Data, error) {
	out := new(Data)
	err := c.cc.Invoke(ctx, "/proto.gophkeeper/GetData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GophkeeperServer is the server API for Gophkeeper service.
// All implementations must embed UnimplementedGophkeeperServer
// for forward compatibility
type GophkeeperServer interface {
	SignUp(context.Context, *SignInData) (*UserAuth, error)
	LogIn(context.Context, *SignInData) (*UserAuth, error)
	GetAccessToken(context.Context, *RefreshToken) (*AccessToken, error)
	StorePassword(context.Context, *StorePasswordRequest) (*ItemID, error)
	StoreBlob(context.Context, *StoreBlobRequest) (*ItemID, error)
	StoreText(context.Context, *StoreTextRequest) (*ItemID, error)
	StoreCard(context.Context, *StoreCardRequest) (*ItemID, error)
	GetData(context.Context, *AccessToken) (*Data, error)
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
func (UnimplementedGophkeeperServer) GetAccessToken(context.Context, *RefreshToken) (*AccessToken, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccessToken not implemented")
}
func (UnimplementedGophkeeperServer) StorePassword(context.Context, *StorePasswordRequest) (*ItemID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StorePassword not implemented")
}
func (UnimplementedGophkeeperServer) StoreBlob(context.Context, *StoreBlobRequest) (*ItemID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoreBlob not implemented")
}
func (UnimplementedGophkeeperServer) StoreText(context.Context, *StoreTextRequest) (*ItemID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoreText not implemented")
}
func (UnimplementedGophkeeperServer) StoreCard(context.Context, *StoreCardRequest) (*ItemID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoreCard not implemented")
}
func (UnimplementedGophkeeperServer) GetData(context.Context, *AccessToken) (*Data, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetData not implemented")
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

func _Gophkeeper_GetAccessToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshToken)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).GetAccessToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.gophkeeper/GetAccessToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).GetAccessToken(ctx, req.(*RefreshToken))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_StorePassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StorePasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).StorePassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.gophkeeper/StorePassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).StorePassword(ctx, req.(*StorePasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_StoreBlob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoreBlobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).StoreBlob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.gophkeeper/StoreBlob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).StoreBlob(ctx, req.(*StoreBlobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_StoreText_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoreTextRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).StoreText(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.gophkeeper/StoreText",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).StoreText(ctx, req.(*StoreTextRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_StoreCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoreCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).StoreCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.gophkeeper/StoreCard",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).StoreCard(ctx, req.(*StoreCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_GetData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AccessToken)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).GetData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.gophkeeper/GetData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).GetData(ctx, req.(*AccessToken))
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
			MethodName: "GetAccessToken",
			Handler:    _Gophkeeper_GetAccessToken_Handler,
		},
		{
			MethodName: "StorePassword",
			Handler:    _Gophkeeper_StorePassword_Handler,
		},
		{
			MethodName: "StoreBlob",
			Handler:    _Gophkeeper_StoreBlob_Handler,
		},
		{
			MethodName: "StoreText",
			Handler:    _Gophkeeper_StoreText_Handler,
		},
		{
			MethodName: "StoreCard",
			Handler:    _Gophkeeper_StoreCard_Handler,
		},
		{
			MethodName: "GetData",
			Handler:    _Gophkeeper_GetData_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/gophkeeper.proto",
}
