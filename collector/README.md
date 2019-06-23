# Collector

## Configuration

create a file: settings.json

```
{
	"name": "My Collector",
	"timeout": "15s",
	"maxAttempts": 20,
	"bulkLoadMax": 3,
	"bulkLoadWait": "1m",
	"monitoring": true,
	"database": {
		"user": "root",
		"password": "password",
		"address": "127.0.0.1:3306"
	},
	"servers": [
		{
			"name": "TDM US",
			"tags": [
				"TDM",
				"US"
			],
			"address": "localhost:50301",
			"password": "password"
		},
		{
			"name": "TDM EU",
			"tags": [
				"TDM",
				"EU"
			],
			"address": "localhost:50302",
			"password": "password"
		}
	]
}
```
