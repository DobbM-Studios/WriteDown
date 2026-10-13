export function markdownToHTML(textoMarkdown) {
    if (typeof window.fromMD !== "function") {
        console.warn("Waiting WASM GO to load");
        return textoMarkdown;
    }

    const htmlResult = window.fromMD(textoMarkdown);

    return htmlResult;
}

window.markdownToHTML = markdownToHTML;
