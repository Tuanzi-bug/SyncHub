package grpc_errs

import (
	"github.com/Tuanzi-bug/SyncHub/common/errs"
)

var (
	//Success             common.BusinessCode = 200
	//LoginMobileNotLegal common.BusinessCode = 2001 //手机号不合法
	NoLegalMobile = errs.NewError(2001, "手机号不合法")
)
