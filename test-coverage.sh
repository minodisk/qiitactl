#!/bin/bash

echo "mode: set" > $CIRCLE_ARTIFACTS/coverage.out

for Dir in $(find ./* -maxdepth 10 -type d );
do
  if ls $Dir/*.go &> /dev/null; then
    go test -v -coverprofile=$CIRCLE_ARTIFACTS/tmp.out $Dir
    if [ -f $CIRCLE_ARTIFACTS/tmp.out ]; then
      cat $CIRCLE_ARTIFACTS/tmp.out | grep -v "mode: set" >> $CIRCLE_ARTIFACTS/coverage.out
    fi
  fi
done

rm -rf $CIRCLE_ARTIFACTS/tmp.out
