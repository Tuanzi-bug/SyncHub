package router

import (
	"github.com/Tuanzi-bug/SyncHub/common/discovery"
	"github.com/Tuanzi-bug/SyncHub/common/logs"
	"github.com/Tuanzi-bug/SyncHub/grpc/user/login"
	"github.com/Tuanzi-bug/SyncHub/user/config"
	loginServiceV1 "github.com/Tuanzi-bug/SyncHub/user/pkg/service/login.service.v1"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"log"
	"net"
)

type Router interface {
	Register(r *gin.Engine)
}
type RegisterRouter struct {
}

func New() RegisterRouter {
	return RegisterRouter{}
}
func (RegisterRouter) Route(router Router, r *gin.Engine) {
	router.Register(r)
}

var routers []Router

func InitRouter(r *gin.Engine) {
	for _, ro := range routers {
		RegisterRouter{}.Route(ro, r)
	}
}

type gRPCConfig struct {
	Addr         string
	RegisterFunc func(*grpc.Server)
}

func RegisterGrpc() *grpc.Server {
	c := gRPCConfig{
		Addr: config.AppConfig.GrpcConfig.Addr,
		RegisterFunc: func(server *grpc.Server) {
			login.RegisterLoginServiceServer(server, loginServiceV1.New())
		},
	}
	server := grpc.NewServer()
	c.RegisterFunc(server)
	lis, err := net.Listen("tcp", config.AppConfig.GrpcConfig.Addr)
	if err != nil {
		log.Println("cannot listen")
	}
	go func() {
		log.Printf("grpc server started as: %s \n", c.Addr)
		err = server.Serve(lis)
		if err != nil {
			log.Println("server started error", err)
			return
		}
	}()
	return server
}

func RegisterEtcdServer() {
	// 1. 创建 服务发现解析器
	etcdRegister := discovery.NewResolver(
		config.AppConfig.EtcdConfig.Addrs, // etcd 服务器地址列表
		logs.LG,                           // 日志组件
	)
	// 2. 注册解析器到 gRPC
	resolver.Register(etcdRegister) // 让 gRPC 知道如何解析 etcd scheme 的服务地址

	// 3. 构造服务信息
	info := discovery.Server{
		Name:    config.AppConfig.GrpcConfig.Name,    // 服务名称
		Addr:    config.AppConfig.GrpcConfig.Addr,    // 服务地址（IP:Port）
		Version: config.AppConfig.GrpcConfig.Version, // 服务版本
		Weight:  config.AppConfig.GrpcConfig.Weight,  // 服务权重，用于负载均衡
	}

	// 4. 创建服务注册器
	r := discovery.NewRegister(
		config.AppConfig.EtcdConfig.Addrs, // etcd 服务器地址列表
		10,                                // TTL（生存时间），单位秒
		logs.LG,                           // 日志组件
	)

	// 5. 注册服务到 etcd
	_, err := r.Register(
		info, // 服务信息
		2,    // 续约时间间隔，单位秒
	)

	// 6. 错误处理
	if err != nil {
		log.Fatalln(err) // 如果注册失败，记录错误并终止程序
	}
}
