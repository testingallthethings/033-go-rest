#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

wait_for_postgres() {
    local PSQL_ATTEMPTS
    local PSQL_ATTEMPTS_DELAY
    local PSQL_ATTEMPTS_MAX
    local PSQL_RET

    PSQL_ATTEMPTS=0
    PSQL_ATTEMPTS_DELAY=3
    PSQL_ATTEMPTS_MAX=10
    PSQL_RET=1

    echo "Waiting for postgres"

    until [ "${PSQL_ATTEMPTS}" -ge "${PSQL_ATTEMPTS_MAX}" ] || [ "${PSQL_RET}" -eq 0 ]; do
        set +e
        echo "Pinging Postgres"
        /usr/bin/psql "${DATABASE_URL}" -c "select 1"
        PSQL_RET=$?
        set -e

        if [ "${PSQL_RET}" -ne 0 ]; then
            echo "ERROR CODE: ${PSQL_RET}"
            echo "Could not ping Postgres"
            sleep $((PSQL_ATTEMPTS * PSQL_ATTEMPTS_DELAY))
        else
            echo "Postgres is up"
            return
        fi

        PSQL_ATTEMPTS=$((PSQL_ATTEMPTS+1))
    done

    return 1
}

main() {
    wait_for_postgres

    /http-service
}


main "$@"