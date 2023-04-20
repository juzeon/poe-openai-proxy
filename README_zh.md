# poe-openai-proxy

这是一个包装器，让你可以使用反向工程的 Python 库 `poe-api` 作为 OpenAI API for ChatGPT 的接口。你可以将你喜欢的基于 OpenAI API 的应用程序连接到这个代理，免费享受 ChatGPT API 的功能！

[Poe.com](https://poe.com/) 是一个免费的网页应用，让你可以和 GPT 模型聊天。`poe-api` 是一个 Python 库，它反向工程了 `poe.com`，所以你可以用 Python 来调用 `poe`。这个项目是一个围绕 `poe-api` 的包装器，让它可以通过一个 HTTP API 来访问，这个 API 模仿了官方的 OpenAI API for ChatGPT，所以它可以和其他使用 OpenAI API for ChatGPT 的程序兼容。

## 安装

1. 将这个仓库克隆到你的本地机器：

```
git clone https://github.com/juzeon/poe-openai-proxy.git
cd poe-openai-proxy/
```

2. 安装 requirements.txt 中的依赖项：

```
pip install -r requirements.txt
```

3. 在项目根目录新建配置文件，并根据注释中的说明修改配置文件：

```
vim config.toml
```

config.toml:

```
# 代理服务的端口号。代理的 OpenAI API 端点为：http://localhost:3700/v1/chat/completions
port = 3700

# poe 令牌的列表。你可以从 poe.com 的 cookies 中获取它们，它们看起来像这样：p-b=fdasac5a1dfa6%3D%3D
tokens = ["fdasac5a1dfa6%3D%3D","d84ef53ad5f132sa%3D%3D"]

# poe-api 的 Python 后端的网关 url。除非你修改了 external/api.py，否则不要改变这个
gateway = "http://127.0.0.1:5000"
```

4. 启动 `poe-api` 的 Python 后端：

```
pip install -r requirements.txt
python external/api.py # 在端口 5000 上运行
```

5. 构建并启动 Go 后端：

```
go build
chmod +x poe-openai-proxy
./poe-openai-proxy
```

## 使用

查看 [OpenAI 文档](https://platform.openai.com/docs/api-reference/chat/create) 了解更多关于如何使用 ChatGPT API 的细节。

只需将你的代码中的 `https://api.openai.com/v1/chat/completions` 替换为 `http://localhost:3700/v1/chat/completions` 即可。

支持的参数：

| 参数     | 说明                                                 |
| -------- | ---------------------------------------------------- |
| model    | 你传递什么都无所谓，poe 总是使用 `gpt-3.5-turbo`。   |
| messages | 你可以像在官方 API 中一样使用这个参数，除了 `name`。 |
| stream   | 你可以像在官方 API 中一样使用这个参数。              |

其他参数将被忽略。

## 感谢

<https://github.com/ading2210/poe-api>