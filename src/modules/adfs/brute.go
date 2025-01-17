package adfs

func CheckTarget(optionsInterface *interface{}) {
	options := (*optionsInterface).(*Options)

	// If the target is not specified, we will try to find the ADFS URL with the endpoint getuserrealm
	if options.Target == "" {
		options.Target = options.findTarget(options.Domain)
		if options.Target == "" {
			log.Error("The ADFS URL was not found")
			return
		}
		options.Log.Verbose("An ADFS instance has been found on " + options.Target)
	}
}

func Authenticate(optionsInterface *interface{}, email, password string) bool {
	options := (*optionsInterface).(*Options)
	return options.brute(email, password)

}
