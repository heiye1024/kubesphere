package gateway

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	virtualizationv1beta1 "kubesphere.io/kubesphere/api/virtualization/v1beta1"
	gatewayclient "kubesphere.io/kubesphere/pkg/virtualization/gateway/client"
	"kubesphere.io/kubesphere/pkg/virtualization/gateway/dto"
)

func respond(c *gin.Context, status int, data any, total int, message string) {
	traceID := c.GetHeader("X-Request-Id")
	if traceID == "" {
		traceID = uuid.NewString()
	}
	auditID := ensureAuditID(c)
	c.Header("X-Trace-Id", traceID)
	c.Header("X-Audit-Id", auditID)
	c.JSON(status, dto.Envelope{
		Data:    data,
		Total:   total,
		TraceID: traceID,
		AuditID: auditID,
		Message: message,
	})
}

func respondError(c *gin.Context, status int, err error) {
	respond(c, status, nil, 0, err.Error())
}

func ensureAuditID(c *gin.Context) string {
	if audit, ok := c.Get("auditID"); ok {
		if s, okCast := audit.(string); okCast && s != "" {
			return s
		}
	}
	generated := uuid.NewString()
	c.Set("auditID", generated)
	return generated
}

func queryOptions(c *gin.Context) gatewayclient.QueryOptions {
	clusters := c.QueryArray("cluster")
	if len(clusters) == 0 {
		if single := c.Query("cluster"); single != "" {
			clusters = []string{single}
		}
	}
	return gatewayclient.QueryOptions{
		Namespace: c.Param("namespace"),
		Clusters:  clusters,
	}
}

func applyClusterFallback(obj metav1.Object, opts gatewayclient.QueryOptions) {
	if obj == nil {
		return
	}
	if len(opts.Clusters) == 0 {
		return
	}
	labels := obj.GetLabels()
	if labels == nil {
		labels = map[string]string{}
	}
	if labels[dto.LabelCluster] == "" {
		cluster := opts.Clusters[0]
		if cluster != "" && cluster != "all" {
			labels[dto.LabelCluster] = cluster
		}
	}
	obj.SetLabels(labels)
}

func (r *Router) listVMs(c *gin.Context) {
	opts := queryOptions(c)
	vms, err := r.Store.ListVMs(c.Request.Context(), opts)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respond(c, http.StatusOK, dto.FromVMList(vms), len(vms.Items), "")
}

func (r *Router) createVM(c *gin.Context) {
	ns := c.Param("namespace")
	opts := queryOptions(c)
	var payload dto.VirtualMachineRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}
	vm := dto.ToVirtualMachine(ns, payload)
	applyClusterFallback(vm, opts)
	if err := r.Store.CreateVM(c.Request.Context(), vm); err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respond(c, http.StatusCreated, vm, 1, "")
}

func (r *Router) powerOn(c *gin.Context)  { r.togglePower(c, virtualizationv1beta1.PowerStateRunning) }
func (r *Router) powerOff(c *gin.Context) { r.togglePower(c, virtualizationv1beta1.PowerStateStopped) }

func (r *Router) togglePower(c *gin.Context, state virtualizationv1beta1.PowerState) {
	ns := c.Param("namespace")
	name := c.Param("name")
	if err := r.Store.UpdatePowerState(c.Request.Context(), ns, name, state); err != nil {
		if gatewayclient.IsForbidden(err) {
			c.Header("X-Deny-Reason", err.Error())
			respondError(c, http.StatusForbidden, err)
			return
		}
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respond(c, http.StatusAccepted, nil, 0, "power state update scheduled")
}

func (r *Router) migrate(c *gin.Context) {
	ns := c.Param("namespace")
	name := c.Param("name")
	if err := r.Store.MigrateVM(c.Request.Context(), ns, name); err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respond(c, http.StatusAccepted, nil, 0, "migration scheduled")
}

func (r *Router) console(c *gin.Context) {
	ns := c.Param("namespace")
	name := c.Param("name")
	url, err := r.Store.Console(c.Request.Context(), ns, name)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respond(c, http.StatusOK, gin.H{"consoleURL": url}, 1, "")
}

func (r *Router) listDisks(c *gin.Context) {
	opts := queryOptions(c)
	disks, err := r.Store.ListDisks(c.Request.Context(), opts)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respond(c, http.StatusOK, dto.FromDiskList(disks), len(disks.Items), "")
}

func (r *Router) createDisk(c *gin.Context) {
	ns := c.Param("namespace")
	opts := queryOptions(c)
	var payload dto.VirtualDiskRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}
	disk := dto.ToVirtualDisk(ns, payload)
	applyClusterFallback(disk, opts)
	if err := r.Store.CreateDisk(c.Request.Context(), disk); err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respond(c, http.StatusCreated, disk, 1, "")
}

func (r *Router) listNets(c *gin.Context) {
	opts := queryOptions(c)
	nets, err := r.Store.ListNetworks(c.Request.Context(), opts)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respond(c, http.StatusOK, dto.FromNetList(nets), len(nets.Items), "")
}

func (r *Router) createNet(c *gin.Context) {
	ns := c.Param("namespace")
	opts := queryOptions(c)
	var payload dto.VirtualNetRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}
	netObj := dto.ToVirtualNet(ns, payload)
	applyClusterFallback(netObj, opts)
	if err := r.Store.CreateNetwork(c.Request.Context(), netObj); err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respond(c, http.StatusCreated, netObj, 1, "")
}

func (r *Router) listSnapshots(c *gin.Context) {
	opts := queryOptions(c)
	snaps, err := r.Store.ListSnapshots(c.Request.Context(), opts)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respond(c, http.StatusOK, dto.FromSnapshotList(snaps), len(snaps.Items), "")
}

func (r *Router) createSnapshot(c *gin.Context) {
	ns := c.Param("namespace")
	opts := queryOptions(c)
	var payload dto.VMSnapshotRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}
	snapshot := dto.ToVMSnapshot(ns, payload)
	applyClusterFallback(snapshot, opts)
	if err := r.Store.CreateSnapshot(c.Request.Context(), snapshot); err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respond(c, http.StatusCreated, snapshot, 1, "")
}

func (r *Router) listTemplates(c *gin.Context) {
	opts := queryOptions(c)
	templates, err := r.Store.ListTemplates(c.Request.Context(), opts)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respond(c, http.StatusOK, dto.FromTemplateList(templates), len(templates.Items), "")
}

func (r *Router) createTemplate(c *gin.Context) {
	ns := c.Param("namespace")
	opts := queryOptions(c)
	var payload dto.VMTemplateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}
	template := dto.ToVMTemplate(ns, payload)
	applyClusterFallback(template, opts)
	if err := r.Store.CreateTemplate(c.Request.Context(), template); err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respond(c, http.StatusCreated, template, 1, "")
}
