package ws

type IHandler interface {
	Router(msg_t int, msg_obj any)
}
