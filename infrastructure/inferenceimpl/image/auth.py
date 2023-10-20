from flask import Flask, request, Response
import requests
import os

app = Flask(__name__)

@app.route('/auth', methods=['GET'])
def auth():
    access_token = request.headers.get('X-Forwarded-Access-Token')
    if access_token is None:
        print("missing access token")
        return Response(status=401)

    user = request.headers.get('X-Forwarded-User')
    if user is None:
        print("missing user")
        return Response(status=401)

    auth_url = os.getenv("ACCESS_TOKEN_ENDPOINT")
    if auth_url is None:
        print("missing auth url")
        return Response(status=401)

    res = requests.get(auth_url, headers={"Authorization": f"{access_token}"})
    if res.status_code != 200:
        print("get user failed: %s %s" % (res.status_code, res.text))
        return Response(status=401)

    env_name = os.getenv("USER")
    if user != env_name:
        print(f"user name not match, {user} vs {env_name}")
        return Response(status=401)

    return Response(status=200)

if __name__ == '__main__':
    app.run(port=5000)