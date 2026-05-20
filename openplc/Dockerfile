FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=Europe/Moscow

RUN apt-get update && apt-get install -y git dos2unix sqlite3 && rm -rf /var/lib/apt/lists/*

RUN git clone https://github.com/thiagoralves/OpenPLC_v3.git /workdir

WORKDIR /workdir

RUN find /workdir -type f -name "*.sh" -exec dos2unix {} \; 2>/dev/null || true \
    && find /workdir -type f -name "*.py" -exec dos2unix {} \; 2>/dev/null || true

RUN mkdir /docker_persistent

VOLUME /docker_persistent

RUN ./install.sh docker \
    && touch /docker_persistent/mbconfig.cfg \
    && touch /docker_persistent/persistent.file \
    && mkdir -p /docker_persistent/st_files \
    && cp /workdir/webserver/openplc.db /docker_persistent/openplc.db \
    && mv /workdir/webserver/openplc.db /workdir/webserver/openplc_default.db \
    && cp /workdir/webserver/dnp3.cfg /docker_persistent/dnp3.cfg \
    && mv /workdir/webserver/dnp3.cfg /workdir/webserver/dnp3_default.cfg \
    && cp -r /workdir/webserver/st_files/ /docker_persistent/st_files/ \
    && mv /workdir/webserver/st_files /workdir/webserver/st_files_default \
    && cp /workdir/webserver/active_program /docker_persistent/active_program \
    && mv /workdir/webserver/active_program /workdir/webserver/active_program_default \
    && ln -s /docker_persistent/mbconfig.cfg /workdir/webserver/mbconfig.cfg \
    && ln -s /docker_persistent/persistent.file /workdir/webserver/persistent.file \
    && ln -s /docker_persistent/openplc.db /workdir/webserver/openplc.db \
    && ln -s /docker_persistent/dnp3.cfg /workdir/webserver/dnp3.cfg \
    && ln -s /docker_persistent/st_files /workdir/webserver/st_files \
    && ln -s /docker_persistent/active_program /workdir/webserver/active_program

# Patch openplc.py: retry _rpc up to 10 times so Modbus enable survives slow runtime start
COPY patch_openplc.py /tmp/patch_openplc.py
RUN python3 /tmp/patch_openplc.py

# Add furnace.st to the default (non-volume) files and compile it
COPY furnace.st /workdir/webserver/st_files_default/furnace.st

RUN echo "furnace.st" > /workdir/webserver/active_program_default \
    && sqlite3 /workdir/webserver/openplc_default.db \
       "INSERT OR IGNORE INTO Programs (Name, Description, File, Date_upload) \
        VALUES ('Furnace','Industrial Furnace Control','furnace.st',strftime('%s','now'))" \
    && sqlite3 /workdir/webserver/openplc_default.db \
       "UPDATE Settings SET Value='true' WHERE Key='Start_run_mode'" \
    && cp /workdir/webserver/st_files_default/furnace.st /workdir/webserver/core/furnace.st \
    && cp /workdir/webserver/st_files_default/furnace.st /docker_persistent/st_files/furnace.st \
    && cd /workdir/webserver && ./scripts/compile_program.sh furnace.st \
    && echo "Furnace compiled OK"

# Custom entrypoint guarantees correct DB state at every startup
COPY entrypoint.sh /workdir/entrypoint.sh
RUN dos2unix /workdir/entrypoint.sh && chmod +x /workdir/entrypoint.sh

EXPOSE 8080 502

ENTRYPOINT ["/workdir/entrypoint.sh"]
