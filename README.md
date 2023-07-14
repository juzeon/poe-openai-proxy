# poe-openai-proxy

![GitHub repo size](https://img.shields.io/github/repo-size/caoyunzhou/poe-openai-proxy)
![Docker Image Size (tag)](https://img.shields.io/docker/image-size/caoyunzhou/poe-openai-proxy/latest)
![GitHub go.mod Go version (branch & subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/caoyunzhou/poe-openai-proxy/master)
![Docker Pulls](https://img.shields.io/docker/pulls/caoyunzhou/poe-openai-proxy)
[![GitHub Repo stars](https://img.shields.io/github/stars/caoyunzhou/poe-openai-proxy?style=social)](https://github.com/caoyunzhou/poe-openai-proxy/stargazers)

- [Poe.com](https://poe.com/) 是一个免费的网页应用，让你可以和 GPT 模型聊天。`poe-api` 它反向工程了 `poe.com`，让它可以通过一个 HTTP API 来访问，这个 API 模仿了官方的 OpenAI API for ChatGPT，所以它可以和其他使用 OpenAI API for ChatGPT 的程序兼容。

## 声明：

- 此项目只发布于 GitHub，基于 MIT 协议，免费且作为开源学习使用。并且不会有任何形式的卖号、付费服务、讨论群、讨论组等行为。谨防受骗。

## 安装

1. 将这个仓库克隆到你的本地机器：

```
git clone https://github.com/caoyunzhou/poe-openai-proxy.git
cd poe-openai-proxy/
```

2. 在项目的根目录创建配置文件。说明写在注释里：

```bash
cp .env.example .env
vim .env
source .env
```

3. 构建并启动Go后端：

```bash
go build
chmod +x poe-openai-proxy
./poe-openai-proxy
```

4. POE token 获取方式
```
- 在浏览器上登录`Poe.com`，然后打开浏览器的开发工具，在以下菜单中查找p-b cookie的值：
- Chromium：开发工具>应用程序>Cookie>poe.com
- Firefox：开发工具>存储>Cookie
- Safari:开发工具>存储>Cookie
- 类似于这样的值[ p-b : 12zNTxAdieuXXszMWYt93g%3D%3D]
```

### Docker一键部署

- docker run快速开始：
  ```
  docker run -d \
  --name poe-openai-proxy \
  -p 8080:8080 \
  -e PORT=8080 \
  -e AuthKey=sk-123456 \
  -e TOKENS=12zNTxAdieuXXszMWYt93g%3D%3D \
  -e SIMULATE_ROLES=0 \
  -e RATE_LIMIT=10 \
  -e COOL_DOWN=10 \
  -e TIMEOUT=60 \
  caoyunzhou/poe-openai-proxy
  ```

- 配置多个POE-Token示例:
  ```
  docker run -d \
  --name poe-openai-proxy \
  -p 8080:8080 \
  -e PORT=8080 \
  -e AuthKey=sk-123456 \
  -e TOKENS="12zNTxAdieuXXszMWYt93g%3D%3D,13zNTxAdieuXXszMWYt93g%3D%3D,14zNTxAdieuXXszMWYt93g%3D%3D" \
  -e SIMULATE_ROLES=0 \
  -e RATE_LIMIT=10 \
  -e COOL_DOWN=10 \
  -e TIMEOUT=60 \
  caoyunzhou/poe-openai-proxy
  ```

- docker-compose部署:
  ```
  docker-compose -f docker-compose.yaml up -d
  ```

- nginx反向代理自己配置


###  使用 Railway 部署

- [注册Railway](https://railway.app?referralCode=CG56Re)

- [![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/new/template/nFaU1x)

#### Railway 环境变量

| 环境变量名称          | 必填                   | 备注                                                                                               |
| --------------------- | ---------------------- | -------------------------------------------------------------------------------------------------- |
| `PORT`                | 必填        | 默认 `8080`
| `AuthKey`             | 必填        | 默认:`sk-123456` , `适配openai的请求，用户自定义的秘钥，放在Authorization header里面做认证` |
| `TOKENS`              | 必填[list]  | poe-Token密钥[token1,token2,token3]                                        |
| `SIMULATE_ROLES`      | `0`        | 角色                                                                       |
| `RATE_LIMIT`          | `10`       | 速率[默认为1分钟内每个令牌调用10个api]    |
| `COOL_DOWN`           | `10`       | 冷却令牌[#冷却几秒钟。同一个令牌在n秒内不能多次使用] |
| `TIMEOUT`             | `60`       | 超时  |

> 注意: `Railway` 修改环境变量会重新 `Deploy`

## 使用

参见[OpenAI文档](https://platform.openai.com/docs/api-reference/chat/create)了解更多关于如何使用ChatGPT API的细节。

只需要把你的代码里的`https://api.openai.com`替换成`http://localhost:8080`或者替换成你的域名`https://api.example.com`就可以了。

## IP参考访问方式[127.0.0.1替换成你的IP]
```
curl --location 'http://127.0.0.1:8080/v1/chat/completions' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer sk-123456' \
--data '{
  "model": "gpt-4",
  "messages": [{"role": "user", "content": "你好"}]
}'
```

## 配置域名参考访问方式
```
curl --location 'https://poe.aivvm.com/v1/chat/completions' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer sk-123456' \
--data '{
  "model": "gpt-4",
  "messages": [{"role": "user", "content": "你好"}]
}'
```

支持的路由：

- /models
- /chat/completions
- /v1/models
- /v1/chat/completions

支持的models对应列表:
| 模型名称                  | poe模型名称  |
| ------------------------ | ---------- |
| sage                     | capybara   |
| claude-instant           | a2         |
| claude-2-100k            | a2_2       |
| claude-instant-100k      | a2_100k    |
| gpt-3.5-turbo-0613       | chinchilla |
| gpt-3.5-turbo-16k-0613   | agouti     |
| gpt-4                    | beaver     |
| gpt-4-0613               | beaver     |
| gpt-4-32k                | vizcacha   |
| google-palm              | acouchy    |

支持的参数：

| 参数     | 说明                                                         |
| -------- | ------------------------------------------------------------ |
| model    | 参见`.env.example`里的`[bot]`部分。模型名字对应着机器人昵称。 |
| messages | 你可以像在官方API里一样使用这个参数，除了`name`。            |
| stream   | 你可以像在官方API里一样使用这个参数。                               |

其他参数会被忽略。


## License
- MIT © [poe-openai-proxy](./license)

## 致谢

- <https://github.com/ading2210/poe-api>
- <https://github.com/juzeon/poe-openai-proxy>
- <https://github.com/lwydyby/poe-api>

## 参与贡献

感谢所有做过贡献的人!

<a href="https://github.com/caoyunzhou/poe-openai-proxy/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=caoyunzhou/poe-openai-proxy" />
</a>



## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=caoyunzhou/poe-openai-proxy&type=Date)](https://star-history.com/#caoyunzhou/poe-openai-proxy&Date)
