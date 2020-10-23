# Kademlia-lab-D7024E-G10
Lab in the course D7024E



## Run nodes

### Locally
Comment/uncomment code in main to run locally

go run main [port] - e.g. go run main 8080

go run main [port] [known port of other node] - e.g. go run main 8081 8080

### Docker
build: docker build --tag kademlia .

create network: docker network create mynetwork

to start on windows: powershell -ExecutionPolicy ByPass -File dockerrun.ps1

to close on windows: powershell -ExecutionPolicy ByPass -File dockerclose.ps1

to start on linux: ./dockerrun.sh

to close on linux: ./dockerclose.sh

read dockerrun.ps1/sh to see how a single node is started. -d flag is not needed then.

run all commands in top level folder of this project

.sh scripts might require chmod

## CLI

step1: docker exec -it NAME /bin/sh

step2: ./cliapp
