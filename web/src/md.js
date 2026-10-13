export function markdownToHTML(textoMarkdown) {
    if (typeof window.toHTML !== "function") {
        console.warn("Waiting WASM GO to load");
        return textoMarkdown;
    }

    const htmlResult = window.toHTML(textoMarkdown);

    return htmlResult;
}

window.markdownToHTML = markdownToHTML;
