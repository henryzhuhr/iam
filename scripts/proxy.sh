#!/bin/bash

export proxy=http://host.docker.internal:7897; export http_proxy=$proxy; export https_proxy=$proxy; export no_proxy="localhost,127.0.0.1,::1"
echo "proxy:        $proxy"; echo "http_proxy:   $http_proxy"; echo "https_proxy:  $https_proxy"; echo "no_proxy:     $no_proxy"