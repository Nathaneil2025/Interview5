pipeline {
    agent any

    options {
        buildDiscarder(logRotator(numToKeepStr: '10'))
        timestamps()
        timeout(time: 30, unit: 'MINUTES')
        disableConcurrentBuilds()
    }

    triggers {
        pollSCM('H/5 * * * *')
    }

    environment {
        AWS_REGION = 'eu-central-1'
        ECR_REGISTRY = '518394500999.dkr.ecr.eu-central-1.amazonaws.com'
        ECR_REPO_PREFIX = 'jenkins-monorepo-ci'
        IMAGE_TAG = "ci-${env.GIT_COMMIT?.take(7) ?: 'latest'}"
        SERVICES = 'user-service,transaction-service,notification-service'
    }

    stages {
        stage('Detect Changes') {
            steps {
                script {
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

        stage('ECR Login') {
            when {
                expression { env.CHANGED_SERVICES?.trim() }
            }
            steps {
                sh '''
                    aws ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin ${ECR_REGISTRY}
                '''
            }
        }

        stage('Docker Build & Push') {
            when {
                expression { env.CHANGED_SERVICES?.trim() }
            }
            parallel {
                stage('Build & Push User Service') {
                    when {
                        expression { env.CHANGED_SERVICES?.contains('user-service') }
                    }
                    steps {
                        buildAndPushImage('user-service')
                    }
                }

                stage('Build & Push Transaction Service') {
                    when {
                        expression { env.CHANGED_SERVICES?.contains('transaction-service') }
                    }
                    steps {
                        buildAndPushImage('transaction-service')
                    }
                }

                stage('Build & Push Notification Service') {
                    when {
                        expression { env.CHANGED_SERVICES?.contains('notification-service') }
                    }
                    steps {
                        buildAndPushImage('notification-service')
                    }
                }
            }
        }

        stage('Manual Approval') {
            when {
                expression { env.CHANGED_SERVICES?.trim() }
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
        def changes = ''
        if (env.GIT_PREVIOUS_COMMIT) {
            changes = sh(script: "git diff --name-only ${env.GIT_PREVIOUS_COMMIT} ${env.GIT_COMMIT}", returnStdout: true).trim()
        } else {
            changes = sh(script: "git diff --name-only origin/main...HEAD 2>/dev/null || git diff --name-only HEAD~1 HEAD", returnStdout: true).trim()
        }
        
        echo "Changed files:\n${changes}"
        
        allServices.each { service ->
            if (changes.contains(service) || changes.contains('shared/')) {
                changedServices.add(service)
            }
        }
        
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
    junit allowEmptyResults: true, testResults: "${service}/**/junit.xml"
    publishHTML(target: [
        allowMissing: true,
        alwaysLinkToLastBuild: false,
        keepAll: true,
        reportDir: "${service}/coverage",
        reportFiles: 'index.html',
        reportName: "${service} Coverage Report"
    ])
}

def buildAndPushImage(String service) {
    echo "Building and pushing Docker image for ${service}..."
    retry(2) {
        sh """
            docker build -t ${ECR_REGISTRY}/${ECR_REPO_PREFIX}/${service}:${IMAGE_TAG} ./${service}
            docker tag ${ECR_REGISTRY}/${ECR_REPO_PREFIX}/${service}:${IMAGE_TAG} ${ECR_REGISTRY}/${ECR_REPO_PREFIX}/${service}:latest
            docker push ${ECR_REGISTRY}/${ECR_REPO_PREFIX}/${service}:${IMAGE_TAG}
            docker push ${ECR_REGISTRY}/${ECR_REPO_PREFIX}/${service}:latest
        """
    }
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
}