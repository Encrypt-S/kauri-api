package daemon


type Daemon interface {

	GetTransForAddresses(addresses []string) ([]Address, error)

}
