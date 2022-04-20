#!/bin/bash
MX_ATTEMPTS=3
ATTEMPTS=0
SUCCESS=0
BACKUP_FILE="1650232071_2022_04_17_14.9.2-ee_gitlab_backup.tar"
GITLAB_CONTAINER="$2"
COMMAND="$1"

function run_command() {
    status_code=$(curl -I --insecure https://gitlab.integration-tests.com/ -ocurl -s -o /dev/null -I -w "%{http_code}" https://gitlab.integration-tests.com/users/sign_in)
    echo $COMMAND
    if [[ "$status_code" -eq 302200 ]] ; then
        if [[ $COMMAND = "create" ]];then 
            echo "gitlab data backup in progress..."
            docker exec -t $GITLAB_CONTAINER gitlab-backup create
            # docker cp $GITLAB_CONTAINER:/var/opt/gitlab/backups $HOME/
        fi

        if [[ $COMMAND = "restore" ]];then
            echo "checking if gitlab service is running..."
            if [ -z `docker ps -q --no-trunc | grep $(docker-compose ps -q web)` ]; then
                echo "No gitlab service running."
            else
                echo "Restoring gitlab backup with pre-configured repository test..."
                docker exec -it $GITLAB_CONTAINER gitlab-ctl stop puma
                docker exec -it $GITLAB_CONTAINER gitlab-ctl stop sidekiq
                docker cp ./srv/gitlab/backups/$BACKUP_FILE $GITLAB_CONTAINER:/var/opt/gitlab/backups/$BACKUP_FILE
                docker exec -it $GITLAB_CONTAINER gitlab-backup restore BACKUP=1650232071_2022_04_17_14.9.2-ee force=yes
                docker restart $GITLAB_CONTAINER
                # docker exec -it  $GITLAB_CONTAINER gitlab-rake gitlab:check SANITIZE=true
            fi
        fi
        SUCCESS=1
    else
        echo -ne "\nwaiting gitlab service to start"
        sleep 5
    fi
}

while [[ "$SUCCESS" -eq 0 ]]
do
    echo -ne "\nattempt: $ATTEMPTS"
    echo -ne "\ntrying to run command: $COMMAND"
    run_command
    ATTEMPTS=$(($ATTEMPTS+1))
    if [[ $ATTEMPTS -gt $MX_ATTEMPTS ]]; then
    echo -ne "\nexceeded number of attemps: $MX_ATTEMPTS"
        break
    fi
done
