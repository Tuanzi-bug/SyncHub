package discovery

import (
	"context"
	"encoding/json"
	"errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Register for grpc server
type Register struct {
	EtcdAddrs   []string
	DialTimeout int

	closeCh     chan struct{}
	leasesID    clientv3.LeaseID
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse

	srvInfo Server
	srvTTL  int64
	cli     *clientv3.Client
	logger  *zap.Logger
}

func NewRegister(etcdAddrs []string, dialTimeout int, logger *zap.Logger) *Register {
	return &Register{
		EtcdAddrs:   etcdAddrs,
		DialTimeout: dialTimeout,
		logger:      logger,
	}
}

func (r *Register) Register(srvInfo Server, ttl int64) (chan<- struct{}, error) {
	var err error
	if strings.Split(srvInfo.Addr, ":")[0] == "" {
		return nil, errors.New("invalid ip")
	}

	// 连接etcd，创建etcd client
	r.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   r.EtcdAddrs,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	})
	if err != nil {
		return nil, err
	}
	// 3. 保存服务信息和TTL
	r.srvInfo = srvInfo
	r.srvTTL = ttl

	// 开始注册服务
	if err = r.register(); err != nil {
		return nil, err
	}
	// 返回关闭通道，用于关闭服务注册
	r.closeCh = make(chan struct{})

	// 监听续约情况
	go r.keepAlive()

	return r.closeCh, nil
}

// Stop 关闭服务注册
func (r *Register) Stop() {
	r.closeCh <- struct{}{}
}

// register 注册节点
func (r *Register) register() error {
	leaseCtx, cancel := context.WithTimeout(context.Background(), time.Duration(r.DialTimeout)*time.Second)
	defer cancel()

	// 创建租约
	leaseResp, err := r.cli.Grant(leaseCtx, r.srvTTL)
	if err != nil {
		return err
	}

	r.leasesID = leaseResp.ID
	// 设置租约续约
	if r.keepAliveCh, err = r.cli.KeepAlive(context.Background(), leaseResp.ID); err != nil {
		return err
	}

	data, err := json.Marshal(r.srvInfo)
	if err != nil {
		return err
	}
	// 注册节点服务
	_, err = r.cli.Put(context.Background(), BuildRegPath(r.srvInfo), string(data), clientv3.WithLease(r.leasesID))
	return err
}

// unregister 删除节点
func (r *Register) unregister() error {
	_, err := r.cli.Delete(context.Background(), BuildRegPath(r.srvInfo))
	return err
}

// keepAlive 保持租约
func (r *Register) keepAlive() {
	// 设置续约时间
	ticker := time.NewTicker(time.Duration(r.srvTTL) * time.Second)
	for {
		select {
		case <-r.closeCh:
			// 关闭续约
			if err := r.unregister(); err != nil {
				r.logger.Error("unregister failed", zap.Error(err))
			}
			// 删除租约
			if _, err := r.cli.Revoke(context.Background(), r.leasesID); err != nil {
				r.logger.Error("revoke failed", zap.Error(err))
			}
			return
		case res := <-r.keepAliveCh:
			// 处理保活响应
			if res == nil { // 如果收到空响应，说明连接可能断开
				if err := r.register(); err != nil {
					r.logger.Error("register keepAliveCh failed", zap.Error(err))
				}
			}
		case <-ticker.C:
			// 定时检查
			if r.keepAliveCh == nil {
				if err := r.register(); err != nil {
					r.logger.Error("register ticker failed", zap.Error(err))
				}
			}
		}
	}
}

func (r *Register) UpdateHandler() http.HandlerFunc {
	/*
		允许在服务运行时动态调整服务权重，这在负载均衡场景下特别有用：
		可以动态调整服务流量分配
		支持服务平滑上下线
		实现动态负载均衡调整
	*/
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// 1. 获取并解析权重参数
		wi := req.URL.Query().Get("weight")
		weight, err := strconv.Atoi(wi)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// 2. 定义更新函数
		var update = func() error {
			// 更新服务信息中的权重值
			r.srvInfo.Weight = int64(weight)
			// 将服务信息序列化为JSON
			data, err := json.Marshal(r.srvInfo)
			if err != nil {
				return err
			}
			// 将更新后的数据写入etcd
			_, err = r.cli.Put(context.Background(), BuildRegPath(r.srvInfo),
				string(data), clientv3.WithLease(r.leasesID))
			return err
		}
		// 3. 执行更新操作
		if err := update(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		// 4. 返回成功信息
		w.Write([]byte("update server weight success"))
	})
}

func (r *Register) GetServerInfo() (Server, error) {
	/*
		使用场景：
		查看服务当前状态
		验证服务注册是否成功
		获取服务最新配置信息
	*/
	// 1. 从etcd获取服务信息
	resp, err := r.cli.Get(context.Background(), BuildRegPath(r.srvInfo))
	if err != nil {
		return r.srvInfo, err
	}

	// 2. 创建一个空的Server对象
	info := Server{}

	// 3. 如果找到了服务信息
	if resp.Count >= 1 {
		// 将etcd中的数据解析到Server结构体
		if err := json.Unmarshal(resp.Kvs[0].Value, &info); err != nil {
			return info, err
		}
	}

	return info, nil
}
