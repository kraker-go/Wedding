package handler

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	tgbot "github.com/go-telegram/bot"
	"wedding/internal/domain"
)

type Notification struct {
	Guest  models.Guest
	Action string // "add" или "update"
}

type Notifier struct {
	ch     chan Notification
	msgCh  chan string
	wg     sync.WaitGroup
	bot    *tgbot.Bot
	chatID string
}

func Telegramm(botToken, chatID string) *Notifier {
	bot, err := tgbot.New(botToken)
	if err != nil {
		log.Fatal("Telegram bot init failed:", err)
	}
	n := &Notifier{
		ch:     make(chan Notification, 100),
		msgCh:  make(chan string, 100),
		bot:    bot,
		chatID: chatID,
	}
	n.wg.Add(1)
	go n.worker() // запускаем воркер
	return n
}

// worker читает из канала и отправляет
func (n *Notifier) worker() {
	defer n.wg.Done()
	for {
		select {
		case notif := <-n.ch:
			msg := formatMessage(notif.Guest, notif.Action)
			n.sendMessage(msg)
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

func formatMessage(g models.Guest, action string) string {
	switch action {
	case "add":
		return fmt.Sprintf("🎉 Новый гость!\nИмя: %s %s", g.FirstName, g.LastName)
	case "delete":
		return fmt.Sprintf("Гость удалился!\nИмя: %s %s", g.FirstName, g.LastName)
	case "update":
		return fmt.Sprintf("🔄 Гость обновлён!\nИмя: %s %s", g.FirstName, g.LastName)
	default:
		return fmt.Sprintf("ℹ️ Гость: %s %s", g.FirstName, g.LastName)
	}
}

// Notify отправляет данные в канал (неблокирующий)))
func (n *Notifier) Notify(guest models.Guest, action string) {
	select {
	case n.ch <- Notification{Guest: guest, Action: action}:
	default:
		log.Println("Канал уведомлений переполнен, сообщение потеряно")
	}
}

// Shutdown останавливает воркер (для graceful shutdown)))
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
