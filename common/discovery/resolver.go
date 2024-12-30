package discovery

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
	"time"
)

const (
	schema = "etcd"
	// 定义属性 key
	weightKey = "weight"
)

type Resolver struct {
	schema       string              // 解析器的协议模式，这里固定为"etcd"
	EtcdAddrs    []string            // etcd服务器地址列表
	DialTimeout  int                 // 连接超时时间
	closeCh      chan struct{}       // 关闭通道，用于优雅退出
	watchCh      clientv3.WatchChan  // etcd监听通道
	cli          *clientv3.Client    // etcd客户端
	keyPrefix    string              // 服务注册的key前缀
	srvAddrsList []resolver.Address  // 服务地址列表
	cc           resolver.ClientConn // gRPC客户端连接
	logger       *zap.Logger         // 日志组件
}

func NewResolver(etcdAddrs []string, logger *zap.Logger) *Resolver {
	// 创建解析器实例，设置基础配置
	return &Resolver{
		schema:      schema,
		EtcdAddrs:   etcdAddrs, // 设置etcd地址
		DialTimeout: 3,         // 默认3秒超时
		logger:      logger,    // 设置日志组件
	}
}

// Scheme 返回解析器支持的协议名
func (r *Resolver) Scheme() string {
	return r.schema
}

// Build 为给定目标创建新的解析器
func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r.cc = cc // 保存gRPC客户端连接
	// target.Endpoint 是服务标识，可能是 "serviceName/version" 的格式
	service := target.URL.Path
	if service[0] == '/' {
		service = service[1:]
	}
	// 构建服务key前缀，用于etcd查询
	r.keyPrefix = BuildPrefix(Server{
		Name:    service,
		Version: target.URL.Host}) // 如果版本信息存在的话
	// 启动解析器
	if _, err := r.start(); err != nil {
		return nil, err
	}
	return r, nil
}

// ResolveNow 是resolver.Resolver接口要求的方法，这里不需要实现
func (r *Resolver) ResolveNow(o resolver.ResolveNowOptions) {}

// Close 用于关闭解析器
func (r *Resolver) Close() {
	r.closeCh <- struct{}{} // 发送关闭信号
}

func (r *Resolver) start() (chan<- struct{}, error) {
	var err error
	// 创建etcd客户端
	r.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   r.EtcdAddrs,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	})
	if err != nil {
		return nil, err
	}

	// 注册解析器
	resolver.Register(r)
	// 创建关闭通道
	r.closeCh = make(chan struct{})

	// 同步获取现有可用服务列表
	if err = r.sync(); err != nil {
		return nil, err
	}

	// 启动监听协程
	go r.watch()
	return r.closeCh, nil
}

func (r *Resolver) watch() {
	// 创建定时器，每分钟同步一次
	ticker := time.NewTicker(time.Minute)
	// 监听key前缀的变更
	r.watchCh = r.cli.Watch(context.Background(), r.keyPrefix, clientv3.WithPrefix())

	for {
		select {
		case <-r.closeCh:
			return
		case res, ok := <-r.watchCh:
			// 处理监听到的变更事件
			if ok {
				r.update(res.Events)
			}
		case <-ticker.C:
			// 定时同步
			if err := r.sync(); err != nil {
				r.logger.Error("sync failed", zap.Error(err))
			}
		}
	}
}

func (r *Resolver) update(events []*clientv3.Event) {
	// 遍历所有事件
	for _, ev := range events {
		var info Server
		var err error
		switch ev.Type {
		case mvccpb.PUT: // 新增或更新事件
			// 解析事件值获取服务信息
			info, err = ParseValue(ev.Kv.Value)
			if err != nil {
				continue
			}
			// 构造服务地址，包含权重信息
			addr := resolver.Address{Addr: info.Addr, Attributes: attributes.New(weightKey, info.Weight)} // 使用 Attributes 替代}
			// 如果地址不存在则添加到列表
			if !Exist(r.srvAddrsList, addr) {
				r.srvAddrsList = append(r.srvAddrsList, addr)
				// 更新gRPC客户端的服务列表状态
				r.cc.UpdateState(resolver.State{Addresses: r.srvAddrsList})
			}

		case mvccpb.DELETE: // 删除事件
			// 从key中解析服务信息
			info, err = SplitPath(string(ev.Kv.Key))
			if err != nil {
				continue
			}
			addr := resolver.Address{Addr: info.Addr}
			// 从列表中移除地址
			if s, ok := Remove(r.srvAddrsList, addr); ok {
				r.srvAddrsList = s
				// 更新gRPC客户端的服务列表状态
				r.cc.UpdateState(resolver.State{Addresses: r.srvAddrsList})
			}
		}
	}
}

func (r *Resolver) sync() error {
	// 创建3秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 查询指定前缀的所有key
	res, err := r.cli.Get(ctx, r.keyPrefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	// 清空当前服务地址列表
	r.srvAddrsList = []resolver.Address{}

	// 遍历结果,解析每个服务信息
	for _, v := range res.Kvs {
		info, err := ParseValue(v.Value)
		if err != nil {
			continue
		}
		// 构造服务地址并添加到列表
		addr := resolver.Address{Addr: info.Addr, Metadata: info.Weight}
		r.srvAddrsList = append(r.srvAddrsList, addr)
	}

	// 更新gRPC客户端的服务列表状态
	r.cc.UpdateState(resolver.State{Addresses: r.srvAddrsList})
	return nil
}
