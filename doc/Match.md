# Mapping objects to users

* Rules to display objects:
	- workflows with matching labels (oidc fields of logged in person, e.g email and groups)
	- pods bearing workflow labels
	- example workflow labels:
		- argovue.io/oidc/group: admin
		- argovue.io/oidc/id: oidc-subect
	- user subscribes to:
		- argovue.io/oidc/group in (admin, user)
		- argovue.io/oidc/id=oidc-subject
