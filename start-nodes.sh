#!/bin/bash
for i in {0..10}
do
    ./memcache-node -port "$(($i + 100))" -type "string" -index $i &
done
for i in {0..10}
do
    ./memcache-node -port "$(($i + 1))" -type "list" -index $i &
done
for i in {0..10}
do
    ./memcache-node -port "$(($i + 1000))" -type "dictionary" -index $i &
done