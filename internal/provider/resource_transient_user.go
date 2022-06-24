package provider

/*
	transient_user is a user that will NOT represent a resource, nor a datasource.
	transient_user queries if a user with set email already exists
	if not, it creates and sends password reset mail
	else if, it returns the user matching the email
	because the user is referenced by email, and not unique_id, we consider this
	transient. There is no guarantee that the user with this email actually exists
	I'll also try to check if the user exists
*/
