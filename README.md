# Hogosuru brotli decoder


Hogosuru brotli decoder provide a live stream transform decoder wasm using brotli algorithm and can be used like the other stock CompressionStream (https://developer.mozilla.org/en-US/docs/Web/API/CompressionStream)


## Why using brotli decoder?

Wasm compiled with official go compiler is big. You can use tinygo to produce smaller files, but produce binary is more complicated to debug 
and TinyGo does not support finalizers (see wasm_exec.js) while the official GO version supported it.

I tried to find the best compromise between the two technologies: Use the official GO version and have the smallest file possible.

Also note that brotli is already supported by most of the browser, only if you set the accept encoding header is set and the webserver supported it (most of the time they not)

Hogosuru Brotli decoder allow you to use official GO wasm ,but with a high compression level thanks to hogosuru and brotli https://github.com/andybalholm/brotli.

## How it works?
hogosuru-brotli-decoder provide a stream transform decompression object and is compiled with tinygo to produce a smaller wasm. This wasm is also gzip to produce a 125ko file.

When hogosuru-brotli-decoder is ready , we can also download the official wasm encoded with brotli and decompress it in real time and execute it.


## Are we really a winner?

Yes! Of course we add an additional payload of 125ko , but the difference between a gzip wasm and a brotli wasm is is most of the time much higher than 125ko

A little example with the hogosuru-jwtdecoder (https://github.com/realPy/hogosuru-jwtdecoder)

|  compiler |  compression|  size | total payload|
|-------------|:--------------------:|:----------:|---------|
| official go | none | 5.4mo | 5.4mo | 
| official go | gzip | 1.3 mo | 1.3 mo | 
| official go | brotli | 948k | 1.04 mo | 

The bigger the project, the wider the gap

## How it works?


Just run  ./buildtinygo.sh to produce the brotlidecoder.wasm.gz or use the precompiled brotlidecoder.wasm.gz in docs/ 
copy the docs/wasm_boot.js in the same place. 

Replace the old official load wasm js script (in index.html)
```
 <script>


const wasmBrowserInstantiate = async (wasmModuleUrl, importObject) => {
  let response = undefined;

  // Check if the browser supports streaming instantiation
  if (WebAssembly.instantiateStreaming) {
    // Fetch the module, and instantiate it as it is downloading
    response = await WebAssembly.instantiateStreaming(
      fetch(wasmModuleUrl),
      importObject
    );
 
  } else {
    // Fallback to using fetch to download the entire module
    // And then instantiate the module
    const fetchAndInstantiateTask = async () => {
      const wasmArrayBuffer = await fetch(wasmModuleUrl).then(response =>
        response.arrayBuffer()
      );
      return WebAssembly.instantiate(wasmArrayBuffer, importObject);
    };
    response = await fetchAndInstantiateTask();
  }

  return response;
};

const go = new Go();
const runWasmAdd = async () => {
  // Get the importObject from the go instance.
  const importObject = go.importObject;

  // Instantiate our wasm module
  const wasmModule = await wasmBrowserInstantiate("jwtdecode.wasm", importObject);

  // Allow the wasm_exec go instance, bootstrap and execute our wasm module
  go.run(wasmModule.instance);

};
runWasmAdd();

        </script>
```

by this new one

```
        <script>




const wasmBrowserInstantiate = async (wasmModuleUrl, importObject, decoder) => {
  let response = undefined;

  // Check if the browser supports streaming instantiation
  if (WebAssembly.instantiateStreaming) {
    // Fetch the module, and instantiate it as it is downloading
    response = await WebAssembly.instantiateStreaming(
      fetch(wasmModuleUrl).then(
        (response) => new Response(response.body.pipeThrough(decoder), { status: 200,statusText:"OK",bodyUsed:false, encodeBody: 'manual',headers:{"content-type":"application/wasm"}})
      ),
      importObject
    );
 
  } else {
    // Fallback to using fetch to download the entire module
    // And then instantiate the module
    const fetchAndInstantiateTask = async () => {
      rsp = await fetch(wasmModuleUrl).then(response =>
      new Response(response.body.pipeThrough(decoder), { status: 200,statusText:"OK",bodyUsed:false,headers:{"content-type":"application/wasm"}})
      );
      wasmArrayBuffer=await rsp.arrayBuffer()
      return WebAssembly.instantiate(wasmArrayBuffer, importObject);
    };
    response = await fetchAndInstantiateTask();
  }

  return response;
};




const go = new Go();
const runWasmAdd = async () => {

  tinygo = new tinyGo();
  // Instantiate our wasm module
  const bootloader = await wasmBrowserInstantiate("brotlidecoder.wasm.gz", tinygo.importObject,new DecompressionStream("gzip"));
  
  tinygo.run(bootloader.instance);

  const wasmModule = await wasmBrowserInstantiate("jwtdecode.wasm.br", go.importObject, hogosuru_brotli_decoder);
  go.run(wasmModule.instance);

};
runWasmAdd();

        </script>
```


Compress your wasm with brotli and replace the name in the javascript code.

You can see result here :

https://realpy.github.io/hogosuru-brotlidecoder/index.html