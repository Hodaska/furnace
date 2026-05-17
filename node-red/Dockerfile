FROM nodered/node-red:latest

RUN npm install \
    node-red-contrib-modbus \
    node-red-contrib-iiot-opcua \
    node-red-contrib-influxdb

# Fix node-opcua-pki template: OpenSSL 3.x forbids authorityKeyIdentifier
# in v3_selfsigned section during CSR creation (no issuer cert present).
# Remove it from the template so certificate generation succeeds.
RUN sed -i 's/authorityKeyIdentifier.*=.*keyid,issuer/subjectKeyIdentifier = hash/' \
    /usr/src/node-red/node_modules/node-opcua-pki/dist/pki/templates/simple_config_template.cnf.js
