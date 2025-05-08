FROM postgres:15

ENV POSTGRES_USER=admin
ENV POSTGRES_PASSWORD=admin
ENV POSTGRES_DB=test

EXPOSE 5434

# docker run -d --name test -p 5434:5432 test
