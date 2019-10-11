import requests
import sys

def main():
    url = "https://notify-api.line.me/api/notify"
    token = "WcudzQXjoEgLad8EA68AkLe98Tl5mxEjbVhgOjdBIZH"
    headers = {"Authorization" : "Bearer "+ token}
    payload = {"message" :  "finish"}
    files = {"imageFile": open("test0.dat")}

    r = requests.post(url ,headers = headers ,params=payload)
    print(r)

if __name__ == '__main__':
    main()