package user

import (
	"github.com/Tuanzi-bug/SyncHub/api/config"
	"github.com/Tuanzi-bug/SyncHub/common/discovery"
	"github.com/Tuanzi-bug/SyncHub/common/logs"
	loginServiceV1 "github.com/Tuanzi-bug/SyncHub/grpc/user/login"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
)

var LoginServiceClient loginServiceV1.LoginServiceClient

func InitRpcUserClient() {
	// 初始化etcd服务发现，注册解析器
	etcdRegister := discovery.NewResolver(config.AppConfig.EtcdConfig.Addrs, logs.LG)
	resolver.Register(etcdRegister)

	conn, err := grpc.NewClient(etcdRegister.Scheme()+":///user_proto", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	LoginServiceClient = loginServiceV1.NewLoginServiceClient(conn)
}
