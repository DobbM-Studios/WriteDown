const go = new Go();

window.wasmDone = WebAssembly.instantiateStreaming(fetch("../main.wasm"), go.importObject)
    .then((result) => {
        go.run(result.instance);
        console.log("GO's WASM loaded.");
    })
    .catch((err) => {
        console.error("WASM critical error:", err);
    });
