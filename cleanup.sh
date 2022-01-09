#!/bin/bash

CREATE_DOMAIN="_acme-challenge.${CERTBOT_DOMAIN}"
dnspod-ycli del ${CREATE_DOMAIN}