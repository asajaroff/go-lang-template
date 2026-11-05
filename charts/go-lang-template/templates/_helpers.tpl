{{- define "go-lang-template.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "go-lang-template.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{- define "go-lang-template.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "go-lang-template.labels" -}}
helm.sh/chart: {{ include "go-lang-template.chart" . }}
{{ include "go-lang-template.selectorLabels" . }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{- define "go-lang-template.selectorLabels" -}}
app.kubernetes.io/name: {{ include "go-lang-template.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- define "go-lang-template.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
{{ default (include "go-lang-template.fullname" .) .Values.serviceAccount.name }}
{{- else -}}
{{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{- define "go-lang-template.image" -}}
{{- printf "%s:%s" .Values.image.repository (default .Chart.AppVersion .Values.image.tag) -}}
{{- end -}}

{{/*
Env shared by the Deployment and the migration Job — keeps the app
container and the one-off admin process (Factor XII) configured identically.
Non-secret settings come from the chart's ConfigMap (envFrom); the DB URL
is the one value pulled from a Secret the environment attaches (Factor IV).
*/}}
{{- define "go-lang-template.envFrom" -}}
- configMapRef:
    name: {{ include "go-lang-template.fullname" . }}
{{- end -}}

{{- define "go-lang-template.secretEnv" -}}
- name: DATABASE_URL
  valueFrom:
    secretKeyRef:
      name: {{ .Values.env.existingSecret }}
      key: {{ .Values.env.existingSecretDatabaseUrlKey }}
{{- end -}}
