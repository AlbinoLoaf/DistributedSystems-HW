package server
import (
    
    // This has to be the same as the go.mod module,
    // followed by the path to the folder the proto file is in.
    gRPC "HW03/proto"

    "google.golang.org/grpc"
)

type Server struct {
	// an interface that the server type needs to have
	gRPC.UnimplementedTemplateServer

	// here you can implement other fields that you want
}

func (s *Server) proto(msgStream gRPC.<service name>_<protocol.proto>Server) error {
    // get all the messages from the stream
    for {
        msg, err := msgStream.Recv()
        if err == io.EOF {
            break
        }
    }

    ack := // make an instance of your return type
    msgStream.SendAndClose(ack)

    return nil
}
func main() {
    list, _ := net.Listen("tcp", "localhost:5400")
	grpcServer := grpc.NewServer(opts...)
	server := &Server{
		// set fields here 
	}
	gRPC.RegisterTwitterServer(grpcServer, server)
	grpcServer.Serve(list)
// Code here will not run as .Serve() blocks the thread
}