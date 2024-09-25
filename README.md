# I am still alive - Dead Person Switch - FlameIT - Immersion Cooling (FIT DPS)

I am still alive - Dead Person Switch Software

Send message to recipients after Your death/inactivity.

```bash
HOST=localhost
HOST=$(who am i | awk '{print $6}' | tr -d '()')
PORT=8080
USER=user1
PASSWORD=password123
MESSAGE="Very Secret Message"

curl -X POST http://$HOST:$PORT/signup -d '{"username":"'$USER'", "password":"'$PASSWORD'"}'
curl -X POST http://$HOST:$PORT/login -d '{"username":"'$USER'", "password":"'$PASSWORD'"}'

TOKEN=$(curl -X POST http://$HOST:$PORT/login -d '{"username":"'$USER'", "password":"'$PASSWORD'"}' | jq -r '.token')

# Duration in hours (21 days * 24 hours -> 504 hours)
curl -X POST http://$HOST:$PORT/switch/create -H "Authorization: $TOKEN" -d '{"duration": 1, "message": "'$MESSAGE'", "recipients": ["recipient1@test.net", "recipient2@test.net"]}'

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
