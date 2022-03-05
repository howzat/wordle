.PHONY: build-AdminBuildDictionary \
 clean

build-AdminBuildDictionary:
	env GOOS=linux go build -ldflags="-X main.CommitID=${CODEBUILD_SOURCE_VERSION} -s -w" -o "${ARTIFACTS_DIR}/admin-build-dictionary" github.com/howzat/wordle/cmd/admin-build-dictionary

build-Wordle:
	env GOOS=linux go build -ldflags="-X main.CommitID=${CODEBUILD_SOURCE_VERSION} -s -w" -o "${ARTIFACTS_DIR}/wordle" github.com/howzat/wordle/cmd/wordle

clean:
	rm -rf
