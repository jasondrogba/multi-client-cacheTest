#!/bin/bash

for ((i=0;i<=1;i++)); do
  dd if=/dev/zero of=./"$i".txt bs=1M count=50
done
