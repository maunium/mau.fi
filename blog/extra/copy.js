function copy(evt) {
    const codeBlock = evt.target.parentElement.querySelector("pre.chroma > code")
    navigator.clipboard.writeText(codeBlock.innerText).then(
        () => {
            evt.target.classList.add("success")
            setTimeout(() => evt.target.classList.remove("success"), 3000)
        },
        err => {
            console.error("Failed to copy code block to clipboard", err)
            evt.target.classList.add("failed")
            setTimeout(() => evt.target.classList.remove("failed"), 3000)
        }
    )
}

for (const elem of document.querySelectorAll("div.chroma")) {
    const copyButton = document.createElement("button")
    copyButton.classList.add("copy-button")
    copyButton.onclick = copy

    const parent = elem.parentNode
    const wrapper = document.createElement("div")
    wrapper.classList.add("chroma-parent")
    parent.replaceChild(wrapper, elem)
    wrapper.appendChild(elem)
    wrapper.appendChild(copyButton)
}
