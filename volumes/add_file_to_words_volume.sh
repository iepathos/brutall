#!/bin/bash
# ${1} = filename
docker run --rm -v words:/words -v $(pwd):/fbackup ubuntu cp -a /fbackup/${1} /words
