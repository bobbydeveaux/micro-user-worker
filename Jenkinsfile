node('master') {
  stage('Unit Tests') {
    git url: "https://github.com/bobbydeveaux/micro-user-worker.git"
  }
  stage('Build Bin') {
    sh "go get -v -d ./..."
    sh "CGO_ENABLED=0 GOOS=linux go build -o micro-user-worker ."
  }
  stage('Build Image') {
    sh "oc start-build micro-user-worker --from-file=. --follow"
  }
  stage('Deploy') {
    openshiftDeploy depCfg: 'micro-user-worker', namespace: 'fbac'
    openshiftVerifyDeployment depCfg: 'micro-user-worker', replicaCount: 1, verifyReplicaCount: true
  }
}