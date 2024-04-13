package markup

import (
	"fmt"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/tg/button"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	StartMenu = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Управление конкурсом", "contest_setting")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Управление данными пользователей", "user_setting")),
	)

	SuperAdminSetting = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назначить администратором", "create_admin")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назначить супер администратором", "create_super_admin")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Забрать права администратора", "delete_admin")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Список администраторов", "all_admin")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", "user_setting")),
		tgbotapi.NewInlineKeyboardRow(button.MainMenuButton),
	)

	ContestSetting = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Получить список конкурсов", "get_all_contest")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Создать конкурс", "create_contest")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Удалить конкурс", "delete_contest")),
		tgbotapi.NewInlineKeyboardRow(button.MainMenuButton),
	)

	CancelState = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Отменить выполнение команды", "cancel_command")))
)

func QuestionSetting(contestID int) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Получить список вопросов", fmt.Sprintf("get_all_question_%d", contestID))),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Создать вопрос", fmt.Sprintf("create_question_%d", contestID))),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Удалить вопрос", fmt.Sprintf("delete_question_%d", contestID))),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", fmt.Sprintf("contest_get_%d", contestID))),
		tgbotapi.NewInlineKeyboardRow(button.MainMenuButton),
	)
}

func QuestionByIDSetting(questionID int, contestID int) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить название", fmt.Sprintf("question_change_name_%d", questionID))),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить ответ", fmt.Sprintf("question_add_answer_%d", questionID))),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Удалить ответ", fmt.Sprintf("question_delete_answer_%d", questionID))),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить дедлайн", fmt.Sprintf("question_change_deadline_%d", questionID))),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", fmt.Sprintf("get_all_question_%d", contestID))),
		tgbotapi.NewInlineKeyboardRow(button.MainMenuButton),
	)
}

func ContestByIDSetting(contestID int) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Управление вопросами конкурса", fmt.Sprintf("question_setting_%d", contestID))),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Скачать рейтинг", fmt.Sprintf("download_rating_%d", contestID))),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить рейтинг (не готово)", fmt.Sprintf("update_rating_%d", contestID))),
		tgbotapi.NewInlineKeyboardRow(button.BackToGetAllContest),
		tgbotapi.NewInlineKeyboardRow(button.MainMenuButton),
	)
}
