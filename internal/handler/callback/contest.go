package callback

import (
	"context"
	"fmt"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/boterror"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/entity"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/handler"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/handler/tgbot"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/service"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/excel"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/logger"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/store"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/tg"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/tg/markup"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sync"
)

type CallbackContest struct {
	contestService service.ContestService
	store          *store.Store
	log            *logger.Logger
	tgMsg          *tg.TelegramMsg
	excel          *excel.Excel

	mu sync.RWMutex
}

func NewCallbackContest(contestService service.ContestService, store *store.Store, log *logger.Logger, tgMsg *tg.TelegramMsg, excel *excel.Excel) *CallbackContest {
	return &CallbackContest{
		contestService: contestService,
		log:            log,
		tgMsg:          tgMsg,
		store:          store,
		excel:          excel,
	}
}

// CallbackContestSetting - contest_setting
func (c *CallbackContest) CallbackContestSetting() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {

		if _, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			&markup.ContestSetting,
			"CallbackContestSetting"); err != nil {
			return err
		}

		return nil
	}
}

// CallbackGetAllContest - get_all_contest
// Выводятся кнопки вида contest_get_{contest_id}
func (c *CallbackContest) CallbackGetAllContest() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		_, markupContest, err := c.contestService.GetAllContestsButtons(ctx, "get")
		if err != nil {
			c.log.Error("contestService.GetAllContests: failed to get contest: %v", err)
			return err
		}

		if _, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			markupContest,
			"CallbackGetAllContest"); err != nil {
			return err
		}

		return nil
	}
}

// CallbackCreateContest - create_contest
func (c *CallbackContest) CallbackCreateContest() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		msg, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			&markup.CancelState,
			`Отправить через json`+
				"\n{\n"+
				`"название_конкурса":"сюда вписать название"`+"\n"+
				`"дедлайн":"определенный формат данных"`+
				"\n}")
		if err != nil {
			return err
		}

		c.store.Delete(update.FromChat().ID)
		c.store.Set(&store.ContestStore{
			MsgID:              msg,
			UserID:             update.FromChat().ID,
			TypeCommandContest: store.ContestCreate,
		}, update.FromChat().ID)

		return nil
	}
}

// CallbackDeleteContest - delete_contest
func (c *CallbackContest) CallbackDeleteContest() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		_, markupContest, err := c.contestService.GetAllContestsButtons(ctx, "delete")
		if err != nil {
			c.log.Error("contestService.GetAllContests: failed to get contest: %v", err)
			return err
		}

		if _, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			markupContest,
			"CallbackDeleteContest"); err != nil {
			return err
		}

		return nil
	}
}

// CallbackGetContestByID - contest_get_{contest_id}
func (c *CallbackContest) CallbackGetContestByID() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		contestID := entity.GetContestID(update.CallbackData())

		contest, err := c.contestService.GetContestByID(ctx, contestID)
		if err != nil {
			c.log.Error("contestService.GetContestByID: failed to get contest: %v", err)
			return err
		}

		markupContest := markup.ContestByIDSetting(contestID)
		if _, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			&markupContest,
			fmt.Sprintf("%v:", contest),
		); err != nil {
			return err
		}

		return nil
	}
}

// CallbackDownloadRating - download_rating_{contest_id}
func (c *CallbackContest) CallbackDownloadRating() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		contestID := entity.GetContestID(update.CallbackData())

		userResult, err := c.contestService.GetAllUserResultsByContest(ctx, contestID)
		if err != nil {
			c.log.Error("contestService.GetAllUserResultsByContest: failed to get contest: %v", err)
			return err
		}

		c.mu.Lock()
		fileName, err := c.excel.GenerateUserResultsExcelFile(userResult, contestID, update.CallbackQuery.From.UserName)
		if err != nil {
			c.log.Error("Excel.GenerateExcelFile: failed to generate excel file: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		fileIDBytes, err := c.excel.GetExcelFile(fileName)
		if err != nil {
			c.log.Error("Excel.GetExcelFile: failed to get excel file: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}
		c.mu.Unlock()

		if fileIDBytes == nil {
			c.log.Error("fileIDBytes: %v", boterror.ErrNil)
			handler.HandleError(bot, update, boterror.ParseErrToText(boterror.ErrNil))
			return nil
		}

		if _, err := c.tgMsg.SendDocument(update.FromChat().ID,
			fileName,
			fileIDBytes,
			"CallbackDownloadRating",
		); err != nil {
			c.log.Error("tgMsg.SendDocument: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		return nil
	}
}

// CallbackContestDelete - contest_delete_{contest_id}
func (c *CallbackContest) CallbackContestDelete() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		contestID := entity.GetContestID(update.CallbackData())

		if err := c.contestService.DeleteContest(ctx, contestID); err != nil {
			c.log.Error("contestService.DeleteContest: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		if _, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			nil,
			"Удалено успешно",
		); err != nil {
			return err
		}

		return nil
	}
}
