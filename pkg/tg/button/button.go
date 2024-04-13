package button

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	MainMenuButton = tgbotapi.NewInlineKeyboardButtonData("Вернуться в главное меню", "main_menu")

	BackToContestSetting = tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", "contest_setting")

	BackToGetAllContest = tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", "get_all_contest")
)

func BackToQuestionSetting(contestID int) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", fmt.Sprintf("question_setting_%d", contestID))
}

//func BackToGetAllQuestion(contestID int) tgbotapi.InlineKeyboardButton {
//	return tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", fmt.Sprintf("get_all_question_%d", contestID))
//}
