package callback

import (
	"context"
	"encoding/json"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/handler"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/handler/tgbot"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/service"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/logger"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/store"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackUser struct {
	userService service.UserService
	log         *logger.Logger
	store       *store.Store
	tgMsg       *tg.TelegramMsg
}

func NewCallbackUser(
	userService service.UserService,
	log *logger.Logger,
	store *store.Store,
	tgMsg *tg.TelegramMsg,
) *CallbackUser {
	return &CallbackUser{
		userService: userService,
		log:         log,
		store:       store,
		tgMsg:       tgMsg,
	}
}

func (c *CallbackUser) AdminRoleSetting() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		text := "Управление администраторами"

		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Назначить роль администратора", "admin_set_role"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Отозвать роль администратора", "admin_delete_role"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Посмотреть список администраторов", "admin_look_up"),
			),
		)

		if _, err := bot.Send(msg); err != nil {
			c.log.Error("failed to send message: %v", err)
			return err
		}

		return nil
	}
}

// AdminLookUp - admin_look_up
func (c *CallbackUser) AdminLookUp() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		admin, err := c.userService.GetAllAdmin(ctx)
		if err != nil {
			c.log.Error("AdminLookUp: UserRepo.GetAllAdmin: %v", err)
			handler.HandleError(bot, update, "Временные неполадки на сервере")
			return nil
		}

		adminByte, err := json.MarshalIndent(admin, "", "\t")
		if err != nil {
			c.log.Error("AdminLookUp: json.MarshalIndent: %v", err)
			handler.HandleError(bot, update, "Временные неполадки на сервере")
			return nil
		}

		msg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, string(adminByte))
		if _, err := bot.Send(msg); err != nil {
			c.log.Error("failed to send message: %v", err)
			return err
		}

		return nil
	}
}

// AdminDeleteRole - admin_delete_role
func (c *CallbackUser) AdminDeleteRole() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		text := "Напишите никнейм пользователя, у которого вы хотите отозвать права администратором.\nДля отмены команды" +
			"отправьте /cancel"

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)

		_, err := bot.Send(msg)
		if err != nil {
			c.log.Error("failed to send message: %v", err)
			return err
		}

		c.store.Delete(update.CallbackQuery.Message.Chat.ID)
		c.store.Set(store.AdminStore{
			TypeCommand: store.UserAdminDelete,
		}, update.CallbackQuery.Message.Chat.ID)

		return nil
	}
}

// AdminSetRole - admin_set_role
func (c *CallbackUser) AdminSetRole() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		text := "Напишите никнейм пользователя, которого вы хотите назначить администратором.\nДля отмены команды" +
			"отправьте /cancel"

		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)

		_, err := bot.Send(msg)
		if err != nil {
			c.log.Error("failed to send message: %v", err)
			return err
		}

		c.store.Delete(update.CallbackQuery.Message.Chat.ID)
		c.store.Set(store.AdminStore{
			TypeCommand: store.UserAdminCreate,
		}, update.CallbackQuery.Message.Chat.ID)

		return nil
	}
}
