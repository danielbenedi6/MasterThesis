type in_comm struct {
	Req   <-chan cmn.Request
	Graph <-chan cmn.Graph
}

type out_comm struct {
	Req   chan<- cmn.Request
	Graph chan<- cmn.Graph
}