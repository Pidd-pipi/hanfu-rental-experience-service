# BUG 复现说明（hanfu-rental__002）

## Bug 是什么
用户查询/登录/押金流多处联合失效：未找到用户返回 (nil,nil)，登录 nil 用户 panic，押金充值/退款方向颠倒。

## 如何触发
```bash
go test ./internal/service -run 'TestGetProfileNilReturnsError|TestLoginNilReturnsUnauthorized|TestDepositFlow' -count=1
```

## 错误信息
```
--- FAIL: TestGetProfileNilReturnsError
    user_combination_test.go:34: expected error for nil user, got user=<nil>
panic: runtime error: invalid memory address or nil pointer dereference
```
