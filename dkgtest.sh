#!/bin/bash

while getopts ":-n:-f:-s:-p:" opt
do
    case $opt in
        n|--nodeNum)
            nodeNum=$OPTARG;;
        f|--degree)
            degree=$OPTARG;;
        s|--secretNum)
            secretNum=$OPTARG;;
        p|--path)
            path=$OPTARG;;
    esac
done

for((i=0;i<nodeNum;i++));
do
    gnome-terminal --window -- bash -c $"go run main.go -n=$nodeNum -id=$i -f=$degree -s=$secretNum -path=$path; exec bash;"&
done
wait