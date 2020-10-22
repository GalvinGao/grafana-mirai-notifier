package main

import (
	"github.com/Logiase/gomirai/bot"
	"github.com/Logiase/gomirai/message"
	"github.com/avast/retry-go"
	"github.com/sirupsen/logrus"
)

func dialBot() {
	if Bot != nil && Bot.FetchMessages() == nil {
		return
	}

	c := bot.NewClient("default", Conf.MiraiHTTP.Address, Conf.MiraiHTTP.AuthKey)
	c.Logger.Level = logrus.TraceLevel
	key, err := c.Auth()
	if err != nil {
		c.Logger.Fatal(err)
	}
	Bot, err = c.Verify(Conf.MiraiHTTP.QQNumber, key)
	if err != nil {
		c.Logger.Fatal(err)
	}
}

func botSendGroupMessage(group, quote uint, msg ...message.Message) error {
	return retry.Do(
		func() error {
			_, err := Bot.SendGroupMessage(group, quote, msg...)
			if err != nil {
				Log.Println("failed to send message due to", err, ". retrying...")
			}
			return err
		},
		retry.OnRetry(func(n uint, err error) {
			dialBot()
		}),
		retry.Attempts(5),
		retry.DelayType(retry.RandomDelay),
	)
}
