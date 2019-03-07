AWS Auth Map Operator
=====================

Requirements
------------
- Configurable cycle time (how often to scrap the IAM configuration from Amazon)
- helm chart support for kiam/kube2iam annotation
- Configmap entry for which groups to translate to which roles
- Optional RBAC for helm
- Optional aws region overwrite?  (Not sure that's neccessary for IAM)
- Optional by tag?


CLI Spec
--------
auth-map-updater

Global Flags:
--verbose
--aws-key
--aws-secret
--aws-region

Commands:
- update
  --dry-run
  arn:group(*)
- help


Pseudocode
----------
Configure AWS by (ranked in order) flag, environment, nodemetadata service

during execution every n seconds:
if usergroups defined:
	map user groups to users to kubernetes groups
look for kubernetes.io/cluster/so-nonprod tags (value is groupARN:group)
	map user groups to users to kubernetes groups
push configmap to Kubernetes ConfigMap in kube-system
sleep