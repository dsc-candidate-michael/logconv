package logconv

type ReqDetail struct {
	remoteAddr string
	statusCode int
	route      string
}

func (rd *ReqDetail) RemoteAddr() string {
	return rd.remoteAddr
}

func (rd *ReqDetail) StatusCode() int {
	return rd.statusCode
}

func (rd *ReqDetail) Route() string {
	return rd.route
}
