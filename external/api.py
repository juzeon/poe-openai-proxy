import poe
from flask import Flask, request
from flask_sock import Sock
from poe import Client

poe.user_agent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"


def get_client(token) -> Client:
    try:
        print("Connecting to poe...")
        client_poe = poe.Client(token)
    except Exception as excp:
        print("Failed to connect to poe due to " + str(excp))
        exit(1)

    print("Connected to poe successfully")

    return client_poe


app = Flask(__name__)
sock = Sock(app)
sock.init_app(app)
client_dict = {}


@app.route('/add_token', methods=['GET', 'POST'])
def add_token():
    token = request.form['token']
    if token not in client_dict.keys():
        c = get_client(token)
        client_dict[token] = c
        return "ok"
    else:
        return "exist"


@app.route('/ask', methods=['GET', 'POST'])
def ask():
    token = request.form['token']
    bot = request.form['bot']
    content = request.form['content']
    for chunk in client_dict[token].send_message(bot, content, with_chat_break=True):
        pass
    return chunk["text"].strip()


@sock.route('/stream')
def stream(ws):
    token = ws.receive()
    bot = ws.receive()
    content = ws.receive()
    for chunk in client_dict[token].send_message(bot, content, with_chat_break=True):
        ws.send(chunk["text_new"])
    ws.close()


if __name__ == '__main__':
    app.run()
