package presentation

import (
	"github.com/chujieyang/commonops/ops/presentation/batch"
	"github.com/chujieyang/commonops/ops/presentation/cd"
	"github.com/chujieyang/commonops/ops/presentation/ci"
	"github.com/chujieyang/commonops/ops/presentation/cloud_account"
	"github.com/chujieyang/commonops/ops/presentation/daily_job"
	"github.com/chujieyang/commonops/ops/presentation/data"
	"github.com/chujieyang/commonops/ops/presentation/dms"
	"github.com/chujieyang/commonops/ops/presentation/dns"
	"github.com/chujieyang/commonops/ops/presentation/ecs"
	"github.com/chujieyang/commonops/ops/presentation/k8s"
	"github.com/chujieyang/commonops/ops/presentation/kv"
	"github.com/chujieyang/commonops/ops/presentation/middleware"
	"github.com/chujieyang/commonops/ops/presentation/monitor"
	"github.com/chujieyang/commonops/ops/presentation/nacos"
	"github.com/chujieyang/commonops/ops/presentation/other"
	"github.com/chujieyang/commonops/ops/presentation/permission"
	"github.com/chujieyang/commonops/ops/presentation/rds"
	"github.com/chujieyang/commonops/ops/presentation/slb"
	"github.com/chujieyang/commonops/ops/presentation/user"
	"github.com/gin-gonic/gin"
)

// RegisterRouter is a global router register
func RegisterRouter(engine *gin.Engine) {
	engine.POST("/user/login", user.ILogin)
	engine.GET("/hc", user.ILogin)

	wsRoute := engine.Group("/ws")
	{
		wsRoute.GET("/kubernetes/container_terminal", middleware.KubernetesMiddleWare(), k8s.GetContainerTerminal)
		wsRoute.GET("/cloud/ssh", ecs.WsSsh)
	}

	userRoute := engine.Group("/user", middleware.AuthMiddleWare())
	{
		userRoute.GET("/tokenRefresh", user.ITokenRefresh)
		userRoute.POST("/updatePassword", user.IUpdatePassword)
		userRoute.POST("/create", user.IPostUserCreate)
		userRoute.PUT("/active", user.IPutUserActive)
		userRoute.GET("/list", user.IGetUsersList)
		userRoute.GET("/permissions", user.IGetUserPermissions)
		userRoute.POST("/feedback", user.IPostUserFeedback)
		userRoute.GET("/feedback", user.IGetUserFeedback)
	}

	cloudRoute := engine.Group("/cloud", middleware.AuthMiddleWare())
	{
		cloudRoute.GET("/accounts", cloud_account.IGetCloudAccounts)
		cloudRoute.POST("/accounts", cloud_account.IPostCloudAccounts)
		cloudRoute.PUT("/accounts", cloud_account.IPutCloudAccounts)
		cloudRoute.DELETE("/accounts", cloud_account.IDeleteCloudAccounts)

		cloudRoute.GET("/servers", ecs.IGetCloudServers)
		cloudRoute.GET("/servers/treedata", ecs.IGetCloudServersTreeData)
		cloudRoute.POST("/servers", ecs.IPostCloudServers)
		cloudRoute.PUT("/servers", ecs.IPutCloudServers)
		cloudRoute.DELETE("/servers", ecs.IDeleteCloudServers)
		cloudRoute.GET("/server", ecs.IGetCloudServerDetail)
		cloudRoute.POST("/servers/batch/ssh", ecs.IPostCloudServersBatchSSH)

		cloudRoute.GET("/rds/detail", rds.IGetCloudRdsDetail)
		cloudRoute.GET("/rds", rds.IGetCloudRds)
		cloudRoute.PUT("/rds", rds.IPutCloudRds)
		cloudRoute.POST("/rds", rds.IPostCloudRds)
		cloudRoute.DELETE("/rds", rds.IDeleteCloudRds)

		cloudRoute.GET("/kv/detail", kv.IGetCloudKvDetail)
		cloudRoute.GET("/kv", kv.IGetCloudKv)
		cloudRoute.PUT("/kv", kv.IPutCloudKv)
		cloudRoute.POST("/kv", kv.IPostCloudKv)
		cloudRoute.DELETE("/kv", kv.IDeleteCloudKv)

		cloudRoute.GET("/slb/detail", slb.IGetCloudSlbDetail)
		cloudRoute.GET("/slb", slb.IGetCloudSlb)
		cloudRoute.DELETE("/slb", slb.IDeleteCloudSlb)

		cloudRoute.GET("/other", other.IGetCloudOtherRes)
		cloudRoute.POST("/other", other.IPostCloudOtherRes)
		cloudRoute.DELETE("/other", other.IDeleteCloudOtherRes)
	}

	cloudMonitor := engine.Group("/cloud/monitor", middleware.AuthMiddleWare())
	{
		cloudMonitor.GET("/ecs", monitor.IGetCloudEcsMonitor)
		cloudMonitor.GET("/rds", monitor.IGetCloudRdsMonitor)
		cloudMonitor.GET("/kv", monitor.IGetCloudKvMonitor)
		cloudMonitor.GET("/slb", monitor.IGetCloudSlbMonitor)
	}

	roleRoute := engine.Group("/roles", middleware.AuthMiddleWare())
	{
		roleRoute.GET("/list", user.IGetRolesList)
		roleRoute.POST("/addRole", user.ICreateRole)
		roleRoute.PUT("/updateRole", user.IUpdateRole)
		roleRoute.DELETE("/deleteRole", user.IDeleteRole)
		roleRoute.GET("/users", user.IGetRoleUserList)
		roleRoute.POST("/users", user.IPostRoleUserList)
		roleRoute.GET("/resources", user.IGetRoleResourceList)
		roleRoute.POST("/resources", user.IPostRoleResourcesList)
		roleRoute.GET("/authLink", user.IGetRoleAuthLinkList)
		roleRoute.POST("/authLinks", user.ICreateRoleAuthLink)
	}

	permissionRoute := engine.Group("/permissions", middleware.AuthMiddleWare())
	{
		permissionRoute.GET("/list", permission.IGetPermissionsList)
		permissionRoute.GET("/authLink", permission.IGetAuthLink)
		permissionRoute.PUT("/authLink", permission.IPutAuthLink)
		permissionRoute.POST("/authLink", permission.ICreateAuthLink)
		permissionRoute.DELETE("/authLink", permission.IDeleteAuthLink)
	}

	jobRoute := engine.Group("/dailyJob", middleware.AuthMiddleWare())
	{
		jobRoute.POST("/", daily_job.IPostDailyJob)
		jobRoute.GET("/list", daily_job.IGetDailyJobs)
		jobRoute.GET("/info", daily_job.IGetDailyJobDetail)
		jobRoute.PUT("/", daily_job.IPutDailyJob)
		jobRoute.PUT("/executorUser", daily_job.IPutDailyJobExecutorUser)
	}

	dataRoute := engine.Group("/data", middleware.AuthMiddleWare())
	{
		dataRoute.GET("/syncAliyunEcs", data.IGetAliyunEcsSync)
		dataRoute.GET("/syncAliyunRds", data.IGetAliyunRdsSync)
		dataRoute.GET("/syncAliyunKv", data.IGetAliyunKvSync)
		dataRoute.GET("/syncAliyunSlb", data.IGetAliyunSlbSync)
		dataRoute.GET("/syncAliyunStatisData", data.IGetAliyunStatisData)
	}

	configCenterRoute := engine.Group("/configCenter", middleware.AuthMiddleWare())
	{
		configCenterRoute.POST("/nacos", nacos.IPostNacos)
		configCenterRoute.GET("/nacos/list", nacos.IGetNacosList)
		configCenterRoute.GET("/nacos/namespaces", nacos.IGetNacosNamespaceList)
		configCenterRoute.GET("/nacos/configs", nacos.IGetNacosConfigList)
		configCenterRoute.GET("/nacos/config", nacos.IGetNacosConfigDetail)
		configCenterRoute.GET("/nacos/configs/all", nacos.IGetNacosAllConfigs)
		configCenterRoute.POST("/nacos/config", nacos.IPostNacosConfig)
		configCenterRoute.PUT("/nacos/config", nacos.IPutNacosConfig)
		configCenterRoute.DELETE("/nacos/config", nacos.IDeleteNacosConfig)
		configCenterRoute.POST("/nacos/config/copy", nacos.IPostNacosConfigCopy)
		configCenterRoute.POST("/nacos/config/sync", nacos.IPostNacosConfigSync)
		configCenterRoute.GET("/configTemplates", nacos.IGetConfigTemplates)
		configCenterRoute.GET("/configTemplates/all", nacos.IGetConfigTemplatesAll)
		configCenterRoute.POST("/configTemplate", nacos.IPostConfigTemplate)
		configCenterRoute.PUT("/configTemplate", nacos.IPutConfigTemplate)
		configCenterRoute.DELETE("/configTemplate", nacos.IDeleteConfigTemplate)
	}

	k8sRoute := engine.Group("/kubernetes", middleware.AuthMiddleWare(), middleware.KubernetesMiddleWare())
	{
		k8sRoute.GET("/cluster", k8s.GetK8sCluster)
		k8sRoute.POST("/cluster", k8s.PostK8sCluster)
		k8sRoute.DELETE("/cluster", k8s.DeleteK8sCluster)
		k8sRoute.GET("/namespaces", k8s.GetNamespaces)
		k8sRoute.POST("/namespaces", k8s.PostNamespaces)
		k8sRoute.DELETE("/namespaces", k8s.DeleteNamespaces)
		k8sRoute.GET("/deployments", k8s.GetDeployments)
		k8sRoute.PUT("/deployment/restart", k8s.RestartDeployments)
		k8sRoute.GET("/replication_controllers", k8s.GetReplicationControllers)
		k8sRoute.GET("/replica_sets", k8s.GetReplicaSets)
		k8sRoute.GET("/pods", k8s.GetPods)
		k8sRoute.GET("/pod/log", k8s.GetPodLogs)
		k8sRoute.GET("/nodes", k8s.GetNodes)
		k8sRoute.GET("/services", k8s.GetServices)
		k8sRoute.GET("/ingress", k8s.GetIngress)
		k8sRoute.GET("/config_dict", k8s.GetConfigDict)
		k8sRoute.GET("/secret_dict", k8s.GetSecretDict)
		k8sRoute.POST("/yaml_resource", k8s.CreateResourceByYaml)
		k8sRoute.PUT("/yaml_resource", k8s.UpdateResourceByYaml)
		k8sRoute.GET("/yaml", k8s.GetResourceYaml)
		k8sRoute.PUT("/scale", k8s.PutResourceScale)
		k8sRoute.DELETE("/resource", k8s.DeleteResource)
		k8sRoute.DELETE("/config_map", k8s.DeleteConfigMap)
		k8sRoute.DELETE("/secret", k8s.DeleteSecret)
		k8sRoute.GET("/metrics/node", k8s.GetNodesMetrics)
		k8sRoute.GET("/metrics/pod", k8s.GetPodMetrics)
	}

	dmsRoute := engine.Group("/dms", middleware.AuthMiddleWare())
	{
		dmsRoute.GET("/instances", dms.IGetDmsInstanceData)
		dmsRoute.POST("/instance", dms.IPostDmsInstance)
		dmsRoute.DELETE("/instance", dms.IDeleteDmsInstance)
		dmsRoute.GET("/instances/all", dms.IGetAllDmsInstancesData)
		dmsRoute.GET("/authData", dms.IGetDmsAuthData)
		dmsRoute.GET("/instanceData", dms.IGetDmsInstanceData)
		dmsRoute.GET("/databaseData", dms.IGetDmsDatabaseByInstanceId)
		dmsRoute.DELETE("/databaseData", dms.IDeleteDmsDatabaseByDatabaseId)
		dmsRoute.POST("/databaseData", dms.IPostDmsDatabase)
		dmsRoute.POST("/auth", dms.IPostDmsUserAuth)
		dmsRoute.DELETE("/auth", dms.IDeleteDmsUserAuth)
		dmsRoute.GET("/userInstanceData", dms.IGetUserDmsInstanceData)
		dmsRoute.GET("/userDatabaseData", dms.IGetUserDmsDatabaseData)
		dmsRoute.POST("/userExecSQL", dms.IPostDmsUserExecSQL)
		dmsRoute.GET("/userLog", dms.IGetDmsQueryLogData)
	}

	dnsRoute := engine.Group("/dns", middleware.AuthMiddleWare())
	{
		dnsRoute.POST("/domain", dns.IPostDnsDomain)
		dnsRoute.GET("/domainList", dns.IGetDnsDomainList)
		dnsRoute.GET("/domainHistoryList", dns.IGetDnsDomainHistoryList)
		dnsRoute.GET("/domainRecordsList", dns.IGetDnsDomainRecordsList)
		dnsRoute.POST("/domainRecord", dns.IPostDnsDomainRecord)
		dnsRoute.POST("/domainRecordUpdate", dns.IUpdateDnsDomainRecord)
		dnsRoute.DELETE("/domainRecord", dns.IDeleteDnsDomainRecord)
		dnsRoute.POST("/domainRecordStatus", dns.IPostDnsDomainRecordStatus)
	}

	ciRoute := engine.Group("/ci", middleware.AuthMiddleWare())
	{
		ciRoute.POST("/job", ci.IPostJenkinsJob)
		ciRoute.GET("/job", ci.IGetJenkinsJobConfig)
		ciRoute.PUT("/job", ci.IPutJenkinsJob)
		ciRoute.DELETE("/job", ci.IDeleteJenkinsJob)
		ciRoute.GET("/jobList", ci.IGetJenkinsJobList)
		ciRoute.GET("/buildList", ci.IGetJenkinsBuildList)
		ciRoute.GET("/buildInfo", ci.IGetJenkinsBuildInfo)
		ciRoute.GET("/build/stages", ci.IGetJenkinsBuildStageList)
		ciRoute.GET("/build/stage/log", ci.IGetJenkinsBuildStageDetailLog)
		ciRoute.POST("/build", ci.IPostJenkinsBuild)
		ciRoute.DELETE("/build", ci.IDeleteJenkinsBuildLog)
		ciRoute.GET("/buildLog", ci.IGetJenkinsBuildLog)
		ciRoute.GET("/credentials/list", ci.IGetJenkinsCredentialsList)
		ciRoute.POST("/credential", ci.IPostJenkinsCredentials)
		ciRoute.GET("/build/archiveArtifactsInfo", ci.IGetJenkinsBuildArchiveArtifactsInfo)
	}

	cdRoute := engine.Group("/cd", middleware.AuthMiddleWare())
	{
		cdRoute.GET("/processTemplateList", cd.IGetCdProcessTemplateList)
		cdRoute.POST("/processTemplate", cd.IPostCdProcessTemplate)
		cdRoute.GET("/processLog", cd.IGetCdProcessLog)
		cdRoute.POST("/processLog", cd.IPostCdProcessLog)
	}

	batchRoute := engine.Group("/batch", middleware.AuthMiddleWare())
	{
		batchRoute.GET("/allNodeTree", batch.IGetAllNodeTree)
	}

}
