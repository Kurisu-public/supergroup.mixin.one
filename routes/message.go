package routes

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/MixinNetwork/supergroup.mixin.one/middlewares"
	"github.com/MixinNetwork/supergroup.mixin.one/models"
	"github.com/MixinNetwork/supergroup.mixin.one/session"
	"github.com/MixinNetwork/supergroup.mixin.one/views"
	"github.com/dimfeld/httptreemux"
)

type messageImpl struct{}

func registerMesseages(router *httptreemux.TreeMux) {
	impl := messageImpl{}

	router.GET("/messages", impl.index)
	router.POST("/messages/:id/recall", impl.recall)
}

func (impl *messageImpl) index(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	user := middlewares.CurrentUser(r)
	if user.GetRole() != "admin" {
		views.RenderErrorResponse(w, r, session.ForbiddenError(r.Context()))
	} else if messages, err := models.LastestMessageWithUser(r.Context(), 200); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderMessages(w, r, messages)
	}
}

func (impl *messageImpl) recall(w http.ResponseWriter, r *http.Request, params map[string]string) {
	message, err := models.FindMessage(r.Context(), params["id"])
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	switch message.Category {
	case models.MessageCategoryPlainText,
		models.MessageCategoryPlainImage,
		models.MessageCategoryPlainVideo,
		models.MessageCategoryPlainData,
		models.MessageCategoryPlainSticker,
		models.MessageCategoryPlainContact,
		models.MessageCategoryPlainAudio,
		models.MessageCategoryAppButtonGroup:
	default:
		views.RenderErrorResponse(w, r, session.ForbiddenError(r.Context()))
		return
	}
	data, err := json.Marshal(map[string]string{"message_id": params["id"]})
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	str := base64.StdEncoding.EncodeToString(data)
	t := time.Now()
	id := models.UniqueConversationId(params["id"], middlewares.CurrentUser(r).UserId)
	_, err = models.CreateMessage(r.Context(), middlewares.CurrentUser(r), id, models.MessageCategoryMessageRecall, "", str, t, t)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderBlankResponse(w, r)
	}
}
