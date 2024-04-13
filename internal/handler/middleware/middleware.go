package middleware

import (
	"context"
	"errors"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/boterror"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/handler/tgbot"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ChatAdminMiddleware(channelID []int64, next tgbot.ViewFunc) tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		for _, chatID := range channelID {
			admins, err := bot.GetChatAdministrators(
				tgbotapi.ChatAdministratorsConfig{
					ChatConfig: tgbotapi.ChatConfig{
						ChatID: chatID,
					},
				})

			if err != nil {
				return err
			}

			for _, admin := range admins {
				if admin.User.ID == update.Message.From.ID {
					return next(ctx, bot, update)
				}
			}
		}
		return boterror.ErrIsNotAdmin
	}
}

func AdminMiddleware(service service.UserService, next tgbot.ViewFunc) tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		user, err := service.GetUserByID(ctx, update.FromChat().ID)
		if err != nil {
			if errors.Is(err, boterror.ErrNoRows) {
				return nil
			}
			return err
		}

		if user.UserRole == "admin" || user.UserRole == "superAdmin" {
			return next(ctx, bot, update)
		}

		return boterror.ErrIsNotAdmin
	}
}

func SuperAdminMiddleware(service service.UserService, next tgbot.ViewFunc) tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		user, err := service.GetUserByID(ctx, update.FromChat().ID)
		if err != nil {
			if errors.Is(err, boterror.ErrNoRows) {
				return nil
			}
			return err
		}

		if user.UserRole == "superAdmin" {
			return next(ctx, bot, update)
		}

		return boterror.ErrIsNotSuperAdmin
	}
}
