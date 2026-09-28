package handler

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbot "github.com/go-telegram/bot"
	tgbotmodels "github.com/go-telegram/bot/models"
	models "wedding/internal/domain"
)

type Notification struct {
	Guest    models.Guest
	NewGuest models.Guest
	Action   string // "add", "delete", "update"
}

type Notifier struct {
	ch      chan Notification
	msgCh   chan string
	wg      sync.WaitGroup
	bot     *tgbot.Bot
	chatID  string
	handler *UserHandler
}

func (n *Notifier) SetHandler(handler *UserHandler) {
	n.handler = handler
}

func (n *Notifier) callbackHandler(
	ctx context.Context,
	bot *tgbot.Bot,
	update *tgbotmodels.Update,
) {
	if update.CallbackQuery == nil {
		return
	}

	data := update.CallbackQuery.Data

	log.Printf("🔥 CALLBACK ПОЛУЧЕН: %s", data)

	// Сообщаем Telegram, что нажатие кнопки получено.
	_, err := bot.AnswerCallbackQuery(
		ctx,
		&tgbot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		},
	)

	if err != nil {
		log.Printf("❌ Ошибка AnswerCallbackQuery: %v", err)
	}

	// Ожидаем формат:
	// approve_delete:123
	// reject_delete:123
	// approve_update:123
	// reject_update:123
	parts := strings.SplitN(data, ":", 2)

	if len(parts) != 2 {
		log.Printf("❌ Некорректный callback: %s", data)
		return
	}

	action := parts[0]

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Printf(
			"❌ Неверный ID гостя: %s: %v",
			parts[1],
			err,
		)
		return
	}

	log.Printf(
		"🔥 CALLBACK ACTION=%s ID=%d",
		action,
		id,
	)

	if n.handler == nil {
		log.Println("❌ UserHandler не установлен")
		return
	}

	switch action {

	case "approve_delete":

		log.Printf(
			"🔥 ПОДТВЕРЖДЕНО УДАЛЕНИЕ ГОСТЯ ID=%d",
			id,
		)

		err := n.handler.DeleteUserHandler(ctx, id)
		if err != nil {
			log.Printf(
				"❌ Ошибка удаления гостя %d: %v",
				id,
				err,
			)

			go n.NotifyMessage(
				fmt.Sprintf(
					"❌ Ошибка удаления пользователя №%d.",
					id,
				),
			)

			return
		}

		log.Printf(
			"✅ Гость успешно удалён: %d",
			id,
		)

	case "reject_delete":

		log.Printf(
			"❌ Удаление отклонено: %d",
			id,
		)

		go n.NotifyMessage(
			"❌ Удаление пользователя отклонено.",
		)

	case "approve_update":

		log.Printf(
			"🔥 ПОДТВЕРЖДЕНО ОБНОВЛЕНИЕ ГОСТЯ ID=%d",
			id,
		)

		err := n.handler.UpdateUserHandler(ctx, id)
		if err != nil {
			log.Printf(
				"❌ Ошибка обновления гостя %d: %v",
				id,
				err,
			)

			go n.NotifyMessage(
				fmt.Sprintf(
					"❌ Ошибка обновления пользователя №%d.",
					id,
				),
			)

			return
		}

		log.Printf(
			"✅ Гость успешно обновлён: %d",
			id,
		)

	case "reject_update":

		log.Printf(
			"❌ Обновление отклонено: %d",
			id,
		)

		go n.NotifyMessage(
			"❌ Изменение пользователя отклонено.",
		)

	default:

		log.Printf(
			"❌ НЕИЗВЕСТНЫЙ CALLBACK: %s",
			data,
		)
	}
}

func Telegramm(
	botToken string,
	chatID string,
	handler *UserHandler,
) *Notifier {

	bot, err := tgbot.New(
		botToken,
		tgbot.WithAllowedUpdates([]string{
			"message",
			"callback_query",
		}),
	)
	if err != nil {
		log.Fatal("Telegram bot init failed:", err)
	}

	n := &Notifier{
		ch:      make(chan Notification, 100),
		msgCh:   make(chan string, 100),
		bot:     bot,
		chatID:  chatID,
		handler: handler,
	}

	// Один обработчик для всех callback-кнопок.
	//
	// Он будет получать:
	// approve_delete:123
	// reject_delete:123
	// approve_update:123
	// reject_update:123
	bot.RegisterHandler(
		tgbot.HandlerTypeCallbackQueryData,
		"",
		tgbot.MatchTypePrefix,
		n.callbackHandler,
	)

	log.Println("🤖 Telegram bot запускается")

	go bot.Start(context.Background())

	n.wg.Add(1)
	go n.worker()

	return n
}

func (n *Notifier) worker() {
	defer n.wg.Done()

	for {
		select {

		case notif := <-n.ch:
			n.sendNotification(notif)

		case text := <-n.msgCh:
			n.sendMessage(text)
		}
	}
}

func (n *Notifier) sendMessage(text string) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	_, err := n.bot.SendMessage(
		ctx,
		&tgbot.SendMessageParams{
			ChatID: n.chatID,
			Text:   text,
		},
	)

	if err != nil {
		log.Printf(
			"❌ Не удалось отправить уведомление: %v",
			err,
		)

		return
	}

	log.Println("✅ Telegram сообщение отправлено")
}

func (n *Notifier) sendNotification(notif Notification) {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	text := formatMessage(
		notif.Guest,
		notif.NewGuest,
		notif.Action,
	)

	var keyboard *tgbotmodels.InlineKeyboardMarkup

	if notif.Action == "delete" || notif.Action == "update" {

		keyboard = &tgbotmodels.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgbotmodels.InlineKeyboardButton{
				{
					{
						Text: "✅ Подтвердить",

						CallbackData: fmt.Sprintf(
							"approve_%s:%d",
							notif.Action,
							notif.Guest.ID,
						),
					},

					{
						Text: "❌ Отклонить",

						CallbackData: fmt.Sprintf(
							"reject_%s:%d",
							notif.Action,
							notif.Guest.ID,
						),
					},
				},
			},
		}
	}

	params := &tgbot.SendMessageParams{
		ChatID: n.chatID,
		Text:   text,
	}

	if keyboard != nil {
		params.ReplyMarkup = keyboard
	}

	_, err := n.bot.SendMessage(ctx, params)

	if err != nil {
		log.Printf(
			"❌ Не удалось отправить уведомление: %v",
			err,
		)

		return
	}

	log.Printf(
		"✅ Telegram уведомление отправлено: action=%s guestID=%d",
		notif.Action,
		notif.Guest.ID,
	)
}

func formatMessage(
	g models.Guest,
	newG models.Guest,
	action string,
) string {

	switch action {

	case "add":

		return fmt.Sprintf(
			"\n🎉 Добавлен новый гость на свадьбу!!! ✅\n\nИмя: %s\nФамилия: %s\n",
			g.FirstName,
			g.LastName,
		)

	case "delete":

		return fmt.Sprintf(
			"Запрос на удаление гостя! ⚠️\nИмя: %d %s %s",
			g.ID,
			g.FirstName,
			g.LastName,
		)

	case "update":

		return fmt.Sprintf(
			"🔄 Запрос на изменение гостя №%d\n👤 Гость: %s %s\n✏️ Редактируем: %s %s",
			g.ID,
			g.FirstName,
			g.LastName,
			newG.FirstName,
			newG.LastName,
		)

	default:

		return fmt.Sprintf(
			"ℹ️ Гость: %d %s %s",
			g.ID,
			g.FirstName,
			g.LastName,
		)
	}
}

func (n *Notifier) Notify(
	guest models.Guest,
	newG models.Guest,
	action string,
) {
	select {

	case n.ch <- Notification{
		Guest:    guest,
		Action:   action,
		NewGuest: newG,
	}:

	default:

		log.Println(
			"❌ Канал уведомлений переполнен, сообщение потеряно",
		)
	}
}

func (n *Notifier) NotifyMessage(text string) {
	select {

	case n.msgCh <- text:

	default:

		log.Println(
			"❌ Канал уведомлений переполнен, сообщение потеряно",
		)
	}
}

func (n *Notifier) Shutdown() {
	close(n.ch)
	n.wg.Wait()
}
