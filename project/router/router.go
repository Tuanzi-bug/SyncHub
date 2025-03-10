package router

import (
	"github.com/Tuanzi-bug/SyncHub/common/discovery"
	"github.com/Tuanzi-bug/SyncHub/common/logs"
	"github.com/Tuanzi-bug/SyncHub/grpc/account"
	"github.com/Tuanzi-bug/SyncHub/grpc/auth"
	"github.com/Tuanzi-bug/SyncHub/grpc/department"
	"github.com/Tuanzi-bug/SyncHub/grpc/menu"
	"github.com/Tuanzi-bug/SyncHub/grpc/project"
	"github.com/Tuanzi-bug/SyncHub/grpc/task"
	"github.com/Tuanzi-bug/SyncHub/project/config"
	"github.com/Tuanzi-bug/SyncHub/project/internal/interceptor"
	"github.com/Tuanzi-bug/SyncHub/project/internal/rpc"
	account_service_v1 "github.com/Tuanzi-bug/SyncHub/project/pkg/service/account.service.v1"
	auth_service_v1 "github.com/Tuanzi-bug/SyncHub/project/pkg/service/auth.service.v1"
	department_service_v1 "github.com/Tuanzi-bug/SyncHub/project/pkg/service/department.service.v1"
	menu_service_v1 "github.com/Tuanzi-bug/SyncHub/project/pkg/service/menu.service.v1"
	project_service_v1 "github.com/Tuanzi-bug/SyncHub/project/pkg/service/project.service.v1"
	task_service_v1 "github.com/Tuanzi-bug/SyncHub/project/pkg/service/task.service.v1"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"log"
	"net"
)

// Router 接口
type Router interface {
	Route(r *gin.Engine)
}

type RegisterRouter struct {
}

func New() *RegisterRouter {
	return &RegisterRouter{}
}

func (*RegisterRouter) Route(ro Router, r *gin.Engine) {
	ro.Route(r)
}

var routers []Router

func InitRouter(r *gin.Engine) {
	//rg := New()
	//rg.Route(&user.RouterUser{}, r)
	for _, ro := range routers {
		ro.Route(r)
	}
}

func Register(ro ...Router) {
	routers = append(routers, ro...)
}

type gRPCConfig struct {
	Addr         string
	RegisterFunc func(*grpc.Server)
}

func RegisterGrpc() *grpc.Server {
	c := gRPCConfig{
		Addr: config.AppConfig.GrpcConfig.Addr,
		RegisterFunc: func(g *grpc.Server) {
			project.RegisterProjectServiceServer(g, project_service_v1.New())
			task.RegisterTaskServiceServer(g, task_service_v1.New())
			account.RegisterAccountServiceServer(g, account_service_v1.New())
			auth.RegisterAuthServiceServer(g, auth_service_v1.New())
			department.RegisterDepartmentServiceServer(g, department_service_v1.New())
			menu.RegisterMenuServiceServer(g, menu_service_v1.New())
		}}
	s := grpc.NewServer(interceptor.New().Cache())
	c.RegisterFunc(s)
	lis, err := net.Listen("tcp", c.Addr)
	if err != nil {
		log.Println("cannot listen")
	}
	go func() {
		log.Printf("grpc server started as: %s \n", c.Addr)
		err = s.Serve(lis)
		if err != nil {
			log.Println("server started error", err)
			return
		}
	}()
	return s
}

func RegisterEtcdServer() {
	etcdRegister := discovery.NewResolver(config.AppConfig.EtcdConfig.Addrs, logs.LG)
	resolver.Register(etcdRegister)

	info := discovery.Server{
		Name:    config.AppConfig.GrpcConfig.Name,
		Addr:    config.AppConfig.GrpcConfig.Addr,
		Version: config.AppConfig.GrpcConfig.Version,
		Weight:  config.AppConfig.GrpcConfig.Weight,
	}
	r := discovery.NewRegister(config.AppConfig.EtcdConfig.Addrs, 10, logs.LG)
	_, err := r.Register(info, 2)
	if err != nil {
		log.Fatalln(err)
	}
}

func InitUserRpc() {
	rpc.InitRpcUserClient()
}
