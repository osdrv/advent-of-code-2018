#!/bin/bash

set -e;

day=$(find . -name 'day-*' | sort -V | tail -n 1 | sed 's/[^0-9]//g');

day=$((day+1));

daydir="day-${day}";

mkdir "${daydir}";

cp main.go.tmpl "${daydir}/main.go";
cp common.go.tmpl "${daydir}/common.go";
cp Makefile.tmpl "${daydir}/Makefile";

touch "${daydir}/INPUT";

touch "${daydir}/INPUT-TST";

echo "new day folder created: ${daydir}";

cd "${daydir}"

go mod init
go mod tidy
