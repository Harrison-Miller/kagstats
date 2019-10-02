#!/bin/bash
echo "starting loadtest against staging.kagstats.com"
for i in {1..10000}
do
    echo "$i"
    curl -k "http://staging.kagstats.com/#/players"
    curl -k "http://staging.kagstats.com/#/servers"
    curl -k "http://staging.kagstats.com/#/kills"
    curl -k "http://staging.kagstats.com/api/players"
    curl -k "http://staging.kagstats.com/api/servers"
    curl -k "http://staging.kagstats.com/api/kills"
done