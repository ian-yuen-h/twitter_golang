package requests

type Requests interface{
	getCommand() string
	getID()	int
}

type addRequest struct{
	Command string
	Id int
	Body string
	Timestamp int
}


type removeRequest struct{
	Command string
	Id int
	Timestamp int
}


type containsRequest struct{
	Command string
	Id int
	Timestamp int
}

type feedRequest struct{
	Command string
	Id int
}

type doneRequest struct{
	Command string
	Id int
}

func (r addRequest) getCommand() string {
    return r.Command
}
func (r removeRequest) getCommand() string {
    return r.Command
}
func (r containsRequest) getCommand() string {
    return r.Command
}
func (r feedRequest) getCommand() string {
    return r.Command
}
func (r doneRequest) getCommand() string {
    return r.Command
}

func (r addRequest) getID() int {
    return r.Id
}
func (r removeRequest) getID()int {
    return r.Id
}
func (r containsRequest) getID() int {
    return r.Id
}
func (r feedRequest) getID() int {
    return r.Id
}
func (r doneRequest) getID() int {
    return r.Id
}

func commandType(g Requests) string{
	return g.getCommand()
}

func identify(g Requests) int{
	return g.getID()
}