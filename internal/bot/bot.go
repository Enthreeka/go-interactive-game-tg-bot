package bot

import (
	"context"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/config"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/handler/callback"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/handler/middleware"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/handler/tgbot"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/handler/view"
	pgRepo "github.com/Entreeka/go-interactive-game-tg-bot/internal/repo/postgres"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/service"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/excel"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/logger"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/postgres"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/store"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"os/signal"
	"syscall"
)

type Bot struct {
	bot   *tgbotapi.BotAPI
	psql  *postgres.Postgres
	store *store.Store
	tgMsg *tg.TelegramMsg
	excel *excel.Excel

	answersService   service.AnswersService
	questionsService service.QuestionsService
	contestService   service.ContestService
	userService      service.UserService

	viewGeneral *view.ViewGeneral

	callbackGeneral  *callback.CallbackGeneral
	callbackContest  *callback.CallbackContest
	callbackQuestion *callback.CallbackQuestion
	callbackUser     *callback.CallbackUser
}

func NewBot() *Bot {
	return &Bot{}
}

func (b *Bot) initServices(log *logger.Logger) {
	userRepo := pgRepo.NewUserRepo(b.psql)
	answerRepo := pgRepo.NewAnswerRepo(b.psql)
	contestRepo := pgRepo.NewContestRepo(b.psql)
	questionRepo := pgRepo.NewQuestionRepo(b.psql)
	historyPointsRepo := pgRepo.NewHistoryPointsRepo(b.psql)
	questionAnswerRepo := pgRepo.NewQuestionAnswerRepo(b.psql)
	userResultRepo := pgRepo.NewUserResultRepo(b.psql)

	b.userService = service.NewUserService(userRepo, userResultRepo, log)
	b.answersService = service.NewAnswersService(answerRepo, questionRepo, questionAnswerRepo, historyPointsRepo, log)
	b.questionsService = service.NewQuestionsService(questionRepo, questionAnswerRepo, historyPointsRepo, log)
	b.contestService = service.NewContestService(contestRepo, userResultRepo, log)
}

func (b *Bot) initHandlers(log *logger.Logger) {
	b.viewGeneral = view.NewViewGeneral(log)

	b.callbackGeneral = callback.NewCallbackGeneral(log, b.store, b.tgMsg)
	b.callbackContest = callback.NewCallbackContest(b.contestService, b.userService, b.store, log, b.tgMsg, b.excel)
	b.callbackQuestion = callback.NewCallbackQuestion(b.questionsService, b.answersService, b.userService, log, b.store, b.tgMsg, b.excel, b.psql)
	b.callbackUser = callback.NewCallbackUser(b.userService, log, b.store, b.tgMsg)
}

func (b *Bot) initExcel(log *logger.Logger) {
	b.excel = excel.NewExcel(log)
}

func (b *Bot) initStore() {
	b.store = store.NewStore()
}

func (b *Bot) initMessageImplement(log *logger.Logger) {
	b.tgMsg = tg.NewMessageSetting(b.bot, log)
}

func (b *Bot) initialize(ctx context.Context, log *logger.Logger) {
	b.initExcel(log)
	b.initMessageImplement(log)
	b.initStore()
	b.initServices(log)
	b.initHandlers(log)
}

func (b *Bot) Run(log *logger.Logger, cfg *config.Config) error {
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Fatal("failed to load token %v", err)
	}
	bot.Debug = false
	b.bot = bot

	log.Info("Authorized on account %s", bot.Self.UserName)

	psql, err := postgres.New(context.Background(), 5, cfg.Postgres.URL)
	if err != nil {
		log.Fatal("failed to connect PostgreSQL: %v", err)
	}
	defer psql.Close()
	b.psql = psql

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	b.initialize(ctx, log)

	newBot := tgbot.NewBot(bot, log, b.store, b.tgMsg, psql, b.userService, b.answersService, b.questionsService, b.contestService)

	//--admin
	newBot.RegisterCommandView("admin", middleware.AdminMiddleware(b.userService, b.viewGeneral.CallbackStartAdminPanel()))
	//--user
	newBot.RegisterCommandView("start", b.viewGeneral.ViewFirstMessage())

	//--user
	newBot.RegisterCommandCallback("answer_get", b.callbackQuestion.CallbackAnswerGet())

	//--admin
	newBot.RegisterCommandCallback("cancel_command", middleware.AdminMiddleware(b.userService, b.callbackGeneral.CallbackCancelCommand()))
	newBot.RegisterCommandCallback("main_menu", middleware.AdminMiddleware(b.userService, b.callbackGeneral.CallbackMainMenu()))

	newBot.RegisterCommandCallback("user_setting", middleware.AdminMiddleware(b.userService, b.callbackUser.AdminRoleSetting()))
	newBot.RegisterCommandCallback("admin_look_up", middleware.AdminMiddleware(b.userService, b.callbackUser.AdminLookUp()))
	newBot.RegisterCommandCallback("admin_delete_role", middleware.AdminMiddleware(b.userService, b.callbackUser.AdminDeleteRole()))
	newBot.RegisterCommandCallback("admin_set_role", middleware.AdminMiddleware(b.userService, b.callbackUser.AdminSetRole()))

	newBot.RegisterCommandCallback("contest_delete", middleware.AdminMiddleware(b.userService, b.callbackContest.CallbackContestDelete()))
	newBot.RegisterCommandCallback("contest_setting", middleware.AdminMiddleware(b.userService, b.callbackContest.CallbackContestSetting()))
	newBot.RegisterCommandCallback("get_all_contest", middleware.AdminMiddleware(b.userService, b.callbackContest.CallbackGetAllContest()))
	newBot.RegisterCommandCallback("create_contest", middleware.AdminMiddleware(b.userService, b.callbackContest.CallbackCreateContest()))
	newBot.RegisterCommandCallback("delete_contest", middleware.AdminMiddleware(b.userService, b.callbackContest.CallbackDeleteContest()))
	newBot.RegisterCommandCallback("contest_get", middleware.AdminMiddleware(b.userService, b.callbackContest.CallbackGetContestByID()))
	newBot.RegisterCommandCallback("download_rating", middleware.AdminMiddleware(b.userService, b.callbackContest.CallbackDownloadRating()))
	newBot.RegisterCommandCallback("contest_reminder", middleware.AdminMiddleware(b.userService, b.callbackContest.CallbackContestReminder()))
	newBot.RegisterCommandCallback("send_rating", middleware.AdminMiddleware(b.userService, b.callbackContest.CallbackSendRating()))
	newBot.RegisterCommandCallback("pick_random", middleware.AdminMiddleware(b.userService, b.callbackContest.CallbackPickRandom()))
	newBot.RegisterCommandCallback("send_message", middleware.AdminMiddleware(b.userService, b.callbackContest.CallbackSendMessage()))
	newBot.RegisterCommandCallback("update_rating", middleware.AdminMiddleware(b.userService, b.callbackContest.CallbackUpdateRating()))

	newBot.RegisterCommandCallback("question_setting", middleware.AdminMiddleware(b.userService, b.callbackQuestion.CallbackQuestionSetting()))
	newBot.RegisterCommandCallback("create_question", middleware.AdminMiddleware(b.userService, b.callbackQuestion.CallbackCreateQuestion()))
	newBot.RegisterCommandCallback("delete_question", middleware.AdminMiddleware(b.userService, b.callbackQuestion.CallbackDeleteQuestion()))
	newBot.RegisterCommandCallback("question_get", middleware.AdminMiddleware(b.userService, b.callbackQuestion.CallbackGetQuestionByID()))
	newBot.RegisterCommandCallback("question_change_name", middleware.AdminMiddleware(b.userService, b.callbackQuestion.CallbackQuestionChangeName()))
	newBot.RegisterCommandCallback("question_add_answer", middleware.AdminMiddleware(b.userService, b.callbackQuestion.CallbackQuestionAddAnswer()))
	newBot.RegisterCommandCallback("question_delete_answer", middleware.AdminMiddleware(b.userService, b.callbackQuestion.CallbackQuestionDeleteAnswer()))
	newBot.RegisterCommandCallback("answer_delete", middleware.AdminMiddleware(b.userService, b.callbackQuestion.CallbackAnswerDelete()))
	newBot.RegisterCommandCallback("question_change_deadline", middleware.AdminMiddleware(b.userService, b.callbackQuestion.CallbackQuestionChangeDeadline()))
	newBot.RegisterCommandCallback("get_all_question", middleware.AdminMiddleware(b.userService, b.callbackQuestion.CallbackGetAllQuestionByContestID()))
	newBot.RegisterCommandCallback("question_admin_view", middleware.AdminMiddleware(b.userService, b.callbackQuestion.CallbackQuestionAdminView()))
	newBot.RegisterCommandCallback("question_send_user", middleware.AdminMiddleware(b.userService, b.callbackQuestion.CallbackQuestionSendUser()))
	newBot.RegisterCommandCallback("close_rating", middleware.AdminMiddleware(b.userService, b.callbackQuestion.CallbackCloseRating()))
	newBot.RegisterCommandCallback("question_delete", middleware.AdminMiddleware(b.userService, b.callbackQuestion.CallbackQuestionDelete()))
	newBot.RegisterCommandCallback("top_10", middleware.AdminMiddleware(b.userService, b.callbackQuestion.CallbackGetTop10Users()))

	if err := newBot.Run(ctx); err != nil {
		log.Error("failed to run tgbot: %v", err)
	}
	return nil
}
