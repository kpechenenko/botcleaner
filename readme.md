# Botcleaner

Telegram bot for removing duplicate forwarded message from channels.

## Demo

I use a personal chat in the demo, but the bot also works for group chat when different people forward same messages.

![demo.gif](assets/demo.gif)

## Launch

1. Configuration

    ```bash
    cp env-example .env
    # fill in variables in .env
    ```

2. Start

    ```bash
    make run-img
    ```