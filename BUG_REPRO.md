# BUG 复现说明（hanfu-rental__004）

## Bug 是什么
租金计算多处联合失效：租期天数少算，月卡/年卡折扣反了，押金计算错误。

## 如何触发
```bash
go test ./internal/util -run 'TestCalcRentalDays|TestCalcRentalFee|TestCalcDeposit' -count=1
```

## 错误信息
```
--- FAIL: TestCalcRentalDays/two_days
    rental_fee_calculator_test.go:22: CalcRentalDays = 1, want 2
--- FAIL: TestCalcRentalFee/month_card_15%
    rental_fee_calculator_test.go:45: total = 150.000000, want 255.000000
--- FAIL: TestCalcDeposit
    rental_fee_calculator_test.go:56: CalcDeposit(300) = 300.000000, want 600
```
