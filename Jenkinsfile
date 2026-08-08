pipeline {
    agent any

    options {
        timestamps()
        disableConcurrentBuilds()
        skipDefaultCheckout(true)
    }

    parameters {
        booleanParam(
            name: 'RUN_INTEGRATION',
            defaultValue: false,
            description: 'Run Docker-backed integration and Kafka broker tests.'
        )
        booleanParam(
            name: 'BUILD_CONTAINERS',
            defaultValue: true,
            description: 'Build production runtime images without publishing them.'
        )
        booleanParam(
            name: 'DEPLOY_STAGING',
            defaultValue: false,
            description: 'Promote an existing immutable image tag on the staging agent.'
        )
        string(
            name: 'IMAGE_TAG',
            defaultValue: '',
            description: 'Existing GHCR image tag to promote when DEPLOY_STAGING is enabled.'
        )
        string(
            name: 'IMAGE_REGISTRY',
            defaultValue: 'ghcr.io/ORG/REPO',
            description: 'GHCR repository containing the promoted runtime images.'
        )
    }

    environment {
        GOTOOLCHAIN = 'local'
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Format and vet') {
            steps {
                sh 'make ci-format'
                sh 'make ci-vet'
            }
        }

        stage('Unit tests') {
            steps {
                sh 'make ci-unit'
            }
        }

        stage('Race detector') {
            steps {
                sh 'make ci-race'
            }
        }

        stage('Generated contracts') {
            steps {
                sh 'make ci-generated'
            }
        }

        stage('Container builds') {
            when {
                expression { params.BUILD_CONTAINERS }
            }
            steps {
                sh 'make ci-containers'
            }
        }

        stage('Staging delivery') {
            when {
                expression { params.DEPLOY_STAGING }
            }
            agent {
                label 'notezy-staging'
            }
            steps {
                checkout scm
                withCredentials([
                    usernamePassword(
                        credentialsId: 'notezy-ghcr',
                        usernameVariable: 'GHCR_USERNAME',
                        passwordVariable: 'GHCR_TOKEN'
                    )
                ]) {
                    sh '''
                        set -eu
                        test -n "${IMAGE_TAG}"
                        : "${COMPOSE_ENV_FILE:=/etc/notezy/staging.env}"
                        echo "${GHCR_TOKEN}" | docker login ghcr.io --username "${GHCR_USERNAME}" --password-stdin
                        make staging-deploy
                        make staging-smoke
                    '''
                }
            }
        }

        stage('Integration tests') {
            when {
                expression { params.RUN_INTEGRATION }
            }
            steps {
                sh 'make test-integration'
                sh 'make test-integration-kafka'
            }
        }
    }

    post {
        always {
            archiveArtifacts artifacts: '**/coverage.out, **/test-results/**, **/integration-results/**, staging-compose.log', allowEmptyArchive: true
        }
    }
}
