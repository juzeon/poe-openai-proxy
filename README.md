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

3. Create the configuration file in the root folder of the project. Instructions are written in the comments:

```bash
cp config.example.toml comfig.toml
vim config.toml
```

4. Start the Python backend for `poe-api`:

```bash
python external/api.py # Running on port 5100
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

Just replace `https://api.openai.com` in your code with `http://localhost:3700` and you're good to go.

Supported routes:

- /models
- /chat/completions
- /v1/models
- /v1/chat/completions

Supported parameters:

| Parameter | Note                                                         |
| --------- | ------------------------------------------------------------ |
| model     | See `[bot]` section of `config.example.toml`. Model names are mapped to bot nicknames. |
| messages  | You can use this as in the official API, except for `name`.            |
| stream    | You can use this as in the official API.                               |

Other parameters will be ignored.

## Credit

<https://github.com/ading2210/poe-api>