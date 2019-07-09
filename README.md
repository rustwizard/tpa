# Places Autocomplete Service

To build and start service you should do only
```
docker-compose build && docker-compose up
```

After that you could test it via curl

```
curl -vd '{"term":"Flo", "locale":"en", "types": ["city", "airport"]}' \ 
-H "Content-Type: application/json" \ 
-X POST http://127.0.0.1:8080
```

There is no tests due to lack of time for the exercise. If that's the issue I can write tests shortly.
