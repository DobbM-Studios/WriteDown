import { markdownToHTML } from "./langs/md.js";
import { dmlToHTML } from "./langs/dml.js";

const result = document.getElementById("htmlresult");
const lang   = document.getElementById("language-select");
const input  = document.getElementById("code");

export function fromPLAIN(t) {
    result.srcdoc = "<pre>" + t + "</pre>";
}

export function fromMARKDOWN(t) {
    result.srcdoc = markdownToHTML(t);   
}

export function fromHTML(t) {
    result.srcdoc = t;
}

export function fromDOBBMML(t) {
    result.srcdoc = dmlToHTML(t);   
}


window.wasmDone.then(() => {
    input.addEventListener("input", () => {
        update();
    })
});

function update() {
    const val = input.value;

    switch (lang.value) {
        case "Markdown":
            fromMARKDOWN(val);
            break;
        case "HTML":
            fromHTML(val);
            break;
        case "Dobbm Markup":
            fromDOBBMML(val);
            break;
        case "Plain":
            fromPLAIN(val);
            break;
    }
}
