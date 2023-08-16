import requests

endpoint = "http://127.0.0.1:9001/api/"

token = "myToken"

json_test = {
    "Info1": "Data",
}

response = requests.post(endpoint, json=json_test, headers={"Authorization": "Bearer {}".format(token)})


print("\nPost Request:")
print(json_test)

print("Status code: ", response.status_code)
print("\nPost Response:")
print(response.json())
