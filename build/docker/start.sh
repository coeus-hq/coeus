#!/bin/sh

# Docker entrypoint script

# Configure and start postgres server as postgres user
su - postgres -c "chmod 0700 /var/lib/postgresql/data &&\
initdb /var/lib/postgresql/data &&\
echo \"host all  all    0.0.0.0/0  md5\" >> /var/lib/postgresql/data/pg_hba.conf &&\
echo \"listen_addresses='*'\" >> /var/lib/postgresql/data/postgresql.conf &&\
pg_ctl start -D /var/lib/postgresql/data &&\
psql -U postgres -c \"CREATE DATABASE coeus\" &&\
psql -U postgres -c \"CREATE USER coeus WITH PASSWORD 'coeus';\" &&\
psql -U postgres -c \"grant all privileges on database coeus to coeus;\""

# Install and setup bash
apk update \
&& apk add bash \
&& echo "" > /root/.ash_history
sed -i -e "s/bin\/ash/bin\/bash/" /etc/passwd

# Start coeus
/root/coeus-bin
