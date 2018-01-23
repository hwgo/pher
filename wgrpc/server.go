package wgrpc

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/hwgo/pher/tracing"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/hwgo/pher/log"
	"github.com/hwgo/pher/metrics"
	"github.com/hwgo/pher/otgrpc"
)

type Server struct {
	hostPort string
	tracer   opentracing.Tracer
	logger   log.Factory
	Gs       *grpc.Server
}

// Run starts the Customer server
func (s *Server) Run() error {
	bg := s.logger.Bg()
	lis, err := net.Listen("tcp", s.hostPort)

	if err != nil {
		bg.Fatal("Unable to start server", zap.Error(err))
		return err
	}

	bg.Info("Starting", zap.String("address", "tcp://"+s.hostPort))
	return s.Gs.Serve(lis)
}

func NewServerWithTracing(name string, hostPort string) *Server {
	logger := log.NewFactory(log.DefaultLogger.With(zap.String("service", name)))
	metricsFactory := metrics.DefaultMetricsFactory()
	tracer := tracing.Init(name, metricsFactory.Namespace(name, nil), logger)

	th := otgrpc.NewTraceHandler(tracer)
	s := grpc.NewServer(grpc.StatsHandler(th))

	// Register reflection service on gRPC server.
	reflection.Register(s)

	return &Server{
		hostPort: hostPort,
		tracer:   tracer,
		logger:   logger,
		Gs:       s,
	}
}
