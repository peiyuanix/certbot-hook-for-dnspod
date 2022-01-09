#!/bin/bash

CREATE_DOMAIN="_acme-challenge.${CERTBOT_DOMAIN}"
dnspod-ycli add ${CREATE_DOMAIN} ${CERTBOT_VALIDATION}
sleep 10