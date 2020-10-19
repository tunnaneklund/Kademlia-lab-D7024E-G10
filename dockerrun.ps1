#for($i = 0; $i -lt 50; $i++) {$port = 8080 + $i; docker run -d -p ${port}:80 --name "cont$i" --network mynetwork nginx:alpine}
docker run -d -p 8080:8080 --name "cont0" --network mynetwork kademlia /app/main "cont0"
for($i = 1; $i -lt 50; $i++) {
    $port = 8080 + $i;
    docker run -d -p ${port}:8080 --name "cont$i" --network mynetwork kademlia /app/main "cont$i" "cont0"
}