package routers

type PaymentGroup struct {
	MoMo MoMoRouter
}

var PaymentServiceRouterGroup = new(PaymentGroup)
