docker image build -f Dockerfile -t forum-img .
docker container run -p 8080:8080 --detach --name forum-container forum-img
