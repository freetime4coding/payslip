package tax

type TaxBracket struct {
    Limit int
    Rate  float64
}

// Range upper limit salary and tax rate
var brackets = []TaxBracket{
    {20000, 0.0},       // 1st index   : 0%
    {40000, 0.10},      // 2nd index   : 10%
    {80000, 0.20},      // 3rd index   : 20%
    {180000, 0.30},     // 4th index   : 30%
    {999999999, 0.40},  // final index : 40%
}

func GenerateMonthlyPayslip(salary int, index int) float64 {
    // Condition after final index tax OR the salary hit 0 or less
    if index >= len(brackets) || salary <= 0 {
        return 0
    }

    var lower int
    if index == 0 {
        lower = 0
    } else {
        // Set amount to reduced next rate tax
        lower = brackets[index-1].Limit
    }

    // Taxable amount in this bracket
    taxable := salary - lower
    if taxable > brackets[index].Limit - lower {
        taxable = brackets[index].Limit - lower
    }
    if taxable < 0 {
        taxable = 0
    }

    // Tax for current bracket + tax for next brackets
    return float64(taxable)*brackets[index].Rate +
        GenerateMonthlyPayslip(salary, index+1)
}