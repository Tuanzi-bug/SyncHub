//protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative  login_service.proto

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.20.3
// source: login_service.proto

package login

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	LoginService_GetCaptcha_FullMethodName      = "/login.LoginService/GetCaptcha"
	LoginService_Register_FullMethodName        = "/login.LoginService/Register"
	LoginService_Login_FullMethodName           = "/login.LoginService/Login"
	LoginService_TokenVerify_FullMethodName     = "/login.LoginService/TokenVerify"
	LoginService_MyOrgList_FullMethodName       = "/login.LoginService/MyOrgList"
	LoginService_FindMemInfoById_FullMethodName = "/login.LoginService/FindMemInfoById"
)

// LoginServiceClient is the client API for LoginService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LoginServiceClient interface {
	GetCaptcha(ctx context.Context, in *CaptchaMessage, opts ...grpc.CallOption) (*CaptchaResponse, error)
	Register(ctx context.Context, in *RegisterMessage, opts ...grpc.CallOption) (*RegisterResponse, error)
	Login(ctx context.Context, in *LoginMessage, opts ...grpc.CallOption) (*LoginResponse, error)
	TokenVerify(ctx context.Context, in *LoginMessage, opts ...grpc.CallOption) (*LoginResponse, error)
	MyOrgList(ctx context.Context, in *UserMessage, opts ...grpc.CallOption) (*OrgListResponse, error)
	FindMemInfoById(ctx context.Context, in *UserMessage, opts ...grpc.CallOption) (*MemberMessage, error)
}

type loginServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLoginServiceClient(cc grpc.ClientConnInterface) LoginServiceClient {
	return &loginServiceClient{cc}
}

func (c *loginServiceClient) GetCaptcha(ctx context.Context, in *CaptchaMessage, opts ...grpc.CallOption) (*CaptchaResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CaptchaResponse)
	err := c.cc.Invoke(ctx, LoginService_GetCaptcha_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loginServiceClient) Register(ctx context.Context, in *RegisterMessage, opts ...grpc.CallOption) (*RegisterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, LoginService_Register_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loginServiceClient) Login(ctx context.Context, in *LoginMessage, opts ...grpc.CallOption) (*LoginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, LoginService_Login_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loginServiceClient) TokenVerify(ctx context.Context, in *LoginMessage, opts ...grpc.CallOption) (*LoginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, LoginService_TokenVerify_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loginServiceClient) MyOrgList(ctx context.Context, in *UserMessage, opts ...grpc.CallOption) (*OrgListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OrgListResponse)
	err := c.cc.Invoke(ctx, LoginService_MyOrgList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loginServiceClient) FindMemInfoById(ctx context.Context, in *UserMessage, opts ...grpc.CallOption) (*MemberMessage, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MemberMessage)
	err := c.cc.Invoke(ctx, LoginService_FindMemInfoById_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LoginServiceServer is the server API for LoginService service.
// All implementations must embed UnimplementedLoginServiceServer
// for forward compatibility.
type LoginServiceServer interface {
	GetCaptcha(context.Context, *CaptchaMessage) (*CaptchaResponse, error)
	Register(context.Context, *RegisterMessage) (*RegisterResponse, error)
	Login(context.Context, *LoginMessage) (*LoginResponse, error)
	TokenVerify(context.Context, *LoginMessage) (*LoginResponse, error)
	MyOrgList(context.Context, *UserMessage) (*OrgListResponse, error)
	FindMemInfoById(context.Context, *UserMessage) (*MemberMessage, error)
	mustEmbedUnimplementedLoginServiceServer()
}

// UnimplementedLoginServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedLoginServiceServer struct{}

func (UnimplementedLoginServiceServer) GetCaptcha(context.Context, *CaptchaMessage) (*CaptchaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCaptcha not implemented")
}
func (UnimplementedLoginServiceServer) Register(context.Context, *RegisterMessage) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedLoginServiceServer) Login(context.Context, *LoginMessage) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedLoginServiceServer) TokenVerify(context.Context, *LoginMessage) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TokenVerify not implemented")
}
func (UnimplementedLoginServiceServer) MyOrgList(context.Context, *UserMessage) (*OrgListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MyOrgList not implemented")
}
func (UnimplementedLoginServiceServer) FindMemInfoById(context.Context, *UserMessage) (*MemberMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindMemInfoById not implemented")
}
func (UnimplementedLoginServiceServer) mustEmbedUnimplementedLoginServiceServer() {}
func (UnimplementedLoginServiceServer) testEmbeddedByValue()                      {}

// UnsafeLoginServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LoginServiceServer will
// result in compilation errors.
type UnsafeLoginServiceServer interface {
	mustEmbedUnimplementedLoginServiceServer()
}

func RegisterLoginServiceServer(s grpc.ServiceRegistrar, srv LoginServiceServer) {
	// If the following call pancis, it indicates UnimplementedLoginServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&LoginService_ServiceDesc, srv)
}

func _LoginService_GetCaptcha_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CaptchaMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoginServiceServer).GetCaptcha(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LoginService_GetCaptcha_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoginServiceServer).GetCaptcha(ctx, req.(*CaptchaMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoginService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoginServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LoginService_Register_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoginServiceServer).Register(ctx, req.(*RegisterMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoginService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoginServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LoginService_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoginServiceServer).Login(ctx, req.(*LoginMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoginService_TokenVerify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoginServiceServer).TokenVerify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LoginService_TokenVerify_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoginServiceServer).TokenVerify(ctx, req.(*LoginMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoginService_MyOrgList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoginServiceServer).MyOrgList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LoginService_MyOrgList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoginServiceServer).MyOrgList(ctx, req.(*UserMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoginService_FindMemInfoById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoginServiceServer).FindMemInfoById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LoginService_FindMemInfoById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoginServiceServer).FindMemInfoById(ctx, req.(*UserMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// LoginService_ServiceDesc is the grpc.ServiceDesc for LoginService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LoginService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "login.LoginService",
	HandlerType: (*LoginServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCaptcha",
			Handler:    _LoginService_GetCaptcha_Handler,
		},
		{
			MethodName: "Register",
			Handler:    _LoginService_Register_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _LoginService_Login_Handler,
		},
		{
			MethodName: "TokenVerify",
			Handler:    _LoginService_TokenVerify_Handler,
		},
		{
			MethodName: "MyOrgList",
			Handler:    _LoginService_MyOrgList_Handler,
		},
		{
			MethodName: "FindMemInfoById",
			Handler:    _LoginService_FindMemInfoById_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "login_service.proto",
}
