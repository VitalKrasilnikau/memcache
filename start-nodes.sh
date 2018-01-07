#!/bin/bash
for i in {0..10}
do
    ./memcache-node -port "$(($i + 9000))" -type "string" -index $i &
done
for i in {0..10}
do
    ./memcache-node -port "$(($i + 8000))" -type "list" -index $i &
done
for i in {0..10}
do
    ./memcache-node -port "$(($i + 10000))" -type "dictionary" -index $i &
done