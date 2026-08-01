package handler

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
	models "wedding/internal/domain"
	"wedding/internal/service"
)

type UserHandler struct {
	hand *service.UserService
	logg *zap.Logger
}

func NewUserHandler(hand *service.UserService, logg *zap.Logger) *UserHandler {
	return &UserHandler{hand: hand, logg: logg}
}

func (uh *UserHandler) AddUserHandler(w http.ResponseWriter, r *http.Request) {
	uh.logg.Info("Добавляем гостя")

	ctx := r.Context()

	var user models.Guest

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "не корректный ввод", http.StatusBadRequest)
		uh.logg.Error(err.Error())
		return
	}

	u, err := uh.hand.AddUser(ctx, user.FirstName, user.LastName)
	if err != nil {
		if errors.Is(err, models.ErrorAlreadyExists) {
			http.Error(w, models.ErrorAlreadyExists.Error(), http.StatusBadRequest)
			uh.logg.Error("попытка повторной регистрации", zap.Error(err))
			return
		}
		http.Error(w, "ошибка сервера, попробуйте позже", http.StatusBadRequest)
		uh.logg.Error("ошибка сервера", zap.Error(err))
		return
	}

	uh.logg.Info("Гость успешно добавлен !")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, "ошибка сервера, не удалось ответить", http.StatusInternalServerError)
		uh.logg.Error("ошибка ответа", zap.Error(err))
		return
	}

}

func (uh *UserHandler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	uh.logg.Info("Поиск всех гостей...")

	ctx := r.Context()

	users, err := uh.hand.GetAllUsers(ctx)
	if err != nil {
		http.Error(w, "ошибка поиска гостей", http.StatusBadRequest)
		uh.logg.Error(err.Error())
		return
	}

	uh.logg.Info("список гостей собран успешно")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(users); err != nil {
		uh.logg.Error(err.Error())
		return
	}
}

func (uh *UserHandler) GetCountUsersHandler(w http.ResponseWriter, r *http.Request) {

	uh.logg.Info("загружаем счетчик количества гостей")

	ctx := r.Context()

	count, err := uh.hand.GetCountUsers(ctx)
	if err != nil {
		uh.logg.Error(err.Error())
		return
	}

	uh.logg.Info("количество гостей определено")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(count); err != nil {
		uh.logg.Error(err.Error())
		return
	}
}
