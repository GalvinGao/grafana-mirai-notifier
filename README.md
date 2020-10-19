# Grafana QQ Notifier
[![Project Status](https://img.shields.io/badge/status-not%20finished%20yet-red)](#)
[![Powered by Mirai](https://img.shields.io/badge/powered%20by-Mirai-%237bbfb9)](https://github.com/mamoe/mirai)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/GalvinGao/grafana-mirai-notifier)](#)


This notifier allows the one to send notify message to a QQ group when a Grafana alert is triggered, by using the builtin webhook support that Grafana provides.

This project is based on [mamoe/mirai](https://github.com/mamoe/mirai), [mamoe/mirai-console](https://github.com/mamoe/mirai-console), [LXY1226/MiraiOK](https://github.com/LXY1226/MiraiOK), [project-mirai/mirai-api-http](https://github.com/project-mirai/mirai-api-http) and [Logiase/gomirai](https://github.com/Logiase/gomirai) 

## Get Started
1. [Setup mirai-console](https://github.com/LXY1226/MiraiOK)
2. Move `config.example.yml` to `config.yml`
3. Modify content of `config.yml` to correspond your settings
4. Build and run this project by using `go run .`
5. Go to Grafana and head to `https://YOUR_GRAFANA_URL/alerting/notification/new` and configure a New Notification Channel by using the config as following:
  - Name: (name of your choice)
  - Type: `webhook`
  - Url: `http://localhost:PORT_SET_IN_CONFIG/webhook`
  - Method: `POST`
  - Username and Password: _(leave as blank)_
6. Go to corresponding alert rules and add the notification channel you've just created

## License
This project is using `AGPLv3` as stated in its `LICENSE` file. If conflict occurs between those two statements, the `LICENSE` file should supersede.
