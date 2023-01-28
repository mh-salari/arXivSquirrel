#!/bin/bash
# Source: https://www.baeldung.com/linux/shell-retry-failed-command

max_iteration=5

# Check if an argument was passed
if [ $# -eq 0 ]
then
  echo "Please provide the path to arxiv program as an argument."
  exit 1
fi

for i in $(seq 1 $max_iteration)
do
  $1
  result=$?
  if [[ $result -eq 0 ]]
  then
    echo "Result successful"
    break
  else
    echo "Result unsuccessful"
    sleep 1
  fi
done

if [[ $result -ne 0 ]]
then
  echo "All of the trials failed!!!"
fi
