#!/bin/bash
for i in {0..10}
do
    ./memcache-node -port "$(($i + 59000))" -type "string" -index $i &
done
for i in {0..10}
do
    ./memcache-node -port "$(($i + 58000))" -type "list" -index $i &
done
for i in {0..10}
do
    ./memcache-node -port "$(($i + 60000))" -type "dictionary" -index $i &
done