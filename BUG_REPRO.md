# BUG 复现说明（hanfu-rental__005）

## Bug 是什么
租赁订单状态枚举、状态文案、会员卡年卡类型文案/月数/价格多处错位。

## 如何触发
```bash
go test ./internal/util -run 'TestRentalStatusAndCardBoundary'
```

## 错误信息
```
--- FAIL: TestRentalStatusAndCardBoundary
    rental_status_diagnosis_test.go:11: IsRentalOrderStatus(returned) should be true
```
