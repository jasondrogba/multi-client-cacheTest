#!/bin/bash

for ((i=0;i<=1;i++)); do
  dd if=/dev/zero of=./"$i".txt bs=1K count=4
done
