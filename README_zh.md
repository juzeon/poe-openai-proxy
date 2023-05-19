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
pip install -r external/requirements.txt
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

# 连接到 poe.com 所用的代理，留空代表不使用代理
proxy = "socks5h://127.0.0.1:7890"

# poe-api 的 Python 后端的网关 url
# 注意，如果使用了docker，请将此值修改为: http://external:5000
gateway = "http://127.0.0.1:5000"

# poe-api 的 Python 后端的网关端口
# 必须和上面 gateway 中指定的端口一致
gateway-port = 5000

# poe 上 bot 的名字，capybara 即 Sage
bot = "capybara"

# 若设置为true，将用前缀prompt来管理角色。如果您使用类似于 https://github.com/TheR1D/shell_gpt 的工具，最好禁用它
# 0:启用, 1:禁用, 2:自动识别
# 例如：
# ||>User:
# Hello!
# ||Assistant:
# Hello! How can I assist you today?
simulate-roles = 2

# API调用速率限制，默认对于每个 token，每分钟最多请求10次
rate-limit = 10

# API调用冷却时间，同一个 token 在 n 秒内不能被重复使用
cool-down = 3

# 单个对话请求超时时间秒数
timeout = 200
```

4. 启动 `poe-api` 的 Python 后端：

```
python external/api.py # 在端口 5000 上运行
```

5. 构建并启动 Go 后端：

```
go build
chmod +x poe-openai-proxy
./poe-openai-proxy
```

### Docker 支持

如果要使用 Docker，只需要在按照上面的步骤创建`config.toml`后运行`docker-compose up -d`即可。

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