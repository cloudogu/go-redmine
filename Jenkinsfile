#!groovy
@Library(['github.com/cloudogu/zalenium-build-lib@v2.1.0', 'github.com/cloudogu/ces-build-lib@1.44.3']) _
import com.cloudogu.ces.zaleniumbuildlib.*
import com.cloudogu.ces.cesbuildlib.*

node('docker') {
    timestamps {
        branch = "${env.BRANCH_NAME}"

        stage('Checkout') {
            checkout scm
        }

        stage('SonarQube') {
            def scannerHome = tool name: 'sonar-scanner', type: 'hudson.plugins.sonar.SonarRunnerInstallation'
            withSonarQubeEnv {
                if (branch == "main") {
                    echo "This branch has been detected as the main branch."
                    sh "${scannerHome}/bin/sonar-scanner"
                } else if (branch == "develop") {
                    echo "This branch has been detected as the develop branch."
                    sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.BRANCH_NAME} -Dsonar.branch.target=master"
                } else if (env.CHANGE_TARGET) {
                    echo "This branch has been detected as a pull request."
                    sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.CHANGE_BRANCH}-PR${env.CHANGE_ID} -Dsonar.branch.target=${env.CHANGE_TARGET}"
                } else if (branch.startsWith("feature/")) {
                    echo "This branch has been detected as a feature branch."
                    sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.BRANCH_NAME} -Dsonar.branch.target=develop"
                }
            }
            timeout(time: 2, unit: 'MINUTES') { // Needed when there is no webhook for example
                def qGate = waitForQualityGate()
                if (qGate.status != 'OK') {
                    unstable("Pipeline unstable due to SonarQube quality gate failure")
                }
            }
        }
    }
}
