package main

import (
	"bytes"
	"io"

	"github.com/andybalholm/brotli"
	"github.com/realPy/hogosuru"
	"github.com/realPy/hogosuru/base/stream"
	"github.com/realPy/hogosuru/base/typedarray"
)

type BrotliEncoder struct {
	controller stream.TransformStreamDefaultController
	b          *bytes.Buffer
	reader     *brotli.Reader
}

func main() {
	var optimalGoBuffer [8 * 1024 * 1024]byte
	var outputdec [8 * 1024 * 1024]byte
	hogosuru.Init()
	var enc *BrotliEncoder

	ts, _ := stream.NewTransformStream(func(controller stream.TransformStreamDefaultController) {
		buffer := bytes.NewBuffer([]byte{})
		enc = &BrotliEncoder{controller: controller, b: buffer, reader: brotli.NewReader(buffer)}

	},
		func(chunk interface{}, controller stream.TransformStreamDefaultController) {

			if buffer, ok := chunk.(typedarray.Uint8Array); ok {
				size, _ := buffer.Length()

				buffer.CopyBytes(optimalGoBuffer[:size])
				enc.b.Write(optimalGoBuffer[:size])
				for {
					if n, err := enc.reader.Read(outputdec[:]); err == nil {

						decarray, _ := typedarray.NewUint8Array(n)
						decarray.CopyFromBytes(outputdec[:n])
						controller.Enqueue(decarray.BaseObject)

					} else {
						if err == io.EOF {
							break
						}
						break
					}

				}

			}
		},
		func(controller stream.TransformStreamDefaultController) {
			controller.Terminate()
		},
	)

	ts.Export("hogosuru_brotli_decoder")
	select {}

}
