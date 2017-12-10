#!/bin/bash

# create words volume
docker volume create words

# add and .txt files in words/ to the words volume
pushd words
	echo "Merge words lists into allwords.txt"
	./merge_words.sh > allwords.txt
popd

for wordfile in words/*.txt ; do
	echo "Adding file $wordfile to words docker volume"
    ./add_file_to_words_volume.sh $wordfile
done
