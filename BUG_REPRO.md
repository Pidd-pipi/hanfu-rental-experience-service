# BUG 复现说明（hanfu-rental__001）

## Bug 是什么
汉服目录生命周期多处联合失效：非法朝代可录入，朝代文案错误，库存更新被仓储层强制置 0，筛选条件被丢弃。

## 如何触发
```bash
go test ./internal/service -run 'TestHanfuLifecycleValidation|TestHanfuUpdateStock' -count=1
```

## 错误信息
```
--- FAIL: TestHanfuLifecycleValidation
    hanfu_lifecycle_test.go:17: invalid dynasty should be rejected
--- FAIL: TestHanfuUpdateStock
    hanfu_lifecycle_test.go:29: stock = 0, want 5
```
