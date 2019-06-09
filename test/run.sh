#!/bin/sh

# generic section

cd ..
bin=${PWD##*/}
go build -o ./test/${bin}
cd test

# specific section

mkfifo input_0.pipe
mkfifo input_1.pipe
mkfifo output.pipe

echo
echo 'Open a shell and paste:'
echo '  1$ cat < output.pipe'
echo
echo 'Please open 2 shells @'${PWD}
echo '  2$ ls > input_0.pipe'
echo '  3$ ls > input_1.pipe'
echo
echo 'Check if ls intputs are being printed in 1$'

PICK=as_comes \
IN_COUNT=2 \
IN_0=input_0.pipe \
IN_1=input_1.pipe \
OUT_COUNT=1 \
OUT_0=output.pipe \
./${bin}

rm ${bin}
rm input_0.pipe
rm input_1.pipe
rm output.pipe