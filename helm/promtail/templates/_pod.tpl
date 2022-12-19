{{/*
Pod template used in Daemonset and Deployment
*/}}
{{- define "promtail.podTemplate" -}}
metadata:
  labels:
    {{- include "promtail.selectorLabels" . | nindent 4 }}
    {{- with .Values.podLabels }}
    {{- toYaml . | nindent 4 }}
    {{- end }}
  annotations:
    checksum/config: {{ include (print .Template.BasePath "/secret.yaml") . | sha256sum }}
    {{- with .Values.podAnnotations }}
    {{- toYaml . | nindent 4 }}
    {{- end }}
spec:
  serviceAccountName: {{ include "promtail.serviceAccountName" . }}
  {{- include "promtail.enableServiceLinks" . | nindent 2 }}
  {{- with .Values.priorityClassName }}
  priorityClassName: {{ . }}
  {{- end }}
  {{- with .Values.initContainer }}
  initContainers:
    {{- toYaml . | nindent 4 }}
  {{- end }}
  {{- with .Values.imagePullSecrets }}
  imagePullSecrets:
    {{- toYaml . | nindent 4 }}
  {{- end }}
  containers:
    - name: promtail
      image: "{{ .Values.image.registry }}/{{ .Values.image.repository }}:{{ .Values.image.tag | default .Chart.AppVersion }}"
      imagePullPolicy: {{ .Values.image.pullPolicy }}
      args:
        - "-config.file=/etc/promtail/promtail.yaml"
        {{- with .Values.extraArgs }}
        {{- toYaml . | nindent 8 }}
        {{- end }}
      volumeMounts:
        - name: config
          mountPath: /etc/promtail
        - name: positions
          mountPath: /run/promtail
        {{- with .Values.extraVolumeMounts }}
        {{- toYaml . | nindent 8 }}
        {{- end }}
      env:
        - name: HOSTNAME
          valueFrom:
            fieldRef:
              fieldPath: spec.nodeName
      {{- with .Values.extraEnv }}
        {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .Values.extraEnvFrom }}
      envFrom:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      ports:
        - name: http-metrics
          containerPort: {{ .Values.config.serverPort }}
          protocol: TCP
        {{- range $key, $values := .Values.extraPorts }}
        - name: {{ .name | default $key }}
          containerPort: {{ $values.containerPort }}
          protocol: {{ $values.protocol | default "TCP" }}
        {{- end }}
      {{- with .Values.livenessProbe }}
      livenessProbe:
        {{- tpl (toYaml .) $ | nindent 8 }}
      {{- end }}
      {{- with .Values.readinessProbe }}
      readinessProbe:
        {{- tpl (toYaml .) $ | nindent 8 }}
      {{- end }}
      {{- with .Values.resources }}
      resources:
        {{- toYaml . | nindent 8 }}
      {{- end }}
    {{- if .Values.deployment.enabled }}
    {{- range $name, $values := .Values.extraContainers }}
    - name: {{ $name }}
      {{ toYaml $values | nindent 6 }}
    {{- end }}
    {{- end }}
  {{- with .Values.affinity }}
  affinity:
    {{- toYaml . | nindent 4 }}
  {{- end }}
  {{- with .Values.nodeSelector }}
  nodeSelector:
    {{- toYaml . | nindent 4 }}
  {{- end }}
  {{- with .Values.tolerations }}
  tolerations:
    {{- toYaml . | nindent 4 }}
  {{- end }}
  volumes:
    - name: config
      {{- if .Values.configmap.enabled }}
      configMap:
        name: {{ include "promtail.fullname" . }}
      {{- else }}
      secret:
        secretName: {{ include "promtail.fullname" . }}
      {{- end }}
    - name: positions
      persistentVolumeClaim:
        claimName: {{ include "promtail.fullname" . }}-positions
    {{- with .Values.extraVolumes }}
    {{- toYaml . | nindent 4 }}
    {{- end }}
{{- end }}
