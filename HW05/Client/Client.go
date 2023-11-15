package main

import (
	proto "/grpc"
	"context"
	"fmt"
	"log"
	"time"
)

type BidClient struct {
}

func main() {

	for i := 0; i < 4; i++ {
		port := int64(5000) + int64(i/2)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		var conn *proto.ClientConn
		defer cancel()
		conn, err = proto.DialContext(ctx, fmt.Sprintf(":%v", port), opts...)
		if err != nil {
			// If there's an error, sleep for a while before retrying
			log.Printf("Attempt %d: did not connect: %v", i+1, err)
			time.Sleep(time.Second * 5)
		} else {
			// If connection is successful, break out of the loop
			break
		}
	}

}
