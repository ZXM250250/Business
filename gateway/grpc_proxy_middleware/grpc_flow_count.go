package grpc_proxy_middleware

import (
	"gateway/dao"
	"gateway/tools"
	"google.golang.org/grpc"
	"log"
)

func GrpcFlowCountMiddleware(serviceDetail *dao.ServiceDetail) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		totalCounter, err := tools.FlowCounterHandler.GetCounter(tools.FlowTotal)
		if err != nil {
			return err
		}
		totalCounter.Increase()
		serviceCounter, err := tools.FlowCounterHandler.GetCounter(tools.FlowServicePrefix + serviceDetail.Info.ServiceName)
		if err != nil {
			return err
		}
		serviceCounter.Increase()

		if err := handler(srv, ss); err != nil {
			log.Printf("GrpcFlowCountMiddleware failed with error %v\n", err)
			return err
		}
		return nil
	}
}
