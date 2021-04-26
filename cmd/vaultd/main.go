package main

import (
	"flag"
	"fmt"
	vault "github.com/akolybelnikov/grpcvault"
	pb "github.com/akolybelnikov/grpcvault/pb"
	"github.com/go-kit/kit/ratelimit"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var (
		httpAddr = flag.String("http", ":8080", "http listen address")
		gRPCAddr = flag.String("grpc", ":8081", "gRPC listen address")
	)
	flag.Parse()
	srv := vault.NewService()
	errChan := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	limit := rate.NewLimiter(rate.Every(time.Second), 5)

	hashEndpoint := ratelimit.NewDelayingLimiter(limit)(vault.MakeHashEndpoint(srv))
	validateEndpoint := ratelimit.NewDelayingLimiter(limit)(vault.MakeValidateEndpoint(srv))

	endpoints := vault.Endpoints{HashEndpoint: hashEndpoint, ValidateEndpoint: validateEndpoint}

	// HTTP transport
	go func() {
		log.Println("http:", *httpAddr)
		handler := vault.NewHTTPServer(endpoints)
		errChan <- http.ListenAndServe(*httpAddr, handler)
	}()

	// gRPC transport
	go func() {
		listener, err := net.Listen("tcp", *gRPCAddr)
		if err != nil {
			errChan <- err
			return
		}
		log.Println("grpc:", *gRPCAddr)
		handler := vault.NewGRPCServer(endpoints)
		gRPCServer := grpc.NewServer()
		pb.RegisterVaultServer(gRPCServer, handler)
		errChan <- gRPCServer.Serve(listener)
	}()

	log.Fatalln(<-errChan)
}
