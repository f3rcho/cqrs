FROM postgres:12

COPY up.sql /docker-entrypoint-initdb.d
EXPOSE 5432
CMD [ "postgres" ]