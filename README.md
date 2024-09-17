# I am still alive - Dead Person Switch - FlameIT - Immersion Cooling

I am still alive - Dead Man Switch Software

Send message to recipients after Your death/inactivity.

```bash

HOST=192.168.1.25
PORT=8080

curl -X POST http://$HOST:$PORT/signup -d '{"username":"user1", "password":"password123"}'
curl -X POST http://$HOST:$PORT/login -d '{"username":"user1", "password":"password123"}'

TOKEN=$(curl -X POST http://$HOST:$PORT/login -d '{"username":"user1", "password":"password123"}' | jq -r '.token')

# Duration in hours (21 days * 24 hours -> 504 hours)
curl -X POST http://$HOST:$PORT/switch/create -H "Authorization: $TOKEN" -d '{"duration": 1, "message": "Wiadomość zza światów", "recipients": ["recipient1@test.net", "recipient2@test.net"]}'

# Switch timer update
curl -X GET http://$HOST:$PORT/switch/checkin -H "Authorization: $TOKEN"

# Remove switch
curl -X DELETE http://$HOST:$PORT/switch/delete -H "Authorization: $TOKEN"

```


## Other projects

* https://instantiator.dev/post/dead-person-switch/
* https://www.deadmansswitch.net/
* https://yankeguo.github.io/lastwill/

## Author

* Paweł 'felixd' Wojciechowski - FlameIT - Immersion Cooling (https://flameit.io)
