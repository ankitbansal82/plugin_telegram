![Abot Telegram Plugin](https://telegram.org/img/t_logo.png)
![Abot](http://i.imgur.com/WBACSyP.png)

A plugin for [Abot](https://github.com/itsabot/abot) platform that integrates it with [Telegram bot api](https://core.telegram.org/bots)

## Integrating with Telegram

Abot makes it easy to add support for multiple communication tools, including SMS, phone, email, Slack, etc. In this guide, we'll learn how to set up Telegram bot, so we can communicate with this new digital assistant via Telegram.

First we will need to create a [Telegram bot](https://core.telegram.org/bots#6-botfather). Take a note of your API token. 
 
You'll want to set the following environment variables in your ~/.bash_profile or ~/.bashrc:

```bash
export TELEGRAM_API_KEY="REPLACE"
```

Telegram bot api sends updates to you using either long polling or webhook. If you interested in using webhook then you will need a publicly exposed web host, a domain name, public and private keys for SSL certificate. The plugin takes care of uploading the public certificate to Telegram and setting up the webhook.

Guide on how to generate certificate keys can be found [here](https://core.telegram.org/bots/self-signed)

You'll also want to set the following environment variables in your ~/.bash_profile or ~/.bashrc:

```bash
export TELEGRAM_API_KEY="REPLACE"
export TELEGRAM_USE_WEBHOOK="true"
export TELEGRAM_WEBHOOK_HOST="REPLACE"
export TELEGRAM_PRIVATE_KEY="private.key"
export TELEGRAM_PUBLIC_KEY="public.pem"
```

Now we'll add the Telegram driver. Since this is a plugin just like any other Abot plugin, you can simply add it to your plugins.json like so:

```json
{
    "Name": "abot",
    "Version": "0.2",
    "Dependencies": {
        "github.com/itsabot/plugin_onboard": "*",
        "github.com/ankitbansal82/plugin_telegram": "*"
    }
}
```

Then from your terminal run:

```
$ abot plugin install
Installing 2 plugins...
Success!
```

### Testing it out

To try it out, start abot server. And send a message from Telegram to your Telegram Bot. 

```
> Say hi
Hello World!
```