#!/bin/bash

function ol { ss|awk -F ":" '{print $3":"$4":"$5}'|grep ffff:|sort |uniq; }
