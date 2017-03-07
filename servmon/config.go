package servmon

//Configuration is the basis for configuring the software
type Configuration struct {
	Username string
	Password string
	Servers  []string
}
