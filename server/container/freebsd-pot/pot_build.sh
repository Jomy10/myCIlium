POT_NAME=mycilium-orchestrator
FREEBSD_VERSION=14.1

pot create \
  -p $POT_NAME \
  -b $FREEBSD_VERSION \
  -t single

pot copy-in \
  -p $POT_NAME \
  -s pot_init.sh \
  -d /pot_init.sh

pot copy-in \
  -p $POT_NAME \
  -s ../.. \
  -d /server

pot start $POT_NAME

pot exec -p $POT_NAME sh /pot_init.sh

pot set-cmd -p $POT_NAME -c "/mycilium-orchestrator 6001"

pot stop $POT_NAME

pot snap -p $POT_NAME -r

# TODO: store database on local storage instead of inside of pot (would be destroyed on upgrade)
