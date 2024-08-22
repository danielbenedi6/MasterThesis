func generator(in in_comm, out out_comm) {
	for {
		r, ok := <-in.Req
		if !ok {
			break
		}
		switch r.Op {
		case cmn.Insert, cmn.Update:
			out_req := make(chan cmn.Request, channelSize)
			out_grph := make(chan cmn.Graph, channelSize)
			new_out := out_comm{Req: out_req, Graph: out_grph}
			new_in := in_comm{Req: out_req, Graph: out_grph}

			go filter(in, new_out, r.E)
			in = new_in
		case cmn.Delete:
			// Do nothing, asked to delete unexistent edge
		case cmn.KMST:
			g, _ := <-in.Graph

			out.Req <- r
			out.Graph <- g
		case cmn.EOF:
			out.Req <- r
			break
		default: //something's wrong
			fmt.Println("Unknown operation in generator")
			break
		}
	}
	close(out.Req)
	close(out.Graph)
}