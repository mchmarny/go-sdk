package client

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/dapr/go-sdk/dapr/proto/runtime/v1"
)

const (
	daprPortDefault    = "50001"
	daprPortEnvVarName = "DAPR_GRPC_PORT"
)

var (
	logger = log.New(os.Stdout, "", 0)
)

// NewClient instantiates dapr client locally using port from DAPR_GRPC_PORT env var
func NewClient() (client *Client, err error) {
	port := os.Getenv(daprPortEnvVarName)
	if port == "" {
		port = daprPortDefault
	}
	return NewClientWithPort(port)
}

// NewClientWithPort instantiates dapr client locally for the specific port
func NewClientWithPort(port string) (client *Client, err error) {
	if port == "" {
		return nil, errors.New("nil port")
	}
	return NewClientWithAddress(net.JoinHostPort("127.0.0.1", port))
}

// NewClientWithAddress instantiates dapr client configured for the specific address
func NewClientWithAddress(address string) (client *Client, err error) {
	if address == "" {
		return nil, errors.New("nil address")
	}
	logger.Printf("dapr client initializing for: %s", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrapf(err, "error creating connection to '%s': %v", address, err)
	}
	client = &Client{
		connection:  conn,
		protoClient: pb.NewDaprClient(conn),
	}
	return
}

// NewClientWithConnection instantiates dapr client configured for the specific connection
func NewClientWithConnection(conn *grpc.ClientConn) (client *Client, err error) {
	client = &Client{
		connection:  conn,
		protoClient: pb.NewDaprClient(conn),
	}
	return
}

// Client is the dapr client
type Client struct {
	connection  *grpc.ClientConn
	protoClient pb.DaprClient
}

// Close cleans up all resources created by the client
func (c *Client) Close() {
	if c.connection != nil {
		c.connection.Close()
	}
}

func authContext(ctx context.Context) context.Context {
	token := os.Getenv("DAPR_API_TOKEN")
	if token == "" {
		return ctx
	}
	md := metadata.Pairs("dapr-api-token", token)
	return metadata.NewOutgoingContext(ctx, md)
}
