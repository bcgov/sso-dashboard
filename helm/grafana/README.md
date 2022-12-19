# Grafana Helm Chart

* Installs the web dashboarding system [Grafana](http://grafana.org/)

## Installing/upgrading the Chart

To install/upgrade the chart with the release name `sso-grafana`:

```console
make upgrade NAMESPACE=<namespace>
```

## Uninstalling the Chart

To uninstall/delete the my-release deployment:

```console
make uninstall NAMESPACE=<namespace>
```

## Configuration

| Parameter                                 | Description                                   | Default                                                 |
|-------------------------------------------|-----------------------------------------------|---------------------------------------------------------|
| `replicas`                                | Number of nodes                               | `1`                                                     |
| `podDisruptionBudget.minAvailable`        | Pod disruption minimum available              | `nil`                                                   |
| `podDisruptionBudget.maxUnavailable`      | Pod disruption maximum unavailable            | `nil`                                                   |
| `deploymentStrategy`                      | Deployment strategy                           | `{ "type": "RollingUpdate" }`                           |
| `livenessProbe`                           | Liveness Probe settings                       | `{ "httpGet": { "path": "/api/health", "port": 3000 } "initialDelaySeconds": 60, "timeoutSeconds": 30, "failureThreshold": 10 }` |
| `readinessProbe`                          | Readiness Probe settings                      | `{ "httpGet": { "path": "/api/health", "port": 3000 } }`|
| `securityContext`                         | Deployment securityContext                    | `{"runAsUser": 472, "runAsGroup": 472, "fsGroup": 472}`  |
| `priorityClassName`                       | Name of Priority Class to assign pods         | `nil`                                                   |
| `image.repository`                        | Image repository                              | `grafana/grafana`                                       |
| `image.tag`                               | Overrides the Grafana image tag whose default is the chart appVersion (`Must be >= 5.0.0`) | ``                                                      |
| `image.sha`                               | Image sha (optional)                          | ``                                                      |
| `image.pullPolicy`                        | Image pull policy                             | `IfNotPresent`                                          |
| `image.pullSecrets`                       | Image pull secrets (can be templated)         | `[]`                                                    |
| `service.enabled`                         | Enable grafana service                        | `true`                                                  |
| `service.type`                            | Kubernetes service type                       | `ClusterIP`                                             |
| `service.port`                            | Kubernetes port where service is exposed      | `80`                                                    |
| `service.portName`                        | Name of the port on the service               | `service`                                               |
| `service.appProtocol`                     | Adds the appProtocol field to the service     | ``                                                      |
| `service.targetPort`                      | Internal service is port                      | `3000`                                                  |
| `service.nodePort`                        | Kubernetes service nodePort                   | `nil`                                                   |
| `service.annotations`                     | Service annotations (can be templated)        | `{}`                                                    |
| `service.labels`                          | Custom labels                                 | `{}`                                                    |
| `service.clusterIP`                       | internal cluster service IP                   | `nil`                                                   |
| `service.loadBalancerIP`                  | IP address to assign to load balancer (if supported) | `nil`                                            |
| `service.loadBalancerSourceRanges`        | list of IP CIDRs allowed access to lb (if supported) | `[]`                                             |
| `service.externalIPs`                     | service external IP addresses                 | `[]`                                                    |
| `headlessService`                         | Create a headless service                     | `false`                                                 |
| `extraExposePorts`                        | Additional service ports for sidecar containers| `[]`                                                   |
| `hostAliases`                             | adds rules to the pod's /etc/hosts            | `[]`                                                    |
| `ingress.enabled`                         | Enables Ingress                               | `false`                                                 |
| `ingress.annotations`                     | Ingress annotations (values are templated)    | `{}`                                                    |
| `ingress.labels`                          | Custom labels                                 | `{}`                                                    |
| `ingress.path`                            | Ingress accepted path                         | `/`                                                     |
| `ingress.pathType`                        | Ingress type of path                          | `Prefix`                                                |
| `ingress.hosts`                           | Ingress accepted hostnames                    | `["chart-example.local"]`                                                    |
| `ingress.extraPaths`                      | Ingress extra paths to prepend to every host configuration. Useful when configuring [custom actions with AWS ALB Ingress Controller](https://kubernetes-sigs.github.io/aws-alb-ingress-controller/guide/ingress/annotation/#actions). Requires `ingress.hosts` to have one or more host entries. | `[]`                                                    |
| `ingress.tls`                             | Ingress TLS configuration                     | `[]`                                                    |
| `resources`                               | CPU/Memory resource requests/limits           | `{}`                                                    |
| `nodeSelector`                            | Node labels for pod assignment                | `{}`                                                    |
| `tolerations`                             | Toleration labels for pod assignment          | `[]`                                                    |
| `affinity`                                | Affinity settings for pod assignment          | `{}`                                                    |
| `extraInitContainers`                     | Init containers to add to the grafana pod     | `{}`                                                    |
| `extraContainers`                         | Sidecar containers to add to the grafana pod  | `""`                                                    |
| `extraContainerVolumes`                   | Volumes that can be mounted in sidecar containers | `[]`                                                |
| `extraLabels`                             | Custom labels for all manifests               | `{}`                                                    |
| `schedulerName`                           | Name of the k8s scheduler (other than default) | `nil`                                                  |
| `persistence.enabled`                     | Use persistent volume to store data           | `false`                                                 |
| `persistence.type`                        | Type of persistence (`pvc` or `statefulset`)  | `pvc`                                                   |
| `persistence.size`                        | Size of persistent volume claim               | `10Gi`                                                  |
| `persistence.existingClaim`               | Use an existing PVC to persist data (can be templated) | `nil`                                          |
| `persistence.storageClassName`            | Type of persistent volume claim               | `nil`                                                   |
| `persistence.accessModes`                 | Persistence access modes                      | `[ReadWriteOnce]`                                       |
| `persistence.annotations`                 | PersistentVolumeClaim annotations             | `{}`                                                    |
| `persistence.finalizers`                  | PersistentVolumeClaim finalizers              | `[ "kubernetes.io/pvc-protection" ]`                    |
| `persistence.extraPvcLabels`              | Extra labels to apply to a PVC.               | `{}`                                                    |
| `persistence.subPath`                     | Mount a sub dir of the persistent volume (can be templated) | `nil`                                     |
| `persistence.inMemory.enabled`            | If persistence is not enabled, whether to mount the local storage in-memory to improve performance | `false`                                                   |
| `persistence.inMemory.sizeLimit`          | SizeLimit for the in-memory local storage     | `nil`                                                   |
| `initChownData.enabled`                   | If false, don't reset data ownership at startup | true                                                  |
| `initChownData.image.repository`          | init-chown-data container image repository    | `busybox`                                               |
| `initChownData.image.tag`                 | init-chown-data container image tag           | `1.31.1`                                                |
| `initChownData.image.sha`                 | init-chown-data container image sha (optional)| `""`                                                    |
| `initChownData.image.pullPolicy`          | init-chown-data container image pull policy   | `IfNotPresent`                                          |
| `initChownData.resources`                 | init-chown-data pod resource requests & limits | `{}`                                                   |
| `schedulerName`                           | Alternate scheduler name                      | `nil`                                                   |
| `env`                                     | Extra environment variables passed to pods    | `{}`                                                    |
| `envValueFrom`                            | Environment variables from alternate sources. See the API docs on [EnvVarSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#envvarsource-v1-core) for format details. Can be templated | `{}` |
| `envFromSecret`                           | Name of a Kubernetes secret (must be manually created in the same namespace) containing values to be added to the environment. Can be templated | `""` |
| `envFromSecrets`                          | List of Kubernetes secrets (must be manually created in the same namespace) containing values to be added to the environment. Can be templated | `[]` |
| `envFromConfigMaps`                       | List of Kubernetes ConfigMaps (must be manually created in the same namespace) containing values to be added to the environment. Can be templated | `[]` |
| `envRenderSecret`                         | Sensible environment variables passed to pods and stored as secret | `{}`                               |
| `enableServiceLinks`                      | Inject Kubernetes services as environment variables. | `true`                                           |
| `extraSecretMounts`                       | Additional grafana server secret mounts       | `[]`                                                    |
| `extraVolumeMounts`                       | Additional grafana server volume mounts       | `[]`                                                    |
| `createConfigmap`                         | Enable creating the grafana configmap         | `true`                                                  |
| `extraConfigmapMounts`                    | Additional grafana server configMap volume mounts (values are templated) | `[]`                         |
| `extraEmptyDirMounts`                     | Additional grafana server emptyDir volume mounts | `[]`                                                 |
| `plugins`                                 | Plugins to be loaded along with Grafana       | `[]`                                                    |
| `datasources`                             | Configure grafana datasources (passed through tpl) | `{}`                                               |
| `alerting`                                | Configure grafana alerting (passed through tpl) | `{}`                                                  |
| `notifiers`                               | Configure grafana notifiers                   | `{}`                                                    |
| `dashboardProviders`                      | Configure grafana dashboard providers         | `{}`                                                    |
| `dashboards`                              | Dashboards to import                          | `{}`                                                    |
| `dashboardsConfigMaps`                    | ConfigMaps reference that contains dashboards | `{}`                                                    |
| `grafana.ini`                             | Grafana's primary configuration               | `{}`                                                    |
| `global.imagePullSecrets`                 | Global image pull secrets (can be templated). Allows either an array of {name: pullSecret} maps (k8s-style), or an array of strings (more common helm-style).  | `[]`                                                    |
| `ldap.enabled`                            | Enable LDAP authentication                    | `false`                                                 |
| `ldap.existingSecret`                     | The name of an existing secret containing the `ldap.toml` file, this must have the key `ldap-toml`. | `""` |
| `ldap.config`                             | Grafana's LDAP configuration                  | `""`                                                    |
| `annotations`                             | Deployment annotations                        | `{}`                                                    |
| `labels`                                  | Deployment labels                             | `{}`                                                    |
| `podAnnotations`                          | Pod annotations                               | `{}`                                                    |
| `podLabels`                               | Pod labels                                    | `{}`                                                    |
| `podPortName`                             | Name of the grafana port on the pod           | `grafana`                                               |
| `lifecycleHooks`                          | Lifecycle hooks for podStart and preStop [Example](https://kubernetes.io/docs/tasks/configure-pod-container/attach-handler-lifecycle-event/#define-poststart-and-prestop-handlers)     | `{}`                                                    |
| `sidecar.image.repository`                | Sidecar image repository                      | `quay.io/kiwigrid/k8s-sidecar`                          |
| `sidecar.image.tag`                       | Sidecar image tag                             | `1.19.2`                                                |
| `sidecar.image.sha`                       | Sidecar image sha (optional)                  | `""`                                                    |
| `sidecar.imagePullPolicy`                 | Sidecar image pull policy                     | `IfNotPresent`                                          |
| `sidecar.resources`                       | Sidecar resources                             | `{}`                                                    |
| `sidecar.securityContext`                 | Sidecar securityContext                       | `{}`                                                    |
| `sidecar.enableUniqueFilenames`           | Sets the kiwigrid/k8s-sidecar UNIQUE_FILENAMES environment variable. If set to `true` the sidecar will create unique filenames where duplicate data keys exist between ConfigMaps and/or Secrets within the same or multiple Namespaces. | `false`                           |
| `sidecar.alerts.enabled`             | Enables the cluster wide search for alerts and adds/updates/deletes them in grafana |`false`       |
| `sidecar.alerts.label`               | Label that config maps with alerts should have to be added | `grafana_alert`                               |
| `sidecar.alerts.labelValue`          | Label value that config maps with alerts should have to be added | `""`                                |
| `sidecar.alerts.searchNamespace`     | Namespaces list. If specified, the sidecar will search for alerts config-maps  inside these namespaces. Otherwise the namespace in which the sidecar is running will be used. It's also possible to specify ALL to search in all namespaces. | `nil`                               |
| `sidecar.alerts.watchMethod`         | Method to use to detect ConfigMap changes. With WATCH the sidecar will do a WATCH requests, with SLEEP it will list all ConfigMaps, then sleep for 60 seconds. | `WATCH` |
| `sidecar.alerts.resource`            | Should the sidecar looks into secrets, configmaps or both. | `both`                               |
| `sidecar.alerts.reloadURL`           | Full url of datasource configuration reload API endpoint, to invoke after a config-map change | `"http://localhost:3000/api/admin/provisioning/alerting/reload"` |
| `sidecar.alerts.skipReload`          | Enabling this omits defining the REQ_URL and REQ_METHOD environment variables | `false` |
| `sidecar.alerts.initDatasources`     | Set to true to deploy the datasource sidecar as an initContainer in addition to a container. This is needed if skipReload is true, to load any alerts defined at startup time. | `false` |
| `sidecar.dashboards.enabled`              | Enables the cluster wide search for dashboards and adds/updates/deletes them in grafana | `false`       |
| `sidecar.dashboards.SCProvider`           | Enables creation of sidecar provider          | `true`                                                  |
| `sidecar.dashboards.provider.name`        | Unique name of the grafana provider           | `sidecarProvider`                                       |
| `sidecar.dashboards.provider.orgid`       | Id of the organisation, to which the dashboards should be added | `1`                                   |
| `sidecar.dashboards.provider.folder`      | Logical folder in which grafana groups dashboards | `""`                                                |
| `sidecar.dashboards.provider.disableDelete` | Activate to avoid the deletion of imported dashboards | `false`                                       |
| `sidecar.dashboards.provider.allowUiUpdates` | Allow updating provisioned dashboards from the UI | `false`                                          |
| `sidecar.dashboards.provider.type`        | Provider type                                 | `file`                                                  |
| `sidecar.dashboards.provider.foldersFromFilesStructure`        | Allow Grafana to replicate dashboard structure from filesystem.                                 | `false`                                                  |
| `sidecar.dashboards.watchMethod`          | Method to use to detect ConfigMap changes. With WATCH the sidecar will do a WATCH requests, with SLEEP it will list all ConfigMaps, then sleep for 60 seconds. | `WATCH` |
| `sidecar.skipTlsVerify`                   | Set to true to skip tls verification for kube api calls | `nil`                                         |
| `sidecar.dashboards.label`                | Label that config maps with dashboards should have to be added | `grafana_dashboard`                                |
| `sidecar.dashboards.labelValue`                | Label value that config maps with dashboards should have to be added | `""`                                |
| `sidecar.dashboards.folder`               | Folder in the pod that should hold the collected dashboards (unless `sidecar.dashboards.defaultFolderName` is set). This path will be mounted. | `/tmp/dashboards`    |
| `sidecar.dashboards.folderAnnotation`     | The annotation the sidecar will look for in configmaps to override the destination folder for files | `nil`                                                  |
| `sidecar.dashboards.defaultFolderName`    | The default folder name, it will create a subfolder under the `sidecar.dashboards.folder` and put dashboards in there instead | `nil`                                |
| `sidecar.dashboards.searchNamespace`      | Namespaces list. If specified, the sidecar will search for dashboards config-maps  inside these namespaces. Otherwise the namespace in which the sidecar is running will be used. It's also possible to specify ALL to search in all namespaces. | `nil`                                |
| `sidecar.dashboards.script`               | Absolute path to shell script to execute after a configmap got reloaded. | `nil`                                |
| `sidecar.dashboards.resource`             | Should the sidecar looks into secrets, configmaps or both. | `both`                               |
| `sidecar.dashboards.extraMounts`          | Additional dashboard sidecar volume mounts. | `[]`                               |
| `sidecar.datasources.enabled`             | Enables the cluster wide search for datasources and adds/updates/deletes them in grafana |`false`       |
| `sidecar.datasources.label`               | Label that config maps with datasources should have to be added | `grafana_datasource`                               |
| `sidecar.datasources.labelValue`          | Label value that config maps with datasources should have to be added | `""`                                |
| `sidecar.datasources.searchNamespace`     | Namespaces list. If specified, the sidecar will search for datasources config-maps  inside these namespaces. Otherwise the namespace in which the sidecar is running will be used. It's also possible to specify ALL to search in all namespaces. | `nil`                               |
| `sidecar.datasources.watchMethod`         | Method to use to detect ConfigMap changes. With WATCH the sidecar will do a WATCH requests, with SLEEP it will list all ConfigMaps, then sleep for 60 seconds. | `WATCH` |
| `sidecar.datasources.resource`            | Should the sidecar looks into secrets, configmaps or both. | `both`                               |
| `sidecar.datasources.reloadURL`           | Full url of datasource configuration reload API endpoint, to invoke after a config-map change | `"http://localhost:3000/api/admin/provisioning/datasources/reload"` |
| `sidecar.datasources.skipReload`          | Enabling this omits defining the REQ_URL and REQ_METHOD environment variables | `false` |
| `sidecar.datasources.initDatasources`     | Set to true to deploy the datasource sidecar as an initContainer in addition to a container. This is needed if skipReload is true, to load any datasources defined at startup time. | `false` |
| `sidecar.notifiers.enabled`               | Enables the cluster wide search for notifiers and adds/updates/deletes them in grafana | `false`        |
| `sidecar.notifiers.label`                 | Label that config maps with notifiers should have to be added | `grafana_notifier`                               |
| `sidecar.notifiers.labelValue`            | Label value that config maps with notifiers should have to be added | `""`                                |
| `sidecar.notifiers.searchNamespace`       | Namespaces list. If specified, the sidecar will search for notifiers config-maps (or secrets) inside these namespaces. Otherwise the namespace in which the sidecar is running will be used. It's also possible to specify ALL to search in all namespaces. | `nil`                               |
| `sidecar.notifiers.watchMethod`           | Method to use to detect ConfigMap changes. With WATCH the sidecar will do a WATCH requests, with SLEEP it will list all ConfigMaps, then sleep for 60 seconds. | `WATCH` |
| `sidecar.notifiers.resource`              | Should the sidecar looks into secrets, configmaps or both. | `both`                               |
| `sidecar.notifiers.reloadURL`             | Full url of notifier configuration reload API endpoint, to invoke after a config-map change | `"http://localhost:3000/api/admin/provisioning/notifications/reload"` |
| `sidecar.notifiers.skipReload`            | Enabling this omits defining the REQ_URL and REQ_METHOD environment variables | `false` |
| `sidecar.notifiers.initNotifiers`         | Set to true to deploy the notifier sidecar as an initContainer in addition to a container. This is needed if skipReload is true, to load any notifiers defined at startup time. | `false` |
| `smtp.existingSecret`                     | The name of an existing secret containing the SMTP credentials. | `""`                                  |
| `smtp.userKey`                            | The key in the existing SMTP secret containing the username. | `"user"`                                 |
| `smtp.passwordKey`                        | The key in the existing SMTP secret containing the password. | `"password"`                             |
| `admin.existingSecret`                    | The name of an existing secret containing the admin credentials (can be templated). | `""`                                 |
| `admin.userKey`                           | The key in the existing admin secret containing the username. | `"admin-user"`                          |
| `admin.passwordKey`                       | The key in the existing admin secret containing the password. | `"admin-password"`                      |
| `serviceAccount.autoMount`                | Automount the service account token in the pod| `true`                                                  |
| `serviceAccount.annotations`              | ServiceAccount annotations                    |                                                         |
| `serviceAccount.create`                   | Create service account                        | `true`                                                  |
| `serviceAccount.labels`                   | ServiceAccount labels                         | `{}`                                                    |
| `serviceAccount.name`                     | Service account name to use, when empty will be set to created account if `serviceAccount.create` is set else to `default` | `` |
| `serviceAccount.nameTest`                 | Service account name to use for test, when empty will be set to created account if `serviceAccount.create` is set else to `default` | `nil` |
| `rbac.create`                             | Create and use RBAC resources                 | `true`                                                  |
| `rbac.namespaced`                         | Creates Role and Rolebinding instead of the default ClusterRole and ClusteRoleBindings for the grafana instance  | `false` |
| `rbac.useExistingRole`                    | Set to a rolename to use existing role - skipping role creating - but still doing serviceaccount and rolebinding to the rolename set here. | `nil` |
| `rbac.pspEnabled`                         | Create PodSecurityPolicy (with `rbac.create`, grant roles permissions as well) | `true`                 |
| `rbac.pspUseAppArmor`                     | Enforce AppArmor in created PodSecurityPolicy (requires `rbac.pspEnabled`)  | `true`                    |
| `rbac.extraRoleRules`                     | Additional rules to add to the Role           | []                                                      |
| `rbac.extraClusterRoleRules`              | Additional rules to add to the ClusterRole    | []                                                      |
| `command`                     | Define command to be executed by grafana container at startup  | `nil`                                              |
| `testFramework.enabled`                   | Whether to create test-related resources      | `true`                                                  |
| `testFramework.image`                     | `test-framework` image repository.            | `bats/bats`                                             |
| `testFramework.tag`                       | `test-framework` image tag.                   | `v1.4.1`                                                |
| `testFramework.imagePullPolicy`           | `test-framework` image pull policy.           | `IfNotPresent`                                          |
| `testFramework.securityContext`           | `test-framework` securityContext              | `{}`                                                    |
| `downloadDashboards.env`                  | Environment variables to be passed to the `download-dashboards` container | `{}`                        |
| `downloadDashboards.envFromSecret`        | Name of a Kubernetes secret (must be manually created in the same namespace) containing values to be added to the environment. Can be templated | `""` |
| `downloadDashboards.resources`            | Resources of `download-dashboards` container  | `{}`                                                    |
| `downloadDashboardsImage.repository`      | Curl docker image repo                        | `curlimages/curl`                                       |
| `downloadDashboardsImage.tag`             | Curl docker image tag                         | `7.73.0`                                                |
| `downloadDashboardsImage.sha`             | Curl docker image sha (optional)              | `""`                                                    |
| `downloadDashboardsImage.pullPolicy`      | Curl docker image pull policy                 | `IfNotPresent`                                          |
| `namespaceOverride`                       | Override the deployment namespace             | `""` (`Release.Namespace`)                              |
| `serviceMonitor.enabled`                  | Use servicemonitor from prometheus operator   | `false`                                                 |
| `serviceMonitor.namespace`                | Namespace this servicemonitor is installed in |                                                         |
| `serviceMonitor.interval`                 | How frequently Prometheus should scrape       | `1m`                                                    |
| `serviceMonitor.path`                     | Path to scrape                                | `/metrics`                                              |
| `serviceMonitor.scheme`                   | Scheme to use for metrics scraping            | `http`                                                  |
| `serviceMonitor.tlsConfig`                | TLS configuration block for the endpoint      | `{}`                                                    |
| `serviceMonitor.labels`                   | Labels for the servicemonitor passed to Prometheus Operator      |  `{}`                                |
| `serviceMonitor.scrapeTimeout`            | Timeout after which the scrape is ended       | `30s`                                                   |
| `serviceMonitor.relabelings`              | MetricRelabelConfigs to apply to samples before ingestion.  | `[]`                                      |
| `revisionHistoryLimit`                    | Number of old ReplicaSets to retain           | `10`                                                    |
| `imageRenderer.enabled`                    | Enable the image-renderer deployment & service                                     | `false`                          |
| `imageRenderer.image.repository`           | image-renderer Image repository                                                    | `grafana/grafana-image-renderer` |
| `imageRenderer.image.tag`                  | image-renderer Image tag                                                           | `latest`                         |
| `imageRenderer.image.sha`                  | image-renderer Image sha (optional)                                                | `""`                             |
| `imageRenderer.image.pullPolicy`           | image-renderer ImagePullPolicy                                                     | `Always`                         |
| `imageRenderer.env`                        | extra env-vars for image-renderer                                                  | `{}`                             |
| `imageRenderer.serviceAccountName`         | image-renderer deployment serviceAccountName                                       | `""`                             |
| `imageRenderer.securityContext`            | image-renderer deployment securityContext                                          | `{}`                             |
| `imageRenderer.hostAliases`                | image-renderer deployment Host Aliases                                             | `[]`                             |
| `imageRenderer.priorityClassName`          | image-renderer deployment priority class                                           | `''`                             |
| `imageRenderer.service.enabled`            | Enable the image-renderer service                                                  | `true`                           |
| `imageRenderer.service.portName`           | image-renderer service port name                                                   | `http`                           |
| `imageRenderer.service.port`               | image-renderer port used by deployment                                             | `8081`                           |
| `imageRenderer.service.targetPort`         | image-renderer service port used by service                                        | `8081`                           |
| `imageRenderer.appProtocol`                | Adds the appProtocol field to the service                                          | ``                               |
| `imageRenderer.grafanaSubPath`             | Grafana sub path to use for image renderer callback url                            | `''`                             |
| `imageRenderer.podPortName`                | name of the image-renderer port on the pod                                         | `http`                           |
| `imageRenderer.revisionHistoryLimit`       | number of image-renderer replica sets to keep                                      | `10`                             |
| `imageRenderer.networkPolicy.limitIngress` | Enable a NetworkPolicy to limit inbound traffic from only the created grafana pods | `true`                           |
| `imageRenderer.networkPolicy.limitEgress`  | Enable a NetworkPolicy to limit outbound traffic to only the created grafana pods  | `false`                          |
| `imageRenderer.resources`                  | Set resource limits for image-renderer pdos                                        | `{}`                             |
| `imageRenderer.nodeSelector`               | Node labels for pod assignment                | `{}`                                                    |
| `imageRenderer.tolerations`                | Toleration labels for pod assignment          | `[]`                                                    |
| `imageRenderer.affinity`                   | Affinity settings for pod assignment          | `{}`                                                    |
| `networkPolicy.enabled`                    | Enable creation of NetworkPolicy resources.                                                                              | `false`             |
| `networkPolicy.allowExternal`              | Don't require client label for connections                                                                               | `true`              |
| `networkPolicy.explicitNamespacesSelector` | A Kubernetes LabelSelector to explicitly select namespaces from which traffic could be allowed                           | `{}`                |
| `networkPolicy.ingress`                    | Enable the creation of an ingress network policy             | `true`    |
| `networkPolicy.egress.enabled`             | Enable the creation of an egress network policy              | `false`   |
| `networkPolicy.egress.ports`               | An array of ports to allow for the egress                    | `[]`    |
| `enableKubeBackwardCompatibility`          | Enable backward compatibility of kubernetes where pod's defintion version below 1.13 doesn't have the enableServiceLinks option  | `false`     |

* see https://github.com/grafana/helm-charts/tree/main/charts/grafana for more detail
