pipeline {
    agent { label 'slave' }
    stages{
        stage('Trigger neofs-node CI') {
            build job: 'neofs_node_ci', parameters: [string(name: 'branch', value: env.BRANCH_NAME)]
        }
    }
}
