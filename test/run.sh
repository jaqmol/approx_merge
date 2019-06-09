#!/bin/sh

# Setup

cd ..
bin=${PWD##*/}

# Clenup

cleanup() {
  echo
  echo Cleaning up ...
  rm ${bin}
  rm approx_check
  rm input_0.pipe
  rm input_1.pipe
  echo Done.
  trap - SIGINT SIGTERM # clear the trap
  # kill -- -$$ # Sends SIGTERM to child/sub processes
}
trap cleanup SIGINT SIGTERM

# Building

go build -o ./test/${bin}
cd ../approx_check
go build -o ../${bin}/test/approx_check
cd ../${bin}/test

mkfifo input_0.pipe
mkfifo input_1.pipe

# Running

echo Starting producers ...
echo One is producing much faster than the other.
MODE=produce SPEED=moderate ./approx_check > input_0.pipe &
MODE=produce SPEED=slow ./approx_check > input_1.pipe &

echo Starting MERGE in round_robin mode ...
echo MERGE is then throttling the faster producer to the speed of the slower one.
PICK=round_robin \
IN_COUNT=2 \
IN_0=input_0.pipe \
IN_1=input_1.pipe \
./${bin} | MODE=consume ./approx_check
