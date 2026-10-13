import { markdownToHTML } from "./md.js";

const result = document.getElementById("htmlresult");
const lang   = document.getElementById("language-select");
const input  = document.getElementById("code");

export function fromMARKDOWN(t) {
    result.srcdoc = markdownToHTML(t);   
}

export function fromHTML(t) {
    result.srcdoc = t;
}

window.wasmDone.then(() => {
    input.addEventListener("input", () => {
        update();
    })
});

function update() {
    if (lang.value == "Markdown") {
        fromMARKDOWN(input.value);
    } else if (lang.value == "HTML") {
        fromHTML(input.value);
    }
}