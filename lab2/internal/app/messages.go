package app

import (
	"encoding/json"
	"errors"
	"net/http"

	"lab2/internal/model"
	"lab2/internal/service"
)

func listMessagesHandler(messageService *service.MessageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		messages, err := messageService.List(r.Context())
		if err != nil {
			writeError(w, http.StatusInternalServerError, "list messages")
			return
		}

		writeJSON(w, http.StatusOK, messages)
	}
}

func createMessageHandler(messageService *service.MessageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var input model.CreateMessageInput
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&input); err != nil {
			writeError(w, http.StatusBadRequest, "invalid json body")
			return
		}

		message, err := messageService.Create(r.Context(), input.Text)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrEmptyMessage), errors.Is(err, service.ErrMessageTooLong):
				writeError(w, http.StatusBadRequest, err.Error())
			default:
				writeError(w, http.StatusInternalServerError, "create message")
			}
			return
		}

		writeJSON(w, http.StatusCreated, message)
	}
}
