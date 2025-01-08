package project

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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
)

var ProjectServiceClient project.ProjectServiceClient
var TaskServiceClient task.TaskServiceClient
var AccountServiceClient account.AccountServiceClient
var DepartmentServiceClient department.DepartmentServiceClient
var AuthServiceClient auth.AuthServiceClient
var MenuServiceClient menu.MenuServiceClient

func InitRpcProjectClient() {
	etcdRegister := discovery.NewResolver(config.AppConfig.EtcdConfig.Addrs, logs.LG)
	resolver.Register(etcdRegister)

	conn, err := grpc.NewClient(etcdRegister.Scheme()+":///project_proto", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	ProjectServiceClient = project.NewProjectServiceClient(conn)
	TaskServiceClient = task.NewTaskServiceClient(conn)
	AccountServiceClient = account.NewAccountServiceClient(conn)
	DepartmentServiceClient = department.NewDepartmentServiceClient(conn)
	AuthServiceClient = auth.NewAuthServiceClient(conn)
	MenuServiceClient = menu.NewMenuServiceClient(conn)
}
