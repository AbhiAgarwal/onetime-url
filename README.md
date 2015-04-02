# onetime-request

A simple onetime URL generator. Created for a personal project. 

The webserver is ran by `go run main.go`. You have to set your SendGridUser and SendGridKey as environment variables. `export SendGridUser=___`, and `export SendGridKey=___`.

Onetime URL is sent to someone's email. To send it to someone's email set your environment variable OTUSecretKey to the secret key. For example: `export OTUSecretKey=a0sd0s9sd9sds91`, and then make a POST request to `http://localhost:3000/email` with the JSON being:

```
{
	"sk": "a0sd0s9sd9sds91",
	"email": "myemail@email.com"
}
```

You're able to pick your own secret key. Then a key will arrive in that email. You're able to go to `http://localhost:3000/_ID IN EMAIL HERE_`.