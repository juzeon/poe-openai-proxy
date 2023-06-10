# poe-openai-proxy

这是一个包装器，让你可以使用反向工程的 Python 库 `poe-api` 作为 OpenAI API for ChatGPT 的接口。你可以将你喜欢的基于 OpenAI API 的应用程序连接到这个代理，免费享受 ChatGPT API 的功能！

[Poe.com](https://poe.com/) 是一个免费的网页应用，让你可以和 GPT 模型聊天。`poe-api` 是一个 Python 库，它反向工程了 `poe.com`，所以你可以用 Python 来调用 `poe`。这个项目是一个围绕 `poe-api` 的包装器，让它可以通过一个 HTTP API 来访问，这个 API 模仿了官方的 OpenAI API for ChatGPT，所以它可以和其他使用 OpenAI API for ChatGPT 的程序兼容。

## 安装

1. 将这个仓库克隆到你的本地机器：

```
git clone https://github.com/juzeon/poe-openai-proxy.git
cd poe-openai-proxy/
```

2. 从requirements.txt安装依赖：

```bash
pip install -r external/requirements.txt
```

3. 在项目的根目录创建配置文件。说明写在注释里：

```bash
cp config.example.toml comfig.toml
vim config.toml
```

4. 启动`poe-api`的Python后端：

```bash
python external/api.py # 运行在5100端口上
```

5. 构建并启动Go后端：

```bash
go build
chmod +x poe-openai-proxy
./poe-openai-proxy
```

### Docker支持

如果你想使用docker，只需要在按照上面的说明创建好`config.toml`之后运行`docker-compose up -d`即可。

## 使用

参见[OpenAI文档](https://platform.openai.com/docs/api-reference/chat/create)了解更多关于如何使用ChatGPT API的细节。

只需要把你的代码里的`https://api.openai.com`替换成`http://localhost:3700`就可以了。

支持的路由：

- /models
- /chat/completions
- /v1/models
- /v1/chat/completions

支持的参数：

| 参数     | 说明                                                         |
| -------- | ------------------------------------------------------------ |
| model    | 参见`config.example.toml`里的`[bot]`部分。模型名字对应着机器人昵称。 |
| messages | 你可以像在官方API里一样使用这个参数，除了`name`。            |
| stream   | 你可以像在官方API里一样使用这个参数。                               |

其他参数会被忽略。

## 致谢

<https://github.com/ading2210/poe-api>