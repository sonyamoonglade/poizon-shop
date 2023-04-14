package router

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"domain"

	"household_bot/internal/telegram/callback"
	"household_bot/internal/telegram/tg_errors"
	"logger"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

var (
	ErrNoHandler = errors.New("no handler was found")
)

type RouteHandler interface {
	Start(ctx context.Context, message *tg.Message) error
	Menu(ctx context.Context, chatID int64) error
	Catalog(ctx context.Context, chatID int64, prevMsgID *int) error

	GetCart(ctx context.Context, chatID int64) error
	EditCart(ctx context.Context, chatID int64, cartMsgID int) error
	DeletePositionFromCart(ctx context.Context, chatID int64, buttonsMsgID int, args []string) error

	Categories(ctx context.Context, chatID int64, prevMsgID int, onlyAvailableInStock bool) error
	Subcategories(ctx context.Context, chatID int64, prevMsgID int, args []string) error
	ProductsNew(ctx context.Context, chatID int64, msgIDForDeletion int, args []string) error
	Products(ctx context.Context, chatID int64, prevMsgID int, args []string) error
	ProductCard(ctx context.Context, chatID int64, prevMsgID int, args []string) error
	AddToCart(ctx context.Context, chatID int64, args []string) error

	AskForOrderType(ctx context.Context, chatID int64) error
	HandleOrderTypeInput(ctx context.Context, chatID int64, args []string) error
	HandleFIOInput(ctx context.Context, m *tg.Message) error
	HandlePhoneNumberInput(ctx context.Context, m *tg.Message) error
	HandleDeliveryAddressInput(ctx context.Context, m *tg.Message) error
	HandlePayment(ctx context.Context, c *tg.CallbackQuery, args []string) error

	AnswerCallback(c *tg.CallbackQuery) error
	Sorry(chatID int64) error
}

type StateProvider interface {
	GetState(ctx context.Context, telegramID int64) (domain.State, error)
}

type Router struct {
	updates       <-chan tg.Update
	handler       RouteHandler
	stateProvider StateProvider

	wg       *sync.WaitGroup
	timeout  time.Duration
	shutdown chan struct{}
}

func NewRouter(updates <-chan tg.Update, h RouteHandler, s StateProvider, timeout time.Duration) *Router {
	return &Router{
		updates:       updates,
		handler:       h,
		wg:            new(sync.WaitGroup),
		timeout:       timeout,
		shutdown:      make(chan struct{}),
		stateProvider: s,
	}
}

func (r *Router) Bootstrap() error {
	logger.Get().Info("router is listening for updates")
	for {
		select {
		case <-r.shutdown:
			logger.Get().Info("router is shutting down")
			return nil
		case update, ok := <-r.updates:
			if !ok {
				return nil
			}

			ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
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
					r.handleError(ctx, err, update)
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
		return ErrNoHandler
	}
}

const (
	Start = "/start"
	Menu  = "Меню"
)

func (r *Router) mapToCommandHandler(ctx context.Context, m *tg.Message) error {
	var (
		chatID = m.Chat.ID
		cmd    = r.command(m.Text)
	)
	logger.Get().Debug("message info",
		zap.String("text", m.Text),
		zap.String("from", r.getUsername(m.From)),
		zap.String("date", m.Time().Format(time.RFC822)))

	switch {
	case cmd(Start):
		return r.handler.Start(ctx, m)
	case cmd(Menu):
		return r.handler.Menu(ctx, chatID)
	default:
		state, err := r.stateProvider.GetState(ctx, chatID)
		if err != nil {
			return tg_errors.New(tg_errors.Config{
				OriginalErr: err,
				Handler:     "mapToCommandHandler",
				CausedBy:    "GetCustomerState",
			})
		}
		switch state {
		case domain.StateWaitingForFIO:
			return r.handler.HandleFIOInput(ctx, m)
		case domain.StateWaitingForPhoneNumber:
			return r.handler.HandlePhoneNumberInput(ctx, m)
		case domain.StateWaitingForDeliveryAddress:
			return r.handler.HandleDeliveryAddressInput(ctx, m)
		default:
			return ErrNoHandler
		}
	}
}

func (r *Router) mapToCallbackHandler(ctx context.Context, c *tg.CallbackQuery) error {
	defer r.handler.AnswerCallback(c)

	var (
		data     = c.Data
		chatID   = c.From.ID
		username = r.getUsername(c.From)
		date     = c.Message.Time().Format(time.RFC822)
		msgID    = c.Message.MessageID
	)

	logger.Get().Debug("callback info",
		zap.String("data", c.Data),
		zap.String("from", username),
		zap.String("date", date),
	)
	intCallback, parsedArgs, err := callback.ParseButtonData(data)
	if err != nil {
		return fmt.Errorf("parse button data: %w", err)
	}
	_ = parsedArgs

	switch intCallback {
	case callback.NoOpCallback:
		return nil
	case callback.Catalog:
		return r.handler.Catalog(ctx, chatID, &msgID)
	case callback.MyCart:
		return r.handler.GetCart(ctx, chatID)
	case callback.CTypeInStock:
		return r.handler.Categories(ctx, chatID, msgID, true)
	case callback.CTypeOrder:
		return r.handler.Categories(ctx, chatID, msgID, false)
	case callback.SelectCategory:
		return r.handler.Subcategories(ctx, chatID, msgID, parsedArgs)
	case callback.FromProductCardToProducts:
		return r.handler.ProductsNew(ctx, chatID, msgID, parsedArgs)
	case callback.SelectSubcategory:
		return r.handler.Products(ctx, chatID, msgID, parsedArgs)
	case callback.SelectProduct:
		return r.handler.ProductCard(ctx, chatID, msgID, parsedArgs)
	case callback.AddToCart:
		return r.handler.AddToCart(ctx, chatID, parsedArgs)
	case callback.SelectOrderType:
		return r.handler.HandleOrderTypeInput(ctx, chatID, parsedArgs)
	case callback.EditCart:
		return r.handler.EditCart(ctx, chatID, msgID)
	case callback.DeletePositionFromCart:
		return r.handler.DeletePositionFromCart(ctx, chatID, msgID, parsedArgs)
	case callback.MakeOrder:
		// Initial step to make order
		return r.handler.AskForOrderType(ctx, chatID)
	case callback.AcceptPayment:
		return r.handler.HandlePayment(ctx, c, parsedArgs)
	default:
		return ErrNoHandler
	}

}

const (
	defaultName = ""
)

func (r *Router) getUsername(u *tg.User) string {
	if u.UserName == "" {
		return defaultName
	}
	return u.UserName
}

func (r *Router) command(actual string) func(string) bool {
	return func(want string) bool {
		return actual == want
	}
}

func (r *Router) handleError(ctx context.Context, err error, u tg.Update) {
	var (
		telegramID = u.FromChat().ID
		from       = r.getUsername(u.SentFrom())
	)
	var telegramError *tg_errors.Error

	if errors.As(err, &telegramError) {
		errAsJson, err := telegramError.ToJSON()
		if err != nil {
			logger.Get().Error("telegramError.ToJSON", zap.Error(err))
		}
		logger.Get().Error("error in handler occurred",
			zap.Any("error", errAsJson),
			zap.String("from", from),
			zap.Int64("telegramId", telegramID),
		)
	} else {
		logger.Get().Error("non-telegram error in handler occurred",
			zap.String("from", from),
			zap.Int64("telegramId", telegramID),
			zap.Error(err),
		)
	}

	r.handler.Sorry(telegramID)
}
