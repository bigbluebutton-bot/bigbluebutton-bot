// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: changeset.proto

package changesetproto

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

// ChangesetClient is the client API for Changeset service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChangesetClient interface {
	// Sends a greeting
	Generate(ctx context.Context, in *GenerateRequest, opts ...grpc.CallOption) (*GenerateReply, error)
}

type changesetClient struct {
	cc grpc.ClientConnInterface
}

func NewChangesetClient(cc grpc.ClientConnInterface) ChangesetClient {
	return &changesetClient{cc}
}

func (c *changesetClient) Generate(ctx context.Context, in *GenerateRequest, opts ...grpc.CallOption) (*GenerateReply, error) {
	out := new(GenerateReply)
	err := c.cc.Invoke(ctx, "/changset.Changeset/Generate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChangesetServer is the server API for Changeset service.
// All implementations must embed UnimplementedChangesetServer
// for forward compatibility
type ChangesetServer interface {
	// Sends a greeting
	Generate(context.Context, *GenerateRequest) (*GenerateReply, error)
	mustEmbedUnimplementedChangesetServer()
}

// UnimplementedChangesetServer must be embedded to have forward compatible implementations.
type UnimplementedChangesetServer struct {
}

func (UnimplementedChangesetServer) Generate(context.Context, *GenerateRequest) (*GenerateReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Generate not implemented")
}
func (UnimplementedChangesetServer) mustEmbedUnimplementedChangesetServer() {}

// UnsafeChangesetServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChangesetServer will
// result in compilation errors.
type UnsafeChangesetServer interface {
	mustEmbedUnimplementedChangesetServer()
}

func RegisterChangesetServer(s grpc.ServiceRegistrar, srv ChangesetServer) {
	s.RegisterService(&Changeset_ServiceDesc, srv)
}

func _Changeset_Generate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChangesetServer).Generate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/changset.Changeset/Generate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChangesetServer).Generate(ctx, req.(*GenerateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Changeset_ServiceDesc is the grpc.ServiceDesc for Changeset service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Changeset_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "changset.Changeset",
	HandlerType: (*ChangesetServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Generate",
			Handler:    _Changeset_Generate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "changeset.proto",
}