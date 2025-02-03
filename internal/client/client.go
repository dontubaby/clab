/* Данный пакет является gRPC клиентов сервиса game-actions*/

package client

import (
	"context"
	"log"
	"time"

	pb "cyberball/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GameLogic Обрабатывает gRPC соединения с другими сервисами
type GameLogicClient struct {
	pb.UnimplementedGameLogicServiceServer
}

func (s *GameLogic) GetAction(ctx context.Context, rexq *pb.ActionRequest) (*pb.ActionResponse, error) {

}

func (s *GameLogic) AddAction(ctx context.Context, rexq *pb.LogicRequest) (*pb.ActionResponse, error) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInesecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("cant connect game-logic client by gRPC:%v\n", err)
	}
	defer conn.Close()

	client := pb.NewGameLogicServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	actionResponse, err := client.AddAction(ctx)
}
