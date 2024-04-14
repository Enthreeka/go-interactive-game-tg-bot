package tgbot

import (
	"context"
	"encoding/json"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/boterror"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/handler"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/service"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/logger"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/postgres"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/store"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"runtime/debug"
	"sync"
	"time"
)

type ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error

type Bot struct {
	bot   *tgbotapi.BotAPI
	log   *logger.Logger
	store *store.Store
	tgMsg *tg.TelegramMsg
	pg    *postgres.Postgres

	userService      service.UserService
	answersService   service.AnswersService
	questionsService service.QuestionsService
	contestService   service.ContestService

	cmdView      map[string]ViewFunc
	callbackView map[string]ViewFunc

	mu      sync.RWMutex
	isDebug bool
}

func NewBot(bot *tgbotapi.BotAPI,
	log *logger.Logger,
	store *store.Store,
	tgMsg *tg.TelegramMsg,
	pg *postgres.Postgres,
	userService service.UserService,
	answersService service.AnswersService,
	questionsService service.QuestionsService,
	contestService service.ContestService,
) *Bot {
	return &Bot{
		bot:              bot,
		log:              log,
		store:            store,
		tgMsg:            tgMsg,
		pg:               pg,
		userService:      userService,
		answersService:   answersService,
		questionsService: questionsService,
		contestService:   contestService,
	}
}

func (b *Bot) RegisterCommandView(cmd string, view ViewFunc) {
	if b.cmdView == nil {
		b.cmdView = make(map[string]ViewFunc)
	}

	b.cmdView[cmd] = view
}

func (b *Bot) RegisterCommandCallback(callback string, view ViewFunc) {
	if b.callbackView == nil {
		b.callbackView = make(map[string]ViewFunc)
	}

	b.callbackView[callback] = view
}

func (b *Bot) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)
	for {
		select {
		case update := <-updates:
			updateCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)

			b.isDebug = false
			b.jsonDebug(update)

			b.handlerUpdate(updateCtx, &update)

			cancel()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (b *Bot) jsonDebug(update any) {
	if b.isDebug {
		updateByte, err := json.MarshalIndent(update, "", " ")
		if err != nil {
			b.log.Error("%v", err)
		}
		b.log.Info("%s", updateByte)
	}
}

func (b *Bot) handlerUpdate(ctx context.Context, update *tgbotapi.Update) {
	defer func() {
		if p := recover(); p != nil {
			b.log.Error("panic recovered: %v, %s", p, string(debug.Stack()))
		}
	}()

	// if write message
	if update.Message != nil {
		b.log.Info("[%s] %s", update.Message.From.UserName, update.Message.Text)

		isProcessing, err := b.isStoreProcessing(ctx, update)
		if err != nil {
			b.log.Error("failed in isStoreProcessing: %v", err)
			handler.HandleError(b.bot, update, err.Error())
			return
		}
		if isProcessing {
			return
		}

		if err := b.userService.CreateUserIfNotExist(ctx, userUpdateToModel(update)); err != nil {
			b.log.Error("userService.CreateUserIfNotExist: failed to create user: %v", err)
			return
		}

		var view ViewFunc

		cmd := update.Message.Command()

		cmdView, ok := b.cmdView[cmd]
		if !ok {
			return
		}

		view = cmdView

		if err := view(ctx, b.bot, update); err != nil {
			b.log.Error("failed to handle update: %v", err)
			handler.HandleError(b.bot, update, boterror.ParseErrToText(err))
			return
		}
		//  if press button
	} else if update.CallbackQuery != nil {
		b.log.Info("[%s] %s", update.CallbackQuery.From.UserName, update.CallbackData())

		var callback ViewFunc

		err, callbackView := b.CallbackStrings(update.CallbackData())
		if err != nil {
			b.log.Error("%v", err)
			return
		}

		callback = callbackView

		if err := callback(ctx, b.bot, update); err != nil {
			b.log.Error("failed to handle update: %v", err)
			handler.HandleError(b.bot, update, boterror.ParseErrToText(err))
			return
		}
		// if request on join chat
	} else if update.ChatJoinRequest != nil {
		b.log.Info("[%s] %s", update.ChatJoinRequest.From.UserName, update.ChatJoinRequest.InviteLink.InviteLink)

		// if bot update/delete from channel
	} else if update.MyChatMember != nil {

		if update.MyChatMember.NewChatMember.WasKicked() {
			if err := b.userService.UpdateBlockedBotStatus(ctx, update.MyChatMember.From.ID, true); err != nil {
				b.log.Error("userService.UpdateBlockedBotStatus: %v", err)
				return
			}
		}

	}
}
