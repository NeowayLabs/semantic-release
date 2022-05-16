GITLAB_CONTAINER_NAME="$1"
COMPOSE_COMMAND="$2"

if [ -z `docker ps -q --no-trunc | grep $(docker-compose ps -q gitlab)` ]; then
    echo -ne "\nStarting gitlab environment."
    cd ..
    GITLAB_CONTAINER_NAME=$GITLAB_CONTAINER_NAME $COMPOSE_COMMAND gitlab
    make gitlab-restore
else
    echo -ne "\nGitlab service is already up."
fi