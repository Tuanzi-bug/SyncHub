package login_service_v1

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Tuanzi-bug/SyncHub/common"
	"github.com/Tuanzi-bug/SyncHub/common/encrypts"
	"github.com/Tuanzi-bug/SyncHub/common/errs"
	"github.com/Tuanzi-bug/SyncHub/common/jwts"
	"github.com/Tuanzi-bug/SyncHub/common/tms"
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
	"strings"
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
	memMsg.Code, _ = encrypts.EncryptInt64(mem.Id, model.AESKey)
	memMsg.LastLoginTime = tms.FormatByMill(mem.LastLoginTime)
	memMsg.CreateTime = tms.FormatByMill(mem.CreateTime)
	//2.根据用户id查组织
	orgs, err := h.organizationRepo.FindOrganizationByMemId(c, mem.Id)
	if err != nil {
		zap.L().Error("Login db FindMember error", zap.Error(err))
		return nil, errs.GrpcError(grpc_errs.DBError)
	}
	var orgsMessage []*login.OrganizationMessage
	err = copier.Copy(&orgsMessage, orgs)
	for _, v := range orgsMessage {
		v.Code, _ = encrypts.EncryptInt64(v.Id, model.AESKey)
		v.OwnerCode = memMsg.Code
		o := organization.ToMap(orgs)[v.Id]
		v.CreateTime = tms.FormatByMill(o.CreateTime)
	}
	if len(orgs) > 0 {
		memMsg.OrganizationCode, _ = encrypts.EncryptInt64(orgs[0].Id, model.AESKey)
	}
	//3.用jwt生成token
	memIdStr := strconv.FormatInt(mem.Id, 10)
	exp := time.Duration(config.AppConfig.JwtConfig.AccessExp*3600*24) * time.Second
	rExp := time.Duration(config.AppConfig.JwtConfig.RefreshExp*3600*24) * time.Second
	token := jwts.CreateToken(memIdStr, exp, config.AppConfig.JwtConfig.AccessSecret, rExp, config.AppConfig.JwtConfig.RefreshSecret, msg.Ip)
	tokenList := &login.TokenMessage{
		AccessToken:    token.AccessToken,
		RefreshToken:   token.RefreshToken,
		AccessTokenExp: token.AccessExp,
		TokenType:      "bearer",
	}
	// 优化：放入缓存 member orgs
	go func() {
		marshal, _ := json.Marshal(mem)
		h.cache.Put(c, model.Member+"::"+memIdStr, string(marshal), exp)
		orgsJson, _ := json.Marshal(orgs)
		h.cache.Put(c, model.MemberOrganization+"::"+memIdStr, string(orgsJson), exp)
	}()
	return &login.LoginResponse{
		Member:           memMsg,
		OrganizationList: orgsMessage,
		TokenList:        tokenList,
	}, nil
}

func (h *LoginService) TokenVerify(ctx context.Context, msg *login.LoginMessage) (*login.LoginResponse, error) {
	// 获取token
	token := msg.Token
	if strings.Contains(token, "bearer") {
		token = strings.ReplaceAll(token, "bearer ", "")
	}
	// 解析token
	parseToken, err := jwts.ParseToken(token, config.AppConfig.JwtConfig.AccessSecret, msg.Ip)
	if err != nil {
		zap.L().Error("Login  TokenVerify error", zap.Error(err))
		return nil, errs.GrpcError(grpc_errs.NoLogin)
	}
	//从缓存中查询 如果没有 直接返回认证失败
	memJson, err := h.cache.Get(ctx, model.Member+"::"+parseToken)
	if err != nil {
		zap.L().Error("TokenVerify cache get member error", zap.Error(err))
		return nil, errs.GrpcError(grpc_errs.NoLogin)
	}
	if memJson == "" {
		zap.L().Error("TokenVerify cache get member expire")
		return nil, errs.GrpcError(grpc_errs.NoLogin)
	}
	memberById := &member.Member{}
	json.Unmarshal([]byte(memJson), memberById)
	memMsg := &login.MemberMessage{}
	copier.Copy(memMsg, memberById)
	memMsg.Code, _ = encrypts.EncryptInt64(memberById.Id, model.AESKey)
	//数据库查询 优化点 登录之后 应该把用户信息缓存起来
	orgsJson, err := h.cache.Get(context.Background(), model.MemberOrganization+"::"+parseToken)
	if err != nil {
		zap.L().Error("TokenVerify cache get organization error", zap.Error(err))
		return nil, errs.GrpcError(grpc_errs.NoLogin)
	}
	if orgsJson == "" {
		zap.L().Error("TokenVerify cache get organization expire")
		return nil, errs.GrpcError(grpc_errs.NoLogin)
	}
	var orgs []*organization.Organization
	json.Unmarshal([]byte(orgsJson), &orgs)

	if len(orgs) > 0 {
		memMsg.OrganizationCode, _ = encrypts.EncryptInt64(orgs[0].Id, model.AESKey)
	}
	memMsg.CreateTime = tms.FormatByMill(memberById.CreateTime)
	return &login.LoginResponse{Member: memMsg}, nil
}

func (h *LoginService) MyOrgList(ctx context.Context, msg *login.UserMessage) (*login.OrgListResponse, error) {
	memId := msg.MemId
	// 根据用户id查组织
	orgs, err := h.organizationRepo.FindOrganizationByMemId(ctx, memId)
	if err != nil {
		zap.L().Error("MyOrgList FindOrganizationByMemId err", zap.Error(err))
		return nil, errs.GrpcError(grpc_errs.DBError)
	}
	var orgsMessage []*login.OrganizationMessage
	err = copier.Copy(&orgsMessage, orgs)
	for _, org := range orgsMessage {
		org.Code, _ = encrypts.EncryptInt64(org.Id, model.AESKey)
	}
	return &login.OrgListResponse{OrganizationList: orgsMessage}, nil
}

// FindMemInfoById 根据用户id查用户信息和组织信息
func (h *LoginService) FindMemInfoById(ctx context.Context, msg *login.UserMessage) (*login.MemberMessage, error) {
	// 根据用户id查用户信息
	memberById, err := h.memberRepo.FindMemberById(ctx, msg.MemId)
	if err != nil {
		zap.L().Error("TokenVerify db FindMemberById error", zap.Error(err))
		return nil, errs.GrpcError(grpc_errs.DBError)
	}
	memMsg := &login.MemberMessage{}
	copier.Copy(memMsg, memberById)
	memMsg.Code, _ = encrypts.EncryptInt64(memberById.Id, model.AESKey)
	// 根据用户id查组织
	orgs, err := h.organizationRepo.FindOrganizationByMemId(ctx, memberById.Id)
	if err != nil {
		zap.L().Error("TokenVerify db FindMember error", zap.Error(err))
		return nil, errs.GrpcError(grpc_errs.DBError)
	}
	if len(orgs) > 0 {
		memMsg.OrganizationCode, _ = encrypts.EncryptInt64(orgs[0].Id, model.AESKey)
	}
	memMsg.CreateTime = tms.FormatByMill(memberById.CreateTime)
	return memMsg, nil
}

func (h *LoginService) FindMemInfoByIds(ctx context.Context, msg *login.UserMessage) (*login.MemberMessageList, error) {
	// 根据用户id列表查用户信息
	memberList, err := h.memberRepo.FindMemberByIds(ctx, msg.MIds)
	if err != nil {
		zap.L().Error("FindMemInfoByIds db memberRepo.FindMemberByIds error", zap.Error(err))
		return nil, errs.GrpcError(grpc_errs.DBError)
	}
	if memberList == nil || len(memberList) <= 0 {
		return &login.MemberMessageList{List: nil}, nil
	}
	// todo: 脱了裤子放屁的感觉：处理返回信息
	mMap := make(map[int64]*member.Member)
	for _, v := range memberList {
		mMap[v.Id] = v
	}
	var memMsgs []*login.MemberMessage
	copier.Copy(&memMsgs, memberList)
	for _, v := range memMsgs {
		m := mMap[v.Id]
		v.CreateTime = tms.FormatByMill(m.CreateTime)
		v.Code = encrypts.EncryptNoErr(v.Id)
	}
	return &login.MemberMessageList{List: memMsgs}, nil
}
