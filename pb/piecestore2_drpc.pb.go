// Code generated by protoc-gen-go-drpc. DO NOT EDIT.
// protoc-gen-go-drpc version: v0.0.33-0.20230222225142-dbfc88394a82
// source: piecestore2.proto

package pb

import (
	bytes "bytes"
	context "context"
	errors "errors"

	jsonpb "github.com/gogo/protobuf/jsonpb"
	proto "github.com/gogo/protobuf/proto"

	drpc "drpc"
	drpcerr "drpc/drpcerr"
)

type drpcEncoding_File_piecestore2_proto struct{}

func (drpcEncoding_File_piecestore2_proto) Marshal(msg drpc.Message) ([]byte, error) {
	return proto.Marshal(msg.(proto.Message))
}

func (drpcEncoding_File_piecestore2_proto) Unmarshal(buf []byte, msg drpc.Message) error {
	return proto.Unmarshal(buf, msg.(proto.Message))
}

func (drpcEncoding_File_piecestore2_proto) JSONMarshal(msg drpc.Message) ([]byte, error) {
	var buf bytes.Buffer
	err := new(jsonpb.Marshaler).Marshal(&buf, msg.(proto.Message))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (drpcEncoding_File_piecestore2_proto) JSONUnmarshal(buf []byte, msg drpc.Message) error {
	return jsonpb.Unmarshal(bytes.NewReader(buf), msg.(proto.Message))
}

type DRPCPiecestoreClient interface {
	DRPCConn() drpc.Conn

	Upload(ctx context.Context) (DRPCPiecestore_UploadClient, error)
	Download(ctx context.Context) (DRPCPiecestore_DownloadClient, error)
	Delete(ctx context.Context, in *PieceDeleteRequest) (*PieceDeleteResponse, error)
	DeletePieces(ctx context.Context, in *DeletePiecesRequest) (*DeletePiecesResponse, error)
	Retain(ctx context.Context, in *RetainRequest) (*RetainResponse, error)
	RestoreTrash(ctx context.Context, in *RestoreTrashRequest) (*RestoreTrashResponse, error)
	Exists(ctx context.Context, in *ExistsRequest) (*ExistsResponse, error)
}

type drpcPiecestoreClient struct {
	cc drpc.Conn
}

func NewDRPCPiecestoreClient(cc drpc.Conn) DRPCPiecestoreClient {
	return &drpcPiecestoreClient{cc}
}

func (c *drpcPiecestoreClient) DRPCConn() drpc.Conn { return c.cc }

func (c *drpcPiecestoreClient) Upload(ctx context.Context) (DRPCPiecestore_UploadClient, error) {
	stream, err := c.cc.NewStream(ctx, "/piecestore.Piecestore/Upload", drpcEncoding_File_piecestore2_proto{})
	if err != nil {
		return nil, err
	}
	x := &drpcPiecestore_UploadClient{stream}
	return x, nil
}

type DRPCPiecestore_UploadClient interface {
	drpc.Stream
	Send(*PieceUploadRequest) error
	CloseAndRecv() (*PieceUploadResponse, error)
}

type drpcPiecestore_UploadClient struct {
	drpc.Stream
}

func (x *drpcPiecestore_UploadClient) GetStream() drpc.Stream {
	return x.Stream
}

func (x *drpcPiecestore_UploadClient) Send(m *PieceUploadRequest) error {
	return x.MsgSend(m, drpcEncoding_File_piecestore2_proto{})
}

func (x *drpcPiecestore_UploadClient) CloseAndRecv() (*PieceUploadResponse, error) {
	if err := x.CloseSend(); err != nil {
		return nil, err
	}
	m := new(PieceUploadResponse)
	if err := x.MsgRecv(m, drpcEncoding_File_piecestore2_proto{}); err != nil {
		return nil, err
	}
	return m, nil
}

func (x *drpcPiecestore_UploadClient) CloseAndRecvMsg(m *PieceUploadResponse) error {
	if err := x.CloseSend(); err != nil {
		return err
	}
	return x.MsgRecv(m, drpcEncoding_File_piecestore2_proto{})
}

func (c *drpcPiecestoreClient) Download(ctx context.Context) (DRPCPiecestore_DownloadClient, error) {
	stream, err := c.cc.NewStream(ctx, "/piecestore.Piecestore/Download", drpcEncoding_File_piecestore2_proto{})
	if err != nil {
		return nil, err
	}
	x := &drpcPiecestore_DownloadClient{stream}
	return x, nil
}

type DRPCPiecestore_DownloadClient interface {
	drpc.Stream
	Send(*PieceDownloadRequest) error
	Recv() (*PieceDownloadResponse, error)
}

type drpcPiecestore_DownloadClient struct {
	drpc.Stream
}

func (x *drpcPiecestore_DownloadClient) GetStream() drpc.Stream {
	return x.Stream
}

func (x *drpcPiecestore_DownloadClient) Send(m *PieceDownloadRequest) error {
	return x.MsgSend(m, drpcEncoding_File_piecestore2_proto{})
}

func (x *drpcPiecestore_DownloadClient) Recv() (*PieceDownloadResponse, error) {
	m := new(PieceDownloadResponse)
	if err := x.MsgRecv(m, drpcEncoding_File_piecestore2_proto{}); err != nil {
		return nil, err
	}
	return m, nil
}

func (x *drpcPiecestore_DownloadClient) RecvMsg(m *PieceDownloadResponse) error {
	return x.MsgRecv(m, drpcEncoding_File_piecestore2_proto{})
}

func (c *drpcPiecestoreClient) Delete(ctx context.Context, in *PieceDeleteRequest) (*PieceDeleteResponse, error) {
	out := new(PieceDeleteResponse)
	err := c.cc.Invoke(ctx, "/piecestore.Piecestore/Delete", drpcEncoding_File_piecestore2_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcPiecestoreClient) DeletePieces(ctx context.Context, in *DeletePiecesRequest) (*DeletePiecesResponse, error) {
	out := new(DeletePiecesResponse)
	err := c.cc.Invoke(ctx, "/piecestore.Piecestore/DeletePieces", drpcEncoding_File_piecestore2_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcPiecestoreClient) Retain(ctx context.Context, in *RetainRequest) (*RetainResponse, error) {
	out := new(RetainResponse)
	err := c.cc.Invoke(ctx, "/piecestore.Piecestore/Retain", drpcEncoding_File_piecestore2_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcPiecestoreClient) RestoreTrash(ctx context.Context, in *RestoreTrashRequest) (*RestoreTrashResponse, error) {
	out := new(RestoreTrashResponse)
	err := c.cc.Invoke(ctx, "/piecestore.Piecestore/RestoreTrash", drpcEncoding_File_piecestore2_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcPiecestoreClient) Exists(ctx context.Context, in *ExistsRequest) (*ExistsResponse, error) {
	out := new(ExistsResponse)
	err := c.cc.Invoke(ctx, "/piecestore.Piecestore/Exists", drpcEncoding_File_piecestore2_proto{}, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type DRPCPiecestoreServer interface {
	Upload(DRPCPiecestore_UploadStream) error
	Download(DRPCPiecestore_DownloadStream) error
	Delete(context.Context, *PieceDeleteRequest) (*PieceDeleteResponse, error)
	DeletePieces(context.Context, *DeletePiecesRequest) (*DeletePiecesResponse, error)
	Retain(context.Context, *RetainRequest) (*RetainResponse, error)
	RestoreTrash(context.Context, *RestoreTrashRequest) (*RestoreTrashResponse, error)
	Exists(context.Context, *ExistsRequest) (*ExistsResponse, error)
}

type DRPCPiecestoreUnimplementedServer struct{}

func (s *DRPCPiecestoreUnimplementedServer) Upload(DRPCPiecestore_UploadStream) error {
	return drpcerr.WithCode(errors.New("Unimplemented"), drpcerr.Unimplemented)
}

func (s *DRPCPiecestoreUnimplementedServer) Download(DRPCPiecestore_DownloadStream) error {
	return drpcerr.WithCode(errors.New("Unimplemented"), drpcerr.Unimplemented)
}

func (s *DRPCPiecestoreUnimplementedServer) Delete(context.Context, *PieceDeleteRequest) (*PieceDeleteResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), drpcerr.Unimplemented)
}

func (s *DRPCPiecestoreUnimplementedServer) DeletePieces(context.Context, *DeletePiecesRequest) (*DeletePiecesResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), drpcerr.Unimplemented)
}

func (s *DRPCPiecestoreUnimplementedServer) Retain(context.Context, *RetainRequest) (*RetainResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), drpcerr.Unimplemented)
}

func (s *DRPCPiecestoreUnimplementedServer) RestoreTrash(context.Context, *RestoreTrashRequest) (*RestoreTrashResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), drpcerr.Unimplemented)
}

func (s *DRPCPiecestoreUnimplementedServer) Exists(context.Context, *ExistsRequest) (*ExistsResponse, error) {
	return nil, drpcerr.WithCode(errors.New("Unimplemented"), drpcerr.Unimplemented)
}

type DRPCPiecestoreDescription struct{}

func (DRPCPiecestoreDescription) NumMethods() int { return 7 }

func (DRPCPiecestoreDescription) Method(n int) (string, drpc.Encoding, drpc.Receiver, interface{}, bool) {
	switch n {
	case 0:
		return "/piecestore.Piecestore/Upload", drpcEncoding_File_piecestore2_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return nil, srv.(DRPCPiecestoreServer).
					Upload(
						&drpcPiecestore_UploadStream{in1.(drpc.Stream)},
					)
			}, DRPCPiecestoreServer.Upload, true
	case 1:
		return "/piecestore.Piecestore/Download", drpcEncoding_File_piecestore2_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return nil, srv.(DRPCPiecestoreServer).
					Download(
						&drpcPiecestore_DownloadStream{in1.(drpc.Stream)},
					)
			}, DRPCPiecestoreServer.Download, true
	case 2:
		return "/piecestore.Piecestore/Delete", drpcEncoding_File_piecestore2_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCPiecestoreServer).
					Delete(
						ctx,
						in1.(*PieceDeleteRequest),
					)
			}, DRPCPiecestoreServer.Delete, true
	case 3:
		return "/piecestore.Piecestore/DeletePieces", drpcEncoding_File_piecestore2_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCPiecestoreServer).
					DeletePieces(
						ctx,
						in1.(*DeletePiecesRequest),
					)
			}, DRPCPiecestoreServer.DeletePieces, true
	case 4:
		return "/piecestore.Piecestore/Retain", drpcEncoding_File_piecestore2_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCPiecestoreServer).
					Retain(
						ctx,
						in1.(*RetainRequest),
					)
			}, DRPCPiecestoreServer.Retain, true
	case 5:
		return "/piecestore.Piecestore/RestoreTrash", drpcEncoding_File_piecestore2_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCPiecestoreServer).
					RestoreTrash(
						ctx,
						in1.(*RestoreTrashRequest),
					)
			}, DRPCPiecestoreServer.RestoreTrash, true
	case 6:
		return "/piecestore.Piecestore/Exists", drpcEncoding_File_piecestore2_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCPiecestoreServer).
					Exists(
						ctx,
						in1.(*ExistsRequest),
					)
			}, DRPCPiecestoreServer.Exists, true
	default:
		return "", nil, nil, nil, false
	}
}

func DRPCRegisterPiecestore(mux drpc.Mux, impl DRPCPiecestoreServer) error {
	return mux.Register(impl, DRPCPiecestoreDescription{})
}

type DRPCPiecestore_UploadStream interface {
	drpc.Stream
	SendAndClose(*PieceUploadResponse) error
	Recv() (*PieceUploadRequest, error)
}

type drpcPiecestore_UploadStream struct {
	drpc.Stream
}

func (x *drpcPiecestore_UploadStream) SendAndClose(m *PieceUploadResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_piecestore2_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

func (x *drpcPiecestore_UploadStream) Recv() (*PieceUploadRequest, error) {
	m := new(PieceUploadRequest)
	if err := x.MsgRecv(m, drpcEncoding_File_piecestore2_proto{}); err != nil {
		return nil, err
	}
	return m, nil
}

func (x *drpcPiecestore_UploadStream) RecvMsg(m *PieceUploadRequest) error {
	return x.MsgRecv(m, drpcEncoding_File_piecestore2_proto{})
}

type DRPCPiecestore_DownloadStream interface {
	drpc.Stream
	Send(*PieceDownloadResponse) error
	Recv() (*PieceDownloadRequest, error)
}

type drpcPiecestore_DownloadStream struct {
	drpc.Stream
}

func (x *drpcPiecestore_DownloadStream) Send(m *PieceDownloadResponse) error {
	return x.MsgSend(m, drpcEncoding_File_piecestore2_proto{})
}

func (x *drpcPiecestore_DownloadStream) Recv() (*PieceDownloadRequest, error) {
	m := new(PieceDownloadRequest)
	if err := x.MsgRecv(m, drpcEncoding_File_piecestore2_proto{}); err != nil {
		return nil, err
	}
	return m, nil
}

func (x *drpcPiecestore_DownloadStream) RecvMsg(m *PieceDownloadRequest) error {
	return x.MsgRecv(m, drpcEncoding_File_piecestore2_proto{})
}

type DRPCPiecestore_DeleteStream interface {
	drpc.Stream
	SendAndClose(*PieceDeleteResponse) error
}

type drpcPiecestore_DeleteStream struct {
	drpc.Stream
}

func (x *drpcPiecestore_DeleteStream) SendAndClose(m *PieceDeleteResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_piecestore2_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCPiecestore_DeletePiecesStream interface {
	drpc.Stream
	SendAndClose(*DeletePiecesResponse) error
}

type drpcPiecestore_DeletePiecesStream struct {
	drpc.Stream
}

func (x *drpcPiecestore_DeletePiecesStream) SendAndClose(m *DeletePiecesResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_piecestore2_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCPiecestore_RetainStream interface {
	drpc.Stream
	SendAndClose(*RetainResponse) error
}

type drpcPiecestore_RetainStream struct {
	drpc.Stream
}

func (x *drpcPiecestore_RetainStream) SendAndClose(m *RetainResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_piecestore2_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCPiecestore_RestoreTrashStream interface {
	drpc.Stream
	SendAndClose(*RestoreTrashResponse) error
}

type drpcPiecestore_RestoreTrashStream struct {
	drpc.Stream
}

func (x *drpcPiecestore_RestoreTrashStream) SendAndClose(m *RestoreTrashResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_piecestore2_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCPiecestore_ExistsStream interface {
	drpc.Stream
	SendAndClose(*ExistsResponse) error
}

type drpcPiecestore_ExistsStream struct {
	drpc.Stream
}

func (x *drpcPiecestore_ExistsStream) SendAndClose(m *ExistsResponse) error {
	if err := x.MsgSend(m, drpcEncoding_File_piecestore2_proto{}); err != nil {
		return err
	}
	return x.CloseSend()
}
