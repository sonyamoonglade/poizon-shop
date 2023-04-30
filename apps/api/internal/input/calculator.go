package input

type CallbackCalculatorQueryInput struct {
	Category    string `query:"category"`
	Subcategory string `query:"subcategory"`
	ProductName string `query:"productName"`
}
