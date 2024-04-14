package tgbot

import (
	"errors"
	"strings"
)

var (
	ErrNotFound = errors.New("not found in map")
)

func (b *Bot) CallbackStrings(callbackData string) (error, ViewFunc) {
	switch {

	case strings.HasPrefix(callbackData, "channel_get_"):
		callbackView, ok := b.callbackView["channel_get"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "contest_setting"):
		callbackView, ok := b.callbackView["contest_setting"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "get_all_contest"):
		callbackView, ok := b.callbackView["get_all_contest"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "create_contest"):
		callbackView, ok := b.callbackView["create_contest"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "delete_contest"):
		callbackView, ok := b.callbackView["delete_contest"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "contest_get_"):
		callbackView, ok := b.callbackView["contest_get"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "download_rating_"):
		callbackView, ok := b.callbackView["download_rating"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "question_setting_"):
		callbackView, ok := b.callbackView["question_setting"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "get_all_question"):
		callbackView, ok := b.callbackView["get_all_question"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "create_question"):
		callbackView, ok := b.callbackView["create_question"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "delete_question"):
		callbackView, ok := b.callbackView["delete_question"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "question_get_"):
		callbackView, ok := b.callbackView["question_get"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "question_change_name_"):
		callbackView, ok := b.callbackView["question_change_name"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "question_add_answer_"):
		callbackView, ok := b.callbackView["question_add_answer"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "question_delete_answer_"):
		callbackView, ok := b.callbackView["question_delete_answer"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "answer_delete_"):
		callbackView, ok := b.callbackView["answer_delete"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "question_change_deadline_"):
		callbackView, ok := b.callbackView["question_change_deadline"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "main_menu"):
		callbackView, ok := b.callbackView["main_menu"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "contest_delete_"):
		callbackView, ok := b.callbackView["contest_delete"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "cancel_command"):
		callbackView, ok := b.callbackView["cancel_command"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "get_all_question_"):
		callbackView, ok := b.callbackView["get_all_question"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "question_admin_view_"):
		callbackView, ok := b.callbackView["question_admin_view"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "question_send_user_"):
		callbackView, ok := b.callbackView["question_send_user"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	case strings.HasPrefix(callbackData, "answer_get_"):
		callbackView, ok := b.callbackView["answer_get"]
		if !ok {
			return ErrNotFound, nil
		}
		return nil, callbackView

	default:
		return nil, nil
	}
}
