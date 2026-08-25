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
	Action   string // "add" или "update"
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

	fmt.Println("Callback:", data)

	parts := strings.Split(data, ":")
	if len(parts) != 2 {
		log.Println("Некорректный callback:", data)
		return
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Println("Неверный ID:", err)
		return
	}

	fmt.Println("ID гостя:", id)

	if n.handler == nil {
		log.Println("UserHandler не установлен")
		return
	}

	// ==================================================
	// УДАЛЕНИЕ — ПОДТВЕРЖДЕНО
	// ==================================================

	if strings.HasPrefix(data, "approve_delete:") {

		err := n.handler.DeleteUserHandler(ctx, id)
		if err != nil {
			log.Printf(
				"Ошибка удаления гостя %d: %v",
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

		fmt.Println("Гость успешно удалён:", id)

		return
	}

	// УДАЛЕНИЕ — ОТКЛОНЕНО

	if strings.HasPrefix(data, "reject_delete:") {

		fmt.Println("Удаление отклонено:", id)

		go n.NotifyMessage(
			"❌ Удаление пользователя отклонено.",
		)

		return
	}

	// ОБНОВЛЕНИЕ — ПОДТВЕРЖДЕНО

	if strings.HasPrefix(data, "approve_update:") {

		err := n.handler.UpdateUserHandler(ctx, id)
		if err != nil {
			log.Printf(
				"Ошибка обновления гостя %d: %v",
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

		fmt.Println("Гость успешно обновлён:", id)

		return
	}

	// ОБНОВЛЕНИЕ — ОТКЛОНЕНО

	if strings.HasPrefix(data, "reject_update:") {

		fmt.Println("Обновление отклонено:", id)

		go n.NotifyMessage(
			"❌ Изменение пользователя отклонено.",
		)

		return
	}

	// НЕИЗВЕСТНЫЙ CALLBACK

	log.Println("Неизвестный callback:", data)
}

func Telegramm(botToken, chatID string, handler *UserHandler) *Notifier {
	bot, err := tgbot.New(botToken)
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

	// Регистрируем обработчик ПОСЛЕ создания n
	bot.RegisterHandler(
		tgbot.HandlerTypeCallbackQueryData,
		"approve_delete:",
		tgbot.MatchTypePrefix,
		n.callbackHandler,
	)

	bot.RegisterHandler(
		tgbot.HandlerTypeCallbackQueryData,
		"approve_update:",
		tgbot.MatchTypePrefix,
		n.callbackHandler,
	)

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := n.bot.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: n.chatID,
		Text:   text,
	})
	if err != nil {
		log.Printf("Не удалось отправить уведомление: %v", err)
	}
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
		log.Printf("Не удалось отправить уведомление: %v", err)
	}

	if err != nil {

		log.Printf(
			"Не удалось отправить уведомление: %v",
			err,
		)
	}
}

func formatMessage(g models.Guest, newG models.Guest, action string) string {
	switch action {
	case "add":
		return fmt.Sprintf("\n🎉 Добавлен новый гость на свадьбу!!! ✅ \n\nИмя: %s\nФамилия: %s\n", g.FirstName, g.LastName)
	case "delete":
		return fmt.Sprintf("Запрос на удаление гостя! ⚠️\nИмя: %d %s %s", g.ID, g.FirstName, g.LastName)
	case "update":
		return fmt.Sprintf("🔄 Запрос на изменение гостя №%d\n👤 Гость: %s %s\n✏️ Редактируем: %s  %s", g.ID, g.FirstName, g.LastName, newG.FirstName, newG.LastName)
	default:
		return fmt.Sprintf("ℹ️ Гость: %d %s %s", g.ID, g.FirstName, g.LastName)
	}
}

func (n *Notifier) Notify(guest models.Guest, newG models.Guest, action string) {
	select {
	case n.ch <- Notification{Guest: guest,
		Action:   action,
		NewGuest: newG}:
	default:
		log.Println("Канал уведомлений переполнен, сообщение потеряно")
	}
}

func (n *Notifier) Shutdown() {
	close(n.ch)
	n.wg.Wait()
}
func (n *Notifier) NotifyMessage(text string) {
	select {
	case n.msgCh <- text:
	default:
		log.Println("Канал уведомлений переполнен, сообщение потеряно")
	}
}
