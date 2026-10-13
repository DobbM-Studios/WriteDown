export function dmlToHTML(textoMarkdown) {
    if (typeof window.fromDML !== "function") {
        console.warn("Waiting WASM GO to load");
        return textoMarkdown;
    }

    const htmlResult = window.fromDML(textoMarkdown);

    return htmlResult;
}

window.markdownToHTML = markdownToHTML;
