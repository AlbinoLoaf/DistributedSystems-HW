package Client
import (

    // This has to be the same as the go.mod module,
    // followed by the path to the folder the proto file is in.
    gRPC "HW03/proto"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

var opts []grpc.DialOption
opts = append(
    opts, grpc.WithBlock(), 
    grpc.WithTransportCredentials(insecure.NewCredentials())
)
conn, err := grpc.Dial(":5400", opts...)

server = gRPC.NewTwitterClient(conn)

stream, err := server.SayHi(context.Background())

message := // make an instance of your input ty
stream.Send(message)

farewell, err := stream.CloseAndRecv()


