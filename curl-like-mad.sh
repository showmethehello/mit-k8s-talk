#!/bin/bash
while :
do
  echo -n "$(date +%T) : " && curl http://localhost:30000/
done
