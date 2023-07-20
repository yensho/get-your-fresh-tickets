FROM postgres:15.2-alpine

RUN mkdir -p /tmp/psql_data/

COPY tables.sql /tmp/psql_data/
COPY db-startup.sh /docker-entrypoint-initdb.d/

RUN chmod 777 /tmp