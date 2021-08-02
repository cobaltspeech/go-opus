#!groovy
// Copyright (2021) Cobalt Speech and Language Inc.

// Keep only 10 builds on Jenkins
properties([
	buildDiscarder(logRotator(
		artifactDaysToKeepStr: '', artifactNumToKeepStr: '', daysToKeepStr: '', numToKeepStr: '10'))
])

// commit.setBuildStatus is defined in cobalt's shared jenkins library

// Build libopus for requested libc and architecture
def build_opus_linux(libc, arch) {
	stage("build-$libc-$arch") {
		sh "echo Building for $libc on $arch"
		// build for the given platform
		try {
			commit.setBuildStatus("build-${libc}-${arch}", "PENDING", "")
			echo "Building..."
			sh "source /setup_env.sh ${libc} ${arch} && build/build.sh"
			sh "mv opus/build/libopus.a lib/libopus.${libc}.${arch}.a"
			archiveArtifacts artifacts: "lib/libopus.${libc}.${arch}.a"
			commit.setBuildStatus("build-$libc-$arch", "SUCCESS", "Build succeeded.")
		} catch(err) {
			commit.setBuildStatus("build-$libc-$arch", "ERROR", "Build failed.")
			throw err
		}
	}
}

// Build go-opus (mobile) for ios
def build_go_opus_ios_macos(platform, arch) {
	
	stage("build-$platform-$arch") {
		try {
			checkout scm
			commit.setBuildStatus("build-$platform-$arch", "PENDING", "")
			sh "export CMAKE_TOOLCHAIN_FILE=${env.WORKSPACE}/build/toolchains/${platform}.${arch}.cmake && build/build.sh"
			dir("${env.WORKSPACE}"){
				sh "mv opus/build/libopus.a lib/libopus.${platform}.${arch}.a"
				archiveArtifacts artifacts: "lib/libopus.${platform}.${arch}.a"
				
			}
			
			commit.setBuildStatus("build-$platform-$arch", "SUCCESS", "Build succeeded.")
		} catch(err) {
			commit.setBuildStatus("build-$platform-$arch", "ERROR", "Build failed.")
			throw err
		}
	}
}

if (env.CHANGE_ID || env.TAG_NAME) {
	// building a PR or a tag or a debug build
	try {
		extraMessage = ""

		parallel ios: {
			node('ios') {
				build_go_opus_ios_macos('iphoneos','arm64')
				build_go_opus_ios_macos('iphonesimulator','x86_64')
				build_go_opus_ios_macos('macosx','arm64')
				build_go_opus_ios_macos('macosx','x86_64')
			}
		},
		linux: {
			node {
				sh '$(aws ecr get-login --region us-east-1 --no-include-email)'
				timeout(time: 30, unit: 'MINUTES') {
					ecrRegistry = "https://494415350827.dkr.ecr.us-east-1.amazonaws.com"
					docker.withRegistry("${ecrRegistry}") {
						docker.image("private/xcc-toolchains:20210417").inside('-u root') {
							try {
								checkout scm
								// Build go-opus binary for different platforms
								build_opus_linux("gnu", "x86_64")
								build_opus_linux("musl", "x86_64")
								build_opus_linux("musl", "aarch64")
								build_opus_linux("musl", "arm")
								build_opus_linux("android", "arm")
								build_opus_linux("android", "aarch64")
								build_opus_linux("android", "x86_64")
								
							} finally {
								// Change ownership of everything so it can be cleaned up later
								sh "chown -R 1000:1000 ."
								sh "rm -f ~/.netrc"
							}
						}
					}

					// these two git commands are intentional, to produce output in the jenkins console for debugging.
					sh "git describe --tags --dirty --always"
					sh "git diff"

				}
			}
		}, failFast: true
		mattermostSend channel: 'g-ci-notifications', color: 'good', message: "Build Successful - ${env.JOB_NAME} ${env.BUILD_NUMBER} (<${env.BUILD_URL}|Open>) ${extraMessage}"
	} catch (err) {
		mattermostSend channel: 'g-ci-notifications', color: 'danger', message: "Build Failed - ${env.JOB_NAME} ${env.BUILD_NUMBER} (<${env.BUILD_URL}|Open>)"
		throw err
	}
}


// Emacs configuration
// Local Variables:
// tab-width: 4
// indent-tabs-mode: t
// End:
