FROM grafana/grafana:7.4.2

COPY grafana/grafana.ini ${GF_PATHS_CONFIG}
COPY grafana/testnet.json ${GF_PATHS_HOME}
COPY grafana/provisioning ${GF_PATHS_PROVISIONING}
COPY grafana/grafana_icon.svg ${GF_PATHS_HOME}/public/img/grafana_icon.svg

