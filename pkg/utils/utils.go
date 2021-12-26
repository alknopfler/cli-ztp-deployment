package utils

import (
	"context"
	"encoding/json"
	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
	"github.com/alknopfler/cli-ztp-deployment/pkg/httpd"
	"github.com/alknopfler/cli-ztp-deployment/pkg/resources"
	"k8s.io/client-go/dynamic"

	"log"
)

func GetDomainFromCluster(client dynamic.Interface, ctx context.Context) string {
	d, err := resources.NewGenericGet(ctx, client, httpd.INGRESS_CONTROLLER_GROUP, httpd.INGRESS_CONTROLLER_VERSION, httpd.INGRESS_CONTROLLER_KIND, httpd.INGRESS_CONTROLLER_NS, httpd.INGRESS_CONTROLLER_NAME, httpd.INGRESS_CONTROLLER_JQPATH).GetResourceByJq()
	if err != nil {
		log.Fatalf("[ERROR] Getting resources in GetDomainFromCluster: %e", err)
	}
	b, _ := json.Marshal(d)
	value := jsonvalue.MustUnmarshal(b)
	domain, _ := value.Get("Object", "status", "domain")

	return domain.String()
}
