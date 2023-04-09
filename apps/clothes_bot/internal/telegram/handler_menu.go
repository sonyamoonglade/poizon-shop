package telegram

import (
	"context"
	"errors"
	"fmt"

	"domain"
	"functools"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"utils/ranges"
)

func (h *handler) Start(ctx context.Context, m *tg.Message) error {
	var (
		chatID       = m.Chat.ID
		telegramID   = chatID
		chatUsername = m.From.String()
		username     = domain.MakeUsername(chatUsername)
	)
	// register customer
	_, err := h.customerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		if !errors.Is(err, domain.ErrCustomerNotFound) {
			return err
		}
		// save to db
		if err := h.customerRepo.Save(ctx, domain.NewCustomer(telegramID, username)); err != nil {
			return err
		}
	}

	return h.sendWithKeyboard(chatID, getStartTemplate(username), initialMenuKeyboard)

}

func (h *handler) Menu(ctx context.Context, chatID int64) error {
	if err := h.customerRepo.UpdateState(ctx, chatID, domain.StateDefault); err != nil {
		return err
	}
	rate, err := h.rateProvider.GetYuanRate(ctx)
	if err != nil {
		return fmt.Errorf("get yuan rate: %w", err)
	}
	if err := h.sendMessage(chatID, fmt.Sprintf("Курс юаня на сегодня: %.2f ₽", rate)); err != nil {
		return err
	}

	return h.sendWithKeyboard(chatID, getTemplate().Menu, menuButtons)
}

func (h *handler) MyOrders(ctx context.Context, chatID int64) error {
	var telegramID = chatID

	customer, err := h.customerRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return err
	}

	orders, err := h.orderService.GetAllForCustomer(ctx, customer.CustomerID)
	if err != nil {
		if errors.Is(err, domain.ErrNoOrders) {
			return h.sendMessage(chatID, "У тебя еще нет заказов 🦕")
		}

		return err
	}
	if len(orders) == 0 {
		return h.sendMessage(chatID, "У тебя еще нет заказов 🦕")
	}
	var name string
	if customer.FullName != nil {
		name = *customer.FullName
	} else {
		name = *customer.Username
	}
	out := getMyOrdersStart(name)
	for _, o := range orders {
		out += getSingleOrderPreview(singleOrderArgs{
			shortID:         o.ShortID,
			isExpress:       o.IsExpress,
			isPaid:          o.IsPaid,
			isApproved:      o.IsApproved,
			cartLen:         len(o.Cart),
			deliveryAddress: o.DeliveryAddress,
			comment:         o.Comment,
			status:          o.Status,
			totalYuan:       o.AmountYUAN,
			totalRub:        o.AmountRUB,
		})
		for nCartItem, cartItem := range o.Cart {
			out += getPositionTemplate(cartPositionPreviewArgs{
				n:         nCartItem + 1,
				link:      cartItem.ShopLink,
				size:      cartItem.Size,
				category:  string(cartItem.Category),
				priceRub:  cartItem.PriceRUB,
				priceYuan: cartItem.PriceYUAN,
			})
		}

		out += getTemplate().MyOrdersEnd
	}

	return h.sendMessage(chatID, out)
}

func (h *handler) FAQ(ctx context.Context, chatID int64) error {
	return h.sendWithKeyboard(chatID, faqMenuTemplate, prepareFaqButtons())
}

func (h *handler) AnswerQuestion(chatID int64, n int) error {
	defer h.askForMoreFaq(chatID)
	answers := GetAnswers(n)
	// For questions 1,2,4,5 attach image to msg.
	if ranges.In(n, []int{1, 2, 4, 5}) {
		imageURLs, ok := GetImageURLs(n)
		if !ok {
			return fmt.Errorf("invalid image urls ask")
		}
		return h.answerQuestionWithPhoto(chatID, answers, imageURLs)
	} else {
		return h.answerQuestionWithVideo(chatID, n, answers, GetVideoPath(n))
	}
}

func (h *handler) answerQuestionWithPhoto(chatID int64, answers []string, imageURLs []string) error {
	n_answers := len(answers)

	if n_answers > 1 {
		// Send images with caption (answers[0])
		if err := h.sendAnswerWithImages(chatID, answers[0], imageURLs); err != nil {
			return err
		}

		// Send the rest of answers
		for _, leftAns := range answers[1:] {
			return h.sendMessage(chatID, leftAns)
		}

		return nil
	}

	// If there's only one answer then just attach every image to it and send as a group
	return h.sendAnswerWithImages(chatID, answers[0], imageURLs)
}

func (h *handler) sendAnswerWithImages(chatID int64, answer string, imageURLs []string) error {
	var first bool
	thumbnails := functools.Map(func(url string, i int) interface{} {
		thumbnail := tg.NewInputMediaPhoto(tg.FileURL(url))
		if !first {
			thumbnail.Caption = answer
			thumbnail.ParseMode = parseModeHTML
			first = true
		}
		return thumbnail
	}, imageURLs)

	group := tg.NewMediaGroup(chatID, thumbnails)
	_, err := h.b.SendMediaGroup(group)
	return err
}

func (h *handler) answerQuestionWithVideo(chatID int64, n int, answers []string, videoURL string) error {
	// Case for video file attached to msg.
	var sentBaseImage bool
	for i, ans := range answers {
		// Last answer, so send video with it
		if i == len(answers)-1 {
			// Send answer
			msg := tg.NewMessage(chatID, ans)
			msg.ParseMode = parseModeHTML
			if err := h.cleanSend(msg); err != nil {
				return err
			}
			if err := h.sendMessage(chatID, "Отправляем видео..."); err != nil {
				return err
			}
			// Send video
			return h.cleanSend(tg.NewVideo(chatID, tg.FilePath(videoURL)))
		}
		// In order to prevent default image
		hasLink := AnswerHasLink(n)
		if hasLink && !sentBaseImage {
			if err := h.sendAnswerWithPhoto(chatID, ans, GetBaseImageURL()); err != nil {
				return err
			}
			sentBaseImage = true
			continue
		}
		msg := tg.NewMessage(chatID, ans)
		msg.ParseMode = parseModeHTML
		if err := h.cleanSend(msg); err != nil {
			return err
		}
	}
	return nil
}

func (h *handler) sendAnswerWithPhoto(chatID int64, answer string, imageURL string) error {
	photo := tg.NewPhoto(chatID, tg.FileURL(imageURL))
	photo.ParseMode = parseModeHTML
	photo.Caption = answer
	return h.cleanSend(photo)
}

func (h *handler) askForMoreFaq(chatID int64) error {
	return h.sendWithKeyboard(chatID, "Остались вопросы ❓", askMoreFaqButtons)
}
