package login_service_v1

import (
	"context"
	"errors"
	"github.com/Tuanzi-bug/SyncHub/common"
	"github.com/Tuanzi-bug/SyncHub/common/encrypts"
	"github.com/Tuanzi-bug/SyncHub/common/errs"
	"github.com/Tuanzi-bug/SyncHub/common/jwts"
	"github.com/Tuanzi-bug/SyncHub/grpc/user/login"
	"github.com/Tuanzi-bug/SyncHub/user/config"
	"github.com/Tuanzi-bug/SyncHub/user/internal/dao"
	"github.com/Tuanzi-bug/SyncHub/user/internal/database"
	"github.com/Tuanzi-bug/SyncHub/user/internal/database/tran"
	"github.com/Tuanzi-bug/SyncHub/user/internal/domain/member"
	"github.com/Tuanzi-bug/SyncHub/user/internal/domain/organization"
	"github.com/Tuanzi-bug/SyncHub/user/internal/repo"
	"github.com/Tuanzi-bug/SyncHub/user/pkg/grpc_errs"
	"github.com/Tuanzi-bug/SyncHub/user/pkg/model"
	"github.com/jinzhu/copier"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"log"
	"strconv"
	"time"
)

type LoginService struct {
	login.UnimplementedLoginServiceServer
	cache            repo.Cache
	memberRepo       repo.MemberRepo
	organizationRepo repo.OrganizationRepo
	transaction      tran.Transaction
}

func New() *LoginService {
	return &LoginService{
		cache:            dao.Rc,
		memberRepo:       dao.NewMemberDao(),
		organizationRepo: dao.NewOrganizationDao(),
		transaction:      dao.NewTransaction(),
	}
}

func (h *LoginService) GetCaptcha(ctx context.Context, msg *login.CaptchaMessage) (*login.CaptchaResponse, error) {
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
		err := h.cache.Put(c, "REGISTER_"+mobile, code, 5*time.Minute)
		if err != nil {
			log.Println("验证码存储失败")
		}
	}()
	return &login.CaptchaResponse{Code: code}, nil
}

func (h *LoginService) Register(ctx context.Context, msg *login.RegisterMessage) (*login.RegisterResponse, error) {
	c := context.Background()
	//1.可以校验参数
	redisCode, err := h.cache.Get(c, model.RegisterRedisKey+msg.Mobile)
	if errors.Is(err, redis.Nil) {
		return nil, errs.GrpcError(grpc_errs.CaptchaNotExist)
	}
	if err != nil {
		zap.L().Error("Register redis get error", zap.Error(err))
		return nil, errs.GrpcError(grpc_errs.RedisError)
	}
	//2.校验验证码
	if redisCode != msg.Captcha {
		return nil, errs.GrpcError(grpc_errs.CaptchaError)
	}
	//3.校验业务逻辑（邮箱是否被注册 账号是否被注册 手机号是否被注册）
	exist, err := h.memberRepo.GetMemberByEmail(c, msg.Email)
	if err != nil {
		zap.L().Error("Register db get error", zap.Error(err))
		return nil, errs.GrpcError(grpc_errs.DBError)
	}
	if exist {
		return nil, errs.GrpcError(grpc_errs.EmailExist)
	}
	exist, err = h.memberRepo.GetMemberByAccount(c, msg.Name)
	if err != nil {
		zap.L().Error("Register db get error", zap.Error(err))
		return nil, errs.GrpcError(grpc_errs.DBError)
	}
	if exist {
		return nil, errs.GrpcError(grpc_errs.AccountExist)
	}
	exist, err = h.memberRepo.GetMemberByMobile(c, msg.Mobile)
	if err != nil {
		zap.L().Error("Register db get error", zap.Error(err))
		return nil, errs.GrpcError(grpc_errs.DBError)
	}
	if exist {
		return nil, errs.GrpcError(grpc_errs.MobileExist)
	}
	//4.执行业务 将数据存入member表 生成一个数据 存入组织表 organization
	pwd := encrypts.Md5(msg.Password)
	mem := &member.Member{
		Account:       msg.Name,
		Password:      pwd,
		Name:          msg.Name,
		Mobile:        msg.Mobile,
		Email:         msg.Email,
		CreateTime:    time.Now().UnixMilli(),
		LastLoginTime: time.Now().UnixMilli(),
		Status:        model.Normal,
	}
	err = h.transaction.Action(func(conn database.DbConn) error {
		err = h.memberRepo.SaveMember(conn, c, mem)
		if err != nil {
			zap.L().Error("Register db SaveMember error", zap.Error(err))
			return errs.GrpcError(grpc_errs.DBError)
		}
		//存入组织
		org := &organization.Organization{
			Name:       mem.Name + "个人组织",
			MemberId:   mem.Id,
			CreateTime: time.Now().UnixMilli(),
			Personal:   model.Personal,
			Avatar:     "https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fc-ssl.dtstatic.com%2Fuploads%2Fblog%2F202103%2F31%2F20210331160001_9a852.thumb.1000_0.jpg&refer=http%3A%2F%2Fc-ssl.dtstatic.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1673017724&t=ced22fc74624e6940fd6a89a21d30cc5",
		}
		err = h.organizationRepo.SaveOrganization(conn, c, org)
		if err != nil {
			zap.L().Error("register SaveOrganization db err", zap.Error(err))
			return errs.GrpcError(grpc_errs.DBError)
		}
		return nil
	})

	return &login.RegisterResponse{}, err
}

func (h *LoginService) Login(ctx context.Context, msg *login.LoginMessage) (*login.LoginResponse, error) {
	c := context.Background()
	//1.去数据库查询 账号密码是否正确
	pwd := encrypts.Md5(msg.Password)
	mem, err := h.memberRepo.FindMember(c, msg.Account, pwd)
	if err != nil {
		zap.L().Error("Login db FindMember error", zap.Error(err))
		return nil, errs.GrpcError(grpc_errs.DBError)
	}
	if mem == nil {
		return nil, errs.GrpcError(grpc_errs.AccountAndPwdError)
	}
	memMsg := &login.MemberMessage{}
	err = copier.Copy(memMsg, mem)
	//2.根据用户id查组织
	orgs, err := h.organizationRepo.FindOrganizationByMemId(c, mem.Id)
	if err != nil {
		zap.L().Error("Login db FindMember error", zap.Error(err))
		return nil, errs.GrpcError(grpc_errs.DBError)
	}
	var orgsMessage []*login.OrganizationMessage
	err = copier.Copy(&orgsMessage, orgs)
	//3.用jwt生成token
	memIdStr := strconv.FormatInt(mem.Id, 10)
	exp := time.Duration(config.AppConfig.JwtConfig.AccessExp*3600*24) * time.Second
	rExp := time.Duration(config.AppConfig.JwtConfig.RefreshExp*3600*24) * time.Second
	token := jwts.CreateToken(memIdStr, exp, config.AppConfig.JwtConfig.AccessSecret, rExp, config.AppConfig.JwtConfig.RefreshSecret)
	tokenList := &login.TokenMessage{
		AccessToken:    token.AccessToken,
		RefreshToken:   token.RefreshToken,
		AccessTokenExp: token.AccessExp,
		TokenType:      "bearer",
	}
	return &login.LoginResponse{
		Member:           memMsg,
		OrganizationList: orgsMessage,
		TokenList:        tokenList,
	}, nil
}
