#!/bin/bash
docker run -d -p 8080:8080 --name "cont0" --network mynetwork kademlia /app/main "cont0"
for i in {1..50}
    do
        let PORT=8080+$i
        docker run -d -p ${PORT}:8080 --name "cont$i" --network mynetwork kademlia /app/main "cont$i" "cont0"
    done
