FROM grafana/grafana:7.4.2

COPY grafana/grafana.ini ${GF_PATHS_CONFIG}
COPY grafana/testnet.json ${GF_PATHS_HOME}
COPY grafana/provisioning ${GF_PATHS_PROVISIONING}
COPY grafana/web/ ${GF_PATHS_HOME}
