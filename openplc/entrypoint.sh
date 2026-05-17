#!/bin/bash
set -e

# Initialize /docker_persistent with defaults (uses -n so won't overwrite existing)
if [ -d "/docker_persistent" ]; then
    mkdir -p /docker_persistent/st_files
    cp -n /workdir/webserver/dnp3_default.cfg /docker_persistent/dnp3.cfg 2>/dev/null || true
    cp -n /workdir/webserver/openplc_default.db /docker_persistent/openplc.db 2>/dev/null || true
    cp -n /workdir/webserver/active_program_default /docker_persistent/active_program 2>/dev/null || true
    cp -rn /workdir/webserver/st_files_default/. /docker_persistent/st_files/ 2>/dev/null || true
    touch /docker_persistent/mbconfig.cfg /docker_persistent/persistent.file 2>/dev/null || true
fi

# Always ensure furnace.st is registered and enabled
sqlite3 /docker_persistent/openplc.db \
    "INSERT OR IGNORE INTO Programs (Name, Description, File, Date_upload)
     VALUES ('Furnace','Industrial Furnace Control','furnace.st',strftime('%s','now'))"
sqlite3 /docker_persistent/openplc.db \
    "UPDATE Settings SET Value='true' WHERE Key='Start_run_mode'"
echo "furnace.st" > /docker_persistent/active_program
cp /workdir/webserver/st_files_default/furnace.st /docker_persistent/st_files/furnace.st

cd /workdir/webserver
exec "/workdir/.venv/bin/python3" webserver.py
