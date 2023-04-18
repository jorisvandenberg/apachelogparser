#!/bin/bash

# Fetch the latest changes from the remote repository
git fetch --prune

# Loop over all the branches that are marked as [gone]
for branch in $(git branch -vv | grep ': gone]' | awk '{print $1}'); do
    # Delete the local branch
    git branch -D $branch
done
