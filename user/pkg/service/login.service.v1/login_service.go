package login_service_v1

import (
	"context"
	"github.com/Tuanzi-bug/SyncHub/common"
	"github.com/Tuanzi-bug/SyncHub/common/errs"
	"github.com/Tuanzi-bug/SyncHub/user/pkg/dao"
	"github.com/Tuanzi-bug/SyncHub/user/pkg/grpc_errs"
	"github.com/Tuanzi-bug/SyncHub/user/pkg/repo"
	"log"
	"time"
)

type LoginService struct {
	UnimplementedLoginServiceServer
	Cache repo.Cache
}

func New() *LoginService {
	return &LoginService{
		Cache: dao.Rc,
	}
}

func (h *LoginService) GetCaptcha(ctx context.Context, msg *CaptchaMessage) (*CaptchaResponse, error) {
	//1. 获取参数
	mobile := msg.Mobile
	//2. 验证手机合法性
	if !common.VerifyMobile(mobile) {
		return nil, errs.GrpcError(grpc_errs.NoLegalMobile)
	}
	//3.生成验证码
	code := "123456"
	//4. 发送验证码
	go func() {
		time.Sleep(2 * time.Second)
		log.Println("调用短信平台发送短信")
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := h.Cache.Put(c, "REGISTER_"+mobile, code, 5*time.Minute)
		if err != nil {
			log.Println("验证码存储失败")
		}
	}()
	return &CaptchaResponse{Code: code}, nil
}
