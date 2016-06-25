package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/mrd0ll4r/tbotapi"
	"github.com/itsabot/abot/core"
	"github.com/itsabot/abot/core/log"
	"github.com/itsabot/abot/shared/datatypes"
	"github.com/itsabot/abot/shared/interface/messenger"
	"github.com/itsabot/abot/shared/interface/messenger/driver"
	"github.com/julienschmidt/httprouter"
)

type drv struct{}

func (d *drv) Open(r *httprouter.Router) (driver.Conn, error) {
	var err error
	var api *tbotapi.TelegramBotAPI

	apiKey := os.Getenv("TELEGRAM_API_KEY")	
	if os.Getenv("TELEGRAM_USE_WEBHOOK") == "true"{
		api, err = openWithWebhook(apiKey)
	} else{
		api, err = openWithLongPolling(apiKey)
	}
	if err != nil {
		log.Info(err)
		return nil, err
	}
	c := conn(*api)
	
	// To confirm api connection is working.
	log.Info("User ID: ", c.ID)
	log.Info("Bot Name: ", c.Name)
	log.Info("Bot Username: ", c.Username)	
	
	go func() {
		for {
			update := <-api.Updates
			if update.Error() != nil {
				log.Info("Update error: ", update.Error())
				continue
			}
				handlerTelegram(update.Update(), api)
		}
	}()
	return &c, nil
}

type conn tbotapi.TelegramBotAPI

// Send an message using a TElegram api to a specific open chat
func (c *conn) Send(chatID int, msg string) error {
	return nil
}

// Close the telegram api connection
func (c *conn) Close() error {
	c.Close()
	return nil
}

func init() {
	messenger.Register("telegram", &drv{})
}

func handlerTelegram(update tbotapi.Update, c *tbotapi.TelegramBotAPI) {
	var ret string
	switch update.Type() {
	case tbotapi.MessageUpdate:
		msg := update.Message
		typ := msg.Type()
		if typ != tbotapi.TextMessage {
			// Ignore non-text messages for now.
			log.Info("Ignoring non-text message")
			return
		}
		// Note: Bots cannot receive from channels, at least no text messages. So we don't have to distinguish anything here.
		log.Info("From User: ", msg.Chat)
		log.Info("From USer ID: ", msg.From.ID)
		log.Info("Text: ", *msg.Text)

		tmp := struct {
			CMD        string
			FlexID     string
			FlexIDType dt.FlexIDType
		}{
			CMD:        *msg.Text,
			FlexID:     string(msg.From.ID),
			FlexIDType: 3,
		}
		byt, err := json.Marshal(tmp)
		if err != nil {
			log.Info("failed marshaling req struct.", err)
			ret = "Something went wrong... Please contact support."
		}
		r, err := http.NewRequest("POST", "http://localhost", bytes.NewBuffer(byt))
		if err != nil {
			log.Info("failed building http request.", err)
			ret = "Something went wrong... Please contact support."
		}
		ret, err = core.ProcessText(r)
		if err != nil {
			log.Info("failed processing text.", err)
			ret = "Something went wrong... Please contact support."
		}

		_, err = c.NewOutgoingMessage(tbotapi.NewRecipientFromChat(msg.Chat), ret).Send()

		if err != nil {
			log.Info("Failed sending response to telegram: ", err)
			return
		}
	case tbotapi.InlineQueryUpdate:
		log.Info("Ignoring received inline query: ", update.InlineQuery.Query)
	case tbotapi.ChosenInlineResultUpdate:
		log.Info("Ignoring chosen inline query result (ID): ", update.ChosenInlineResult.ID)
	default:
		log.Info("Ignoring unknown Update type.")
	}
}

func openWithLongPolling(apiKey string)(*tbotapi.TelegramBotAPI, error){
	api,err := tbotapi.New(apiKey)
	if err != nil {
		log.Info(err)
		return nil, err
	}
	return api, err
}

func openWithWebhook(apiKey string)(*tbotapi.TelegramBotAPI, error){
	webhookHost := os.Getenv("TELEGRAM_WEBHOOK_HOST")
	webhookPort := uint16(8443)
	privkey := os.Getenv("TELEGRAM_PRIVATE_KEY")
	pubkey := os.Getenv("TELEGRAM_PUBLIC_KEY")
	u := url.URL{
	Host:   webhookHost + ":" + fmt.Sprint(webhookPort),
	Scheme: "https",
	Path:   apiKey,
	}
	api, handler, err := tbotapi.NewWithWebhook(apiKey, u.String(), pubkey)
	if err != nil {
		log.Info(err)
		return nil, err
	}
	http.HandleFunc("/"+apiKey, handler)
	log.Info("Starting Telegram webhook ", u.String())
	go func() {
		log.Fatal(http.ListenAndServeTLS("0.0.0.0:"+fmt.Sprint(webhookPort), pubkey, privkey, nil))
	}()
	return api, err
}