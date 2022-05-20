notifications:
  enabled: true
  nodeSelector:
    fusion_node_type: infra
  metrics:
    enabled: true
  secret:
    create: true
    items:
      slack-token: ${slackNotificationBotToken}
  templates:
    template.app-deployed: |
      email:
        subject: New version of an application {{.app.metadata.name}} is up and running.
      message: |
        {{if eq .serviceType "slack"}}:white_check_mark:{{end}} Application {{.app.metadata.name}} is now running new version of deployments manifests.
      slack:
        attachments: |
          [{
            "title": "{{ .app.metadata.name}}",
            "title_link":"{{.context.argocdUrl}}/applications/{{.app.metadata.name}}",
            "color": "#18be52",
            "fields": [
            {
              "title": "Sync Status",
              "value": "{{.app.status.sync.status}}",
              "short": true
            },
            {
              "title": "Repository",
              "value": "{{.app.spec.source.repoURL}}",
              "short": true
            },
            {
              "title": "Revision",
              "value": "{{.app.status.sync.revision}}",
              "short": true
            }
            {{range $index, $c := .app.status.conditions}}
            {{if not $index}},{{end}}
            {{if $index}},{{end}}
            {
              "title": "{{$c.type}}",
              "value": "{{$c.message}}",
              "short": true
            }
            {{end}}
            ]
          }]
    template.app-health-degraded: |
      email:
        subject: Application {{.app.metadata.name}} has degraded.
      message: |
        {{if eq .serviceType "slack"}}:exclamation:{{end}} Application {{.app.metadata.name}} has degraded.
        Application details: {{.context.argocdUrl}}/applications/{{.app.metadata.name}}.
      slack:
        attachments: |-
          [{
            "title": "{{ .app.metadata.name}}",
            "title_link": "{{.context.argocdUrl}}/applications/{{.app.metadata.name}}",
            "color": "#f4c030",
            "fields": [
            {
              "title": "Sync Status",
              "value": "{{.app.status.sync.status}}",
              "short": true
            },
            {
              "title": "Repository",
              "value": "{{.app.spec.source.repoURL}}",
              "short": true
            }
            {{range $index, $c := .app.status.conditions}}
            {{if not $index}},{{end}}
            {{if $index}},{{end}}
            {
              "title": "{{$c.type}}",
              "value": "{{$c.message}}",
              "short": true
            }
            {{end}}
            ]
          }]
    template.app-sync-failed: |
      email:
        subject: Failed to sync application {{.app.metadata.name}}.
      message: |
        {{if eq .serviceType "slack"}}:exclamation:{{end}}  The sync operation of application {{.app.metadata.name}} has failed at {{.app.status.operationState.finishedAt}} with the following error: {{.app.status.operationState.message}}
        Sync operation details are available at: {{.context.argocdUrl}}/applications/{{.app.metadata.name}}?operation=true .
      slack:
        attachments: |-
          [{
            "title": "{{ .app.metadata.name}}",
            "title_link":"{{.context.argocdUrl}}/applications/{{.app.metadata.name}}",
            "color": "#E96D76",
            "fields": [
            {
              "title": "Sync Status",
              "value": "{{.app.status.sync.status}}",
              "short": true
            },
            {
              "title": "Repository",
              "value": "{{.app.spec.source.repoURL}}",
              "short": true
            }
            {{range $index, $c := .app.status.conditions}}
            {{if not $index}},{{end}}
            {{if $index}},{{end}}
            {
              "title": "{{$c.type}}",
              "value": "{{$c.message}}",
              "short": true
            }
            {{end}}
            ]
          }]
    template.app-sync-running: |
      email:
        subject: Start syncing application {{.app.metadata.name}}.
      message: |
        The sync operation of application {{.app.metadata.name}} has started at {{.app.status.operationState.startedAt}}.
        Sync operation details are available at: {{.context.argocdUrl}}/applications/{{.app.metadata.name}}?operation=true .
      slack:
        attachments: |-
          [{
            "title": "{{ .app.metadata.name}}",
            "title_link":"{{.context.argocdUrl}}/applications/{{.app.metadata.name}}",
            "color": "#0DADEA",
            "fields": [
            {
              "title": "Sync Status",
              "value": "{{.app.status.sync.status}}",
              "short": true
            },
            {
              "title": "Repository",
              "value": "{{.app.spec.source.repoURL}}",
              "short": true
            }
            {{range $index, $c := .app.status.conditions}}
            {{if not $index}},{{end}}
            {{if $index}},{{end}}
            {
              "title": "{{$c.type}}",
              "value": "{{$c.message}}",
              "short": true
            }
            {{end}}
            ]
          }]
    template.app-sync-status-unknown: |
      email:
        subject: Application {{.app.metadata.name}} sync status is 'Unknown'
      message: |
        {{if eq .serviceType "slack"}}:exclamation:{{end}} Application {{.app.metadata.name}} sync is 'Unknown'.
        Application details: {{.context.argocdUrl}}/applications/{{.app.metadata.name}}.
        {{if ne .serviceType "slack"}}
        {{range $c := .app.status.conditions}}
            * {{$c.message}}
        {{end}}
        {{end}}
      slack:
        attachments: |-
          [{
            "title": "{{ .app.metadata.name}}",
            "title_link":"{{.context.argocdUrl}}/applications/{{.app.metadata.name}}",
            "color": "#E96D76",
            "fields": [
            {
              "title": "Sync Status",
              "value": "{{.app.status.sync.status}}",
              "short": true
            },
            {
              "title": "Repository",
              "value": "{{.app.spec.source.repoURL}}",
              "short": true
            }
            {{range $index, $c := .app.status.conditions}}
            {{if not $index}},{{end}}
            {{if $index}},{{end}}
            {
              "title": "{{$c.type}}",
              "value": "{{$c.message}}",
              "short": true
            }
            {{end}}
            ]
          }]
    template.app-sync-succeeded: |
      email:
        subject: Application {{.app.metadata.name}} has been successfully synced.
      message: |
        {{if eq .serviceType "slack"}}:white_check_mark:{{end}} Application {{.app.metadata.name}} has been successfully synced at {{.app.status.operationState.finishedAt}}.
        Sync operation details are available at: {{.context.argocdUrl}}/applications/{{.app.metadata.name}}?operation=true .
      slack:
        attachments: "[{\n  \"title\": \"{{ .app.metadata.name}}\",\n  \"title_link\":\"{{.context.argocdUrl}}/applications/{{.app.metadata.name}}\",\n  \"color\": \"#18be52\",\n  \"fields\": [\n  {\n    \"title\": \"Sync Status\",\n    \"value\": \"{{.app.status.sync.status}}\",\n    \"short\": true\n  },\n  {\n    \"title\": \"Repository\",\n    \"value\": \"{{.app.spec.source.repoURL}}\",\n    \"short\": true\n  }\n  {{range $index, $c := .app.status.conditions}}\n  {{if not $index}},{{end}}\n  {{if $index}},{{end}}\n  {\n    \"title\": \"{{$c.type}}\",\n    \"value\": \"{{$c.message}}\",\n    \"short\": true\n  }\n  {{end}}\n  ]\n}]    "
  triggers:
    trigger.on-deployed: |
      - description: Application is synced and healthy. Triggered once per commit.
        oncePer: app.status.sync.revision
        send:
        - app-deployed
        when: app.status.operationState.phase in ['Succeeded'] and app.status.health.status == 'Healthy'
    trigger.on-health-degraded: |
      - description: Application has degraded
        send:
        - app-health-degraded
        when: app.status.health.status == 'Degraded'
    trigger.on-sync-failed: |
      - description: Application syncing has failed
        send:
        - app-sync-failed
        when: app.status.operationState.phase in ['Error', 'Failed']
    trigger.on-sync-running: |
      - description: Application is being synced
        send:
        - app-sync-running
        when: app.status.operationState.phase in ['Running']
    trigger.on-sync-status-unknown: |
      - description: Application status is 'Unknown'
        send:
        - app-sync-status-unknown
        when: app.status.sync.status == 'Unknown'
    trigger.on-sync-succeeded: |
      - description: Application syncing has succeeded
        send:
        - app-sync-succeeded
        when: app.status.operationState.phase in ['Succeeded']
    # For more information: https://argocd-notifications.readthedocs.io/en/stable/triggers/#default-triggers
    defaultTriggers: |
      - on-sync-status-unknown
      - app-deployed
      - app-health-degraded
      - app-sync-failed
      - on-sync-running
  # Subscriptions
  subscriptions:
    - recipients:
      - slack:cloud-infrastructure-alerts
      triggers:
      - on-sync-status-unknown
      - app-deployed
      - app-health-degraded
      - app-sync-failed