# EKS IAM Group Mapper

A utility to synchronize Group associated IAM users with Kubernetes groups.

## Purpose
AWS currently allows you to map individual users, accounts, and roles to Kubernetes groups via a configmap
for the AWS IAM Authenticator.  Currently there is no direct support for "Groups" membership (if you're part of an IAM group)
this tool attempts to rectify it.  

## Development Status
I'm currently working on this development actively, further this is the first time I've used Go besides creating some minor
patches.  Please bare with me for the short term, also feel free to file issues/pull requests.
