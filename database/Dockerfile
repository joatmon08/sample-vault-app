FROM postgres:14.4

ENV POSTGRES_DB test

COPY setup.sql /docker-entrypoint-initdb.d/
RUN chmod a+r /docker-entrypoint-initdb.d/*