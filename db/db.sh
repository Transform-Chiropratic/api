#!/bin/bash -e

function usage() {
  echo "Usage: $0 <operation> <database>"
  echo "-"
  echo "operation : [ create, drop, reset ]"
  exit 1
}

function create() {
  echo "CREATE USER jesseokeya SUPERUSER;" | psql -h127.0.0.1 || :
  cat <<EOF | psql -h127.0.0.1 -U jesseokeya
    CREATE USER transform_chiropratic WITH PASSWORD 'supersecret';
    CREATE DATABASE $database ENCODING 'UTF-8' LC_COLLATE='en_US.UTF-8' LC_CTYPE='en_US.UTF-8' TEMPLATE template0 OWNER transform_chiropratic;
EOF

  cat <<EOF | psql -h127.0.0.1 -U jesseokeya $database
    CREATE EXTENSION IF NOT EXISTS pg_trgm;
    CREATE EXTENSION IF NOT EXISTS plpgsql;
    GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO transform_chiropratic;
    GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO transform_chiropratic;
    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON TABLES TO transform_chiropratic;
    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON SEQUENCES TO transform_chiropratic;
EOF
}

function drop() {
  cat <<EOF | psql -h127.0.0.1 -U jesseokeya
    SELECT pg_terminate_backend(pg_stat_activity.pid)
    FROM pg_stat_activity
    WHERE pg_stat_activity.datname = '$database' AND pid <> pg_backend_pid();
    DROP DATABASE IF EXISTS $database;
EOF
}

if [ $# -lt 2 ]; then
  usage
fi

operation="$1"
database="$2"

case "$operation" in
  "create")
    create
    ;;
  "drop")
    drop
    ;;
  "reset")
    drop
    create
    ;;
  *)
    echo "no such operation"
    usage
    ;;
esac
