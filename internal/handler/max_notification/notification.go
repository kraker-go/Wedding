package max_notification

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "embed"

	models "wedding/internal/domain"
)

const apiURL = "https://platform-api2.max.ru"

// UserActions — действия, которые MAX может попросить выполнить.
type UserActions interface {
	DeleteUserHandler(ctx context.Context, id int) error
	UpdateUserHandler(ctx context.Context, id int) error
}

type Notification struct {
	Guest    models.Guest
	NewGuest models.Guest
	Action   string
}

type Notifier struct {
	token   string
	chatIDs []string

	client  *http.Client
	handler UserActions
}

func NewNotifier(
	token string,
	chatID1 string,
	chatID2 string,
	handler UserActions,
) *Notifier {
	return &Notifier{
		token: token,
		chatIDs: []string{
			chatID1,
			chatID2,
		},
		handler: handler,

		client: createMaxHTTPClient(),
	}
}

func (n *Notifier) sendMessage(
	ctx context.Context,
	text string,
	buttons [][]Button,
) error {

	body := MessageRequest{
		Text: text,
	}

	if len(buttons) > 0 {
		body.Attachments = []Attachment{
			{
				Type: "inline_keyboard",
				Payload: InlineKeyboardPayload{
					Buttons: buttons,
				},
			},
		}
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal MAX message: %w", err)
	}

	for _, chatID := range n.chatIDs {

		url := fmt.Sprintf(
			"%s/messages?chat_id=%s",
			apiURL,
			chatID,
		)

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			url,
			bytes.NewReader(data),
		)
		if err != nil {
			log.Printf(
				"❌ MAX: ошибка создания запроса для chat_id=%s: %v",
				chatID,
				err,
			)

			continue
		}

		req.Header.Set("Authorization", n.token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := n.client.Do(req)
		if err != nil {
			log.Printf(
				"❌ MAX: ошибка отправки в chat_id=%s: %v",
				chatID,
				err,
			)

			continue
		}

		resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			log.Printf(
				"❌ MAX: chat_id=%s вернул статус %s",
				chatID,
				resp.Status,
			)

			continue
		}

		log.Printf(
			"✅ MAX: сообщение отправлено chat_id=%s",
			chatID,
		)
	}

	return nil
}

type MessageRequest struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Type    string                `json:"type"`
	Payload InlineKeyboardPayload `json:"payload"`
}

type InlineKeyboardPayload struct {
	Buttons [][]Button `json:"buttons"`
}

type Button struct {
	Type string `json:"type"`
	Text string `json:"text"`

	Payload string `json:"payload"`
}

func (n *Notifier) NotifyDelete(
	guest models.Guest,
) {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	text := fmt.Sprintf(
		"⚠️ Запрос на удаление гостя\n\nID: %d\nИмя: %s %s",
		guest.ID,
		guest.FirstName,
		guest.LastName,
	)

	buttons := [][]Button{
		{
			{
				Type: "callback",
				Text: "✅ Подтвердить",
				Payload: fmt.Sprintf(
					"approve_delete:%d",
					guest.ID,
				),
			},
			{
				Type: "callback",
				Text: "❌ Отклонить",
				Payload: fmt.Sprintf(
					"reject_delete:%d",
					guest.ID,
				),
			},
		},
	}

	if err := n.sendMessage(ctx, text, buttons); err != nil {
		log.Printf(
			"❌ MAX: ошибка отправки удаления: %v",
			err,
		)

		return
	}

	log.Printf(
		"✅ MAX: уведомление об удалении отправлено, guestID=%d",
		guest.ID,
	)
}

func (n *Notifier) HandleCallback(
	ctx context.Context,
	payload string,
) {

	log.Printf(
		"🔥 MAX CALLBACK: %s",
		payload,
	)

	parts := strings.SplitN(payload, ":", 2)

	if len(parts) != 2 {
		log.Printf(
			"❌ MAX: неправильный callback: %s",
			payload,
		)

		return
	}

	action := parts[0]

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Printf(
			"❌ MAX: неправильный ID: %s",
			parts[1],
		)

		return
	}

	switch action {

	case "approve_delete":

		log.Printf(
			"🔥 MAX: подтверждено удаление ID=%d",
			id,
		)

		if n.handler == nil {
			log.Println(
				"❌ MAX: UserActions не установлен",
			)

			return
		}

		err := n.handler.DeleteUserHandler(
			ctx,
			id,
		)

		if err != nil {
			log.Printf(
				"❌ MAX: ошибка удаления ID=%d: %v",
				id,
				err,
			)

			return
		}

		log.Printf(
			"✅ MAX: пользователь удалён ID=%d",
			id,
		)

	case "reject_delete":

		log.Printf(
			"❌ MAX: удаление отклонено ID=%d",
			id,
		)

	case "approve_update":

		log.Printf(
			"🔥 MAX: подтверждено изменение ID=%d",
			id,
		)

		if n.handler == nil {
			log.Println(
				"❌ MAX: UserActions не установлен",
			)

			return
		}

		err := n.handler.UpdateUserHandler(
			ctx,
			id,
		)

		if err != nil {
			log.Printf(
				"❌ MAX: ошибка изменения ID=%d: %v",
				id,
				err,
			)

			return
		}

		log.Printf(
			"✅ MAX: пользователь изменён ID=%d",
			id,
		)

	case "reject_update":

		log.Printf(
			"❌ MAX: изменение отклонено ID=%d",
			id,
		)

	default:

		log.Printf(
			"❌ MAX: неизвестный callback: %s",
			payload,
		)
	}
}

func (n *Notifier) NotifyMessage(text string) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := n.sendMessage(ctx, text, nil); err != nil {
		log.Printf(
			"❌ MAX: ошибка отправки сообщения: %v",
			err,
		)

		return
	}

	log.Println("✅ MAX: сообщение отправлено")
}

type UpdatesResponse struct {
	Updates []Update `json:"updates"`
	Marker  *int64   `json:"marker"`
}

type Update struct {
	UpdateType string    `json:"update_type"`
	Timestamp  int64     `json:"timestamp"`
	Callback   *Callback `json:"callback,omitempty"`
}

type Callback struct {
	CallbackID string `json:"callback_id"`
	Payload    string `json:"payload"`
}

func (n *Notifier) StartPolling(ctx context.Context) {
	go func() {
		var marker *int64

		for {
			select {
			case <-ctx.Done():
				log.Println("MAX polling остановлен")
				return

			default:
			}

			reqURL := apiURL + "/updates?timeout=30&limit=100&types=message_callback"

			if marker != nil {
				reqURL += fmt.Sprintf("&marker=%d", *marker)
			}

			req, err := http.NewRequestWithContext(
				ctx,
				http.MethodGet,
				reqURL,
				nil,
			)
			if err != nil {
				log.Printf("❌ MAX polling: ошибка создания запроса: %v", err)
				time.Sleep(time.Second)
				continue
			}

			req.Header.Set("Authorization", n.token)

			resp, err := n.client.Do(req)
			if err != nil {
				log.Printf("❌ MAX polling: ошибка запроса: %v", err)
				time.Sleep(2 * time.Second)
				continue
			}

			var result UpdatesResponse

			err = json.NewDecoder(resp.Body).Decode(&result)
			resp.Body.Close()

			if err != nil {
				log.Printf("❌ MAX polling: ошибка JSON: %v", err)
				time.Sleep(time.Second)
				continue
			}

			for _, update := range result.Updates {

				if update.UpdateType != "message_callback" {
					continue
				}

				if update.Callback == nil {
					continue
				}

				log.Printf(
					"🔥 MAX CALLBACK: id=%s payload=%s",
					update.Callback.CallbackID,
					update.Callback.Payload,
				)

				n.HandleCallback(
					ctx,
					update.Callback.Payload,
				)

				go n.answerCallback(
					update.Callback.CallbackID,
					"✅ Принято",
				)
			}

			marker = result.Marker
		}
	}()
}

func (n *Notifier) answerCallback(
	callbackID string,
	notification string,
) {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	body := map[string]string{
		"notification": notification,
	}

	data, err := json.Marshal(body)
	if err != nil {
		log.Printf(
			"❌ MAX answer: ошибка JSON: %v",
			err,
		)
		return
	}

	url := fmt.Sprintf(
		"%s/answers?callback_id=%s",
		apiURL,
		callbackID,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewReader(data),
	)
	if err != nil {
		log.Printf(
			"❌ MAX answer: ошибка создания запроса: %v",
			err,
		)
		return
	}

	req.Header.Set("Authorization", n.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		log.Printf(
			"❌ MAX answer: ошибка запроса: %v",
			err,
		)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf(
			"❌ MAX answer: статус %s",
			resp.Status,
		)
		return
	}

	log.Printf(
		"✅ MAX callback подтверждён: %s",
		callbackID,
	)
}

func (n *Notifier) NotifyAdd(guest models.Guest) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	text := fmt.Sprintf(
		"🎉 Добавлен новый гость на свадьбу!!! ✅\n\nИмя: %s\nФамилия: %s",
		guest.FirstName,
		guest.LastName,
	)

	if err := n.sendMessage(ctx, text, nil); err != nil {
		log.Printf(
			"❌ MAX: ошибка отправки добавления гостя: %v",
			err,
		)

		return
	}

	log.Printf(
		"✅ MAX: уведомление о добавлении отправлено",
	)
}

func (n *Notifier) NotifyUpdate(
	guest models.Guest,
	newGuest models.Guest,
) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	text := fmt.Sprintf(
		"🔄 Запрос на изменение гостя №%d\n\n"+
			"👤 Сейчас: %s %s\n"+
			"✏️ Новые данные: %s %s",
		guest.ID,
		guest.FirstName,
		guest.LastName,
		newGuest.FirstName,
		newGuest.LastName,
	)

	buttons := [][]Button{
		{
			{
				Type: "callback",
				Text: "✅ Подтвердить",
				Payload: fmt.Sprintf(
					"approve_update:%d",
					guest.ID,
				),
			},
			{
				Type: "callback",
				Text: "я загрузил наш код на амвера, а там не работают оповещения, что делать Отклонить",
				Payload: fmt.Sprintf(
					"reject_update:%d",
					guest.ID,
				),
			},
		},
	}

	if err := n.sendMessage(ctx, text, buttons); err != nil {
		log.Printf(
			"❌ MAX: ошибка отправки запроса на изменение: %v",
			err,
		)

		return
	}

	log.Printf(
		"✅ MAX: запрос на изменение отправлен, guestID=%d",
		guest.ID,
	)
}

//go:embed certs/max.cer
var maxCert []byte

func createMaxHTTPClient() *http.Client {
	rootCAs, err := x509.SystemCertPool()
	if err != nil || rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// Если сертификат PEM.
	if ok := rootCAs.AppendCertsFromPEM(maxCert); ok {
		log.Println("✅ MAX: сертификат Минцифры добавлен из PEM")
	} else {
		// Если сертификат DER (.cer).
		cert, err := x509.ParseCertificate(maxCert)
		if err != nil {
			log.Printf(
				"❌ MAX: не удалось прочитать сертификат: %v",
				err,
			)
		} else {
			rootCAs.AddCert(cert)

			log.Println(
				"✅ MAX: сертификат Минцифры добавлен из DER",
			)
		}
	}

	return &http.Client{
		Timeout: 95 * time.Second,

		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: rootCAs,
			},
		},
	}
}
