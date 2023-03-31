#!/usr/bin/bash
cp template_config.ini config.ini
sed -i '/^;/d' config.ini