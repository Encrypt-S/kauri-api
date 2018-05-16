package deamonmgr

activeDeamons impletments DeaemonInterface = [navd, ethd, btcd, xrpd]

func GetTransactions(incomingTxs IncomingTransactions) {

	r = resutls

	for activeDeamons as d {

		r.add(d.GetTransactionsForAddresses(incomingTxs))

	}

}