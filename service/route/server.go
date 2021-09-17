package route

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go_im/pkg/logger"
	"go_im/service/pb"
	"go_im/service/rpc"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	ExtraTag        = "rt_extra_tag"
	ExtraSrvUrl     = "rt_extra_srv_url"
	ExtraFrom       = "rt_extra_from"
	ExtraSelectMode = "rt_extra_select_mode"
)

type Server struct {
	*rpc.BaseServer
	rts map[string]*service
}

func NewServer(options *rpc.ServerOptions) *Server {
	s := &Server{
		BaseServer: rpc.NewBaseServer(options),
		rts:        map[string]*service{},
	}
	s.BaseServer.Register(options.Name, s)
	return s
}

func (s *Server) SetTag(ctx context.Context, req *pb.SetTagReq, _ *emptypb.Empty) error {
	rt, ok := s.rts[req.SrvId]
	if !ok {
		return fmt.Errorf("service not found: srvId=%s", req.SrvId)
	}
	rt.addTag(req.GetTag(), req.GetValue())
	return nil
}

func (s *Server) RemoveTag(ctx context.Context, req *pb.ClearTagReq, _ *emptypb.Empty) error {
	rt, ok := s.rts[req.SrvId]
	if !ok {
		return fmt.Errorf("service not found: srvId=%s", req.SrvId)
	}
	rt.removeTag(req.GetTag())
	return nil
}

func (s *Server) Route(ctx context.Context, param *pb.RouteReq, reply *pb.RouteReply) error {
	rt, ok := s.rts[param.SrvId]
	if !ok {
		return fmt.Errorf("service not register: srvId=%s", param.SrvId)
	}
	reply.Success = true
	reply.Msg = "success"
	reply.Reply = &anypb.Any{}

	err := rt.route(ctx, param.Fn, param, reply)
	if err != nil {
		reply.Success = false
		reply.Msg = err.Error()
		return errors.Wrap(err, "service route error")
	}
	return nil
}

func (s *Server) Unregister(ctx context.Context, param *pb.UnRegisterReq, _ *emptypb.Empty) error {
	rv, ok := s.rts[param.SrvId]
	if ok {
		return rv.Close()
	}
	return errors.New("service not register")
}

func (s *Server) Register(ctx context.Context, param *pb.RegisterRtReq, _ *emptypb.Empty) error {
	sv := newService(&rpc.ClientOptions{
		Name:        param.GetSrvName(),
		EtcdServers: param.GetDiscoverySrvUrl(),
	})
	err := sv.BaseClient.Run()
	if err != nil {
		return err
	}
	_, ok := s.rts[param.SrvId]
	if ok {
		// override
	}
	s.rts[param.SrvId] = sv
	logger.D("service registered: %s", param.SrvName)
	return nil
}
