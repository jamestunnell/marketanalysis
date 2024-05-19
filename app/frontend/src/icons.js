import van from "vanjs-core"

const {i} = van.tags

const IconAdd = () => {
    return i({class: "fa-solid fa-plus"});
}

const IconCheck = () => {
    return i({class: "fa-solid fa-check"});
}

const IconCollapsed = () => {
    return i({class: "fa-solid fa-chevron-right"});
}

const IconError = () => {
    return i({class: "fa-solid fa-triangle-exclamation"});
}

const IconExpanded = () => {
    return i({class: "fa-solid fa-chevron-down"});
}

const IconDelete = () => {
    return i({class: "fa-solid fa-trash"});
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

const IconSave = () => {
    return i({class: "fa-solid fa-floppy-disk"});
}

const IconView = () => {
    return i({class: "fa-solid fa-eye"});
}

export {IconAdd, IconCheck, IconCollapsed, IconDelete, IconError, IconExpanded, IconExport, IconImport, IconPlay, IconSave, IconView};