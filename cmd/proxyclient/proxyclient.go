package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sakirsensoy/genv"
	"github.com/sakirsensoy/genv/dotenv"
	"github.com/tkcrm/modules/logger"
	"github.com/tkcrm/modules/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type server struct {
	logger     logger.Logger
	remotePort int
	devKey     string
	remoteHost string

	grpcHealthConn   *grpc.ClientConn
	grpcHealthClient grpc_health_v1.HealthClient
}

func main() {
	l := logger.DefaultLogger(utils.GetDefaultString(os.Getenv("LOG_LEVEL"), "info"), "tracker-proxy-client")

	remotePort := flag.Int("port", 0, "relayd port")

	flag.Parse()

	if *remotePort == 0 {
		l.Fatal("Undefined port")
	}

	dotenv.Load()

	proxyDevKey := genv.Key("PROXY_DEV_KEY").String()
	if proxyDevKey == "" {
		l.Fatal("Undefined dev key")
	}

	proxyRemoteHost := genv.Key("PROXY_REMOTE_HOST").String()
	if proxyRemoteHost == "" {
		proxyRemoteHost = "relay.tracker.tkcrm.ru"
	}

	grpcPort := genv.Key("GRPC_PORT").String()
	if grpcPort == "" {
		l.Fatal("Undefined dev key")
	}

	s := server{
		logger:     l,
		remotePort: *remotePort,
		devKey:     proxyDevKey,
		remoteHost: proxyRemoteHost,
	}

	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", "localhost", grpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		s.logger.Fatalf("GRPC connect error: %v", err)
	}

	s.logger.Info("GRPC: Connected to ms tracker")

	s.grpcHealthConn = conn
	s.grpcHealthClient = grpc_health_v1.NewHealthClient(conn)
	defer func() {
		s.grpcHealthConn.Close()
	}()

	cn := make(chan os.Signal, 1)

	go func() {
		for {
			if err := s.start(); err != nil {
				s.logger.Errorf("Server stopped with error: %v", err)
			}
			time.Sleep(time.Second)
			s.logger.Errorf("Try to restart server")
		}
	}()

	signal.Notify(cn, os.Interrupt, syscall.SIGTERM)
	<-cn
	os.Exit(1)
}

func (s *server) start() error {

	// Устанавливаем локальное соединение
	conn, err := net.DialTimeout("tcp4", fmt.Sprintf(":%d", s.remotePort), time.Second)
	if err != nil {
		return fmt.Errorf("connect to local server error: %v", err)
	}
	defer func() {
		s.logger.Info("Disconnected from local server")
		conn.Close()
	}()
	s.logger.Info("Connected to local server")

	// Устанавливаем удаленное соединение
	remoteConn, err := net.DialTimeout("tcp4", fmt.Sprintf("%s:%d", s.remoteHost, s.remotePort+1), time.Second)
	if err != nil {
		return fmt.Errorf("connect to remote server with host: %s error: %v", s.remoteHost, err)
	}
	defer func() {
		s.logger.Info("Disconnected from remote server")
		remoteConn.Close()
	}()

	s.logger.Info("Try to authenticate on relay server")

	if _, err := remoteConn.Write([]byte(s.devKey)); err != nil {
		s.logger.Fatal(err)
	}

	resultAuth := make([]byte, 1)
	if _, err := remoteConn.Read(resultAuth); err != nil {
		s.logger.Fatal(err)
	}

	if resultAuth[0] != 0x01 {
		s.logger.Fatal("Authentication faield")
	}

	s.logger.Info("Authenticated on relay server")

	errChan := make(chan error, 3)

	go func() {
		if _, err := io.Copy(remoteConn, conn); err != nil {
			s.logger.Error(err)
			errChan <- err
		}
	}()
	go func() {
		if _, err := io.Copy(conn, remoteConn); err != nil {
			s.logger.Error(err)
			errChan <- err
		}
	}()
	go func() {
		for {
			time.Sleep(time.Millisecond * 100)
			_, err := s.grpcHealthClient.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
			if err != nil {
				errChan <- err
			}
		}
	}()

	return <-errChan
}
