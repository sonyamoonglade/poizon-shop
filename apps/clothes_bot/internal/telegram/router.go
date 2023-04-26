package telegram

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"domain"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"logger"
	"utils/ranges"
)

var (
	ErrNoHandler = errors.New("handler not found")
	ErrNoRoute   = errors.New("no route")
)

type StateProvider interface {
	GetState(ctx context.Context, telegramID int64) (domain.State, error)
}

type RouteHandler interface {
	// Main menu
	Start(ctx context.Context, m *tg.Message) error
	Menu(ctx context.Context, chatID int64) error
	Catalog(ctx context.Context, chatID int64) error
	MyOrders(ctx context.Context, chatID int64) error

	AskForPromocode(ctx context.Context, chatID int64) error
	HandlePromocodeInput(ctx context.Context, m *tg.Message) error

	FAQ(ctx context.Context, chatID int64) error
	AnswerQuestion(chatID int64, n int) error

	AskForCalculatorOrderType(ctx context.Context, chatID int64) error
	HandleCalculatorOrderTypeInput(ctx context.Context, chatID int64, typ domain.OrderType) error
	HandleCalculatorCategoryInput(ctx context.Context, chatID int64, cat domain.Category) error
	AskForCalculatorCategory(ctx context.Context, chatID int64) error
	HandleCalculatorInput(ctx context.Context, m *tg.Message) error

	GetCart(ctx context.Context, chatID int64) error
	EditCart(ctx context.Context, chatID int64, cartPreviewMsgID int) error
	RemoveCartPosition(ctx context.Context, chatID int64, callbackData int, originalMsgID, cartPreviewMsgID int) error

	// Add position is like StartmakeOrderGuide but without instruction
	AddPosition(ctx context.Context, chatID int64) error

	HandleOrderTypeInput(ctx context.Context, chatID int64, typ domain.OrderType) error
	HandleCategoryInput(ctx context.Context, chatID int64, cat domain.Category) error

	// StartMakeOrderGuide is initial guide handler
	StartMakeOrderGuide(ctx context.Context, m *tg.Message) error

	MakeOrderGuideStep1(ctx context.Context, chatID int64, controlButtonsMessageID int, guideMsgIDs []int) error
	MakeOrderGuideStep2(ctx context.Context, chatID int64, controlButtonsMessageID int, guideMsgIDs []int) error
	MakeOrderGuideStep3(ctx context.Context, chatID int64, controlButtonsMessageID int, guideMsgIDs []int) error
	MakeOrderGuideStep4(ctx context.Context, chatID int64, controlButtonsMessageID int, guideMsgIDs []int) error
	MakeOrderGuideStep5(ctx context.Context, chatID int64, controlButtonsMessageID int, guideMsgIDs []int) error
	MakeOrderGuideStep6(ctx context.Context, chatID int64, controlButtonsMessageID int, guideMsgIDs []int) error

	AskForFIO(ctx context.Context, chatID int64) error
	// Use tg.Message because user types in and userID is user's
	HandleFIOInput(ctx context.Context, m *tg.Message) error
	HandlePhoneNumberInput(ctx context.Context, m *tg.Message) error
	HandleDeliveryAddressInput(ctx context.Context, m *tg.Message) error
	HandlePayment(ctx context.Context, shortOrderID string, c *tg.CallbackQuery) error

	HandleSizeInput(ctx context.Context, m *tg.Message) error
	// Use tg.CallbackQuery because callback is asosiated with c.User.ID, message is from bot
	HandleButtonSelect(ctx context.Context, c *tg.CallbackQuery, button domain.Button) error
	HandlePriceInput(ctx context.Context, m *tg.Message) error
	HandleLinkInput(ctx context.Context, m *tg.Message) error

	// Catalog manupulations
	HandleCatalogNext(ctx context.Context, chatID int64, controlButtonsMessageID int64, thumnailMsgIDs []int) error
	HandleCatalogPrev(ctx context.Context, chatID int64, controlButtonsMessageID int64, thumnailMsgIDs []int) error

	// Utils
	HandleError(ctx context.Context, err error, m tg.Update)
	AnswerCallback(callbackID string) error
}

type Router struct {
	h              RouteHandler
	updates        <-chan tg.Update
	shutdown       chan struct{}
	handlerTimeout time.Duration
	wg             *sync.WaitGroup
	stateProvider  StateProvider
}

func NewRouter(updates <-chan tg.Update, h RouteHandler, stateProvider StateProvider, timeout time.Duration) *Router {
	return &Router{
		h:              h,
		updates:        updates,
		handlerTimeout: timeout,
		stateProvider:  stateProvider,
		shutdown:       make(chan struct{}),
		wg:             new(sync.WaitGroup),
	}
}

func (r *Router) Bootstrap() {
	logger.Get().Info("router is listening for updates")
	for {
		select {
		case <-r.shutdown:
			logger.Get().Info("router is shutting down")
			return
		case update, ok := <-r.updates:
			if !ok {
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), r.handlerTimeout)
			r.wg.Add(1)
			go func() {
				defer func() {
					if panicMsg := recover(); panicMsg != nil {
						logger.Get().Error("panic in handler",
							zap.Any("msg", panicMsg),
							zap.ByteString("stacktrace", debug.Stack()))
					}
				}()
				if err := r.mapToHandler(ctx, update); err != nil {
					var (
						username   = domain.MakeUsername(update.FromChat().UserName)
						telegramID = update.FromChat().ID
					)

					logger.Get().Error("error in handler occurred",
						zap.String("from", username),
						zap.Int64("telegramId", telegramID),
						zap.Error(err))

					r.h.HandleError(ctx, err, update)
				}
				defer cancel()
				defer r.wg.Done()
			}()
		}
	}
}

func (r *Router) Shutdown() {
	close(r.shutdown)
	r.wg.Wait()
}

func (r *Router) mapToHandler(ctx context.Context, u tg.Update) error {
	switch {
	case u.Message != nil:
		return r.mapToCommandHandler(ctx, u.Message)
	case u.CallbackQuery != nil:
		return r.mapToCallbackHandler(ctx, u.CallbackQuery)
	default:
		return ErrNoRoute
	}
}

const (
	startCommand       = "/start"
	menuCommand        = "Меню"
	getCartCommand     = "Посмотреть корзину"
	addPositionCommand = "Добавить позицию"
)

func (r *Router) mapToCommandHandler(ctx context.Context, m *tg.Message) error {
	var (
		chatID = m.Chat.ID
		cmd    = r.command(m.Text)
	)
	// get state and route accordingly
	logger.Get().Debug("message info",
		zap.String("text", m.Text),
		zap.String("from", domain.MakeUsername(m.From.String())),
		zap.String("date", m.Time().Format(time.RFC822)))
	switch true {
	case cmd(startCommand):
		return r.h.Start(ctx, m)
	case cmd(menuCommand):
		return r.h.Menu(ctx, chatID)
	case cmd(getCartCommand):
		return r.h.GetCart(ctx, chatID)
	case cmd(addPositionCommand):
		return r.h.AddPosition(ctx, chatID)
	default:
		customerState, err := r.stateProvider.GetState(ctx, chatID)
		if err != nil {
			return err
		}
		switch customerState {
		case domain.StateWaitingForSize:
			return r.h.HandleSizeInput(ctx, m)
		case domain.StateWaitingForPrice:
			return r.h.HandlePriceInput(ctx, m)
		case domain.StateWaitingForLink:
			return r.h.HandleLinkInput(ctx, m)
		case domain.StateWaitingForCalculatorInput:
			return r.h.HandleCalculatorInput(ctx, m)
		case domain.StateWaitingForFIO:
			return r.h.HandleFIOInput(ctx, m)
		case domain.StateWaitingForPhoneNumber:
			return r.h.HandlePhoneNumberInput(ctx, m)
		case domain.StateWaitingForDeliveryAddress:
			return r.h.HandleDeliveryAddressInput(ctx, m)
		case domain.StateWaitingForPromocode:
			return r.h.HandlePromocodeInput(ctx, m)
		case domain.StateDefault:
			return ErrNoHandler
		default:
			return ErrNoHandler
		}
	}
}

func (r *Router) mapToCallbackHandler(ctx context.Context, c *tg.CallbackQuery) error {
	logger.Get().Debug("callback info",
		zap.String("data", c.Data),
		zap.String("from", domain.MakeUsername(c.From.String())),
		zap.String("date", c.Message.Time().Format(time.RFC822)))

	defer r.h.AnswerCallback(c.ID)

	var (
		chatID             = c.From.ID
		msgID              = c.Message.MessageID
		intCallbackData    int
		callbackDataMsgIDs []int
		stringData         string
	)

	out, callback, err := parseCallbackData(c.Data)
	if err != nil {
		return fmt.Errorf("parseCallbackData: %w", err)
	}

	intCallbackData = callback

	switch v := out.(type) {
	case []int:
		callbackDataMsgIDs = v
	case string:
		stringData = v
	}

	switch intCallbackData {
	case noopCallback:
		// Do not do anything
		return nil
	case menuMakeOrderCallback:
		return r.h.StartMakeOrderGuide(ctx, c.Message)
	case myCartCallback:
		return r.h.GetCart(ctx, chatID)
	case menuCalculatorCallback:
		return r.h.AskForCalculatorOrderType(ctx, chatID)
	case calculateMoreCallback:
		return r.h.AskForCalculatorOrderType(ctx, chatID)
	case promocodeCallback:
		return r.h.AskForPromocode(ctx, chatID)
	case orderGuideStep0Callback:
		return r.h.MakeOrderGuideStep1(ctx, chatID, msgID, callbackDataMsgIDs)
	case orderGuideStep1Callback:
		return r.h.MakeOrderGuideStep2(ctx, chatID, msgID, callbackDataMsgIDs)
	case orderGuideStep2Callback:
		return r.h.MakeOrderGuideStep3(ctx, chatID, msgID, callbackDataMsgIDs)
	case orderGuideStep3Callback:
		return r.h.MakeOrderGuideStep4(ctx, chatID, msgID, callbackDataMsgIDs)
	case orderGuideStep4Callback:
		return r.h.MakeOrderGuideStep5(ctx, chatID, msgID, callbackDataMsgIDs)
	case orderGuideStep5Callback:
		return r.h.MakeOrderGuideStep6(ctx, chatID, msgID, callbackDataMsgIDs)
	case makeOrderCallback:
		return r.h.AskForFIO(ctx, chatID)
	case menuCatalogCallback:
		return r.h.Catalog(ctx, chatID)
	case menuMyOrdersCallback:
		return r.h.MyOrders(ctx, chatID)
	case menuFaqCallback:
		return r.h.FAQ(ctx, chatID)
	case buttonTorqoiseSelectCallback:
		return r.h.HandleButtonSelect(ctx, c, domain.ButtonTorqoise)
	case buttonGreySelectCallback:
		return r.h.HandleButtonSelect(ctx, c, domain.ButtonGrey)
	case button95SelectCallback:
		return r.h.HandleButtonSelect(ctx, c, domain.Button95)
	case editCartCallback:
		return r.h.EditCart(ctx, chatID, msgID)
	case addPositionCallback:
		return r.h.AddPosition(ctx, chatID)
	case orderTypeNormalCallback:
		return r.h.HandleOrderTypeInput(ctx, chatID, domain.OrderTypeNormal)
	case orderTypeNormalCalculatorCallback:
		return r.h.HandleCalculatorOrderTypeInput(ctx, chatID, domain.OrderTypeNormal)
	case orderTypeExpressCallback:
		return r.h.HandleOrderTypeInput(ctx, chatID, domain.OrderTypeExpress)
	case orderTypeExpressCalculatorCallback:
		return r.h.HandleCalculatorOrderTypeInput(ctx, chatID, domain.OrderTypeExpress)
	case categoryLightCallback:
		return r.h.HandleCategoryInput(ctx, chatID, domain.CategoryLight)
	case categoryLightCalculatorCallback:
		return r.h.HandleCalculatorCategoryInput(ctx, chatID, domain.CategoryLight)
	case categoryHeavyCallback:
		return r.h.HandleCategoryInput(ctx, chatID, domain.CategoryHeavy)
	case categoryHeavyCalculatorCallback:
		return r.h.HandleCalculatorCategoryInput(ctx, chatID, domain.CategoryHeavy)
	case categoryOtherCallback:
		return r.h.HandleCategoryInput(ctx, chatID, domain.CategoryOther)
	case categoryOtherCalculatorCallback:
		return r.h.HandleCalculatorCategoryInput(ctx, chatID, domain.CategoryOther)
	case selectCategoryAgainCallback:
		return r.h.AskForCalculatorCategory(ctx, chatID)
	case paymentCallback:
		// stringData in this case is orderShortID
		return r.h.HandlePayment(ctx, stringData, c)
	default:
		// intCallback > edit
		// Remove position callback
		if ranges.IsBetween(intCallbackData, editCartRemovePositionOffset, catalogOffset) {
			// callbackDataMsgIDs[0] - id of preview cart message
			return r.h.RemoveCartPosition(ctx, chatID, intCallbackData, msgID, callbackDataMsgIDs[0])
		}

		// Prev or next callback
		if ranges.IsBetweenInc(intCallbackData, catalogOffset, faqOffset) {
			switch intCallbackData - catalogOffset {
			case catalogNextCallback:
				return r.h.HandleCatalogNext(ctx, chatID, int64(msgID), callbackDataMsgIDs)
			case catalogPrevCallback:
				return r.h.HandleCatalogPrev(ctx, chatID, int64(msgID), callbackDataMsgIDs)
			}
		}
		// todo: rm faqOffset + 1000
		if ranges.IsBetween(intCallbackData, faqOffset, faqOffset+1000) {
			n_question := intCallbackData - faqOffset
			return r.h.AnswerQuestion(chatID, n_question)
		}

		return ErrNoHandler
	}
}

func (r *Router) command(actual string) func(string) bool {
	return func(want string) bool {
		return actual == want
	}
}
