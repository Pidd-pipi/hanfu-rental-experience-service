package constants

// Messages centralizes user-facing prompts, backend responses and log wording.
const (
	MsgOK                  = "ok"
	MsgValidationFailed    = "参数校验失败"
	MsgUnauthorized        = "未登录或登录已过期"
	MsgForbidden           = "没有操作权限"
	MsgNotFound            = "资源不存在"
	MsgConflict            = "资源状态冲突"
	MsgRateLimited         = "请求过于频繁，请稍后再试"
	MsgInternalError       = "服务器内部错误"
	MsgPhoneAlreadyUsed    = "该手机号已被注册"
	MsgPhoneOrPassword     = "手机号或密码错误"
	MsgHanfuUnavailable    = "该汉服当前不可租赁"
	MsgHanfuNotEnoughStock = "该汉服库存不足"
	MsgOrderStatusInvalid  = "当前订单状态不可操作"
	MsgNotOwner            = "仅订单本人可操作"
	MsgActivityClosed      = "报名已截止或人数已满"
	MsgAlreadySignedUp     = "您已报名该活动"
	MsgCardExists          = "您已持有有效会员卡"
	MsgDepositPaid         = "押金已缴纳"
	MsgDepositRefunded     = "押金已退还"
)
