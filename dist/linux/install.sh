#!/bin/bash

# Check OS and VERSION
OS=$(cat /etc/os-release | grep ^ID= | cut -d'=' -f2)
temp="${OS%\"}"
temp="${temp#\"}"
OS="$temp"

# READ .env file
echo PWD IS $(pwd)
if [ -f ~/sqvs.env ]; then
    echo Reading Installation options from `realpath ~/sqvs.env`
    env_file=~/sqvs.env
elif [ -f ../sqvs.env ]; then
    echo Reading Installation options from `realpath ../sqvs.env`
    env_file=../sqvs.env
fi

if [ -n $env_file ]; then
    source $env_file
    env_file_exports=$(cat $env_file | grep -E '^[A-Z0-9_]+\s*=' | cut -d = -f 1)
    if [ -n "$env_file_exports" ]; then eval export $env_file_exports; fi
else
    echo No .env file found
    SQVS_NOSETUP="true"
fi

SERVICE_USERNAME=sqvs
COMPONENT_NAME=sqvs
  
if [[ $EUID -ne 0 ]]; then
    echo "This installer must be run as root"
    exit 1
fi

# Upgrade if component is already installed
if command -v $COMPONENT_NAME &>/dev/null; then
  echo "$COMPONENT_NAME is installed, proceeding with the upgrade"
  ./${COMPONENT_NAME}_upgrade.sh
  exit $?
fi

echo "Setting up SGX Quote Verification Service Linux User..."
id -u $SERVICE_USERNAME 2> /dev/null || useradd --shell /bin/false $SERVICE_USERNAME

echo "Installing SGX Quote Verification Service..."

PRODUCT_HOME=/opt/$COMPONENT_NAME
BIN_PATH=$PRODUCT_HOME/bin
LOG_PATH=/var/log/$COMPONENT_NAME/
CONFIG_PATH=/etc/$COMPONENT_NAME/
CERTS_PATH=$CONFIG_PATH/certs
CERTDIR_TRUSTEDJWTCERTS=$CERTS_PATH/trustedjwt
CERTDIR_TRUSTEDJWTCAS=$CERTS_PATH/trustedca

for directory in $BIN_PATH $LOG_PATH $CONFIG_PATH $CERTS_PATH $CERTDIR_TRUSTEDJWTCERTS $CERTDIR_TRUSTEDJWTCAS; do
  # mkdir -p will return 0 if directory exists or is a symlink to an existing directory or directory and parents can be created
  mkdir -p $directory
  if [ $? -ne 0 ]; then
    echo "Cannot create directory: $directory"
    exit 1
  fi
  chown -R $SERVICE_USERNAME:$SERVICE_USERNAME $directory
  chmod 700 $directory
done

cp $COMPONENT_NAME $BIN_PATH/ && chown $SERVICE_USERNAME:$SERVICE_USERNAME $BIN_PATH/*
chmod 700 $BIN_PATH/*
ln -sfT $BIN_PATH/$COMPONENT_NAME /usr/bin/$COMPONENT_NAME

# Create logging dir in /var/log
mkdir -p $LOG_PATH && chown sqvs:sqvs $LOG_PATH
chmod 640 $LOG_PATH

# Install systemd script
cp sqvs.service $PRODUCT_HOME && chown $SERVICE_USERNAME:$SERVICE_USERNAME $PRODUCT_HOME/sqvs.service && chown $SERVICE_USERNAME:$SERVICE_USERNAME $PRODUCT_HOME

# Enable systemd service
systemctl disable sqvs.service > /dev/null 2>&1
systemctl enable $PRODUCT_HOME/sqvs.service
systemctl daemon-reload

#Install log rotation
auto_install() {
  local component=${1}
  local cprefix=${2}
  local packages=$(eval "echo \$${cprefix}_PACKAGES")
if [ "$OS" == "rhel" ]; then
  dnf -y install $packages
elif [ "$OS" == "ubuntu" ]; then
  apt -y install $packages
fi
}

# SCRIPT EXECUTION
logRotate_clear() {
  logrotate=""
}

logRotate_detect() {
  local logrotaterc=`ls -1 /etc/logrotate.conf 2>/dev/null | tail -n 1`
  logrotate=`which logrotate 2>/dev/null`
  if [ -z "$logrotate" ] && [ -f "/usr/sbin/logrotate" ]; then
    logrotate="/usr/sbin/logrotate"
  fi
}

logRotate_install() {
  LOGROTATE_PACKAGES="logrotate"
  if [ "$(whoami)" == "root" ]; then
    auto_install "Log Rotate" "LOGROTATE"
    if [ $? -ne 0 ]; then echo "Failed to install logrotate"; exit -1; fi
  fi
  logRotate_clear; logRotate_detect;
    if [ -z "$logrotate" ]; then
      echo "logrotate is not installed"
    else
      echo  "logrotate installed in $logrotate"
    fi
}

logRotate_install

export LOG_ROTATION_PERIOD=${LOG_ROTATION_PERIOD:-weekly}
export LOG_COMPRESS=${LOG_COMPRESS:-compress}
export LOG_DELAYCOMPRESS=${LOG_DELAYCOMPRESS:-delaycompress}
export LOG_COPYTRUNCATE=${LOG_COPYTRUNCATE:-copytruncate}
export LOG_SIZE=${LOG_SIZE:-100M}
export LOG_OLD=${LOG_OLD:-12}

mkdir -p /etc/logrotate.d

if [ ! -a /etc/logrotate.d/sqvs ]; then
 echo "/var/log/sqvs/*.log {
    missingok
    notifempty
    rotate $LOG_OLD
    maxsize $LOG_SIZE
    nodateext
    $LOG_ROTATION_PERIOD
    $LOG_COMPRESS
    $LOG_DELAYCOMPRESS
    $LOG_COPYTRUNCATE
}" > /etc/logrotate.d/sqvs
fi

# check if SQVS_NOSETUP is defined
if [ "${SQVS_NOSETUP,,}" == "true" ]; then
    echo "SQVS_NOSETUP is true, skipping setup"
    echo "Installation completed successfully!"
else
    $COMPONENT_NAME setup all
    SETUPRESULT=$?
    if [ ${SETUPRESULT} == 0 ]; then
        systemctl start $COMPONENT_NAME
        echo "Waiting for daemon to settle down before checking status"
        sleep 3
        systemctl status $COMPONENT_NAME 2>&1 > /dev/null
        if [ $? != 0 ]; then
            echo "Installation completed with Errors - $COMPONENT_NAME daemon not started."
            echo "Please check errors in syslog using \`journalctl -u $COMPONENT_NAME\`"
            exit 1
        fi
        echo "$COMPONENT_NAME daemon is running"
        echo "Installation completed successfully!"
    else
        echo "Installation completed with errors"
    fi
fi
