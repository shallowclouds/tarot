module github.com/shallowclouds/tarot

go 1.18

require (
	github.com/disintegration/imaging v1.6.2
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.3
	golang.org/x/image v0.7.0
)

require github.com/shallowclouds/go-utils v0.0.0-00010101000000-000000000000

require (
	github.com/fogleman/gg v1.3.0
	github.com/sashabaranov/go-openai v1.9.4
	golang.org/x/sys v0.5.0 // indirect
)

replace github.com/shallowclouds/go-utils => /home/ubuntu/repos/go-utils
