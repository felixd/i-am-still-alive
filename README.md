# I am still alive - Dead Person Switch - FlameIT - Immersion Cooling

I am still alive - Dead Man Switch Software

Send message to recipients after Your death/inactivity.

```bash
curl -X POST http://localhost:8080/signup -d '{"username":"user1", "password":"password123"}'
curl -X POST http://localhost:8080/login -d '{"username":"user1", "password":"password123"}'

TOKEN="TOKEN-RETURNED BY /login endpoint"

# Duration in hours (21 days * 24 hours -> 504 hours)
curl -X POST http://localhost:8080/switch/create -H "Authorization: $TOKEN" -d '{"duration": 504}'

# Switch timer update
curl -X GET http://localhost:8080/switch/checkin -H "Authorization: $TOKEN"

# Remove switch
curl -X DELETE http://localhost:8080/switch/delete -H "Authorization: $TOKEN"

```


## Other projects

* https://instantiator.dev/post/dead-person-switch/
* https://www.deadmansswitch.net/
* https://yankeguo.github.io/lastwill/

## Author

* Paweł 'felixd' Wojciechowski - FlameIT - Immersion Cooling (https://flameit.io)
