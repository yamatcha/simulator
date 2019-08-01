#!/bin/sh

if [ $# -ne 1 ]; then
  echo "Usage: $0 [infile]"
  exit 1
fi

if [ ! -f $1 ]; then
  echo "File $1 doesn't exists."
  exit 2
fi


abspathdir=$(cd $(dirname $1) && pwd)

IMAGE_NAME=kyontan/cacti:6.5.0

exec docker run --rm -v $abspathdir:/input $IMAGE_NAME -infile /input/$1
