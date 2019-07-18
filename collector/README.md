# Collector

## Configuration

create a file: settings.json

```
{
	"name": "My Collector",
	"rconTimeout": "15s",
	"batchSize": 10,
	"commitInterval": "1m",
	"monitoring": {
		"enabled": true,
		"refreshRate": "30s",
		"host": ":8080"
	},
	"databaseConnection": "root:password@tcp(127.0.0.1:3306)/kagstats",
	"servers": [
		{
			"address": "127.0.0.1",
			"port": "50301",
			"password": "admin"
		}
	]
}
```
