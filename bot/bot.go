package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/BaboaAtCity/BigNewsMorgan/alerts"
	"github.com/BaboaAtCity/BigNewsMorgan/coingecko"
	"github.com/BaboaAtCity/BigNewsMorgan/watchlist"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api *tgbotapi.BotAPI
}

func New(token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{api: api}, nil
}

func (b *Bot) getMainKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Watchlist", "watchlist"),
			tgbotapi.NewInlineKeyboardButtonData("Add Coin", "add"),
			tgbotapi.NewInlineKeyboardButtonData("Remove Coin", "remove"),
		),
	)
}

func (b *Bot) Start() {
	b.api.Debug = true
	log.Printf("Authorized on account %s", b.api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		var msg tgbotapi.MessageConfig

		if update.Message != nil && update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome! Choose an option:")
			case "add":
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, b.handleAdd(update.Message.CommandArguments()))
			case "remove":
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, b.handleRemove(update.Message.CommandArguments()))
			case "addalert":
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, b.handleAddAlert(update.Message.CommandArguments()))
			default:
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know that command")
			}
		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := b.api.Request(callback); err != nil {
				log.Println("Error sending callback:", err)
			}

			switch update.CallbackQuery.Data {
			case "watchlist":
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, b.handleWatchlist())
			case "add":
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Please send the ticker or coin name to add. Format: /add <ticker/name>")
			case "remove":
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Please send the ticker or coin name to remove. Format: /remove <ticker/name>")
			case "addalert":
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Please send the alert details. Format: /addalert <coin> <price>")
			}
		}

		// Add the inline keyboard to every message
		msg.ReplyMarkup = b.getMainKeyboard()

		if _, err := b.api.Send(msg); err != nil {
			log.Println("Error sending message:", err)
		}
	}
}

func (b *Bot) handleWatchlist() string {
	prices, err := coingecko.GetPrices(watchlist.Get())
	if err != nil {
		return "Error fetching prices. Please try again later."
	}
	return watchlist.FormatPrices(prices)
}

func (b *Bot) handleAdd(args string) string {
	if args == "" {
		return "Please provide a ticker or coin name to add. Usage: /add <ticker/name>"
	}
	return watchlist.Add(args)
}

func (b *Bot) handleRemove(args string) string {
	if args == "" {
		return "Please provide a ticker or coin name to remove. Usage: /remove <ticker/name>"
	}
	return watchlist.Remove(args)
}

func (b *Bot) handleAddAlert(args string) string {
	parts := strings.Fields(args)
	if len(parts) != 2 {
		return "Usage: /addalert <coin> <price>"
	}

	coin := strings.ToUpper(parts[0])
	price, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return "Invalid price. Please enter a valid number."
	}

	err = alerts.AddAlert(coin, price)
	if err != nil {
		return fmt.Sprintf("Error adding alert: %v", err)
	}

	return fmt.Sprintf("Alert added for %s at price $%.2f", coin, price)
}
