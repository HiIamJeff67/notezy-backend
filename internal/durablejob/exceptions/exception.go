package exceptions

type DurableJobException struct {
	Domain string
}

func NewDurableJobException(domain string) DurableJobException {
	return DurableJobException{Domain: domain}
}
