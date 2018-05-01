package lndmobile

import (
	"context"
	"fmt"
	"os"

	"github.com/golang/protobuf/proto"
	flags "github.com/jessevdk/go-flags"
	"github.com/lightningnetwork/lnd/daemon"
	"github.com/lightningnetwork/lnd/lnrpc"
)

// TODO: move to build script.
//go:generate go build gen_bindings.go
//go:generate ./gen_bindings.sh
//go:generate rm ./gen_bindings

func Start(appDir string, callback Callback) {
	// Call the "real" main in a nested manner so the defers will properly
	// be executed in the case of a graceful shutdown.
	go func() {
		if err := daemon.LndMain(appDir, lis); err != nil {
			if e, ok := err.(*flags.Error); ok &&
				e.Type == flags.ErrHelp {
			} else {
				fmt.Fprintln(os.Stderr, err)
			}
			os.Exit(1)
		}
	}()

	// TODO(halseth): callback when RPC server is running.
	callback.OnResponse([]byte("started"))
}

// GetInfo takes a serialized GetInfoRequest and returns single serialized
// GetInfoResponse to the provided callback.
func GetInfo(msg []byte, callback Callback) {
	s := &onceHandler{
		newProto: func() proto.Message {
			return &lnrpc.GetInfoRequest{}
		},
		getSync: func(ctx context.Context, client lnrpc.LightningClient,
			req proto.Message) (proto.Message, error) {
			r := req.(*lnrpc.GetInfoRequest)
			return client.GetInfo(ctx, r)
		},
	}
	s.start(msg, callback)
}

// SubscribeInvoices takes a serialized InvoiceSubscription and returns a
// stream of Invoices to the provided callback.
func SubscribeInvoices(msg []byte, callback Callback) {
	s := &readStreamHandler{
		newProto: func() proto.Message {
			return &lnrpc.InvoiceSubscription{}
		},
		recvStream: func(ctx context.Context,
			client lnrpc.LightningClient,
			req proto.Message) (*receiver, error) {
			r := req.(*lnrpc.InvoiceSubscription)
			stream, err := client.SubscribeInvoices(ctx, r)
			if err != nil {
				return nil, err
			}
			return &receiver{
				recv: func() (proto.Message, error) {
					return stream.Recv()
				},
			}, nil
		},
	}
	s.start(msg, callback)
}

// SendPayment opens a bidirectional payment stream to the server, letting the
// caller send serialized SendRequests on the returned SendStream. Serialized
// SendResponses will be delivered to the provided callback.
func SendPayment(callback Callback) (SendStream, error) {
	b := &biStreamHandler{
		newProto: func() proto.Message {
			return &lnrpc.SendRequest{}
		},
		biStream: func(ctx context.Context,
			client lnrpc.LightningClient) (*receiver, *sender,
			error) {
			stream, err := client.SendPayment(ctx)
			if err != nil {
				return nil, nil, err
			}
			return &receiver{
					recv: func() (proto.Message, error) {
						return stream.Recv()
					},
				},
				&sender{
					send: func(req proto.Message) error {
						r := req.(*lnrpc.SendRequest)
						return stream.Send(r)
					},
					close: stream.CloseSend,
				}, nil
		},
	}
	return b.start(callback)
}
