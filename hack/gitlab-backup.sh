#!/bin/bash
MX_ATTEMPTS=50
ATTEMPTS=1
SUCCESS=0
BACKUP_FILE="$3"
BACKUP=${BACKUP_FILE:0:31}
GITLAB_CONTAINER="$2"
COMMAND="$1"

function run_command() {
    status_code=$(curl -I --insecure https://localhost/ -ocurl -s -o /dev/null -I -w "%{http_code}" https://localhost/users/sign_in)
    if [[ "$status_code" -eq 302200 ]] ; then
        if [[ $COMMAND == "create" ]];then 
            echo -ne "\ngitlab data backup in progress..."
            docker exec -t $GITLAB_CONTAINER gitlab-backup create
            NEW_BACKUP=$(docker exec -t $GITLAB_CONTAINER sh -c "ls -t /var/opt/gitlab/backups/ | head -1" | tr -d '\r')
            cd ../srv/gitlab/
            docker cp $GITLAB_CONTAINER:/var/opt/gitlab/backups/$NEW_BACKUP backups/
            chmod a+rwx backups/$NEW_BACKUP
            cd ../..
            sed -i "s/$BACKUP_FILE/$NEW_BACKUP/1" Makefile
            cd hack/
            rm curl
            SUCCESS=1
        fi

        if [[ $COMMAND == "restore" ]];then
            echo -ne "\nchecking if gitlab service is running..."
            if [ -z `docker ps -q --no-trunc | grep $(docker-compose ps -q gitlab)` ]; then
                echo -ne "\nNo gitlab service running."            
            else
                echo -ne "\nRestoring gitlab backup with pre-configured repository test...\n"
                sleep 240
                echo -ne "\nstop gitlab puma service"
                docker exec -t $GITLAB_CONTAINER gitlab-ctl stop puma force=yes
                echo -ne "\nstop gitlab sidekiq service"
                docker exec -t $GITLAB_CONTAINER gitlab-ctl stop sidekiq force=yes
                echo -ne "\ncopy backup to gitlab container"
                docker cp ../srv/gitlab/backups/$BACKUP_FILE $GITLAB_CONTAINER:/var/opt/gitlab/backups/$BACKUP_FILE
                echo -ne "\nrestore gitlab backup"
                docker exec -t $GITLAB_CONTAINER gitlab-backup restore BACKUP=$BACKUP force=yes
                echo -ne "\nRestarting gitlab container..."
                docker restart -t 15 $GITLAB_CONTAINER
                # docker exec -it  $GITLAB_CONTAINER gitlab-rake gitlab:check SANITIZE=true
                SUCCESS=1
                sleep 120
            fi
        fi
    else
        echo -ne "\nwaiting gitlab service to start..."
        sleep 10
    fi
}

if [[ $COMMAND == "restore" ]];then
    echo -ne "\nRestoring gitlab environment..."
    echo -ne "\nDon't worry! It can take between 5 and 10 minutes until gitlab service start."
fi

while [[ "$SUCCESS" -eq 0 ]]
do
    echo -ne "\nattempt $ATTEMPTS: trying to reach gitlab service"
    run_command
    ATTEMPTS=$(($ATTEMPTS+1))
    if [[ $ATTEMPTS -gt $MX_ATTEMPTS ]]; then
    echo -ne "\nexceeded number of attemps: $MX_ATTEMPTS"
        break
    fi
done