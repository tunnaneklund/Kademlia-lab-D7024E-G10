#!/bin/bash
for i in {0..49}
    do
        docker rm -f cont$i
    done