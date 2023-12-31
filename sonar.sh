#!/bin/bash

sonar-scanner -X \
    -Dsonar.host.url=${SONAR_HOST_URL} \
    -Dsonar.language=go \
    -Dsonar.projectName=Apigee_Agents \
    -Dsonar.projectVersion=1.0 \
    -Dsonar.projectKey=Apigee_Agents \
    -Dsonar.sourceEncoding=UTF-8 \
    -Dsonar.projectBaseDir=${WORKSPACE} \
    -Dsonar.sources=./client/pkg/,./discovery/,./traceability/ \
    -Dsonar.tests=./client/pkg/,./discovery/,./traceability/ \
	-Dsonar.exclusions=**/*.json \
    -Dsonar.test.inclusions=**/*test*.go \
    -Dsonar.go.tests.reportPaths=goreport.json \
    -Dsonar.go.coverage.reportPaths=gocoverage.out \
    -Dsonar.issuesReport.console.enable=true \
    -Dsonar.report.export.path=sonar-report.json \
    -Dsonar.verbose=true
