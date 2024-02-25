# Бот-расшифровщик голосовых сообщений

Телеграм бот, который расшифровывает голосовые сообщения в личных чатах и в группах с помощью Yandex Cloud speech-to-text.


### Использование

##### Подготовка
[Создайте вашего бота и получите для него токен.](https://core.telegram.org/bots/features#creating-a-new-bot)

[Разрашите вашему боту получать сообщения в групповых чатах.](https://core.telegram.org/bots/features#privacy-mode)

[Создайте сервисный аккаунт в yandex-cloud с необходимыми ролями и получите для него токен.](https://cloud.yandex.ru/ru/docs/speechkit/stt/api/stt-language-labels-example#preparations)

Добавить в переменные окружения `TELEGRAM_TOKEN` и `YANDEX_OAUTH` (можно записать их в .env в корне проекта).

Создайте в корне проекта файл `config.yml`.

```yaml
telegram:
    allow_chats: [-12345, -67890]  # id групп, в которых боту можно будет расшифровывать сообщения
    allow_users: [12345, 67890]  # id личных чатов, в которых боту можно будет расшифровывать сообщения
    alert_chat: 12345  # id личного чата или группы, в который будут приходить сообщения об ошибках
```


##### Локальный запуск
```bash
make update-requirements
make run
```


##### Запуск в контейнере
```bash
make run-service
```
