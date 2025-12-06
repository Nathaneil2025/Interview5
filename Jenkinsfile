pipeline {
    agent any

    options {
        buildDiscarder(logRotator(numToKeepStr: '10'))
        timestamps()
        timeout(time: 30, unit: 'MINUTES')
        disableConcurrentBuilds()
    }

    triggers {
        // Trigger on push to main or develop
        pollSCM('H/5 * * * *')
    }

    environment {
        DOCKER_REGISTRY = 'your-registry.com'  // Update with your registry
        IMAGE_TAG = "ci-${env.GIT_COMMIT?.take(7) ?: 'latest'}"
        SERVICES = 'user-service,transaction-service,notification-service'
    }

    stages {
        stage('Detect Changes') {
            steps {
                script {
                    // Get list of changed files
                    def changedServices = detectChangedServices()
                    env.CHANGED_SERVICES = changedServices.join(',')
                    
                    if (changedServices.isEmpty()) {
                        echo "No service changes detected. Skipping pipeline."
                        currentBuild.result = 'SUCCESS'
                    } else {
                        echo "Changed services: ${changedServices}"
                    }
                }
            }
        }

        stage('Quality Gates') {
            when {
                expression { env.CHANGED_SERVICES?.trim() }
            }
            parallel {
                stage('User Service') {
                    when {
                        expression { env.CHANGED_SERVICES?.contains('user-service') }
                    }
                    stages {
                        stage('Lint - User Service') {
                            steps {
                                runLint('user-service')
                            }
                        }
                        stage('Test - User Service') {
                            steps {
                                runTests('user-service')
                            }
                            post {
                                always {
                                    publishTestResults('user-service')
                                }
                            }
                        }
                        stage('Security Scan - User Service') {
                            steps {
                                runSecurityScan('user-service')
                            }
                        }
                    }
                }

                stage('Transaction Service') {
                    when {
                        expression { env.CHANGED_SERVICES?.contains('transaction-service') }
                    }
                    stages {
                        stage('Lint - Transaction Service') {
                            steps {
                                runLint('transaction-service')
                            }
                        }
                        stage('Test - Transaction Service') {
                            steps {
                                runTests('transaction-service')
                            }
                            post {
                                always {
                                    publishTestResults('transaction-service')
                                }
                            }
                        }
                        stage('Security Scan - Transaction Service') {
                            steps {
                                runSecurityScan('transaction-service')
                            }
                        }
                    }
                }

                stage('Notification Service') {
                    when {
                        expression { env.CHANGED_SERVICES?.contains('notification-service') }
                    }
                    stages {
                        stage('Lint - Notification Service') {
                            steps {
                                runLint('notification-service')
                            }
                        }
                        stage('Test - Notification Service') {
                            steps {
                                runTests('notification-service')
                            }
                            post {
                                always {
                                    publishTestResults('notification-service')
                                }
                            }
                        }
                        stage('Security Scan - Notification Service') {
                            steps {
                                runSecurityScan('notification-service')
                            }
                        }
                    }
                }
            }
        }

        stage('Docker Build') {
            when {
                expression { env.CHANGED_SERVICES?.trim() }
            }
            parallel {
                stage('Build User Service Image') {
                    when {
                        expression { env.CHANGED_SERVICES?.contains('user-service') }
                    }
                    steps {
                        buildDockerImage('user-service')
                    }
                }

                stage('Build Transaction Service Image') {
                    when {
                        expression { env.CHANGED_SERVICES?.contains('transaction-service') }
                    }
                    steps {
                        buildDockerImage('transaction-service')
                    }
                }

                stage('Build Notification Service Image') {
                    when {
                        expression { env.CHANGED_SERVICES?.contains('notification-service') }
                    }
                    steps {
                        buildDockerImage('notification-service')
                    }
                }
            }
        }

        stage('Push to Registry') {
            when {
                allOf {
                    expression { env.CHANGED_SERVICES?.trim() }
                    branch pattern: "(main|develop)", comparator: "REGEXP"
                }
            }
            steps {
                script {
                    def services = env.CHANGED_SERVICES.split(',')
                    services.each { service ->
                        pushDockerImage(service)
                    }
                }
            }
        }

        stage('Manual Approval') {
            when {
                allOf {
                    expression { env.CHANGED_SERVICES?.trim() }
                    branch 'main'
                }
            }
            steps {
                script {
                    timeout(time: 1, unit: 'HOURS') {
                        input message: 'Approve deployment?',
                              ok: 'Deploy',
                              submitter: 'admin,deployer',
                              parameters: [
                                  string(name: 'APPROVER_NOTES', defaultValue: '', description: 'Optional notes')
                              ]
                    }
                    echo "Deployment approved. Ready to deploy!"
                    // CD would go here - out of scope for this assignment
                }
            }
        }
    }

    post {
        always {
            cleanWs()
        }
        success {
            script {
                sendNotification('SUCCESS')
            }
        }
        failure {
            script {
                sendNotification('FAILURE')
            }
        }
    }
}

// ==================== REUSABLE FUNCTIONS ====================

def detectChangedServices() {
    def changedServices = []
    def allServices = ['user-service', 'transaction-service', 'notification-service']
    
    try {
        // Get changed files compared to previous commit or main branch
        def changes = ''
        if (env.GIT_PREVIOUS_COMMIT) {
            changes = sh(script: "git diff --name-only ${env.GIT_PREVIOUS_COMMIT} ${env.GIT_COMMIT}", returnStdout: true).trim()
        } else {
            // First build or no previous commit - check against main
            changes = sh(script: "git diff --name-only origin/main...HEAD 2>/dev/null || git diff --name-only HEAD~1 HEAD", returnStdout: true).trim()
        }
        
        echo "Changed files:\n${changes}"
        
        allServices.each { service ->
            if (changes.contains(service) || changes.contains('shared/')) {
                changedServices.add(service)
            }
        }
        
        // If shared CI scripts changed, run all services
        if (changes.contains('shared/ci/') || changes.contains('Jenkinsfile')) {
            changedServices = allServices
        }
        
    } catch (Exception e) {
        echo "Error detecting changes: ${e.message}. Running all services."
        changedServices = allServices
    }
    
    return changedServices
}

def runLint(String service) {
    echo "Running lint for ${service}..."
    retry(2) {
        sh """
            chmod +x shared/ci/lint.sh
            ./shared/ci/lint.sh ${service}
        """
    }
}

def runTests(String service) {
    echo "Running tests for ${service}..."
    retry(2) {
        sh """
            chmod +x shared/ci/test.sh
            ./shared/ci/test.sh ${service}
        """
    }
}

def runSecurityScan(String service) {
    echo "Running security scan for ${service}..."
    sh """
        chmod +x shared/ci/scan.sh
        ./shared/ci/scan.sh ${service}
    """
}

def publishTestResults(String service) {
    echo "Publishing test results for ${service}..."
    
    // Publish JUnit results
    junit allowEmptyResults: true, testResults: "${service}/**/junit.xml"
    
    // Publish coverage if available
    publishHTML(target: [
        allowMissing: true,
        alwaysLinkToLastBuild: false,
        keepAll: true,
        reportDir: "${service}/coverage",
        reportFiles: 'index.html',
        reportName: "${service} Coverage Report"
    ])
}

def buildDockerImage(String service) {
    echo "Building Docker image for ${service}..."
    retry(2) {
        sh """
            docker build -t ${service}:${env.IMAGE_TAG} ./${service}
            docker tag ${service}:${env.IMAGE_TAG} ${service}:latest
        """
    }
}

def pushDockerImage(String service) {
    echo "Pushing Docker image for ${service}..."
    // Uncomment and configure for your registry
    // withCredentials([usernamePassword(credentialsId: 'docker-registry-creds', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
    //     sh """
    //         echo \$DOCKER_PASS | docker login ${env.DOCKER_REGISTRY} -u \$DOCKER_USER --password-stdin
    //         docker tag ${service}:${env.IMAGE_TAG} ${env.DOCKER_REGISTRY}/${service}:${env.IMAGE_TAG}
    //         docker push ${env.DOCKER_REGISTRY}/${service}:${env.IMAGE_TAG}
    //     """
    // }
    echo "Push to registry skipped (configure credentials first)"
}

def sendNotification(String status) {
    def color = status == 'SUCCESS' ? 'good' : 'danger'
    def emoji = status == 'SUCCESS' ? ':white_check_mark:' : ':x:'
    
    echo """
    ============================================
    Pipeline ${status}
    Job: ${env.JOB_NAME}
    Build: #${env.BUILD_NUMBER}
    Changed Services: ${env.CHANGED_SERVICES ?: 'None'}
    ============================================
    """
    
    // Uncomment to enable Slack notifications
    // slackSend(
    //     color: color,
    //     message: "${emoji} *${env.JOB_NAME}* #${env.BUILD_NUMBER} - ${status}\nChanged: ${env.CHANGED_SERVICES ?: 'None'}\n${env.BUILD_URL}"
    // )
    
    // Uncomment to enable MS Teams notifications
    // office365ConnectorSend(
    //     webhookUrl: env.TEAMS_WEBHOOK_URL,
    //     status: status,
    //     message: "Pipeline ${status}: ${env.JOB_NAME} #${env.BUILD_NUMBER}"
    // )
}
