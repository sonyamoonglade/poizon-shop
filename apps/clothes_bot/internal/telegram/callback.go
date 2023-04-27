package telegram

// DO NOT CHANGE ORDER
// LOGIC DEMANDS ON IOTA
// todo: change from iota
const (
	noopCallback = iota
	menuCatalogCallback
	menuFaqCallback
	menuMyOrdersCallback
	menuCalculatorCallback
	myCartCallback
	promocodeCallback
	calculateMoreCallback
	menuMakeOrderCallback
	orderGuideStep0Callback
	orderGuideStep1Callback
	orderGuideStep2Callback
	orderGuideStep3Callback
	orderGuideStep4Callback
	orderGuideStep5Callback
	makeOrderCallback
	buttonTorqoiseSelectCallback
	buttonGreySelectCallback
	button95SelectCallback
	addPositionCallback
	editCartCallback
	orderTypeNormalCallback
	orderTypeNormalCalculatorCallback
	orderTypeExpressCallback
	orderTypeExpressCalculatorCallback
	categoryLightCallback
	categoryLightCalculatorCallback
	categoryHeavyCallback
	categoryHeavyCalculatorCallback
	categoryOtherCallback
	categoryOtherCalculatorCallback
	selectCategoryAgainCallback

	paymentCallback
)

const (
	editCartRemovePositionOffset = 1000
	catalogOffset                = 1200
	faqOffset                    = 1400
)

const (
	catalogPrevCallback = iota + 1
	catalogNextCallback
)
