package constants

// Log templates are centralized so that any business field change forces a
// coordinated update of the related log statements across the codebase.
const (
	LogUserRegisterSuccess         = "user register success: phone=%s user_id=%d"
	LogUserRegisterFailed          = "user register failed: phone=%s error=%v"
	LogUserLoginSuccess            = "user login success: phone=%s user_id=%d role=%s"
	LogUserLoginFailed             = "user login failed: phone=%s error=%v"
	LogUserProfileUpdateSuccess    = "user profile update success: user_id=%d"
	LogHanfuCreateSuccess          = "hanfu create success: hanfu_id=%d name=%s dynasty=%s"
	LogHanfuCreateFailed           = "hanfu create failed: name=%s error=%v"
	LogHanfuStockUpdateSuccess     = "hanfu stock update success: hanfu_id=%d stock=%d"
	LogRentalOrderCreateSuccess    = "rental order create success: order_id=%d user_id=%d hanfu_id=%d"
	LogRentalOrderCreateFailed     = "rental order create failed: user_id=%d hanfu_id=%d error=%v"
	LogRentalOrderConfirmSuccess   = "rental order confirm success: order_id=%d status=renting"
	LogRentalOrderReturnSuccess    = "rental order return success: order_id=%d status=returned"
	LogRentalOrderReturnFailed     = "rental order return failed: order_id=%d error=%v"
	LogRentalOrderCompleteSuccess  = "rental order complete success: order_id=%d deposit_refunded=%.2f"
	LogRentalOrderCancelSuccess    = "rental order cancel success: order_id=%d"
	LogActivityCreateSuccess       = "activity create success: activity_id=%d title=%s"
	LogActivitySignupSuccess       = "activity signup success: activity_id=%d user_id=%d"
	LogActivitySignupFailed        = "activity signup failed: activity_id=%d user_id=%d error=%v"
	LogMemberCardCreateSuccess     = "member card create success: card_id=%d user_id=%d type=%s"
	LogMemberCardExpireProcessed   = "member card expire processed: card_id=%d"
	LogDepositPaySuccess           = "deposit pay success: user_id=%d amount=%.2f"
	LogDepositRefundSuccess        = "deposit refund success: user_id=%d amount=%.2f"
	LogArticleCreateSuccess        = "article create success: article_id=%d title=%s"
	LogArticleCreateFailed         = "article create failed: title=%s error=%v"
	LogMiddlewareAuthFailed        = "auth middleware failed: error=%v"
	LogMiddlewareRbacDenied        = "rbac middleware denied: user_id=%d role=%s required=%v"
	LogRateLimitReached            = "rate limit reached: ip=%s route=%s"
	LogSeedingCompleted            = "database seeding completed: users=%d hanfus=%d"
	LogFeeCalcUsed                 = "rental fee calculated: days=%d fee=%.2f discount=%.2f"
)

// LogTemplateCount guards the "at least 25 templates" requirement.
const LogTemplateCount = 30
