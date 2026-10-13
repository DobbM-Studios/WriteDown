import { markdownToHTML } from "./md.js";

const container = document.getElementById("container");

export function fromMARKDOWN(t) {
    container.innerHTML = markdownToHTML(t);   
}

window.wasmDone.then(() => fromMARKDOWN("**a**"));