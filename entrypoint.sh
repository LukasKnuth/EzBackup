#!/bin/sh -eu

for cmd in "$@"; do
  case "$cmd" in
    safe_init)
      /usr/bin/restic snapshots
      if [ $? -eq 0 ]
      then
        exit 1
      else
        exec /usr/bin/restic init
      fi
      ;;

    run)
      shift # dismiss first script argument
      exec /usr/bin/restic "$@"
      ;;

    shell)
      exec sh -i
      ;;

    *)
      echo "Invalid command: $cmd" >&2
      exit 1
      ;;
  esac