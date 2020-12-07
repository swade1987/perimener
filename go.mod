module github.com/swade1987/perimener

go 1.15

require (
	github.com/Azure/go-autorest/autorest v0.11.13 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.8 // indirect
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/gophercloud/gophercloud v0.14.0 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/rs/zerolog v1.20.0
	golang.org/x/oauth2 v0.0.0-20201207163604-931764155e3f // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	k8s.io/api v0.18.12
	k8s.io/apimachinery v0.18.12
	k8s.io/client-go v0.18.12
	k8s.io/klog v1.0.0 // indirect
	k8s.io/utils v0.0.0-20201110183641-67b214c5f920 // indirect
)

replace (
	"github.com/swade1987/perimener/pkg/pods" => "../pkg/pods"
)
