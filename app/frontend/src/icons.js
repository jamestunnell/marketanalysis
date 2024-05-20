import van from "vanjs-core"

const {i} = van.tags

const IconAdd = () => {
    return i({class: "fa-solid fa-plus"});
}

const IconCheck = () => {
    return i({class: "fa-solid fa-check"});
}

const IconClose = () => {
    return i({class: "fa-solid fa-xmark"});
}

const IconCollapsed = () => {
    return i({class: "fa-solid fa-chevron-right"});
}

const IconDelete = () => {
    return i({class: "fa-solid fa-trash"});
}

const IconDownload = () => {
    return i({class: "fa-solid fa-download"});
}

const IconError = () => {
    return i({class: "fa-solid fa-triangle-exclamation"});
}

const IconExpanded = () => {
    return i({class: "fa-solid fa-chevron-down"});
}

const IconExport = () => {
    return i({class: "fa-solid fa-file-export"});
}

const IconImport = () => {
    return i({class: "fa-solid fa-file-import"});
}

const IconPlay = () => {
    return i({class: "fa-solid fa-play"});
}

const IconPlot = () => {
    return i({class: "fa-solid fa-chart-line"});
}

const IconSave = () => {
    return i({class: "fa-solid fa-floppy-disk"});
}

const IconView = () => {
    return i({class: "fa-solid fa-eye"});
}

export {IconAdd, IconCheck, IconClose, IconCollapsed, IconDelete, IconDownload, IconError, IconExpanded, IconExport, IconImport, IconPlay, IconPlot, IconSave, IconView};