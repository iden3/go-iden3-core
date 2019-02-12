# Backup Server

## Backup Service
```
Wallet                 Backup Service
  +                         +
  |       /register         |
  +------------------------>+
  |        200 OK           |
  +<------------------------+
  |                         |
  |                         |
  |      /backup/upload     |
  +------------------------>+
  |        200 OK           |
  +<------------------------+
  |                         |
  |                         |
  |      /backup/download   |
  +------------------------>+
  |        {backup}         |
  +<------------------------+
  |                         |
  +                         +
```

### Endpoints
- POST /register
	- in:
	```js
	{
	    username: "",
	    password: ""
	}
	```
	- out:
	```
	200 OK
	```

- POST /backup/upload
	- in:
	```js
	{
	    username: "",
	    password: ""
	    backup: "base64"
	}
	```
	- out:
	```
	200 OK
	```

- POST /backup/download
	- in:
	```js
	{
	    username: "",
	    password: ""
	}
	```
	- out:
	```js
	{
	    backup: "base64"
	}
	```

- ERROR
	- out:
	```js
	{
	    error: "msg"
	}
	```
