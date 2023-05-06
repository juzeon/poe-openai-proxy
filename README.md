# poe-openai-proxy

A wrapper that lets you use the reverse-engineered Python library `poe-api` as if it was the OpenAI API for ChatGPT. You can connect your favorite OpenAI API based apps to this proxy and enjoy the ChatGPT API for free!

[Poe.com](https://poe.com/) from Quora is a free web app that lets you chat with GPT models. `poe-api` is a Python library that reverse-engineered `poe.com` so you can use Python to call `poe`. This project is a wrapper around `poe-api` that makes it accessible through an HTTP API, which mimics the official OpenAI API for ChatGPT so it can work with other programs that use OpenAI API for their features.

[简体中文](README_zh.md)

## Installation

1. Clone this repository to your local machine:

```bash
git clone https://github.com/juzeon/poe-openai-proxy.git
cd poe-openai-proxy/
```

2. Install dependencies from requirements.txt:

```bash
pip install -r external/requirements.txt
```

3. Create the configuration file in the root folder of the project according to the instructions in the comments:

```bash
vim config.toml
```

config.toml:

```toml
# The port number for the proxy service. The proxied OpenAI API endpoint will be: http://localhost:3700/v1/chat/completions
port = 3700

# A list of poe tokens. You can get them from the cookies on poe.com, they look like this: p-b=fdasac5a1dfa6%3D%3D
tokens = ["fdasac5a1dfa6%3D%3D","d84ef53ad5f132sa%3D%3D"]

# The proxy that will be used to connect to poe.com. Leave it blank if you do not use a proxy
proxy = "socks5h://127.0.0.1:7890"

# The gateway url for the Python backend of poe-api. Don't change this unless you modify external/api.py
# Note that if you use docker this value should be changed into: http://external:5000
gateway = "http://127.0.0.1:5000"

# The bot name to use from poe. `capybara` stands for `Sage`
bot = "capybara"

# Use leading prompts to indicate roles if enabled. You'd better disable it if you are using tools like https://github.com/TheR1D/shell_gpt
# 0:disable, 1:enable, 2:auto detect
# Example: 
# ||>User:
# Hello!
# ||Assistant:
# Hello! How can I assist you today?
simulate-roles = 2

# Rate limit. Default to 10 api calls per token in 1 minute
rate-limit = 10

# Cool down of seconds. One same token cannot be used more than once in n seconds 
cool-down = 3

# Timeout of seconds per response message
timeout = 200
```

4. Start the Python backend for `poe-api`:

```bash
python external/api.py # Running on port 5000
```

5. Build and start the Go backend:

```bash
go build
chmod +x poe-openai-proxy
./poe-openai-proxy
```

### Docker support

If you would like to use docker, just run `docker-compose up -d` after creating `config.toml` according to the instructions above.

## Usage

See [OpenAI Document](https://platform.openai.com/docs/api-reference/chat/create) for more details on how to use the ChatGPT API.

Just replace `https://api.openai.com/v1/chat/completions` in your code with `http://localhost:3700/v1/chat/completions` and you're good to go.

Supported parameters:

| Parameter | Note                                                         |
| --------- | ------------------------------------------------------------ |
| model     | It doesn't matter what you pass here, poe will always use `gpt-3.5-turbo`. |
| messages  | You can use this as in the official API, except for `name`.            |
| stream    | You can use this as in the official API.                               |

Other parameters will be ignored.

## Credit

<https://github.com/ading2210/poe-api>