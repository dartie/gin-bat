T=myToken

#curl -X POST $E -H "Content-Type: application/json" -H "Authorization: Bearer $T" --data "{\"Info\": [\"data1\"]}"
#curl -X POST $E -H "Content-Type: application/json" -H "Authorization: Bearer $T" --data '{"Info": ["data1"]}'
#curl -X POST $E -H "Authorization: Bearer $T" --data "Info1=data1"

curl http://127.0.0.1:9003/api/test -H "Authorization: Bearer ${T}"
