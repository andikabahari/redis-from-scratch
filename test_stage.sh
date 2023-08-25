#!/bin/sh

case $1 in
  '1') 
    CODECRAFTERS_CURRENT_STAGE_SLUG='init'
    ;;
  '2') 
    CODECRAFTERS_CURRENT_STAGE_SLUG='ping-pong'
    ;;
  '3') 
    CODECRAFTERS_CURRENT_STAGE_SLUG='ping-pong-multiple'
    ;;
  '4') 
    CODECRAFTERS_CURRENT_STAGE_SLUG='concurrent-clients'
    ;;
  '5') 
    CODECRAFTERS_CURRENT_STAGE_SLUG='echo'
    ;;
  '6') 
    CODECRAFTERS_CURRENT_STAGE_SLUG='set_get'
    ;;
  '7') 
    CODECRAFTERS_CURRENT_STAGE_SLUG='expiry'
    ;;
  *)
    echo 'Invalid stage'
    exit
    ;;
esac

cd redis-tester
go build -o ../redis-go/test.out ./cmd/tester

cd ../redis-go
CODECRAFTERS_SUBMISSION_DIR=$(pwd) \
CODECRAFTERS_CURRENT_STAGE_SLUG=${CODECRAFTERS_CURRENT_STAGE_SLUG} \
./test.out
rm ./test.out