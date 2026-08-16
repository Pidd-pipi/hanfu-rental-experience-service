# BUG 复现说明（hanfu-rental__003）

## Bug 是什么
分页归一化、汉服筛选条件与仓库 offset 多处联合失效。

## 如何触发
```bash
go test ./internal/dto -run 'TestPageQueryNormalizeDefaults' -count=1
```

## 错误信息
```
--- FAIL: TestPageQueryNormalizeDefaults
    pagination_combination_test.go:9: Page = 0, want 1
```
